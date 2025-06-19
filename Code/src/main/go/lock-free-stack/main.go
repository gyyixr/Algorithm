package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

// node 是链表中的节点
type node struct {
	value interface{}
	next  *node
}

// LockFreeStack 是一个无锁栈结构体
type LockFreeStack struct {
	// top 指向栈顶元素，我们必须使用原子操作来访问它
	top unsafe.Pointer // 仍然指向 *node
}

// NewLockFreeStack 创建一个新的无锁栈
func NewLockFreeStack() *LockFreeStack {
	return &LockFreeStack{}
}

// loadTop 原子地加载栈顶节点。
// 这个辅助方法封装了 atomic.LoadPointer 和 unsafe.Pointer 的转换。
func (s *LockFreeStack) loadTop() *node {
	p := atomic.LoadPointer(&s.top)
	return (*node)(p)
}

// casTop 原子地执行比较并交换操作来更新栈顶。
// 这个辅助方法封装了 atomic.CompareAndSwapPointer 和 unsafe.Pointer 的转换。
func (s *LockFreeStack) casTop(old, new *node) bool {
	// 将 *node 类型的指针转换为 unsafe.Pointer 以供原子操作使用
	oldPtr := unsafe.Pointer(old)
	newPtr := unsafe.Pointer(new)
	return atomic.CompareAndSwapPointer(&s.top, oldPtr, newPtr)
}

// Push 将一个值压入栈顶
func (s *LockFreeStack) Push(value interface{}) {
	newNode := &node{value: value}
	for {
		// 1. 使用辅助方法读取当前的栈顶
		oldTop := s.loadTop()

		// 2. 将新节点的 next 指向旧的栈顶
		newNode.next = oldTop

		// 3. 使用辅助方法尝试进行 CAS 操作
		//    如果成功，说明栈顶已经被更新，循环结束。
		if s.casTop(oldTop, newNode) {
			return
		}
	}
}

// Pop 从栈顶弹出一个值
func (s *LockFreeStack) Pop() (interface{}, bool) {
	for {
		// 1. 使用辅助方法读取当前的栈顶
		oldTop := s.loadTop()

		// 2. 检查栈是否为空
		if oldTop == nil {
			return nil, false
		}

		// 3. 获取下一个节点作为新的栈顶
		// 增加nil检查以提高安全性
		newTop := oldTop.next

		// 4. 使用辅助方法尝试进行 CAS 操作
		// 再次确认oldTop仍然是当前栈顶
		if s.casTop(oldTop, newTop) {
			// 成功弹出，返回旧栈顶的值
			return oldTop.value, true
		}
		// 如果CAS失败，重试整个过程
	}
}

// IsEmpty 检查栈是否为空
func (s *LockFreeStack) IsEmpty() bool {
	return s.loadTop() == nil
}

func main() {
	// --- 主函数的测试逻辑保持不变 ---
	stack := NewLockFreeStack()
	var wg sync.WaitGroup

	// 测试并发 Push
	// 启动 100 个 goroutine，每个 goroutine 推入 100 个元素
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				stack.Push(fmt.Sprintf("g%d-v%d", i, j))
			}
		}(i)
	}
	wg.Wait()

	count := 0
	results := make(map[interface{}]bool)
	for {
		val, ok := stack.Pop()
		if !ok {
			break
		}
		results[val] = true
		count++
	}

	fmt.Printf("总共推入了 100 * 100 = 10000 个元素\n")
	fmt.Printf("成功弹出了 %d 个元素\n", count)
	if len(results) == 10000 && count == 10000 {
		fmt.Println("测试通过：所有推入的元素都被成功弹出，且没有重复。")
	} else {
		fmt.Println("测试失败！")
	}

	_, ok := stack.Pop()
	fmt.Printf("在空栈上执行 Pop 操作，成功: %v\n", ok)
	fmt.Printf("栈现在是否为空: %v\n", stack.IsEmpty())
}
