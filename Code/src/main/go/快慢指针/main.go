package main

import "fmt"

// ListNode defines a node in a singly linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

// MiddleLeftNode returns the middle node of the list.
// If the number of nodes is even, it returns the right middle node.
// 偶数时返回右中间节点，奇数时返回中间节点。
func MiddleRightNode(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}
	slow := head
	fast := head

	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	return slow
}

// MiddleRightNode returns the middle node of the list.
// If the number of nodes is even, it returns the left middle node.
// 偶数时返回左中间节点，奇数时返回中间节点。
func MiddleLeftNode(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}
	slow := head
	fast := head

	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	return slow
}

func main() {
	root := &ListNode{Val: 1}
	root.Next = &ListNode{Val: 2}
	root.Next.Next = &ListNode{Val: 3}
	root.Next.Next.Next = &ListNode{Val: 4}

	fmt.Println(MiddleLeftNode(root).Val)
	fmt.Println(MiddleRightNode(root).Val)

	//Uncomment to test with 5 nodes
	root.Next.Next.Next.Next = &ListNode{Val: 5}
	fmt.Println(MiddleLeftNode(root).Val)
	fmt.Println(MiddleRightNode(root).Val)
}
