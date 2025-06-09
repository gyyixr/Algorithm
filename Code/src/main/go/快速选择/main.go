package main

import (
	"fmt"
	"math/rand"
	"time"
)

// findKthLargest: 寻找切片中第 k 大的元素
// 使用快速选择算法，平均时间复杂度 O(n)
func findKthLargest(nums []int, k int) int {
	// 在一个升序排序的数组中，第 k 大的元素的索引是 len(nums) - k
	targetIndex := len(nums) - k
	left, right := 0, len(nums)-1

	for {
		// 对 [left, right] 区间进行分区操作，返回基准值的最终索引
		pivotIndex := partition(nums, left, right)

		if pivotIndex == targetIndex {
			// 如果基准值的索引正好是目标索引，说明找到了第 k 大的元素
			return nums[pivotIndex]
		} else if pivotIndex < targetIndex {
			// 如果基准值索引偏小，说明目标元素在右侧部分，更新左边界
			left = pivotIndex + 1
		} else {
			// 如果基准值索引偏大，说明目标元素在左侧部分，更新右边界
			right = pivotIndex - 1
		}
	}
}

// partition: 分区函数 (使用 Lomuto 分区方案)
// 它会为 nums 的 [left, right] 子切片选择一个基准值，
// 并将所有小于基准值的元素放到它左边，大于等于的放到右边。
// 最后返回基准值在新切片中的索引。
func partition(nums []int, left, right int) int {
	// 1. 为了避免最坏情况 (O(n^2))，随机选择一个基准值
	// rand.Intn(n) 返回 [0, n) 之间的随机整数
	pivotIndex := left + rand.Intn(right-left+1)
	pivotValue := nums[pivotIndex]

	// 2. 将基准值移动到区间的末尾，方便处理
	nums[pivotIndex], nums[right] = nums[right], nums[pivotIndex]

	// 3. 遍历区间，将所有小于基准值的元素移动到左侧
	storeIndex := left
	for i := left; i < right; i++ {
		if nums[i] < pivotValue {
			nums[storeIndex], nums[i] = nums[i], nums[storeIndex]
			storeIndex++
		}
	}

	// 4. 遍历结束后，storeIndex 的位置就是基准值应该在的位置
	// 将基准值放回其最终位置
	nums[storeIndex], nums[right] = nums[right], nums[storeIndex]

	return storeIndex
}

// main 函数用于测试
func main() {
	// 初始化随机数种子，确保每次运行的随机性
	// 对于 Go 1.20 及以上版本，这一步可以省略，因为默认种子是安全的。
	// 但为了兼容旧版本和明确起见，通常会加上。
	rand.Seed(time.Now().UnixNano())

	// --- 测试用例 1 ---
	nums1 := []int{3, 2, 1, 5, 6, 4}
	k1 := 2
	fmt.Printf("测试用例 1:\n")
	fmt.Printf("原始数组: %v, k = %d\n", nums1, k1)
	result1 := findKthLargest(nums1, k1)
	fmt.Printf("找到第 %d 大的元素是: %d\n", k1, result1)
	fmt.Println("--------------------")

	// --- 测试用例 2 ---
	nums2 := []int{3, 2, 3, 1, 2, 4, 5, 5, 6}
	k2 := 4
	fmt.Printf("测试用例 2:\n")
	fmt.Printf("原始数组: %v, k = %d\n", nums2, k2)
	result2 := findKthLargest(nums2, k2)
	fmt.Printf("找到第 %d 大的元素是: %d\n", k2, result2)
	fmt.Println("--------------------")

	// --- 额外测试用例 ---
	nums3 := []int{7, 6, 5, 4, 3, 2, 1}
	k3 := 5
	fmt.Printf("额外测试用例:\n")
	fmt.Printf("原始数组: %v, k = %d\n", nums3, k3)
	result3 := findKthLargest(nums3, k3)
	fmt.Printf("找到第 %d 大的元素是: %d\n", k3, result3)
}
