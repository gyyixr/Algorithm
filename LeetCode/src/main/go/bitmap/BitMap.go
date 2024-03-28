package main

import (
	"fmt"
)

// BitMap 是一个结构体，使用一个整数数组来存储位图
type BitMap struct {
	data []uint32 // 使用uint32数组来存储位图，根据需要可以更改为uint64或其他类型
}

// NewBitMap 创建一个新的位图，指定大小（以位为单位）
func NewBitMap(size uint) *BitMap {
	// 计算需要多少个uint32来存储所有的位
	capaity := (size + 31) / 32 * 4
	return &BitMap{
		data: make([]uint32, capaity),
	}
}

// Set 设置指定位置的位为1
func (bm *BitMap) Set(pos uint) {
	if pos >= uint(len(bm.data)*32) {
		panic("position out of range")
	}
	bm.data[pos/32] |= (1 << (pos % 32))
}

// Clear 清除指定位置的位，将其设置为0
func (bm *BitMap) Clear(pos uint) {
	if pos >= uint(len(bm.data)*32) {
		panic("position out of range")
	}
	// a &^ b = a &(^b)
	bm.data[pos/32] &^= 1 << (pos % 32)
}

// Test 测试指定位置的位是否为1
func (bm *BitMap) Test(pos uint) bool {
	if pos >= uint(len(bm.data)*32) {
		panic("position out of range")
	}
	return (bm.data[pos/32] & (1 << (pos % 32))) != 0
}

// Flip 翻转指定位置的位
func (bm *BitMap) Flip(pos uint) {
	if pos >= uint(len(bm.data)*32) {
		panic("position out of range")
	}
	bm.data[pos/32] ^= 1 << (pos % 32)
}

func main() {
	bm := NewBitMap(100) // 创建一个可以存储100个位的位图

	bm.Set(5)  // 设置第5位
	bm.Set(10) // 设置第10位

	fmt.Printf("Bit at position 5: %v\n", bm.Test(5))   // 输出：true
	fmt.Printf("Bit at position 10: %v\n", bm.Test(10)) // 输出：true

	bm.Clear(5) // 清除第5位

	fmt.Printf("Bit at position 5 after clear: %v\n", bm.Test(5)) // 输出：false

	bm.Flip(10) // 翻转第10位

	fmt.Printf("Bit at position 10 after flip: %v\n", bm.Test(10)) // 输出：false
}
