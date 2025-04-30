package main

import "fmt"

/**
给定一个单链表 L 的头节点 head ，单链表 L 表示为：

L0 → L1 → … → Ln - 1 → Ln
请将其重新排列后变为：

L0 → Ln → L1 → Ln - 1 → L2 → Ln - 2 → …
不能只是单纯的改变节点内部的值，而是需要实际的进行节点交换。
*/

// 方法一 数组模拟
/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func reorderList(head *ListNode) {
	vec := make([]*ListNode, 0)
	cur := head
	if cur == nil {
		return
	}
	for cur != nil {
		vec = append(vec, cur)
		cur = cur.Next
	}
	cur = head
	i := 1
	j := len(vec) - 1 // i j为前后的双指针
	count := 0        // 计数，偶数取后面，奇数取前面
	for i <= j {
		if count%2 == 0 {
			cur.Next = vec[j]
			j--
		} else {
			cur.Next = vec[i]
			i++
		}
		cur = cur.Next
		count++
	}
	cur.Next = nil // 注意结尾，把链表节点放入数组中时，每个node的Next指针不为空，所以需要单独处理
}

func main() {
	head := &ListNode{Val: 1}
	head.Next = &ListNode{Val: 2}
	head.Next.Next = &ListNode{Val: 3}
	head.Next.Next.Next = &ListNode{Val: 4}
	head.Next.Next.Next.Next = &ListNode{Val: 5}
	printList(head)
	reorderList(head)
	printList(head)
}

// ListNode 定义链表节点结构
type ListNode struct {
	Val  int
	Next *ListNode
}

// 打印链表
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
