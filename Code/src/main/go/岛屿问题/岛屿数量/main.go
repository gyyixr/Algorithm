package main

import "fmt"

// numIslands 是主函数，用于计算岛屿的数量
func numIslands(grid [][]byte) int {
	// 处理边界情况，如果网格为空，则没有岛屿
	if len(grid) == 0 {
		return 0
	}

	rows := len(grid)    // 网格的行数
	cols := len(grid[0]) // 网格的列数
	numIslands := 0      // 初始化岛屿数量为 0

	// 定义一个 dfs 函数（深度优先搜索）
	// 这个函数会"淹没"（标记为'0'）与 (r, c) 相邻的所有陆地
	var dfs func(r, c int)
	dfs = func(r, c int) {
		// 1. 检查边界条件和当前位置是否为水域
		// 如果行或列越界，或者当前已经是水('0')，则直接返回
		if r < 0 || r >= rows || c < 0 || c >= cols || grid[r][c] == '0' {
			return
		}

		// 2. 将当前陆地 '1' 修改为 '0'，表示已经访问过
		// 这是避免重复计数的关键
		grid[r][c] = '0'

		// 3. 递归地访问当前位置的四个方向：上、下、左、右
		dfs(r+1, c) // 下
		dfs(r-1, c) // 上
		dfs(r, c+1) // 右
		dfs(r, c-1) // 左
	}

	// 遍历整个网格
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			// 如果发现一个 '1'，说明找到了一个新岛屿
			if grid[r][c] == '1' {
				numIslands++ // 岛屿数量加 1
				dfs(r, c)    // 调用 dfs 将整个岛屿"淹没"，防止重复计数
			}
		}
	}

	return numIslands
}

func main() {
	// 示例 1
	grid1 := [][]byte{
		{'1', '1', '1', '1', '0'},
		{'1', '1', '0', '1', '0'},
		{'1', '1', '0', '0', '0'},
		{'0', '0', '0', '0', '0'},
	}
	fmt.Printf("示例 1 网格的岛屿数量是: %d\n", numIslands(grid1)) // 输出: 1

	// 示例 2
	grid2 := [][]byte{
		{'1', '1', '0', '0', '0'},
		{'1', '1', '0', '0', '0'},
		{'0', '0', '1', '0', '0'},
		{'0', '0', '0', '1', '1'},
	}
	fmt.Printf("示例 2 网格的岛屿数量是: %d\n", numIslands(grid2)) // 输出: 3
}
