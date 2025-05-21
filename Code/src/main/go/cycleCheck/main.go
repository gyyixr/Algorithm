package main

import (
	"fmt"
)

// 定义节点状态
const (
	Unvisited = 0 // 未访问
	Visiting  = 1 // 正在访问 (在当前DFS路径上)
	Visited   = 2 // 已访问 (所有子节点已处理)
)

// DAG 表示一个有向图
type DAG struct {
	nodes  map[string]struct{} // 存储所有节点，使用空结构体节省空间
	adj    map[string][]string // 邻接表，key是节点，value是该节点指向的节点列表
	states map[string]int      // 存储每个节点的访问状态，用于DFS
}

// NewDAG 创建一个新的DAG实例
func NewDAG() *DAG {
	return &DAG{
		nodes:  make(map[string]struct{}),
		adj:    make(map[string][]string),
		states: make(map[string]int),
	}
}

// AddNode 向图中添加一个节点
func (g *DAG) AddNode(node string) {
	if _, exists := g.nodes[node]; !exists {
		g.nodes[node] = struct{}{}
		g.adj[node] = []string{} // 初始化邻接列表
	}
}

// AddEdge 向图中添加一条有向边 (from -> to)
func (g *DAG) AddEdge(from, to string) error {
	// 确保节点存在
	if _, ok := g.nodes[from]; !ok {
		// 或者可以选择自动添加节点：g.AddNode(from)
		return fmt.Errorf("起始节点 '%s' 不存在", from)
	}
	if _, ok := g.nodes[to]; !ok {
		// 或者可以选择自动添加节点：g.AddNode(to)
		return fmt.Errorf("目标节点 '%s' 不存在", to)
	}

	// 检查是否重复添加边（可选）
	for _, neighbor := range g.adj[from] {
		if neighbor == to {
			return nil // 边已存在
		}
	}

	g.adj[from] = append(g.adj[from], to)
	return nil
}

// hasCycleDFS 是一个辅助函数，用于执行深度优先搜索
// node: 当前访问的节点
// recursionStack: 记录当前DFS路径上的节点状态，用于检测环
func (g *DAG) hasCycleDFS(node string) bool {
	g.states[node] = Visiting // 标记为正在访问

	for _, neighbor := range g.adj[node] {
		if g.states[neighbor] == Visiting {
			// 如果邻居节点正在被访问（即在当前递归栈中），则发现环
			fmt.Printf("检测到循环：... -> %s -> %s -> ...\n", node, neighbor) // 简单打印环的一部分
			return true
		}
		if g.states[neighbor] == Unvisited {
			// 如果邻居节点未被访问，则递归访问
			if g.hasCycleDFS(neighbor) {
				return true // 在子路径中检测到环
			}
		}
		// 如果 g.states[neighbor] == Visited，则说明该邻居及其子树已处理完毕且无环，跳过
	}

	g.states[node] = Visited // 当前节点及其所有子节点都已处理完毕
	return false
}

// HasCycle 检测图中是否存在循环依赖
func (g *DAG) HasCycle() bool {
	// 初始化所有节点的状态为 Unvisited
	g.states = make(map[string]int) // 重置状态，以便多次调用
	for node := range g.nodes {
		g.states[node] = Unvisited
	}

	// 对图中每一个未访问的节点执行DFS
	for node := range g.nodes {
		if g.states[node] == Unvisited {
			if g.hasCycleDFS(node) {
				return true // 发现环
			}
		}
	}
	return false // 没有发现环
}

func main() {
	// 示例1: 无环图
	fmt.Println("--- 示例 1: 无环图 ---")
	dag1 := NewDAG()
	dag1.AddNode("A")
	dag1.AddNode("B")
	dag1.AddNode("C")
	dag1.AddNode("D")
	dag1.AddNode("E")

	dag1.AddEdge("A", "B")
	dag1.AddEdge("A", "C")
	dag1.AddEdge("B", "D")
	dag1.AddEdge("C", "D")
	dag1.AddEdge("D", "E")

	if dag1.HasCycle() {
		fmt.Println("图1中存在循环依赖")
	} else {
		fmt.Println("图1中不存在循环依赖")
	}
	fmt.Println()

	// 示例2: 有环图
	fmt.Println("--- 示例 2: 有环图 ---")
	dag2 := NewDAG()
	dag2.AddNode("1")
	dag2.AddNode("2")
	dag2.AddNode("3")
	dag2.AddNode("4")

	dag2.AddEdge("1", "2")
	dag2.AddEdge("2", "3")
	dag2.AddEdge("3", "1") // 环: 1 -> 2 -> 3 -> 1
	dag2.AddEdge("3", "4")

	if dag2.HasCycle() {
		fmt.Println("图2中存在循环依赖")
	} else {
		fmt.Println("图2中不存在循环依赖")
	}
	fmt.Println()

	// 示例3: 更复杂的环
	fmt.Println("--- 示例 3: 更复杂的环 ---")
	dag3 := NewDAG()
	dag3.AddNode("X")
	dag3.AddNode("Y")
	dag3.AddNode("Z")
	dag3.AddNode("W")
	dag3.AddNode("V")

	dag3.AddEdge("X", "Y")
	dag3.AddEdge("Y", "Z")
	dag3.AddEdge("Z", "W")
	dag3.AddEdge("W", "Y") // 环: Y -> Z -> W -> Y
	dag3.AddEdge("X", "V")

	if dag3.HasCycle() {
		fmt.Println("图3中存在循环依赖")
	} else {
		fmt.Println("图3中不存在循环依赖")
	}
	fmt.Println()

	// 示例4: 自环
	fmt.Println("--- 示例 4: 自环 ---")
	dag4 := NewDAG()
	dag4.AddNode("P")
	dag4.AddNode("Q")
	dag4.AddEdge("P", "P") // 自环
	dag4.AddEdge("P", "Q")

	if dag4.HasCycle() {
		fmt.Println("图4中存在循环依赖")
	} else {
		fmt.Println("图4中不存在循环依赖")
	}
}
