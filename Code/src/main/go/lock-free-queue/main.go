package main

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type Queue interface {
	Enqueue(v interface{})
	Dequeue() interface{}
}

// LKQueue 是一个无锁的无限队列。
type LKQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

type node struct {
	value interface{}
	next  unsafe.Pointer
}

// NewLKQueue 返回一个空队列。
func NewLKQueue() *LKQueue {
	n := unsafe.Pointer(&node{})
	return &LKQueue{head: n, tail: n}
}

// Enqueue 将给定的值 v 放入队列的尾部。
//  在每次循环中，操作步骤如下：
//	1.	读取当前的尾节点 tail 及其后继节点 next。
//	2.	确认 tail 没有被其他线程修改。
//	3.	如果 next 为 nil，说明 tail 是最后一个节点，尝试将新节点 n 添加到 tail.next。
//	4.	如果添加成功，再尝试将 tail 指针移动到新节点 n。
//	5.	如果 next 不为 nil，说明其他线程已经添加了新节点，尝试将 tail 指针向前移动
func (q *LKQueue) Enqueue(v interface{}) {
	n := &node{value: v}
	for {
		tail := load(&q.tail)
		next := load(&tail.next)
		if tail == load(&q.tail) { // tail 和 next 是否一致？
			if next == nil {
				if cas(&tail.next, next, n) {
					cas(&q.tail, tail, n) // 入队完成。尝试将 tail 指向插入的节点
					return
				}
			} else { // tail 没有指向最后一个节点
				// 尝试将 tail 指向下一个节点
				cas(&q.tail, tail, next)
			}
		}
	}
}

// Dequeue 移除并返回队列头部的值。
// 如果队列为空，则返回零值。
// 在每次循环中，操作步骤如下：
//	1.	读取当前的头节点 head、尾节点 tail 及 head 的后继节点 next。
//	2.	确认 head 没有被其他线程修改。
//	3.	如果 head 与 tail 相同，且 next 为 nil，说明队列为空，返回零值。
//	4.	如果 next 不为 nil，说明有其他线程正在入队，但 tail 尚未更新，尝试将 tail 指针向前移动。
//	5.	如果 head 与 tail 不同，读取 next 的值 v，尝试将 head 指针移动到 next，如果成功，返回值 v。
func (q *LKQueue) Dequeue() interface{} {
	var t interface{}
	for {
		head := load(&q.head)
		tail := load(&q.tail)
		next := load(&head.next)
		if head == load(&q.head) { // head、tail 和 next 是否一致？
			if head == tail { // 队列是否为空或 tail 落后？
				if next == nil { // 队列是否为空？
					return t
				}
				// tail 落后。尝试将其前移
				cas(&q.tail, tail, next)
			} else {
				// 在 CAS 之前读取值，否则另一个出队操作可能会释放下一个节点
				v := next.value
				if cas(&q.head, head, next) {
					return v // 出队完成。返回
				}
			}
		}
	}
}

func load(p *unsafe.Pointer) (n *node) {
	return (*node)(atomic.LoadPointer(p))
}

func cas(p *unsafe.Pointer, old, new *node) (ok bool) {
	return atomic.CompareAndSwapPointer(
		p, unsafe.Pointer(old), unsafe.Pointer(new))
}

//=== slice-queue-start==========================================================
// SliceQueue is an unbounded queue which uses a slice as underlying.
type SliceQueue struct {
	data []interface{}
	mu   sync.Mutex
}

// NewSliceQueue returns an empty queue.
func NewSliceQueu(n int) (q *SliceQueue) {
	return &SliceQueue{data: make([]interface{}, 0, n)}
}

// Enqueue puts the given value v at the tail of the queue.
func (q *SliceQueue) Enqueue(v interface{}) {
	q.mu.Lock()
	q.data = append(q.data, v)
	q.mu.Unlock()
}

// Dequeue removes and returns the value at the head of the queue.
// It returns nil if the queue is empty.
func (q *SliceQueue) Dequeue() interface{} {
	var t interface{}
	q.mu.Lock()
	if len(q.data) == 0 {
		q.mu.Unlock()
		return t
	}
	v := q.data[0]
	q.data = q.data[1:]
	q.mu.Unlock()
	return v
}

//=== slice-queue-end==========================================================

func main() {
	q := NewLKQueue()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			q.Enqueue(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			_ = q.Dequeue()
		}
	}()
	wg.Wait()
}
