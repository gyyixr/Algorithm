package main

import (
	"fmt"
	"sort"
)

var (
	res  [][]int
	path []int
	used []bool
)

func combinationSum2(candidates []int, target int) [][]int {
	res, path = make([][]int, 0), make([]int, 0, len(candidates))
	used = make([]bool, len(candidates))
	sort.Ints(candidates) // 排序，为剪枝做准备
	dfs(candidates, 0, target)
	return res
}

func dfs(candidates []int, start int, target int) {
	if target == 0 { // target 不断减小，如果为0说明达到了目标值
		tmp := make([]int, len(path))
		copy(tmp, path)
		res = append(res, tmp)
		return
	}
	for i := start; i < len(candidates); i++ {
		if candidates[i] > target { // 剪枝，提前返回
			break
		}
		if i > 0 && candidates[i] == candidates[i-1] && used[i-1] == false {
			continue
		}
		path = append(path, candidates[i])
		used[i] = true
		dfs(candidates, i+1, target-candidates[i])
		used[i] = false
		path = path[:len(path)-1]
	}
}

func main() {
	input := []int{10, 1, 2, 7, 6, 1, 5}
	fmt.Println(combinationSum2(input, 8))
}
