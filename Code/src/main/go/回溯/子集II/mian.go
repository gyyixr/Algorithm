package main

import (
	"fmt"
	"sort"
)

var (
	result [][]int
	path   []int
)

func subsetsWithDup(nums []int) [][]int {
	result = make([][]int, 0)
	path = make([]int, 0)
	used := make([]bool, len(nums))
	sort.Ints(nums) // 去重需要排序
	backtracing(nums, 0, used)
	return result
}

func backtracing(nums []int, startIndex int, used []bool) {
	tmp := make([]int, len(path))
	copy(tmp, path)
	result = append(result, tmp)
	for i := startIndex; i < len(nums); i++ {
		if i != startIndex && nums[i] == nums[i-1] && used[i-1] == false {
			continue
		}
		path = append(path, nums[i])
		used[i] = true
		backtracing(nums, i+1, used)
		path = path[:len(path)-1]
		used[i] = false
	}
}

func main() {
	input := []int{1, 1, 2}
	fmt.Println(subsetsWithDup(input))
}
