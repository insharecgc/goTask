package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	fmt.Println("---------------------------mutex加锁增加计数器---------------------------")
	mutexDemo()

	fmt.Println("---------------------------atomic原子操作增加计数器---------------------------")
	atomicDemo()
}

func mutexDemo() {
	// 共享计数器
	var mutexCount = 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 每个goroutine调用1000次自增
			for i := 0; i < 1000; i++ {
				addOne(&mutexCount, &mu)
			}
		}()
	}

	wg.Wait()
	fmt.Println("10个goroutine使用Mutex锁各自自增1000次后count值为：", mutexCount)
}

func addOne(mutexCount *int, mu *sync.Mutex) {
	mu.Lock()         // 对整个代码块加锁
	defer mu.Unlock() // 确保执行完后结束
	*mutexCount++
}

func atomicDemo() {
	// 原子操作的变量必须是int32/int64/uint32/uint64等类型
	var count int64 = 0
	var wg sync.WaitGroup
	var goNum = 10
	wg.Add(goNum)
	for i := 0; i < goNum; i++ {
		go func(){
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				// 原子操作：count += 1（等价于 count++，但线程安全）
				atomic.AddInt64(&count, 1)
			}
		}()
	}

	wg.Wait()
	fmt.Println("10个goroutine使用atomic原子操作自增1000次后count值为：", count)
}
