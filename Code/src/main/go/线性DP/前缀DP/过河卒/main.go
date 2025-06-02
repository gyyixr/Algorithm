package main

import "fmt"

// solveRiverCrossingPawn 计算过河卒问题的方案数
func solveRiverCrossingPawn(numRows, numCols int, knightR, knightC int) int64 {
	if numRows <= 0 || numCols <= 0 {
		return 0
	}

	// 1. 创建并标记被阻挡的格子
	isBlocked := make([][]bool, numRows)
	for i := range isBlocked {
		isBlocked[i] = make([]bool, numCols) // 默认都是false (未阻挡)
	}

	// 仅当马在棋盘上时，才标记其自身位置和攻击位置
	if knightR >= 0 && knightR < numRows && knightC >= 0 && knightC < numCols {
		// 马的8个移动方向 (dr, dc)
		knightMoves := []struct{ dr, dc int }{
			{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2},
			{1, -2}, {1, 2}, {2, -1}, {2, 1},
		}

		// 标记马本身的位置
		isBlocked[knightR][knightC] = true

		// 标记马能攻击到的位置
		for _, move := range knightMoves {
			r, c := knightR+move.dr, knightC+move.dc
			// 确保攻击点在棋盘内
			if r >= 0 && r < numRows && c >= 0 && c < numCols {
				isBlocked[r][c] = true
			}
		}
	}
	// 如果马在棋盘外，isBlocked 数组将保持全为 false，这是正确的。

	// 2. 初始化DP表
	dp := make([][]int64, numRows)
	for i := range dp {
		dp[i] = make([]int64, numCols)
	}

	// 3. 填充DP表
	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			if isBlocked[r][c] { // 如果当前格子被阻挡
				dp[r][c] = 0
				continue // 跳过对此格子的路径计算
			}

			// 对于未被阻挡的格子
			if r == 0 && c == 0 {
				// 起始点，如果未被阻挡，则有1种方式到达（即开始就在这里）
				dp[r][c] = 1
			} else if r == 0 && c != 0 {
				dp[r][c] = dp[r][c-1] // 第一行，只能从左边来
			} else if r != 0 && c == 0 {
				dp[r][c] = dp[r-1][c] // 第一列，只能从上边来
			} else {
				dp[r][c] = dp[r-1][c] + dp[r][c-1] // 其他格子，可以从上边或左边来
			}
		}
	}

	// 4. 结果是到达右下角 (numRows-1, numCols-1) 的方案数
	return dp[numRows-1][numCols-1]
}

func main() {
	// 示例：洛谷 P1002 过河卒 (题目中坐标是1-indexed，且目标点和马点是输入)
	// 此处我们假设棋盘大小固定，马位置固定，目标是右下角
	// 例如，一个 3x3 的棋盘，卒子 (0,0) -> (2,2)
	// 马在 (1,1)

	numRows, numCols := 6, 6
	knightR, knightC := 2, 2
	fmt.Printf("棋盘大小: %d x %d\n", numRows, numCols)
	fmt.Printf("卒子从 (0,0) 到 (%d,%d)\n", numRows-1, numCols-1)
	fmt.Printf("马的位置: (%d,%d)\n", knightR, knightC)

	ways := solveRiverCrossingPawn(numRows, numCols, knightR, knightC)
	fmt.Printf("到达右下角的方案数: %d\n", ways) // 对于马(2,2)在6x6棋盘, 路径被有效阻断，0是可能的
	fmt.Println("---")

	knightR2, knightC2 := 0, 0
	fmt.Printf("棋盘大小: %d x %d\n", numRows, numCols)
	fmt.Printf("卒子从 (0,0) 到 (%d,%d)\n", numRows-1, numCols-1)
	fmt.Printf("马的位置: (%d,%d)\n", knightR2, knightC2)
	ways2 := solveRiverCrossingPawn(numRows, numCols, knightR2, knightC2)
	fmt.Printf("到达右下角的方案数 (马在起点): %d\n", ways2) // 预期 0
	fmt.Println("---")

	targetR, targetC := 5, 5
	kR, kC := 2, 2

	boardRows := targetR + 1
	boardCols := targetC + 1

	fmt.Printf("目标点: (%d,%d), 马位置: (%d,%d)\n", targetR, targetC, kR, kC)
	fmt.Printf("等效棋盘大小: %d x %d\n", boardRows, boardCols)
	ways_luogu_style := solveRiverCrossingPawn(boardRows, boardCols, kR, kC)
	fmt.Printf("到达目标点的方案数: %d\n", ways_luogu_style) // 同第一个例子
	fmt.Println("---")

	fmt.Printf("目标点: (1,1), 马位置: (-1,-1) (棋盘外)\n")
	// 棋盘大小为 2x2 (numRows=targetR+1=2, numCols=targetC+1=2)
	ways_simple := solveRiverCrossingPawn(2, 2, -1, -1)
	fmt.Printf("到达目标点(1,1)的方案数 (无有效障碍): %d\n", ways_simple) // 修正后预期 2
	fmt.Println("---")

	fmt.Printf("目标点: (2,2), 马位置: (1,2)\n") // 棋盘大小 3x3
	// 马(1,2)攻击(0,0)因为(1-1, 2-2)=(0,0)
	ways_start_attacked := solveRiverCrossingPawn(3, 3, 1, 2)
	fmt.Printf("到达目标点(2,2)的方案数 (起点被攻击): %d\n", ways_start_attacked) // 预期 0
	fmt.Println("---")
}
