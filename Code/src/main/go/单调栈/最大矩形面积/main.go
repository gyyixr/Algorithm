package main

import "fmt"

func largestRectangleArea(heights []int) int {
	result := 0
	// 使用切片实现栈
	stack := make([]int, 0)
	// 数组头部加入0
	heights = append([]int{0}, heights...)
	// 数组尾部加入0
	heights = append(heights, 0)
	// 初始化栈，序号从0开始
	stack = append(stack, 0)
	for i := 1; i < len(heights); i++ {
		// 结束循环条件为：当即将入栈元素>top元素，也就是形成非单调递增的趋势
		for heights[stack[len(stack)-1]] > heights[i] {
			// mid 是top
			mid := stack[len(stack)-1]
			// 出栈
			stack = stack[0 : len(stack)-1]
			// left是top的下一位元素，i是将要入栈的元素
			left := stack[len(stack)-1]
			// 高度x宽度
			tmp := heights[mid] * (i - left - 1)
			if tmp > result {
				result = tmp
			}
		}
		stack = append(stack, i)
	}
	return result
}

func main() {
	input := []int{2, 1, 5, 6, 2, 3}
	fmt.Println(largestRectangleArea(input))
	fmt.Println(largestRectangleArea1(input))
}

// 双指针解法
func largestRectangleArea1(heights []int) int {
	n := len(heights)
	result := 0
	left, right := make([]int, n), make([]int, n)

	// 初始化left数组
	//for i := 0; i < n; i++ {
	//	left[i] = -1
	//}
	left[0] = -1

	// 初始化right数组
	//for i := 0; i < n; i++ {
	//	right[i] = n
	//}
	right[n-1] = n

	// 计算每个柱子左边第一个小于当前柱子的下标
	for i := 1; i < n; i++ {
		j := i - 1
		for j >= 0 && heights[j] >= heights[i] {
			j = left[j]
		}
		left[i] = j
	}

	// 计算每个柱子右边第一个小于当前柱子的下标
	for i := n - 2; i >= 0; i-- {
		j := i + 1
		for j < n && heights[j] >= heights[i] {
			j = right[j]
		}
		right[i] = j
	}

	// 计算每个柱子的最大矩形面积
	for i := 0; i < n; i++ {
		width := right[i] - left[i] - 1
		area := heights[i] * width
		if area > result {
			result = area
		}
	}

	return result
}
