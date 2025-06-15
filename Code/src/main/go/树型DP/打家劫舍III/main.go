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
	fmt.Println(rob1(root))
}

// 纯递归解法
var umap = make(map[*TreeNode]int)

func rob1(root *TreeNode) int {
	if root == nil {
		return 0
	}
	if root.Left == nil && root.Right == nil {
		return root.Val
	}
	if val, ok := umap[root]; ok {
		return val // 如果umap里已经有记录则直接返回
	}
	// 偷父节点
	val1 := root.Val
	if root.Left != nil {
		val1 += rob(root.Left.Left) + rob(root.Left.Right) // 跳过root->left，相当于不考虑左孩子了
	}
	if root.Right != nil {
		val1 += rob(root.Right.Left) + rob(root.Right.Right) // 跳过root->right，相当于不考虑右孩子了
	}
	// 不偷父节点
	val2 := rob(root.Left) + rob(root.Right) // 考虑root的左右孩子
	umap[root] = max(val1, val2)             // umap记录一下结果
	return max(val1, val2)
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
