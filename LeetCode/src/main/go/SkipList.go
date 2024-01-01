package main

import (
	"fmt"
	"math/rand"
)

const maxLevel = 16 // 跳表的最大层数

// Node 表示跳表中的节点
type Node struct {
	value   int
	forward []*Node
}

// SkipList 表示跳表
type SkipList struct {
	header *Node
	level  int
}

// NewNode 创建一个新的节点
func NewNode(value, level int) *Node {
	return &Node{value: value, forward: make([]*Node, level)}
}

// NewSkipList 创建一个新的跳表
func NewSkipList() *SkipList {
	header := NewNode(0, maxLevel)
	return &SkipList{header: header, level: 1}
}

// randomLevel 生成一个随机层数
func randomLevel() int {
	level := 1
	for rand.Float64() < 0.5 && level < maxLevel {
		level++
	}
	return level
}

// Insert 向跳表中插入一个元素
func (sl *SkipList) Insert(value int) {
	update := make([]*Node, maxLevel)
	current := sl.header

	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].value < value {
			current = current.forward[i]
		}
		update[i] = current
	}

	level := randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			update[i] = sl.header
		}
		sl.level = level
	}

	newNode := NewNode(value, level)
	for i := 0; i < level; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}
}

// Delete 从跳表中删除一个元素
func (sl *SkipList) Delete(value int) {
	if !sl.Search(value) {
		fmt.Printf("To be deleted: %v not found\n", value)
		return
	}
	update := make([]*Node, sl.level)
	current := sl.header

	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].value < value {
			current = current.forward[i]
		}
		//update[i] 是每层比target小的那个元素
		update[i] = current
	}

	// 要删除一个元素，只能用最底层（元素最全）去判断要删除的元素是否存在
	if current.forward[0] != nil && current.forward[0].value == value {
		for i := 0; i < sl.level; i++ {
			if update[i].forward[i] != nil && update[i].forward[i].value == value {
				update[i].forward[i] = update[i].forward[i].forward[i]
			}
		}
		// 更新层数
		for sl.level > 1 && sl.header.forward[sl.level-1] == nil {
			sl.level--
		}
	}
}

// Search 在跳表中搜索一个元素
func (sl *SkipList) Search(value int) bool {
	current := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].value < value {
			current = current.forward[i]
		}
		if current.forward[i] != nil && current.forward[i].value == value {
			return true
		}
	}
	return false
}

// Display 打印跳表的内容
func (sl *SkipList) Display() {
	fmt.Println("Skip List:")
	for i := sl.level - 1; i >= 0; i-- {
		fmt.Printf("Level %d: ", i)
		node := sl.header.forward[i]
		for node != nil {
			fmt.Printf("%d ", node.value)
			node = node.forward[i]
		}
		fmt.Println()
	}
	fmt.Print("level:")
	fmt.Println(sl.level)
	fmt.Println()
}

// DisplayVisual 打印跳表的可视化内容
func (sl *SkipList) DisplayVisual() {
	fmt.Println("Skip List Visualization:")
	for i := sl.level - 1; i >= 0; i-- {
		node := sl.header.forward[i]
		fmt.Printf("Level %d: ", i)

		for node != nil {
			fmt.Print(node.value)
			if node.forward[i] != nil {
				fmt.Print(" -> ")
			}
			node = node.forward[i]
		}
		fmt.Println()
	}
	fmt.Print("level:")
	fmt.Println(sl.level)
}

func main() {
	skipList := NewSkipList()

	// 插入元素
	skipList.Insert(3)
	skipList.Insert(6)
	skipList.Insert(7)
	skipList.Insert(9)
	skipList.Insert(12)
	skipList.Insert(19)
	skipList.Insert(17)
	skipList.Insert(26)
	skipList.Insert(21)
	skipList.Insert(25)
	skipList.Insert(26)
	skipList.Insert(27)

	skipList.DisplayVisual()

	// 搜索元素
	searchValue := 19
	fmt.Printf("Search %d: %t\n", searchValue, skipList.Search(searchValue))

	// 删除元素
	deleteValue := 17
	fmt.Printf("Delete %d\n", deleteValue)
	skipList.Delete(deleteValue)
	skipList.DisplayVisual()

	// 搜索删除后的元素
	fmt.Printf("Search %d: %t\n", deleteValue, skipList.Search(deleteValue))

	fmt.Printf("Delete %d\n", 23)
	skipList.Delete(23)
	skipList.DisplayVisual()
}
