package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

// 任务结构体
type Task struct {
	Name string
	Job  func() error
}

// 任务执行结果结构体
type TaskResult struct {
	TaskName  string        `json:"taskName"`
	StartTime time.Time     `json:"startTime"`
	Duration  time.Duration `json:"duration"`
	Err       error         `json:"err`
}

// 任务调度器结构体
type Scheduler struct {
	taskQueue chan Task      // 任务队列
	workers   int            // 工作协程数量
	wg        sync.WaitGroup // 用于等待所有任务完成
}

// 创建调度器
func newScheduler(workers int, queueSize int) *Scheduler {
	return &Scheduler{
		taskQueue: make(chan Task, queueSize),
		workers:   workers,
	}
}

// 添加任务到调度器
func (s *Scheduler) addTask(name string, job func() error) {
	s.taskQueue <- Task{
		Name: name,
		Job:  job,
	}
	s.wg.Add(1) // 每添加一个任务，等待组计数加1
}

// 开始调度任务
func (s *Scheduler) start(resultCh chan<- TaskResult) {
	for i := 0; i < s.workers; i++ {
		go func() {
			// taskQueue 通道内没数据时，会阻塞
			for task := range s.taskQueue {
				start := time.Now()
				err := s.run(task)

				resultCh <- TaskResult{
					TaskName:  task.Name,
					StartTime: start,
					Duration:  time.Since(start),
					Err:       err,
				}
			}
		}()
	}
}

// 执行任务，并返回执行结果
func (s *Scheduler) run(task Task) error {
	defer s.wg.Done() // 任务完成后，等待组计数减1
	fmt.Printf("开始执行任务: %s\n", task.Name)
	start := time.Now()
	if err := task.Job(); err == nil {
		fmt.Printf("任务[%s]执行完毕，耗时：%s\n", task.Name, time.Since(start))
		return nil
	} else {
		fmt.Printf("任务[%s]执行异常：%s，耗时：%s\n", task.Name, err, time.Since(start))
		return err
	}
}

// 等待所有任务完成
func (s *Scheduler) wait() {
	close(s.taskQueue) // 关闭任务队列，表示不再添加新任务
	s.wg.Wait()
}

func main() {
	fmt.Println("---------------------------两个goroutine交替打印奇数和偶数---------------------------")
	startTwoGoroutine()

	fmt.Println("---------------------------任务调度演示---------------------------")
	var resultWg sync.WaitGroup
	resultCh := make(chan TaskResult, 10)
	scheduler := newScheduler(3, 10)
	scheduler.start(resultCh)

	scheduler.addTask("任务1", func() error {
		time.Sleep(1 * time.Second)
		fmt.Println("任务1完成")
		return nil
	})

	scheduler.addTask("任务2", func() error {
		time.Sleep(2 * time.Second)
		fmt.Println("任务2完成")
		return nil
	})

	scheduler.addTask("任务3", func() error {
		time.Sleep(time.Second * 3)
		fmt.Println("任务3异常")
		return errors.New("异常错误")
	})

	resultWg.Add(3)
	go func() {
		for result := range resultCh {
			jsonData, _ := json.Marshal(result)
			fmt.Println("处理结果：", string(jsonData))
			resultWg.Done()
		}
	}()

	resultWg.Wait()
	scheduler.wait()
	fmt.Println("所有任务完成")
}

func startTwoGoroutine() {
	var wg sync.WaitGroup
	wg.Add(2) // 添加需要等待的goroutine数量
	go func() {
		defer wg.Done() // 当前goroutine执行完毕后调用，同时减1
		f1()
	}()
	go func() {
		defer wg.Done()
		f2()
	}()
	wg.Wait() // 等待所有goroutine执行完毕（计数器减为0）
	fmt.Println("所有goroutine已完成")
}

func f1() {
	for i := 0; i <= 10; i++ {
		if i%2 != 0 {
			fmt.Println("奇数:", i)
			time.Sleep(100 * time.Millisecond) // 添加延时以更明显地看到输出交替
		}
	}
}

func f2() {
	for i := 0; i <= 10; i++ {
		if i%2 == 0 {
			fmt.Println("偶数:", i)
			time.Sleep(100 * time.Millisecond) // 添加延时以更明显地看到输出交替
		}
	}
}
