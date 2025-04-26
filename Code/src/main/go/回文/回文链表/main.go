package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

// 方法二，快慢指针
func isPalindrome(head *ListNode) bool {
	if head == nil && head.Next == nil {
		return true
	}

	// 1->2>3->2>1
	// 1->2->2>1
	// 偶数时，慢指针1->2,快指针反转后1->2
	// 奇数时，慢指针1->2,快指针反转后1->2->3,快指针多一个元素时最中间的元素，不用管
	//慢指针，找到链表中间分位置，作为分割
	slow := head
	fast := head
	//记录慢指针的前一个节点，用来分割链表
	pre := head
	for fast != nil && fast.Next != nil {
		pre = slow
		slow = slow.Next
		fast = fast.Next.Next
	}
	//分割链表
	pre.Next = nil
	//前半部分
	cur1 := head
	//反转后半部分，总链表长度如果是奇数，cur2比cur1多一个节点
	cur2 := ReverseList(slow)

	//开始两个链表的比较
	for cur1 != nil {
		if cur1.Val != cur2.Val {
			return false
		}
		cur1 = cur1.Next
		cur2 = cur2.Next
	}
	return true
}

//反转链表
func ReverseList(head *ListNode) *ListNode {
	var pre *ListNode
	cur := head
	for cur != nil {
		tmp := cur.Next
		cur.Next = pre
		pre = cur
		cur = tmp
	}
	return pre
}

func main() {
	//	测试isPalindrome
	head := &ListNode{
		Val: 1,
		Next: &ListNode{
			Val: 2,
			Next: &ListNode{
				Val: 3,
				Next: &ListNode{
					Val: 2,
					Next: &ListNode{
						Val:  1,
						Next: nil,
					},
				},
			},
		},
	}

	fmt.Println(isPalindrome(head))

}
