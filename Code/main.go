package main

import "fmt"

// TrieNode 表示前缀树中的一个节点
type TrieNode struct {
	children    map[rune]*TrieNode
	isEndOfWord bool
}

// NewTrieNode 创建一个新的前缀树节点
func NewTrieNode() *TrieNode {
	return &TrieNode{
		children:    make(map[rune]*TrieNode),
		isEndOfWord: false,
	}
}

// Trie 表示前缀树
type Trie struct {
	root *TrieNode
}

// NewTrie 创建一个新的前缀树
func NewTrie() *Trie {
	return &Trie{root: NewTrieNode()}
}

// Insert 向前缀树中插入一个单词
func (t *Trie) Insert(word string) {
	node := t.root
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			node.children[char] = NewTrieNode()
		}
		node = node.children[char]
	}
	node.isEndOfWord = true
}

// Search 在前缀树中搜索一个单词是否存在
func (t *Trie) Search(word string) bool {
	node := t.root
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			return false
		}
		node = node.children[char]
	}
	return node.isEndOfWord
}

// StartsWith 检查前缀树中是否存在以指定前缀开头的单词
func (t *Trie) StartsWith(prefix string) bool {
	node := t.root
	for _, char := range prefix {
		if _, ok := node.children[char]; !ok {
			return false
		}
		node = node.children[char]
	}
	return true
}

func main() {
	trie := NewTrie()

	words := []string{"apple", "app", "apricot", "banana", "bandana"}
	for _, word := range words {
		trie.Insert(word)
		fmt.Printf("插入单词: %s\n", word)
	}

	fmt.Println("\n--- 搜索 ---")
	searchWords := []string{"app", "apple", "apricots", "ban", "banana"}
	for _, word := range searchWords {
		if trie.Search(word) {
			fmt.Printf("单词 '%s' 存在于前缀树中\n", word)
		} else {
			fmt.Printf("单词 '%s' 不存在于前缀树中\n", word)
		}
	}

	fmt.Println("\n--- 前缀搜索 ---")
	prefixes := []string{"ap", "b", "band", "cat"}
	for _, prefix := range prefixes {
		if trie.StartsWith(prefix) {
			fmt.Printf("存在以 '%s' 开头的单词\n", prefix)
		} else {
			fmt.Printf("不存在以 '%s' 开头的单词\n", prefix)
		}
	}
}
