package main

import "fmt"

// findMedianSortedArrays 函数用于计算两个已排序数组的中位数。
func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	n := len(nums1)
	m := len(nums2)

	// left 和 right 分别代表合并后数组的两个中间位置。
	// 如果总长度是奇数，left 和 right 的值会相同。
	// 如果总长度是偶数，它们将指向中间的两个数。
	left := (n + m + 1) / 2
	right := (n + m + 2) / 2

	// 将奇数和偶数的情况合并处理。
	// 如果是奇数，会调用两次 getKth 寻找同一个位置的元素。
	// 如果是偶数，则寻找中间两个位置的元素。
	// 最后取两者的平均值。
	val1 := getKth(nums1, 0, n-1, nums2, 0, m-1, left)
	val2 := getKth(nums1, 0, n-1, nums2, 0, m-1, right)
	return float64(val1+val2) * 0.5
}

// getKth 函数在两个有序数组的指定范围内查找第 k 小的元素。
func getKth(nums1 []int, start1 int, end1 int, nums2 []int, start2 int, end2 int, k int) int {
	len1 := end1 - start1 + 1
	len2 := end2 - start2 + 1

	// 保证 len1 始终是较短的数组长度，简化后续处理。
	// 如果 len1 > len2，则交换两个数组的参数位置进行递归。
	if len1 > len2 {
		return getKth(nums2, start2, end2, nums1, start1, end1, k)
	}

	// 如果短的数组已经没有元素，则直接从长的数组中返回第 k 个元素。
	if len1 == 0 {
		return nums2[start2+k-1]
	}

	// 当 k=1 时，表示要找两个数组当前范围内的第一个（即最小的）元素。
	if k == 1 {
		return min(nums1[start1], nums2[start2])
	}

	// 计算两个数组的比较点。为了防止数组越界，取各自数组长度和 k/2 中的较小值。
	i := start1 + min(len1, k/2) - 1
	j := start2 + min(len2, k/2) - 1

	// 比较两个数组在比较点的元素值。
	if nums1[i] > nums2[j] {
		// 如果 nums1 的元素更大，说明 nums2 在 j 点之前的元素都不可能是第 k 小的，可以排除。
		// 更新 k 的值，减去被排除的元素数量。
		return getKth(nums1, start1, end1, nums2, j+1, end2, k-(j-start2+1))
	} else {
		// 反之，排除 nums1 在 i 点之前的元素。
		return getKth(nums1, i+1, end1, nums2, start2, end2, k-(i-start1+1))
	}
}

func main() {
	nums1 := []int{1, 3}
	nums2 := []int{2}
	fmt.Println(findMedianSortedArrays(nums1, nums2))
}

// min 是一个辅助函数，返回两个整数中的较小者。
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
