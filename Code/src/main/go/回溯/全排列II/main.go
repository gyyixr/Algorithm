package main

import (
	"fmt"
	"sort"
)

var (
	res  [][]int
	path []int
	st   []bool // state的缩写
)

func permuteUnique(nums []int) [][]int {
	res, path = make([][]int, 0), make([]int, 0, len(nums))
	st = make([]bool, len(nums))
	sort.Ints(nums)
	dfs(nums, 0)
	return res
}

func dfs(nums []int, cur int) {
	if cur == len(nums) {
		tmp := make([]int, len(path))
		copy(tmp, path)
		res = append(res, tmp)
	}
	for i := 0; i < len(nums); i++ {
		if i != 0 && nums[i] == nums[i-1] && !st[i-1] { // 去重，用st来判别是深度还是广度
			continue
		}
		if !st[i] {
			path = append(path, nums[i])
			st[i] = true
			dfs(nums, cur+1)
			st[i] = false
			path = path[:len(path)-1]
		}
	}
}

func main() {
	input := []int{1, 1, 2}
	fmt.Println(permuteUnique(input))
}
