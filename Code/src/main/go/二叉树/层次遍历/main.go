package main

import (
	"fmt"
)

// TreeNode 定义二叉树节点
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// levelOrder (DFS 实现的层次遍历 - 保持不变，用于对比)
func levelOrder(root *TreeNode) [][]int {
	arr := [][]int{}
	depth := 0
	var order func(root *TreeNode, depth int)
	order = func(root *TreeNode, depth int) {
		if root == nil {
			return
		}
		if len(arr) == depth {
			arr = append(arr, []int{})
		}
		arr[depth] = append(arr[depth], root.Val)
		order(root.Left, depth+1)
		order(root.Right, depth+1)
	}
	order(root, depth)
	return arr
}

// levelOrderBFS 使用切片作为队列实现层次遍历
func levelOrderBFS(root *TreeNode) [][]int {
	if root == nil { // 1. 处理空树的情况
		return [][]int{}
	}

	result := [][]int{}
	queue := []*TreeNode{}      // 2. 初始化一个 *TreeNode 类型的切片作为队列
	queue = append(queue, root) // 3. 将根节点入队 (追加到切片末尾)

	for len(queue) > 0 { // 4. 当队列不为空时，持续循环
		levelSize := len(queue)       // 5. 获取当前层的节点数量
		currentLevelValues := []int{} // 6. 初始化一个切片，用于存储当前层所有节点的值

		for i := 0; i < levelSize; i++ { // 7. 遍历当前层的所有节点
			// 8. 出队：获取队列的第一个元素
			node := queue[0]
			// 9. 出队：移除队列的第一个元素 (通过切片重新赋值)
			queue = queue[1:]

			currentLevelValues = append(currentLevelValues, node.Val) // 10. 将当前节点的值添加到当前层的值列表中

			// 11. 如果当前节点有左子节点，将其入队
			if node.Left != nil {
				queue = append(queue, node.Left) // 入队：追加到切片末尾
			}
			// 12. 如果当前节点有右子节点，将其入队
			if node.Right != nil {
				queue = append(queue, node.Right) // 入队：追加到切片末尾
			}
		}
		result = append(result, currentLevelValues) // 13. 将当前层的值列表添加到最终结果中
	}

	return result // 14. 返回按层存储节点值的二维切片
}

func main() {
	// 构建一个测试二叉树
	//       3
	//      / \
	//     9  20
	//       /  \
	//      15   7
	//     /
	//    1
	root := &TreeNode{Val: 3}
	root.Left = &TreeNode{Val: 9}
	root.Right = &TreeNode{Val: 20}
	root.Right.Left = &TreeNode{Val: 15}
	root.Right.Right = &TreeNode{Val: 7}
	root.Right.Left.Left = &TreeNode{Val: 1} // 添加一个更深的节点

	fmt.Println("测试 DFS 实现的层次遍历 (levelOrder):")
	resultDFS := levelOrder(root)
	fmt.Println(resultDFS) // 预期输出: [[3] [9 20] [15 7] [1]]

	fmt.Println("\n测试 BFS (基于切片) 实现的层次遍历 (levelOrderBFS):")
	resultBFS := levelOrderBFS(root)
	fmt.Println(resultBFS) // 预期输出: [[3] [9 20] [15 7] [1]]

	// 测试一个空树
	fmt.Println("\n测试空树 (DFS):")
	emptyResultDFS := levelOrder(nil)
	fmt.Println(emptyResultDFS) // 预期输出: []

	fmt.Println("\n测试空树 (BFS - 基于切片):")
	emptyResultBFS := levelOrderBFS(nil)
	fmt.Println(emptyResultBFS) // 预期输出: []

	// 测试只有根节点的树
	fmt.Println("\n测试只有根节点的树 (DFS):")
	singleNodeRoot := &TreeNode{Val: 100}
	singleNodeResultDFS := levelOrder(singleNodeRoot)
	fmt.Println(singleNodeResultDFS) // 预期输出: [[100]]

	fmt.Println("\n测试只有根节点的树 (BFS - 基于切片):")
	singleNodeResultBFS := levelOrderBFS(singleNodeRoot)
	fmt.Println(singleNodeResultBFS) // 预期输出: [[100]]

	// 构建另一个稍微不同的树
	//       1
	//      / \
	//     2   3
	//    /     \
	//   4       5
	root2 := &TreeNode{Val: 1,
		Left: &TreeNode{Val: 2,
			Left: &TreeNode{Val: 4},
		},
		Right: &TreeNode{Val: 3,
			Right: &TreeNode{Val: 5},
		},
	}
	fmt.Println("\n测试另一个树 (DFS):")
	resultDFS2 := levelOrder(root2)
	fmt.Println(resultDFS2) // 预期输出: [[1] [2 3] [4 5]]

	fmt.Println("\n测试另一个树 (BFS - 基于切片):")
	resultBFS2 := levelOrderBFS(root2)
	fmt.Println(resultBFS2) // 预期输出: [[1] [2 3] [4 5]]
}
