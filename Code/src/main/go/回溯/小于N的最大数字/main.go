package main

import (
	"fmt"
	"strconv"
)

var maxNumberPath int

// nums: 可用整数（作为构建块）的切片。
// cur: 当前正在构建的数字字符串（新数字会被前置到此字符串）。
// target: 目标值，构建的数字必须小于此目标值。
func dfs(nums []int, cur string, target int) {
	// 如果当前已构建的字符串 cur 非空，则尝试将其作为数字处理。
	if cur != "" {
		num, _ := strconv.Atoi(cur) // 将字符串 cur 转换为整数 num (这里忽略了 Atoi 可能返回的错误)。

		// 如果转换后的数字 num 小于目标值 target。
		if num < target {
			// 如果 num 比当前记录的 maxNumberPath 更大，
			// 或者 maxNumberPath 尚未被有效数字更新过（仍为初始值 -1），则更新 maxNumberPath。
			if maxNumberPath == -1 || num > maxNumberPath {
				maxNumberPath = num
			}
		}
	}

	// 剪枝优化：获取目标值的字符串长度。
	targetStr := strconv.Itoa(target)
	targetLen := len(targetStr)

	// 如果当前构建的数字字符串 cur 的长度已经等于目标值 target 的长度，
	// 则没有必要再前置更多的数字了（因为那样会使数字更长，或者值更大，此处的逻辑是基于长度限制）。
	if len(cur) == targetLen {
		return
	}

	// 递归步骤：遍历 nums 中的每一个可用数字 n。
	for _, n := range nums {
		//newCur := cur + strconv.Itoa(n) // golang的字符串不可变，这里是隐藏式撤销选择
		dfs(nums, cur+strconv.Itoa(n), target) // 以 newCur 作为当前构建的数字字符串，继续进行递归搜索。
	}
}

// nums: 可用整数（构建块）的切片。
// target: 目标值。
func bruteForceMaxNumberGolang(nums []int, target int) int {
	maxNumberPath = -1 // 初始化/重置全局结果为 -1 (表示尚未找到符合条件的数)。

	dfs(nums, "", target) // 开始 DFS 搜索，初始的当前字符串 cur 为空。

	return maxNumberPath // 返回搜索过程中找到的最大符合条件的数。
}

// 示例用法 (如果需要运行，请取消 fmt 的导入注释以及此部分代码注释)
func main() {
	testCases := []struct {
		nums   []int
		target int
		name   string
	}{
		{[]int{1, 2, 3}, 200, "测试用例 1"},                    // 预期: 133
		{[]int{9}, 100, "测试用例 2"},                          // 预期: 99
		{[]int{1}, 1, "测试用例 3"},                            // 预期: -1 (没有数字小于1)
		{[]int{5, 2}, 230, "测试用例 4"},                       // 预期: 225
		{[]int{8}, 8, "测试用例 5"},                            // 预期: -1
		{[]int{7, 8, 9}, 80, "测试用例 6"},                     // 预期: 79
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 1000, "测试用例 7"}, // 预期: 999,
		{[]int{1, 1, 2}, 3, "测试用例 8"},                      // 预期: 2
	}

	for _, tc := range testCases {
		result := bruteForceMaxNumberGolang(tc.nums, tc.target)
		fmt.Printf("%s: Nums: %v, Target: %d, Max Found: %d\n", tc.name, tc.nums, tc.target, result)
	}
}
