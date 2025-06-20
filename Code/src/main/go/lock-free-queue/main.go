package main

import (
	"sync"
	"sync/atomic"
)

type Queue interface {
	Enqueue(v interface{})
	Dequeue() interface{}
}

// LKQueue 是一个无锁的无限队列。
type LKQueue struct {
	head atomic.Pointer[node]
	tail atomic.Pointer[node]
}

type node struct {
	value interface{}
	next  atomic.Pointer[node]
}

// NewLKQueue 返回一个空队列。
func NewLKQueue() *LKQueue {
	n := &node{}
	q := &LKQueue{}
	q.head.Store(n)
	q.tail.Store(n)
	return q
}

// Enqueue 将给定的值 v 放入队列的尾部。
//
// 在每次循环中，操作步骤如下：
// 1.读取当前的尾节点 tail 及其后继节点 next。
// 2.确认 tail 没有被其他线程修改。
// 3.如果 next 为 nil，说明 tail 是最后一个节点，尝试将新节点 n 添加到 tail.next。
// 4.如果添加成功，再尝试将 tail 指针移动到新节点 n。
// 5.如果 next 不为 nil，说明其他线程已经添加了新节点，尝试将 tail 指针向前移动
func (q *LKQueue) Enqueue(v interface{}) {
	n := &node{value: v}
	for {
		tail := q.tail.Load()
		next := tail.next.Load()
		if tail == q.tail.Load() { // tail 和 next 是否一致？
			if next == nil {
				if tail.next.CompareAndSwap(next, n) {
					q.tail.CompareAndSwap(tail, n) // 入队完成。尝试将 tail 指向插入的节点
					return
				}
			} else { // tail 没有指向最后一个节点
				// 尝试将 tail 指向下一个节点
				q.tail.CompareAndSwap(tail, next)
			}
		}
	}
}

// Dequeue 移除并返回队列头部的值。
// 如果队列为空，则返回零值。
// 在每次循环中，操作步骤如下：
//  1. 读取当前的头节点 head、尾节点 tail 及 head 的后继节点 next。
//  2. 确认 head 没有被其他线程修改。
//  3. 如果 head 与 tail 相同，且 next 为 nil，说明队列为空，返回零值。
//  4. 如果 next 不为 nil，说明有其他线程正在入队，但 tail 尚未更新，尝试将 tail 指针向前移动。
//  5. 如果 head 与 tail 不同，读取 next 的值 v，尝试将 head 指针移动到 next，如果成功，返回值 v。
func (q *LKQueue) Dequeue() interface{} {
	var t interface{}
	for {
		head := q.head.Load()
		tail := q.tail.Load()
		next := head.next.Load()
		if head == q.head.Load() { // head、tail 和 next 是否一致？
			if head == tail { // 队列是否为空或 tail 落后？
				if next == nil { // 队列是否为空？
					return t
				}
				// tail 落后。尝试将其前移
				q.tail.CompareAndSwap(tail, next)
			} else {
				// 在 CAS 之前读取值，否则另一个出队操作可能会释放下一个节点
				v := next.value
				if q.head.CompareAndSwap(head, next) {
					return v // 出队完成。返回
				}
			}
		}
	}
}

//=== slice-queue-end==========================================================

func main() {
	// 测试基本功能
	testBasicOperations()

	// 测试并发安全性
	testConcurrency()
}

func testBasicOperations() {
	println("=== 测试基本操作 ===")
	q := NewLKQueue()

	// 测试空队列
	if v := q.Dequeue(); v != nil {
		println("错误：空队列应该返回 nil")
	}

	// 测试入队出队
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	if v := q.Dequeue(); v != 1 {
		println("错误：期望 1，得到", v)
	}
	if v := q.Dequeue(); v != 2 {
		println("错误：期望 2，得到", v)
	}
	if v := q.Dequeue(); v != 3 {
		println("错误：期望 3，得到", v)
	}

	// 测试队列为空
	if v := q.Dequeue(); v != nil {
		println("错误：队列应该为空")
	}

	println("基本操作测试通过")
}

func testConcurrency() {
	println("=== 测试并发安全性 ===")
	q := NewLKQueue()
	var wg sync.WaitGroup

	// 生产者
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			q.Enqueue(i)
		}
	}()

	// 消费者
	wg.Add(1)
	go func() {
		defer wg.Done()
		count := 0
		for count < 100000 {
			if v := q.Dequeue(); v != nil {
				count++
			}
		}
	}()

	wg.Wait()
	println("并发测试完成")

	// 验证队列最终为空
	if v := q.Dequeue(); v != nil {
		println("错误：最终队列应该为空")
	} else {
		println("并发安全性测试通过")
	}
}
