package main

import (
	"fmt"
	"slices"
)

type TreeNode struct {
	Val   int
	Right *TreeNode
	Left  *TreeNode
}

func rob(root *TreeNode) int {
	res := robTree(root)
	return slices.Max(res)
}

func robTree(cur *TreeNode) []int {
	if cur == nil {
		return []int{0, 0}
	}
	// 后序遍历
	left := robTree(cur.Left)
	right := robTree(cur.Right)

	// 考虑去偷当前的屋子
	robCur := cur.Val + left[0] + right[0]
	// 考虑不去偷当前的屋子
	notRobCur := slices.Max(left) + slices.Max(right)

	// 注意顺序：0:不偷，1:去偷
	return []int{notRobCur, robCur}
}

func main() {

	// 构建测试用例的二叉树
	//      3
	//     / \
	//    2   3
	//     \   \
	//      3   1
	root := &TreeNode{Val: 3}
	root.Left = &TreeNode{Val: 2}
	root.Right = &TreeNode{Val: 3}
	root.Left.Right = &TreeNode{Val: 3}
	root.Right.Right = &TreeNode{Val: 1}

	// 调用rob函数并打印结果
	maxMoney := rob(root)
	fmt.Printf("可以偷窃的最大金额是：%d\n", maxMoney)
}
