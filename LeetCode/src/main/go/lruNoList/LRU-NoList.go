package main

import (
	"fmt"
)

// LRUCache represents a simple LRU cache.
type LRUCache struct {
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
	return &LRUCache{
		capacity: capacity,
		cacheMap: make(map[int]*Node),
	}
}

// moveToHead moves a node to the front of the linked list.
func (lru *LRUCache) moveToHead(node *Node) {
	if node == lru.head {
		return
	}

	if node == lru.tail {
		lru.tail = node.prev
	} else {
		node.next.prev = node.prev
	}

	node.prev.next = node.next
	node.prev = nil
	node.next = lru.head
	lru.head.prev = node
	lru.head = node
}

// removeTail removes the last node from the linked list.
func (lru *LRUCache) removeTail() {
	if lru.tail == nil {
		return
	}

	delete(lru.cacheMap, lru.tail.key)

	if lru.head == lru.tail {
		lru.head, lru.tail = nil, nil
	} else {
		lru.tail = lru.tail.prev
		lru.tail.next = nil
	}
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
		// Update the existing entry
		node.value = value
		lru.moveToHead(node)
	} else {
		// Add a new entry
		newNode := &Node{key, value, nil, nil}
		lru.cacheMap[key] = newNode

		if lru.head == nil {
			lru.head, lru.tail = newNode, newNode
		} else {
			newNode.next = lru.head
			lru.head.prev = newNode
			lru.head = newNode
		}

		if len(lru.cacheMap) > lru.capacity {
			lru.removeTail()
		}
	}
}

// PrintCacheState prints the current state of the cache.
func (lru *LRUCache) PrintCacheState() {
	fmt.Println("Cache State:")
	node := lru.head
	for node != nil {
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
	lruCache.PrintCacheState()

	lruCache.Put(4, 40) // This will trigger eviction of the least recently used element (key 1)
	lruCache.PrintCacheState()
}
