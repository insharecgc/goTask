package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("---------------------------不带缓冲区的通道---------------------------")
	chanDemo1()

	fmt.Println("---------------------------带缓冲区的通道---------------------------")
	chanDemo2()

	fmt.Println("---------------------------一个生产者，多个消费者---------------------------")
	chanDemo3()
}

// 不带缓冲的通道通信
func chanDemo1() {
	// 定义一个不带缓冲的通道，无缓冲通道的发送和接收操作是同步的。也就是说，发送操作会阻塞置顶另一个goroutine执行接收操作，反之依然
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		// 往通道写
		for i := 1; i <= 10; i++ {
			ch <- i
		}
		// 通道写完内容后，就需要关闭，不然读通道中通过for range就永远不会跳出循环，就会导致死锁
		close(ch)
	}()

	go func() {
		defer wg.Done()
		// 从通道中读
		for val := range ch {
			fmt.Printf("从通道中读取的内容：%d\n", val)
		}
	}()
	wg.Wait()
}

// 带缓冲的通道通信
func chanDemo2() {
	// 定义一个带缓冲的通道，当通道已满，写会阻塞；当通道为空，读会阻塞
	ch := make(chan int, 10)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			ch <- i
		}
		// 不再往通道写内容后，要关闭，不然会死锁
		close(ch)
	}()

	go func() {
		defer wg.Done()
		for val := range ch {
			fmt.Println("读取的内容：", val)
		}
	}()
	wg.Wait()
}

// 一个生产者，多个消费者
func chanDemo3() {
	const (
		totalTaskNum = 100 // 100个任务数据
		bufferSize   = 10  // 缓冲区大小
		consumerNum  = 3   // 消费者数量
	)
	var wg sync.WaitGroup
	ch := make(chan int, bufferSize)
	// 启动一个生产者
	wg.Add(1)
	go producer(ch, &wg, totalTaskNum)
	for i := 0; i < consumerNum; i++ {
		wg.Add(1)
		go consumer(ch, &wg, i)
	}

	wg.Wait()
	fmt.Println("程序已执行完成")
}

func producer(ch chan <- int, wg *sync.WaitGroup, totalTaskNum int) {
	defer wg.Done()
	for i := 0; i < totalTaskNum; i++ {
		ch <- i
		fmt.Println("生产者已生成一个数据：", i)
		time.Sleep(100 * time.Millisecond)
	}
	close(ch)
}

func consumer(ch <- chan int, wg *sync.WaitGroup, no int) {
	defer wg.Done()
	for val := range ch {
		fmt.Printf("第%d个消费者，正在读取数据：%d\n", no, val)
		time.Sleep(200 * time.Millisecond)
	}
}
