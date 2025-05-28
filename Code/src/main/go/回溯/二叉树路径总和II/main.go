package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

var (
	result   [][]int
	currPath []int
)

func pathSum(root *TreeNode, targetSum int) [][]int {
	result = make([][]int, 0)
	currPath = make([]int, 0)
	traverse(root, targetSum)
	return result
}

func traverse(node *TreeNode, targetSum int) {
	if node == nil { // 这个判空也可以挪到递归遍历左右子树时去判断
		return
	}

	targetSum -= node.Val                 // 将targetSum在遍历每层的时候都减去本层节点的值
	currPath = append(currPath, node.Val) // 把当前节点放到路径记录里

	if node.Left == nil && node.Right == nil && targetSum == 0 { // 如果剩余的targetSum为0, 则正好就是符合的结果
		pathCopy := make([]int, len(currPath))
		copy(pathCopy, currPath)
		result = append(result, pathCopy) // 将副本放到结果集里
	}

	traverse(node.Left, targetSum)
	traverse(node.Right, targetSum)
	currPath = currPath[:len(currPath)-1] // 当前节点遍历完成, 从路径记录里删除掉
}

func main() {
	root := &TreeNode{Val: 5}
	root.Left = &TreeNode{Val: 4}
	root.Right = &TreeNode{Val: 8}
	root.Left.Left = &TreeNode{Val: 11}
	root.Left.Left.Left = &TreeNode{Val: 7}
	root.Left.Left.Right = &TreeNode{Val: 2}
	root.Right.Left = &TreeNode{Val: 13}
	root.Right.Right = &TreeNode{Val: 4}
	root.Right.Right.Right = &TreeNode{Val: 1}

	targetSum := 22
	resultTmp := pathSum(root, targetSum)
	for _, path := range resultTmp {
		for _, val := range path {
			print(val, " ")
		}
		println()
	}
}
