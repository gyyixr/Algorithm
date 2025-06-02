package main

import (
	"fmt"
)

// maxSumSubarrayDPWithIndices 计算最大子段和，并返回和、起始索引、结束索引
func maxSumSubarrayDPWithIndices(w []int) (maxSum int, startIndex int, endIndex int) {
	n := len(w)
	// dp[i] 表示以 w[i] 结尾的最大子段和
	dp := make([]int, n)
	// currentStartIndices[i] 表示以 w[i] 结尾的最大子段的起始索引
	currentStartIndices := make([]int, n)

	// 初始化第一个元素的情况
	dp[0] = w[0]
	currentStartIndices[0] = 0

	overallMaxSum := dp[0]
	finalStartIndex := 0
	finalEndIndex := 0

	/**
	当 i = 0 时：f(0) = w[0] (以第一个元素结尾的最大子段和就是它本身)
	当 i > 0 时：f(i) = max(f(i-1) + w[i], w[i])
	*/

	for i := 1; i < n; i++ {
		if dp[i-1]+w[i] > w[i] {
			dp[i] = dp[i-1] + w[i]
			currentStartIndices[i] = currentStartIndices[i-1]
		} else {
			dp[i] = w[i]
			currentStartIndices[i] = i
		}

		if dp[i] > overallMaxSum {
			overallMaxSum = dp[i]
			finalStartIndex = currentStartIndices[i]
			finalEndIndex = i
		}
	}

	return overallMaxSum, finalStartIndex, finalEndIndex
}

// 为了对比，这里保留 Kadane's Algorithm 的原始版本（仅返回和）
func maxSumSubarray_kadane(w []int) int {
	n := len(w)

	maxSoFar := w[0]
	currentMax := w[0]

	for i := 1; i < n; i++ {
		if w[i] > currentMax+w[i] {
			currentMax = w[i]
		} else {
			currentMax = currentMax + w[i]
		}
		if currentMax > maxSoFar {
			maxSoFar = currentMax
		}
	}
	return maxSoFar
}

func main() {
	exampleArray := []int{2, -4, 3, -1, 2, -4, 3}
	fmt.Printf("数组: %v\n", exampleArray)

	maxSum, start, end := maxSumSubarrayDPWithIndices(exampleArray)
	fmt.Printf("最大子段和 (DP f(i) table): %d\n", maxSum)
	fmt.Printf("起始索引: %d, 结束索引: %d\n", start, end)
	if start <= end && start >= 0 { // 确保索引有效
		fmt.Printf("最大子段: %v\n", exampleArray[start:end+1])
	}
	fmt.Println("---")
	// 预期:
	// f(0)=2, start(0)=0. OverallMax=2, finalStart=0, finalEnd=0
	// f(1)=max(2-4, -4) = -2 (来自f(0)+w[1]), start(1)=0. OverallMax=2
	// f(2)=max(-2+3, 3) = 3 (来自w[2]), start(2)=2. OverallMax=3, finalStart=2, finalEnd=2
	// f(3)=max(3-1, -1) = 2 (来自f(2)+w[3]), start(3)=2. OverallMax=3
	// f(4)=max(2+2, 2) = 4 (来自f(3)+w[4]), start(4)=2. OverallMax=4, finalStart=2, finalEnd=4
	// f(5)=max(4-4, -4) = 0 (来自f(4)+w[5]), start(5)=2. OverallMax=4
	// f(6)=max(0+3, 3) = 3 (来自w[6]或f(5)+w[6],若f(5)+w[6]>=w[6]则选前者).
	//    若 sumEndingPrev (0+3=3) > w[6] (3) -> false.
	//    else: dp[6]=w[6]=3, currentStartIndices[6]=6. OverallMax=4
	// 结果: maxSum=4, start=2, end=4. 子段: [3, -1, 2]

	testArray2 := []int{-2, -3, -4, -1, -2, -1, -5, -3}
	fmt.Printf("数组: %v\n", testArray2)
	maxSum2, start2, end2 := maxSumSubarrayDPWithIndices(testArray2)
	fmt.Printf("最大子段和 (全负数): %d\n", maxSum2)
	fmt.Printf("起始索引: %d, 结束索引: %d\n", start2, end2)
	if start2 <= end2 && start2 >= 0 {
		fmt.Printf("最大子段: %v\n", testArray2[start2:end2+1]) // 预期: -1, start=3, end=3. 子段: [-1]
	}
	fmt.Println("---")

	testArray3 := []int{1, 2, 3, 4, 5}
	fmt.Printf("数组: %v\n", testArray3)
	maxSum3, start3, end3 := maxSumSubarrayDPWithIndices(testArray3)
	fmt.Printf("最大子段和 (全正数): %d\n", maxSum3)
	fmt.Printf("起始索引: %d, 结束索引: %d\n", start3, end3)
	if start3 <= end3 && start3 >= 0 {
		fmt.Printf("最大子段: %v\n", testArray3[start3:end3+1]) // 预期: 15, start=0, end=4. 子段: [1 2 3 4 5]
	}
	fmt.Println("---")

	testArray4 := []int{-1}
	fmt.Printf("数组: %v\n", testArray4)
	maxSum4, start4, end4 := maxSumSubarrayDPWithIndices(testArray4)
	fmt.Printf("最大子段和 (单个负数): %d\n", maxSum4)
	fmt.Printf("起始索引: %d, 结束索引: %d\n", start4, end4)
	if start4 <= end4 && start4 >= 0 {
		fmt.Printf("最大子段: %v\n", testArray4[start4:end4+1]) // 预期: -1, start=0, end=0. 子段: [-1]
	}
	fmt.Println("---")

	testArray5 := []int{5}
	fmt.Printf("数组: %v\n", testArray5)
	maxSum5, start5, end5 := maxSumSubarrayDPWithIndices(testArray5)
	fmt.Printf("最大子段和 (单个正数): %d\n", maxSum5)
	fmt.Printf("起始索引: %d, 结束索引: %d\n", start5, end5)
	if start5 <= end5 && start5 >= 0 {
		fmt.Printf("最大子段: %v\n", testArray5[start5:end5+1]) // 预期: 5, start=0, end=0. 子段: [5]
	}
	fmt.Println("---")

	// Kadane 算法的调用 (它不返回索引，只是为了完整性)
	// fmt.Printf("Kadane (只返回和) for exampleArray: %d\n", maxSumSubarray_kadane(exampleArray))
}
