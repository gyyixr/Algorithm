package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Node 定义了链表中的节点
type Node struct {
	value int
	next  atomic.Pointer[Node] // 指向下一个节点的原子指针
}

// ConcurrentLinkedList 定义了一个并发安全的链表
type ConcurrentLinkedList struct {
	head atomic.Pointer[Node] // 指向头节点的原子指针
}

// NewConcurrentLinkedList 创建一个新的并发链表实例
func NewConcurrentLinkedList() *ConcurrentLinkedList {
	return &ConcurrentLinkedList{}
}

// Insert 使用 CAS 操作在链表头部并发安全地插入一个新节点
func (l *ConcurrentLinkedList) Insert(value int) {
	newNode := &Node{value: value}

	for {
		// 1. 原子地读取当前的头节点
		oldHead := l.head.Load()

		// 2. 将新节点的 next 指向旧的头节点
		newNode.next.Store(oldHead)

		// 3. 尝试使用 CAS 将头节点从 oldHead 更新为 newNode
		// 如果在 l.head.Load() 和此行代码之间有其他 goroutine 修改了 l.head,
		// 那么 l.head 的当前值就不再是 oldHead，CAS 操作会失败。
		// 此时，循环会继续，我们会获取最新的 head 值重试。
		if l.head.CompareAndSwap(oldHead, newNode) {
			// 如果 CAS 成功，意味着我们成功地插入了节点，可以退出循环
			return
		}
	}
}

// Display 打印链表的所有节点值（这是一个快照）
func (l *ConcurrentLinkedList) Display() {
	fmt.Print("List: ")
	curr := l.head.Load()
	count := 0
	for curr != nil {
		fmt.Printf("%d -> ", curr.value)
		curr = curr.next.Load()
		count++
	}
	fmt.Println("nil")
	fmt.Printf("Total nodes: %d\n", count)
}

func main() {
	// 创建一个并发链表
	cll := NewConcurrentLinkedList()

	// 使用 WaitGroup 等待所有 goroutine 完成
	var wg sync.WaitGroup

	// 定义并发数量
	numGoroutines := 100
	nodesPerGoroutine := 10

	// 启动多个 goroutine 并发地插入节点
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < nodesPerGoroutine; j++ {
				value := goroutineID*100 + j
				cll.Insert(value)
			}
		}(i)
	}

	// 等待所有插入操作完成
	wg.Wait()

	fmt.Println("All goroutines have finished insertion.")

	// 显示最终的链表内容和节点总数
	cll.Display()

	// 验证节点总数
	expectedCount := numGoroutines * nodesPerGoroutine
	fmt.Printf("Expected node count: %d\n", expectedCount)
}
