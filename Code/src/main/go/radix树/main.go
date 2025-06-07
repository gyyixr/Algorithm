package main

import (
	"fmt"
	"strings"
)

// RadixNode 表示基数树中的一个节点
type RadixNode struct {
	// 边上的字符串
	edge string
	// 子节点, key 是子节点 edge 的第一个字符
	children    map[byte]*RadixNode
	isEndOfWord bool
}

// RadixTree 表示基数树
type RadixTree struct {
	root *RadixNode
}

// NewRadixTree 创建一个新的基数树
func NewRadixTree() *RadixTree {
	return &RadixTree{
		root: &RadixNode{
			edge:     "",
			children: make(map[byte]*RadixNode),
		},
	}
}

// Insert 向基数树中插入一个单词
func (t *RadixTree) Insert(word string) {
	node := t.root
	remaining := word

	for {
		// 如果当前剩余的字符串为空，标记节点并返回
		if remaining == "" {
			node.isEndOfWord = true
			return
		}

		// 在子节点中查找匹配的边
		child, ok := node.children[remaining[0]]
		if !ok {
			// 没有匹配的子节点，直接创建一个新节点
			node.children[remaining[0]] = &RadixNode{
				edge:        remaining,
				children:    make(map[byte]*RadixNode),
				isEndOfWord: true,
			}
			return
		}

		// 找到了一个以相同字符开头的子节点，计算公共前缀的长度
		commonPrefixLen := 0
		for commonPrefixLen < len(remaining) && commonPrefixLen < len(child.edge) {
			if remaining[commonPrefixLen] != child.edge[commonPrefixLen] {
				break
			}
			commonPrefixLen++
		}

		// 情况 1: 剩余字符串与子节点边的前缀完全匹配
		if commonPrefixLen == len(remaining) && commonPrefixLen == len(child.edge) {
			// "apple" 已经存在, 现在插入 "apple"
			child.isEndOfWord = true
			return
		}

		// 情况 2: 剩余字符串是子节点边的一个前缀
		// 例如: 树中已有 "apple", 现在插入 "app"
		if commonPrefixLen == len(remaining) {
			// 需要分裂边
			splitNode := &RadixNode{
				edge:        child.edge[commonPrefixLen:], // 新分裂出的节点，边为 "le"
				children:    child.children,
				isEndOfWord: child.isEndOfWord,
			}

			child.edge = remaining // 原有子节点边缩短为 "app"
			child.isEndOfWord = true
			child.children = map[byte]*RadixNode{
				splitNode.edge[0]: splitNode,
			}
			return
		}

		// 情况 3: 子节点边是剩余字符串的一个前缀
		// 例如: 树中已有 "app", 现在插入 "apple"
		if commonPrefixLen == len(child.edge) {
			node = child
			remaining = remaining[commonPrefixLen:]
			// 继续循环查找
			continue
		}

		// 情况 4: 剩余字符串和子节点边有部分共同前缀，但都需要分裂
		// 例如: 树中已有 "apple", 现在插入 "apply"
		// 共同前缀是 "appl"
		splitParent := &RadixNode{
			edge:        child.edge[:commonPrefixLen], // 新的父分裂节点，边为 "appl"
			children:    make(map[byte]*RadixNode),
			isEndOfWord: false, // 这不是一个单词的结尾
		}

		// 原有子节点更新
		child.edge = child.edge[commonPrefixLen:] // 边变为 "e"
		splitParent.children[child.edge[0]] = child

		// 新插入单词的剩余部分
		newChild := &RadixNode{
			edge:        remaining[commonPrefixLen:], // 边变为 "y"
			children:    make(map[byte]*RadixNode),
			isEndOfWord: true,
		}
		splitParent.children[newChild.edge[0]] = newChild

		// 将分裂后的新父节点连接到当前节点
		node.children[splitParent.edge[0]] = splitParent
		return
	}
}

// Search 在基数树中搜索一个单词是否存在
func (t *RadixTree) Search(word string) bool {
	node := t.root
	remaining := word

	for {
		if remaining == "" {
			return node.isEndOfWord
		}

		// 在子节点中查找匹配的边
		child, ok := node.children[remaining[0]]
		if !ok {
			// 没有以该字符开头的边
			return false
		}

		// 检查剩余部分是否以子节点的边为前缀
		if !strings.HasPrefix(remaining, child.edge) {
			return false
		}

		// 如果完全匹配
		if len(remaining) == len(child.edge) {
			return child.isEndOfWord
		}

		// 如果剩余部分比边长，继续向下查找
		node = child
		remaining = remaining[len(child.edge):]
	}
}

func main() {
	tree := NewRadixTree()

	words := []string{"romane", "romanus", "romulus", "rubens", "ruber", "rubicon", "rubicundus"}
	fmt.Println("--- 插入单词 ---")
	for _, word := range words {
		tree.Insert(word)
		fmt.Printf("插入: %s\n", word)
	}

	fmt.Println("\n--- 搜索测试 ---")
	searchWords := []string{"romane", "romanus", "romulus", "ruber", "rubicon", "rom", "rubic"}
	for _, word := range searchWords {
		if tree.Search(word) {
			fmt.Printf("单词 '%s' 存在\n", word)
		} else {
			fmt.Printf("单词 '%s' 不存在\n", word)
		}
	}

	fmt.Println("\n--- 更多插入和分裂测试 ---")
	tree.Insert("test")
	fmt.Println("插入: test")
	tree.Insert("testing")
	fmt.Println("插入: testing (分裂 'test')")
	tree.Insert("team")
	fmt.Println("插入: team (分裂 'te')")

	fmt.Println("\n--- 搜索分裂后的单词 ---")
	searchWords2 := []string{"test", "testing", "team", "tea"}
	for _, word := range searchWords2 {
		if tree.Search(word) {
			fmt.Printf("单词 '%s' 存在\n", word)
		} else {
			fmt.Printf("单词 '%s' 不存在\n", word)
		}
	}
}
