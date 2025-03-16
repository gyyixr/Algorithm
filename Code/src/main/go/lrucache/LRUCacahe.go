package main

import (
	"container/list"
	"fmt"
)

// LRUCache represents a simple LRU cache.
type LRUCache struct {
	capacity int
	cacheMap map[int]*list.Element
	lruList  *list.List
}

// Entry represents a key-value pair in the cache.
type Entry struct {
	key   int
	value int
}

// NewLRUCache creates a new LRUCache with the specified capacity.
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cacheMap: make(map[int]*list.Element),
		lruList:  list.New(),
	}
}

// Get retrieves the value associated with the key from the cache.
func (lru *LRUCache) Get(key int) int {
	if elem, ok := lru.cacheMap[key]; ok {
		lru.lruList.MoveToFront(elem)
		return elem.Value.(Entry).value
	}
	return -1
}

// Put inserts a key-value pair into the cache.
func (lru *LRUCache) Put(key, value int) {
	if elem, ok := lru.cacheMap[key]; ok {
		// Update the existing entry
		elem.Value = Entry{key, value}
		lru.lruList.MoveToFront(elem)
	} else {
		// Add a new entry
		if len(lru.cacheMap) >= lru.capacity {
			// Evict the least recently used entry
			oldest := lru.lruList.Back()
			delete(lru.cacheMap, oldest.Value.(Entry).key)
			lru.lruList.Remove(oldest)
		}

		// Add the new entry to the front
		newElem := lru.lruList.PushFront(Entry{key, value})
		lru.cacheMap[key] = newElem
	}
}

// PrintCacheState prints the current state of the cache.
func (lru *LRUCache) PrintCacheState() {
	fmt.Println("Cache State:")
	for elem := lru.lruList.Front(); elem != nil; elem = elem.Next() {
		fmt.Printf("%d: %d\n", elem.Value.(Entry).key, elem.Value.(Entry).value)
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
