package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func constructMaximumBinaryTree(nums []int) *TreeNode {
	if len(nums) == 0 {
		return nil
	}
	// 找到最大值
	index := findMax(nums)
	// 构造二叉树
	root := &TreeNode{
		Val:   nums[index],
		Left:  constructMaximumBinaryTree(nums[:index]),   //左半边
		Right: constructMaximumBinaryTree(nums[index+1:]), //右半边
	}
	return root
}
func findMax(nums []int) int {
	max := nums[0]
	maxIndex := 0
	for i, _ := range nums {
		if nums[i] > max {
			max = nums[i]
			maxIndex = i
		}
	}
	return maxIndex
}

func main() {
	nums := []int{3, 2, 1, 6, 0, 5}
	root := constructMaximumBinaryTree(nums)
	// 可以在这里添加代码来打印或验证构造的二叉树
	// 例如，使用前序遍历打印树的值
	printPreOrder(root)
}

func printPreOrder(node *TreeNode) {
	if node == nil {
		return
	}
	println(node.Val)         // 打印当前节点的值
	printPreOrder(node.Left)  // 递归打印左子树
	printPreOrder(node.Right) // 递归打印右子树
}
