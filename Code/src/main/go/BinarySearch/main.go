package main

import (
	"fmt"
	"sort" // 导入 sort 包，方便对比和实际使用
)

// BinarySearch 在一个有序的整数切片中查找指定的元素。
// 如果找到，则返回该元素的索引和 true。
// 如果未找到，则返回 -1 和 false。
func BinarySearch(slice []int, target int) (int, bool) {
	low := 0
	high := len(slice) - 1

	for low <= high {
		// 计算中间索引，这种方式可以避免 (low + high) 可能导致的整数溢出
		mid := low + (high-low)/2
		guess := slice[mid]

		if guess == target {
			return mid, true // 找到目标
		}
		if guess > target {
			high = mid - 1 // 目标在左半部分
		} else {
			low = mid + 1 // 目标在右半部分
		}
	}

	return -1, false // 未找到目标
}

func main() {
	// 示例1: 整数切片
	intSlice := []int{2, 5, 7, 8, 11, 12, 15, 18, 22, 25, 30}
	targetInt := 12
	index, found := BinarySearch(intSlice, targetInt)
	if found {
		fmt.Printf("整数切片: 在索引 %d 处找到 %d\n", index, targetInt)
	} else {
		fmt.Printf("整数切片: 未找到 %d\n", targetInt)
	}

	targetIntNotFound := 13
	index, found = BinarySearch(intSlice, targetIntNotFound)
	if found {
		fmt.Printf("整数切片: 在索引 %d 处找到 %d\n", index, targetIntNotFound)
	} else {
		fmt.Printf("整数切片: 未找到 %d\n", targetIntNotFound)
	}

	emptySlice := []int{}
	targetEmpty := 5
	index, found = BinarySearch(emptySlice, targetEmpty)
	if found {
		fmt.Printf("空切片: 在索引 %d 处找到 %d\n", index, targetEmpty)
	} else {
		fmt.Printf("空切片: 未找到 %d\n", targetEmpty)
	}

	singleElementSlice := []int{10}
	targetSingleFound := 10
	targetSingleNotFound := 20

	index, found = BinarySearch(singleElementSlice, targetSingleFound)
	if found {
		fmt.Printf("单元素切片: 在索引 %d 处找到 %d\n", index, targetSingleFound)
	} else {
		fmt.Printf("单元素切片: 未找到 %d\n", targetSingleFound)
	}
	index, found = BinarySearch(singleElementSlice, targetSingleNotFound)
	if found {
		fmt.Printf("单元素切片: 在索引 %d 处找到 %d\n", index, targetSingleNotFound)
	} else {
		fmt.Printf("单元素切片: 未找到 %d\n", targetSingleNotFound)
	}

	// sort.SearchInts 是专门为 []int 设计的
	targetStdLib := 12
	// sort.SearchInts 返回的是第一个大于等于 x 的索引，如果都小于 x，则返回 len(a)
	// 所以如果元素存在，它返回的是该元素的索引
	// 如果元素不存在，它返回的是元素应该插入的位置以保持排序
	i := sort.SearchInts(intSlice, targetStdLib)
	if i < len(intSlice) && intSlice[i] == targetStdLib {
		fmt.Printf("标准库 sort.SearchInts: 在索引 %d 处找到 %d\n", i, targetStdLib)
	} else {
		fmt.Printf("标准库 sort.SearchInts: 未找到 %d，建议插入索引为 %d\n", targetStdLib, i)
	}

	targetStdLibNotFound := 13
	i = sort.SearchInts(intSlice, targetStdLibNotFound)
	if i < len(intSlice) && intSlice[i] == targetStdLibNotFound {
		fmt.Printf("标准库 sort.SearchInts: 在索引 %d 处找到 %d\n", i, targetStdLibNotFound)
	} else {
		fmt.Printf("标准库 sort.SearchInts: 未找到 %d，建议插入索引为 %d\n", targetStdLibNotFound, i)
	}
}
