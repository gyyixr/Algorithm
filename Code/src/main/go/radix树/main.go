package main

import (
	"fmt"
	"strings"
)

// RadixNode, RadixTree, NewRadixTree 结构体和构造函数与之前版本相同
// RadixNode 表示基数树中的一个节点
type RadixNode struct {
	edge        string
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

// --- 递归实现 ---

// Insert 是公开的插入方法，它调用递归辅助函数
func (t *RadixTree) Insert(word string) {
	if word == "" {
		return
	}
	t.insert(t.root, word)
}

// insert 是私有的递归辅助函数
func (t *RadixTree) insert(node *RadixNode, remaining string) {
	// 在子节点中查找匹配的边
	child, ok := node.children[remaining[0]]
	if !ok {
		// --- 基本情况 1: 没有匹配的子节点 ---
		// 直接创建一个新节点并返回
		node.children[remaining[0]] = &RadixNode{
			edge:        remaining,
			children:    make(map[byte]*RadixNode),
			isEndOfWord: true,
		}
		return
	}

	// 计算公共前缀的长度
	commonPrefixLen := 0
	for commonPrefixLen < len(remaining) && commonPrefixLen < len(child.edge) {
		if remaining[commonPrefixLen] != child.edge[commonPrefixLen] {
			break
		}
		commonPrefixLen++
	}

	// 如果公共前缀的长度小于子节点边的长度，说明需要分裂
	if commonPrefixLen < len(child.edge) {
		// 例如: 树中已有 "apple", 现在插入 "apply"
		// 共同前缀是 "appl"
		splitParent := &RadixNode{
			edge:        child.edge[:commonPrefixLen], // 新的父分裂节点，边为 "appl"
			children:    make(map[byte]*RadixNode),
			isEndOfWord: false, // 这不是一个单词的结尾
		}

		// 1. 原有子节点更新
		child.edge = child.edge[commonPrefixLen:] // 边变为 "e"
		splitParent.children[child.edge[0]] = child

		// 2. 将分裂后的新父节点连接到当前节点
		node.children[splitParent.edge[0]] = splitParent

		// 如果插入的单词在分裂点就结束了 (例如: 插入 "app" 而树中有 "apple")
		if commonPrefixLen == len(remaining) {
			splitParent.isEndOfWord = true
			return
		}

		// 3. 为新插入单词的剩余部分创建新节点
		newChild := &RadixNode{
			edge:        remaining[commonPrefixLen:], // 边变为 "y"
			children:    make(map[byte]*RadixNode),
			isEndOfWord: true,
		}
		splitParent.children[newChild.edge[0]] = newChild
		return
	}

	// 如果公共前缀与子节点边完全匹配
	if commonPrefixLen == len(child.edge) {
		// 如果插入的单词也在这里结束
		if commonPrefixLen == len(remaining) {
			// --- 基本情况 2: 单词已存在 ---
			child.isEndOfWord = true
			return
		}

		// --- 递归步骤 ---
		// 边是新单词的前缀，对剩余部分递归插入
		// 例如: 树中已有 "app", 现在插入 "apple"
		t.insert(child, remaining[commonPrefixLen:])
	}
}

// Search 是公开的搜索方法，它调用递归辅助函数
func (t *RadixTree) Search(word string) bool {
	if word == "" {
		return false
	}
	return t.search(t.root, word)
}

// search 是私有的递归辅助函数
func (t *RadixTree) search(node *RadixNode, remaining string) bool {
	// 在子节点中查找匹配的边
	child, ok := node.children[remaining[0]]
	if !ok {
		// --- 基本情况 1: 没找到匹配的子节点 ---
		return false
	}

	// 检查剩余部分是否以子节点的边为前缀
	if !strings.HasPrefix(remaining, child.edge) {
		// --- 基本情况 2: 边不匹配 ---
		return false
	}

	// 如果剩余部分与边的长度完全相同
	if len(remaining) == len(child.edge) {
		// --- 基本情况 3: 找到完全匹配，返回单词标记 ---
		return child.isEndOfWord
	}

	// --- 递归步骤 ---
	// 边是剩余部分的前缀，继续向下一层搜索
	return t.search(child, remaining[len(child.edge):])
}

func main() {
	// main 函数与之前的版本完全相同，用于测试
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
