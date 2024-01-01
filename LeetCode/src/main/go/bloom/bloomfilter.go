package main

import (
	"fmt"
	"hash/fnv"
	"sync"
)

// BloomFilter 结构定义
type BloomFilter struct {
	size      uint
	hashFuncs uint
	bitset    []bool
	mutex     sync.Mutex
}

// 新建一个 BloomFilter
func NewBloomFilter(size uint, hashFuncs uint) *BloomFilter {
	return &BloomFilter{
		size:      size,
		hashFuncs: hashFuncs,
		bitset:    make([]bool, size),
	}
}

// 向 BloomFilter 中添加一个元素
func (bf *BloomFilter) Add(data string) {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	for i := uint(0); i < bf.hashFuncs; i++ {
		index := bf.hash(data, i) % bf.size
		bf.bitset[index] = true
	}
}

// 检查 BloomFilter 是否包含某个元素
func (bf *BloomFilter) Contains(data string) bool {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	for i := uint(0); i < bf.hashFuncs; i++ {
		index := bf.hash(data, i) % bf.size
		if !bf.bitset[index] {
			return false
		}
	}

	return true
}

// 使用 FNV-1a 哈希函数
func (bf *BloomFilter) hash(data string, seed uint) uint {
	h := fnv.New32a()
	h.Write([]byte(data))
	hash := h.Sum32()
	return (uint(hash) + seed) % bf.size
}

func main() {
	// 创建一个 BloomFilter，大小为100，使用3个哈希函数
	bloomFilter := NewBloomFilter(100, 3)

	// 添加一些元素
	elements := []string{"apple", "orange", "banana"}
	for _, element := range elements {
		bloomFilter.Add(element)
	}

	// 检查元素是否存在于 BloomFilter 中
	fmt.Println("Contains apple:", bloomFilter.Contains("apple"))           // true
	fmt.Println("Contains pear:", bloomFilter.Contains("pear"))             // false
	fmt.Println("Contains banana:", bloomFilter.Contains("banana"))         // true
	fmt.Println("Contains watermelon:", bloomFilter.Contains("watermelon")) // false
}
