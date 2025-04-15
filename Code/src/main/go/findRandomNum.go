//给定1千万个随机数，范围在1到1亿之间，要求找出所有在这个范围内但没有出现在随机数中的数字。
package main

import (
	"fmt"
)

//bitmap法
func findMissingNumbersBitmap(randomNumbers []int, maxNum int) []int {
	// 计算需要的字节数 (每个字节8位)
	bitmapSize := (maxNum + 7) / 8
	bitmap := make([]byte, bitmapSize)

	// 标记存在的数字
	for _, num := range randomNumbers {
		if num < 1 || num > maxNum {
			continue
		}
		index := (num - 1) / 8
		bit := (num - 1) % 8
		bitmap[index] |= 1 << bit
	}

	// 收集缺失的数字
	missing := make([]int, 0)
	for num := 1; num <= maxNum; num++ {
		index := (num - 1) / 8
		bit := (num - 1) % 8
		if (bitmap[index] & (1 << bit)) == 0 {
			missing = append(missing, num)
		}
	}

	return missing
}

func main() {
	// 示例使用
	randomNumbers := []int{2, 3, 5, 7, 11} // 假设这是你的1千万个随机数
	maxNum := 20                           // 1亿
	missing := findMissingNumbersBitmap(randomNumbers, maxNum)
	fmt.Println("Missing numbers:", missing)
}
