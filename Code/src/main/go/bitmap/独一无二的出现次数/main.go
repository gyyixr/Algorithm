package main

import "fmt"

// uniqueOccurrences 使用哈希表和位图思想解决问题
func uniqueOccurrences(arr []int) bool {
	// 步骤 1: 使用 map 统计每个数字的出现频率
	freqCounter := make(map[int]int)
	for _, num := range arr {
		freqCounter[num]++
	}

	// 步骤 2: 使用布尔数组作为“位图”来检查频率是否唯一
	// 数组长度最大为 1000，所以频率的最大值也为 1000。
	// 我们创建一个大小为 1001 的布尔数组来记录出现过的频率。
	// 索引代表频率，值代表该频率是否已出现。
	seenFrequencies := make([]bool, 1001)

	// 遍历所有统计出来的频率
	for _, freq := range freqCounter {
		// 如果该频率在“位图”中已被标记为 true，说明频率重复
		if seenFrequencies[freq] {
			return false
		}
		// 否则，在“位图”中标记该频率为已出现
		seenFrequencies[freq] = true
	}

	// 如果所有频率都未重复，返回 true
	return true
}

func main() {
	// 示例 1
	arr1 := []int{1, 2, 2, 1, 1, 3}
	fmt.Printf("输入: arr = %v\n", arr1)
	fmt.Printf("输出: %v\n\n", uniqueOccurrences(arr1)) // 期望输出: true

	// 示例 2
	arr2 := []int{1, 2}
	fmt.Printf("输入: arr = %v\n", arr2)
	fmt.Printf("输出: %v\n\n", uniqueOccurrences(arr2)) // 期望输出: false

	// 示例 3
	arr3 := []int{-3, 0, 1, -3, 1, 1, 1, -3, 10, 0}
	fmt.Printf("输入: arr = %v\n", arr3)
	fmt.Printf("输出: %v\n\n", uniqueOccurrences(arr3)) // 期望输出: true
}
