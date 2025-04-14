package main

import "fmt"

// ListNode 链表节点定义
type ListNode struct {
	Val  int
	Next *ListNode
}

// 反转整个链表
func reverse(a *ListNode) *ListNode {
	var pre, cur, nxt *ListNode
	pre = nil
	cur = a
	nxt = nil
	for cur != nil {
		nxt = cur.Next
		// 逐个节点反转
		cur.Next = pre // 更新指针位置
		pre = cur
		cur = nxt
	}
	return pre
	// 返回反转后的头节点
}

// 反转区间 [a, b) 的元素，左闭右开
func reverseRange(a, b *ListNode) *ListNode {
	var pre, cur, nxt *ListNode
	pre = nil
	cur = a
	nxt = nil
	// while 终止的条件改为 cur != b
	for cur != b {
		nxt = cur.Next
		cur.Next = pre
		pre = cur
		cur = nxt
	}
	// 返回反转后的头节点
	return pre
}

// K个一组反转链表
func reverseKGroup(head *ListNode, k int) *ListNode {
	if head == nil {
		return nil
	}
	// 区间 [a, b) 包含 k 个待反转元素
	a, b := head, head
	for i := 0; i < k; i++ {
		// 不足 k 个，不需要反转，base case
		if b == nil {
			return head
		}
		b = b.Next
	}
	// 反转前 k 个元素
	newHead := reverseRange(a, b)
	// 递归反转后续链表并连接起来
	a.Next = reverseKGroup(b, k)
	return newHead
}

// 辅助函数：创建链表
func createList(nums []int) *ListNode {
	if len(nums) == 0 {
		return nil
	}
	head := &ListNode{Val: nums[0]}
	current := head
	for i := 1; i < len(nums); i++ {
		current.Next = &ListNode{Val: nums[i]}
		current = current.Next
	}
	return head
}

// 辅助函数：打印链表
func printList(head *ListNode) {
	for head != nil {
		fmt.Printf("%d", head.Val)
		if head.Next != nil {
			fmt.Print(" -> ")
		}
		head = head.Next
	}
	fmt.Println()
}

func main() {
	// 示例测试
	nums := []int{1, 2, 3, 4, 5}
	k := 2

	head := createList(nums)
	fmt.Print("原始链表: ")
	printList(head)

	reversed := reverseKGroup(head, k)
	fmt.Printf("每%d个一组反转后: ", k)
	printList(reversed)

	// 另一个测试用例
	nums2 := []int{1, 2, 3, 4, 5, 6, 7, 8}
	k2 := 3

	head2 := createList(nums2)
	fmt.Print("\n原始链表: ")
	printList(head2)

	reversed2 := reverseKGroup(head2, k2)
	fmt.Printf("每%d个一组反转后: ", k2)
	printList(reversed2)
}
