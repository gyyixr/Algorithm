package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// 使用一个无缓冲通道来进行协程间的同步
	ch := make(chan struct{}, 1)
	var wg sync.WaitGroup

	// 启动两个goroutine，一个打印奇数，一个打印偶数
	wg.Add(2)
	// 向通道发送初始信号，让奇数goroutine先执行
	ch <- struct{}{}
	go printOdd(ch, &wg)
	go printEven(ch, &wg)

	// 等待两个goroutine执行完毕
	wg.Wait()
}

func printOdd(ch chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 10; i += 2 {
		// 等待通道接收到信号
		<-ch
		fmt.Println("Odd:", i)
		// 发送信号给偶数goroutine
		ch <- struct{}{}
	}
}

func printEven(ch chan struct{}, wg *sync.WaitGroup) {
	time.Sleep(1 * time.Millisecond)
	defer wg.Done()
	for i := 2; i <= 10; i += 2 {
		// 等待通道接收到信号
		<-ch
		fmt.Println("Even:", i)
		// 发送信号给奇数goroutine
		ch <- struct{}{}
	}
}
