package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// Node 定义了链表中的每个节点
type Node struct {
	value int
	next  *Node
	lock  sync.Mutex // 每个节点都有一把自己的锁
}

// SortedList 定义了有序链表
type SortedList struct {
	head *Node // 链表的哨兵头节点
}

// NewSortedList 创建一个新的有序链表，并初始化哨兵头节点
func NewSortedList() *SortedList {
	// 创建一个哨兵头节点，其值设为整型最小值，确保它总是在最前面
	head := &Node{value: math.MinInt, next: nil}
	return &SortedList{head: head}
}

// Insert 将一个值并发安全地插入到有序链表中
func (l *SortedList) Insert(value int) {
	newNode := &Node{value: value}

	// 使用一个无限循环来处理因并发冲突导致的重试
	for {
		// 1. 乐观遍历：不加锁，找到插入位置
		prev := l.head
		curr := prev.next

		for curr != nil && curr.value < value {
			prev = curr
			curr = curr.next
		}

		// 2. 锁定关键节点 prev
		prev.lock.Lock()

		// 3. 验证：检查 prev 和 curr 的关系是否在遍历后被其他goroutine改变
		if prev.next == curr {
			// 验证成功，可以安全插入
			newNode.next = curr
			prev.next = newNode

			// 操作完成，解锁并退出循环
			prev.lock.Unlock()
			return
		}

		// 4. 验证失败：解锁并重试
		// 如果 prev.next 不等于 curr，说明在我们遍历和加锁之间有其他goroutine插入了节点。
		// 我们必须释放锁，然后从头开始重试整个过程。
		prev.lock.Unlock()
		// continue 关键字会使 for 循环进入下一次迭代
	}
}

// Display 打印链表内容（仅用于调试和演示）
// 注意：在一个高并发系统中，Display本身也需要考虑并发安全问题，
// 但为了突出Insert的逻辑，这里我们简化处理。
func (l *SortedList) Display() {
	fmt.Print("List: head -> ")
	curr := l.head.next
	for curr != nil {
		fmt.Printf("%d -> ", curr.value)
		curr = curr.next
	}
	fmt.Println("nil")
}

func main() {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	list := NewSortedList()
	var wg sync.WaitGroup

	// 模拟10个goroutine并发插入100个随机数
	numGoroutines := 10
	numInsertsPerGoroutine := 10
	totalInserts := numGoroutines * numInsertsPerGoroutine

	wg.Add(totalInserts)

	fmt.Printf("启动 %d 个 goroutine，每个 goroutine 插入 %d 个随机数...\n", numGoroutines, numInsertsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numInsertsPerGoroutine; j++ {
				value := rand.Intn(200) // 插入0-199之间的随机数
				list.Insert(value)
				wg.Done()
			}
		}()
	}

	// 等待所有插入操作完成
	wg.Wait()

	fmt.Println("所有插入操作已完成。")

	// 显示最终的有序链表
	list.Display()

	// 验证链表是否有序且长度正确
	count := 0
	curr := list.head.next
	prev := list.head
	isSorted := true
	for curr != nil {
		if prev != list.head && prev.value > curr.value {
			isSorted = false
		}
		prev = curr
		curr = curr.next
		count++
	}

	fmt.Printf("链表节点总数: %d (预期: %d)\n", count, totalInserts)
	fmt.Printf("链表是否有序: %t\n", isSorted)
}
