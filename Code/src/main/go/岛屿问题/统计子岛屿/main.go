package main

import "fmt"

func main() {
	// --- 示例 1 ---
	fmt.Println("--- 示例 1 ---")
	grid1_ex1 := [][]int{
		{1, 1, 1, 0, 0},
		{0, 1, 1, 1, 1},
		{0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0},
		{1, 1, 0, 1, 1},
	}
	// 注意：因为函数会修改 grid2，所以我们传入一个副本或者重新声明
	grid2_ex1 := [][]int{
		{1, 1, 1, 0, 0},
		{0, 0, 1, 1, 1},
		{0, 1, 0, 0, 0},
		{1, 0, 1, 1, 0},
		{0, 1, 0, 1, 0},
	}
	fmt.Println("输入 grid1: [[1,1,1,0,0],[0,1,1,1,1],[0,0,0,0,0],[1,0,0,0,0],[1,1,0,1,1]]")
	fmt.Println("输入 grid2: [[1,1,1,0,0],[0,0,1,1,1],[0,1,0,0,0],[1,0,1,1,0],[0,1,0,1,0]]")
	result1 := countSubIslands(grid1_ex1, grid2_ex1)
	fmt.Printf("输出: %d\n", result1)
	fmt.Println("预期输出: 3")
	fmt.Println()

	// --- 示例 2 ---
	fmt.Println("--- 示例 2 ---")
	grid1_ex2 := [][]int{
		{1, 0, 1, 0, 1},
		{1, 1, 1, 1, 1},
		{0, 0, 0, 0, 0},
		{1, 1, 1, 1, 1},
		{1, 0, 1, 0, 1},
	}
	grid2_ex2 := [][]int{
		{0, 0, 0, 0, 0},
		{1, 1, 1, 1, 1},
		{0, 1, 0, 1, 0},
		{0, 1, 0, 1, 0},
		{1, 0, 0, 0, 1},
	}
	fmt.Println("输入 grid1: [[1,0,1,0,1],[1,1,1,1,1],[0,0,0,0,0],[1,1,1,1,1],[1,0,1,0,1]]")
	fmt.Println("输入 grid2: [[0,0,0,0,0],[1,1,1,1,1],[0,1,0,1,0],[0,1,0,1,0],[1,0,0,0,1]]")
	result2 := countSubIslands(grid1_ex2, grid2_ex2)
	fmt.Printf("输出: %d\n", result2)
	fmt.Println("预期输出: 2")
}

/**
 * @param grid1 第一个 m x n 的二进制矩阵
 * @param grid2 第二个 m x n 的二进制矩阵
 * @return grid2 中子岛屿的数量
 */
func countSubIslands(grid1 [][]int, grid2 [][]int) int {
	// 获取网格的行数和列数
	m, n := len(grid1), len(grid1[0])
	// 初始化子岛屿计数器
	subIslandCount := 0

	// 遍历 grid2 的每一个格子
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			// 当找到一个岛屿的起点 (值为 1)
			if grid2[r][c] == 1 {
				// 假设当前岛屿是一个子岛屿
				isSub := true
				// 通过深度优先搜索 (DFS) 遍历整个岛屿，并检查其是否为子岛屿
				dfs(grid1, grid2, r, c, m, n, &isSub)
				// 如果遍历结束后，isSub 仍为 true，则计数器加一
				if isSub {
					subIslandCount++
				}
			}
		}
	}

	return subIslandCount
}

/**
 * 深度优先搜索 (DFS) 函数
 * @param grid1 第一个网格
 * @param grid2 第二个网格 (会被修改以标记访问过的节点)
 * @param r 当前行
 * @param c 当前列
 * @param m 总行数
 * @param n 总列数
 * @param isSub 指向布尔值的指针，用于记录当前岛屿是否为子岛屿
 */
func dfs(grid1 [][]int, grid2 [][]int, r, c, m, n int, isSub *bool) {
	// 边界检查：如果超出网格范围，或者是水域 (值为 0)，则返回
	if r < 0 || r >= m || c < 0 || c >= n || grid2[r][c] == 0 {
		return
	}

	// 将 grid2 的当前格子置为 0，标记为已访问，防止重复计算
	grid2[r][c] = 0

	// 关键检查：如果 grid2 的陆地在 grid1 对应位置是水域
	// 那么当前岛屿就不是子岛屿。
	if grid1[r][c] == 0 {
		*isSub = false
	}

	// 向四个方向递归探索
	dfs(grid1, grid2, r+1, c, m, n, isSub) // 下
	dfs(grid1, grid2, r-1, c, m, n, isSub) // 上
	dfs(grid1, grid2, r, c+1, m, n, isSub) // 右
	dfs(grid1, grid2, r, c-1, m, n, isSub) // 左
}
