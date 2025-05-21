package main

import "fmt"

// TreeNode 定义二叉树节点 (假设与之前的定义相同)
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// findBottomLeftValue 使用切片作为队列寻找二叉树最左下角节点的值
func findBottomLeftValue(root *TreeNode) int {
	if root == nil { // 1. 处理空树的情况
		// 根据题目要求，空树可能需要返回特定值或引发错误
		// 这里我们假设题目约定树非空，或返回0作为一种默认值
		// 如果题目有明确定义，请按定义处理
		return 0 // 或者可以 panic("tree is nil")
	}

	var gradation int      // 2. 用于存储最终结果的变量 (每层最左边的值，最后为最深层最左边的值)
	queue := []*TreeNode{} // 3. 初始化一个 *TreeNode 类型的切片作为队列

	queue = append(queue, root) // 4. 将根节点入队 (追加到切片末尾)

	for len(queue) > 0 { // 5. 当队列不为空时，进行层序遍历
		length := len(queue) // 6. 获取当前层的节点数量

		for i := 0; i < length; i++ { // 7. 遍历当前层的所有节点
			// 8. 出队：获取队列的第一个元素
			node := queue[0]
			// 9. 出队：移除队列的第一个元素 (通过切片重新赋值)
			queue = queue[1:]

			// 10. 关键步骤：记录每层最左边的节点值
			if i == 0 {
				gradation = node.Val // 如果是当前层的第一个节点 (i=0)，则更新 gradation
			}

			// 11. 将当前节点的左子节点入队 (如果存在)
			if node.Left != nil {
				queue = append(queue, node.Left) // 入队：追加到切片末尾
			}
			// 12. 将当前节点的右子节点入队 (如果存在)
			if node.Right != nil {
				queue = append(queue, node.Right) // 入队：追加到切片末尾
			}
		}
	}
	return gradation // 13. 返回最后记录的 gradation 值
}

var depth int // 全局变量 最大深度
var res int   // 记录最终结果
func findBottomLeftValueDFS(root *TreeNode) int {
	depth, res = 0, 0 // 初始化
	dfs(root, 1)
	return res
}

func dfs(root *TreeNode, d int) {
	if root == nil {
		return
	}
	// 因为先遍历左边，所以左边如果有值，右边的同层不会更新结果
	if root.Left == nil && root.Right == nil && depth < d {
		depth = d
		res = root.Val
	}
	dfs(root.Left, d+1) // 隐藏回溯
	dfs(root.Right, d+1)
}

func main() {
	// 测试用例 1:
	//       2
	//      / \
	//     1   3
	root1 := &TreeNode{Val: 2}
	root1.Left = &TreeNode{Val: 1}
	root1.Right = &TreeNode{Val: 3}
	fmt.Printf("Test Case 1 (Expected: 1): %d\n", findBottomLeftValue(root1))

	// 测试用例 2:
	//       1
	//      / \
	//     2   3
	//    /   / \
	//   4   5   6
	//  /
	// 7
	root2 := &TreeNode{Val: 1}
	root2.Left = &TreeNode{Val: 2}
	root2.Left.Left = &TreeNode{Val: 4}
	root2.Left.Left.Left = &TreeNode{Val: 7}
	root2.Right = &TreeNode{Val: 3}
	root2.Right.Left = &TreeNode{Val: 5}
	root2.Right.Right = &TreeNode{Val: 6}
	fmt.Printf("层次遍历：Test Case 2 (Expected: 7): %d\n", findBottomLeftValue(root2))
	fmt.Printf("递归遍历：Test Case 2 (Expected: 7): %d\n", findBottomLeftValueDFS(root2))

	// 测试用例 3: 只有一个节点
	root3 := &TreeNode{Val: 100}
	fmt.Printf("Test Case 3 (Expected: 100): %d\n", findBottomLeftValue(root3))

	// 测试用例 4: 树偏向右边
	//   0
	//    \
	//     -1
	root4 := &TreeNode{Val: 0}
	root4.Right = &TreeNode{Val: -1}
	fmt.Printf("Test Case 4 (Expected: -1): %d\n", findBottomLeftValue(root4))

	// 测试用例 5:
	//        10
	//       /
	//      5
	//     /
	//    3
	//   /
	//  1
	root5 := &TreeNode{Val: 10}
	root5.Left = &TreeNode{Val: 5}
	root5.Left.Left = &TreeNode{Val: 3}
	root5.Left.Left.Left = &TreeNode{Val: 1}
	fmt.Printf("Test Case 5 (Expected: 1): %d\n", findBottomLeftValue(root5))

	// 测试用例 6: 空树 (根据函数实现，可能返回0或panic，这里我们添加了nil检查)
	fmt.Printf("Test Case 6 (Empty Tree, Expected based on current code: 0): %d\n", findBottomLeftValue(nil))
}
