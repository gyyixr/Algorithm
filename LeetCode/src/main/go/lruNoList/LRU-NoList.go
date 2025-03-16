package main

import (
	"fmt"
)

// LRUCache represents a simple LRU cache.
type LRUCache struct {
	size     int
	capacity int
	cacheMap map[int]*Node
	head     *Node
	tail     *Node
}

// Node represents a node in the doubly linked list.
type Node struct {
	key   int
	value int
	prev  *Node
	next  *Node
}

// NewLRUCache creates a new LRUCache with the specified capacity.
func NewLRUCache(capacity int) *LRUCache {
	lru := &LRUCache{
		size:     0,
		capacity: capacity,
		cacheMap: make(map[int]*Node, capacity),
	}
	lru.head = &Node{}
	lru.tail = &Node{}
	lru.head.next = lru.tail
	lru.tail.prev = lru.head
	return lru
}

// moveToHead moves a node to the front of the linked list.
func (lru *LRUCache) moveToHead(node *Node) {
	lru.removeNode(node)
	lru.addToHead(node)
}

func (lru *LRUCache) removeNode(node *Node) {
	node.next.prev = node.prev
	node.prev.next = node.next
}

func (lru *LRUCache) addToHead(node *Node) {
	node.prev = lru.head
	node.next = lru.head.next
	lru.head.next.prev = node
	lru.head.next = node
}

func (lru *LRUCache) removeTail() *Node {
	prev := lru.tail.prev
	lru.removeNode(prev)
	return prev
}

// Get retrieves the value associated with the key from the cache.
func (lru *LRUCache) Get(key int) int {
	if node, ok := lru.cacheMap[key]; ok {
		lru.moveToHead(node)
		return node.value
	}
	return -1
}

// Put inserts a key-value pair into the cache.
func (lru *LRUCache) Put(key, value int) {
	if node, ok := lru.cacheMap[key]; ok {
		node.value = value
		lru.moveToHead(node)
	} else {
		newNode := &Node{key: key, value: value}
		lru.cacheMap[key] = newNode
		lru.addToHead(newNode)
		lru.size++
		if lru.size > lru.capacity {
			needRemoveTail := lru.removeTail()
			delete(lru.cacheMap, needRemoveTail.key)
			lru.size--
		}
	}
}

// PrintCacheState prints the current state of the cache.
func (lru *LRUCache) PrintCacheState() {
	fmt.Println("Cache State:")
	node := lru.head.next
	for node.next != nil {
		fmt.Printf("%d: %d\n", node.key, node.value)
		node = node.next
	}
	fmt.Println("---------------")
}

func main() {
	lruCache := NewLRUCache(3)

	lruCache.Put(1, 10)
	lruCache.PrintCacheState()

	lruCache.Put(2, 20)
	lruCache.PrintCacheState()

	lruCache.Put(3, 30)
	lruCache.PrintCacheState()

	fmt.Println("Value for key 2:", lruCache.Get(2))

	lruCache.Put(4, 40) // This will trigger eviction of the least recently used element (key 1)
	lruCache.PrintCacheState()
}
