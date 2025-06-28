package main

import "fmt"

// maxAreaOfIsland 是主函数，计算并返回最大岛屿面积
func maxAreaOfIsland(grid [][]int) int {
	// 边界情况：如果网格为空，则没有岛屿，面积为 0
	if len(grid) == 0 {
		return 0
	}

	rows := len(grid)
	cols := len(grid[0])
	maxArea := 0

	// 定义 dfs 函数，它会返回从 (r, c) 开始的岛屿面积
	var dfs func(r, c int) int
	dfs = func(r, c int) int {
		// 1. 检查边界或遇到水域，此处对岛屿面积无贡献，返回 0
		if r < 0 || r >= rows || c < 0 || c >= cols || grid[r][c] == 0 {
			return 0
		}

		// 2. 将当前陆地 '1' 修改为 '0'，表示已经访问过（淹没）
		// 这是为了防止无限递归和重复计算
		grid[r][c] = 0

		// 3. 计算面积：当前单元格面积为 1，
		// 再加上其四个方向相邻陆地的面积之和
		currentArea := 1 +
			dfs(r+1, c) + // 下
			dfs(r-1, c) + // 上
			dfs(r, c+1) + // 右
			dfs(r, c-1) // 左

		return currentArea
	}

	// 遍历整个网格
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			// 如果发现一个 '1'，说明找到了一个新岛屿
			if grid[r][c] == 1 {
				// 调用 dfs 计算当前岛屿的面积
				currentArea := dfs(r, c)
				// 更新全局最大面积
				if currentArea > maxArea {
					maxArea = currentArea
				}
			}
		}
	}

	return maxArea
}

func main() {
	// 题目中给出的示例
	grid := [][]int{
		{0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0},
		{0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0},
		{0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
	}

	fmt.Printf("示例网格的最大岛屿面积是: %d\n", maxAreaOfIsland(grid)) // 输出: 6
}
