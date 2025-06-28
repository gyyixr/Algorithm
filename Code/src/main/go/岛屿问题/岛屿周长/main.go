package main

import "fmt"

// IslandPerimeter 计算岛屿周长
func IslandPerimeter(grid [][]int) int {
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[0]); c++ {
			if grid[r][c] == 1 {
				// 题目限制只有一个岛屿，计算一个即可
				return dfs(grid, r, c)
			}
		}
	}
	return 0
}

func dfs(grid [][]int, r, c int) int {
	// 函数因为「坐标 (r, c) 超出网格范围」返回，对应一条黄色的边
	if !inArea(grid, r, c) {
		return 1
	}
	// 函数因为「当前格子是海洋格子」返回，对应一条蓝色的边
	if grid[r][c] == 0 {
		return 1
	}
	// 函数因为「当前格子是已遍历的陆地格子」返回，和周长没关系
	if grid[r][c] != 1 {
		return 0
	}
	grid[r][c] = 2
	return dfs(grid, r-1, c) +
		dfs(grid, r+1, c) +
		dfs(grid, r, c-1) +
		dfs(grid, r, c+1)
}

// inArea 判断坐标 (r, c) 是否在网格中
func inArea(grid [][]int, r, c int) bool {
	return 0 <= r && r < len(grid) &&
		0 <= c && c < len(grid[0])
}

func main() {
	// 测试用例1：简单的小岛屿
	fmt.Println("测试用例1：")
	grid1 := [][]int{
		{0, 1, 0, 0},
		{1, 1, 1, 0},
		{0, 1, 0, 0},
		{1, 1, 0, 0},
	}
	result1 := IslandPerimeter(grid1)
	fmt.Printf("岛屿周长: %d\n\n", result1)

	// 测试用例2：单个格子的岛屿
	fmt.Println("测试用例2：")
	grid2 := [][]int{
		{1},
	}
	result2 := IslandPerimeter(grid2)
	fmt.Printf("岛屿周长: %d\n\n", result2)

	// 测试用例3：矩形岛屿
	fmt.Println("测试用例3：")
	grid3 := [][]int{
		{1, 1},
		{1, 1},
	}
	result3 := IslandPerimeter(grid3)
	fmt.Printf("岛屿周长: %d\n\n", result3)

	// 测试用例4：复杂形状的岛屿
	fmt.Println("测试用例4：")
	grid4 := [][]int{
		{1, 1, 1, 1, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 1, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 1, 1, 1, 1},
	}
	result4 := IslandPerimeter(grid4)
	fmt.Printf("岛屿周长: %d\n", result4)
}
