package main

import (
	"fmt"
	"math"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func maxPathSum(root *TreeNode) int {
	maxSum := math.MinInt32

	var dfs func(root *TreeNode) int
	dfs = func(root *TreeNode) int {
		if root == nil {
			return 0
		}
		left := dfs(root.Left)
		right := dfs(root.Right)

		innerMaxSum := left + root.Val + right
		maxSum = max(maxSum, innerMaxSum)
		outputMaxSum := root.Val + max(left, right) // left,right都是非负的，就不用和0比较了
		return max(outputMaxSum, 0)
	}

	dfs(root)
	return maxSum
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	root := &TreeNode{Val: -10}
	root.Right = &TreeNode{Val: 20}
	root.Left = &TreeNode{Val: 9}
	root.Right.Right = &TreeNode{Val: 7}
	root.Right.Left = &TreeNode{Val: 15}
	fmt.Println(maxPathSum(root))
}
