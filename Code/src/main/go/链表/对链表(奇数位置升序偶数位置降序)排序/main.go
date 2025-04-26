package main

import "fmt"

// ListNode 定义链表节点结构
type ListNode struct {
	Val  int
	Next *ListNode
}

// splitList 分离链表为奇数位和偶数位
func splitList(head *ListNode) (*ListNode, *ListNode) {
	var oddHead, oddTail *ListNode   // 奇数位链表的头和尾
	var evenHead, evenTail *ListNode // 偶数位链表的头和尾
	current := head
	isOdd := true // 用布尔值代替索引，更直观表示奇偶

	for current != nil {
		if isOdd {
			if oddHead == nil {
				oddHead = current // 初始化奇数链表
			} else {
				oddTail.Next = current // 添加到奇数链表尾部
			}
			oddTail = current // 更新奇数链表尾部
		} else {
			if evenHead == nil {
				evenHead = current // 初始化偶数链表
			} else {
				evenTail.Next = current // 添加到偶数链表尾部
			}
			evenTail = current // 更新偶数链表尾部
		}

		current = current.Next // 移动到下一个节点
		isOdd = !isOdd         // 切换奇偶状态
	}

	// 切断两个链表的尾部连接
	if oddTail != nil {
		oddTail.Next = nil
	}
	if evenTail != nil {
		evenTail.Next = nil
	}

	return oddHead, evenHead
}

// reverseList 反转链表
// 1->2->3->4->nil
func reverseList(head *ListNode) *ListNode {
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

// mergeLists 合并两个升序链表
// 应该始终保留 dummy 的头节点，并用另一个指针（如 tail）来遍历链表
func mergeLists(l1, l2 *ListNode) *ListNode {
	dummy := &ListNode{Val: -1} // 虚拟头节点
	tail := dummy               // tail 用于构建新链表

	for l1 != nil && l2 != nil {
		if l1.Val < l2.Val {
			tail.Next = l1
			l1 = l1.Next
		} else {
			tail.Next = l2
			l2 = l2.Next
		}
		tail = tail.Next // 移动 tail 到新链表的末尾
	}

	// 将剩余部分直接接上去
	if l1 != nil {
		tail.Next = l1
	} else {
		tail.Next = l2
	}

	return dummy.Next // 返回虚拟头节点的下一个节点
}

// sortLinkedList 主函数：对链表进行排序
func sortLinkedList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	// 保存 splitList 的结果
	oddHead, evenHead := splitList(head)

	// 反转偶数位链表
	evenHead = reverseList(evenHead)

	// 合并两个升序链表
	return mergeLists(oddHead, evenHead)
}

// printList 打印链表
func printList(head *ListNode) {
	for head != nil {
		fmt.Printf("%d -> ", head.Val)
		head = head.Next
	}
	fmt.Println("nil")
}

// 测试代码
func main() {
	head := &ListNode{Val: 1}
	head.Next = &ListNode{Val: 8}
	head.Next.Next = &ListNode{Val: 3}
	head.Next.Next.Next = &ListNode{Val: 6}
	head.Next.Next.Next.Next = &ListNode{Val: 5}
	head.Next.Next.Next.Next.Next = &ListNode{Val: 4}
	head.Next.Next.Next.Next.Next.Next = &ListNode{Val: 7}

	fmt.Print("Original List: ")
	printList(head)

	sortedHead := sortLinkedList(head)

	fmt.Print("Sorted List: ")
	printList(sortedHead)
}
