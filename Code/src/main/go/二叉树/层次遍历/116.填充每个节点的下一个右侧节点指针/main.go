package main

type Node struct {
	Val   int
	Left  *Node
	Right *Node
	Next  *Node
}

func connect(root *Node) *Node {
	if root == nil {
		return root
	}
	queue := make([]*Node, 0)
	queue = append(queue, root)
	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			node := queue[i]
			if i != size-1 {
				queue[i].Next = queue[i+1]
			} else {
				queue[i].Next = nil
			}
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		queue = queue[size:]
	}
	return root
}

func main() {
	root := &Node{Val: 3}
	root.Left = &Node{Val: 9}
	root.Right = &Node{Val: 20}
	root.Right.Left = &Node{Val: 15}
	root.Right.Right = &Node{Val: 7}
	root.Right.Left.Left = &Node{Val: 1} // 添加一个更深的节点
	connect(root)
}
