package main

import "fmt"

func solveQuiz(n int, k int) int {
	// dp[i][j]: 表示回答了 i 个问题后得到分数 j 的方案数。
	// 分数最高可以达到 n。
	dp := make([][]int, n+1)
	for r := 0; r <= n; r++ {
		dp[r] = make([]int, n+1) // 分数列需要存储到 n
	}

	// 基本情况：回答 0 个问题，得到 0 分的方案数为 1 (什么都不做)。
	dp[0][0] = 1

	// currentPowerOf2 用于计算 dp[i][0]，表示 2^(i-1)
	currentPowerOf2 := 1
	for i := 1; i <= n; i++ { // i 是已回答问题的数量
		// 情况 1：第 i 个问题回答错误，分数变为 0。
		// dp[i][0] 是回答 i-1 个问题（任何结果）的方案总数，即 2^(i-1)。
		dp[i][0] = currentPowerOf2
		currentPowerOf2 *= 2 // 为下一次迭代 (计算 2^i) 做准备

		// 情况 2：第 i 个问题回答正确。
		// 要在回答 i 个问题后得到分数 j (且第 i 个问题正确)，
		// 那么在回答 i-1 个问题后分数必须是 j-1。
		for j := 1; j <= i; j++ { // 回答 i 个问题得到的分数 j 最高为 i。
			// dp[i-1][j-1] 总是有效的，因为这里 j-1 >= 0。
			dp[i][j] = dp[i-1][j-1]
		}
	}

	// 通过累加得到分数 k 的方案数来计算最终答案。
	finalAns := 0
	if k == 0 {
		// 如果目标分数是 0，则累加回答 0 到 n 个问题后得到 0 分的方案数。
		for i := 0; i <= n; i++ {
			finalAns += dp[i][0]
		}
	} else if k > 0 {
		// 如果目标分数 k > 0，我们至少需要回答 k 个问题。
		// 并且 k 不能超过 n。
		if k <= n {
			for i := k; i <= n; i++ {
				finalAns += dp[i][k]
			}
		}
		// 如果 k > n (且 k 不为 0)，finalAns 保持为 0，这是正确的。
	}
	// 如果 k < 0 (通常分数不会这样)，finalAns 保持为 0。

	return finalAns
}

func main() {
	// 示例用法：
	// n: 可用的总问题数
	// k: 目标分数

	var n, k int
	fmt.Println("请输入总问题数 (n):") // Enter the total number of questions (n):
	fmt.Scanln(&n)
	fmt.Println("请输入目标分数 (k):") // Enter the target score (k):
	fmt.Scanln(&k)

	result := solveQuiz(n, k)
	fmt.Printf("在 %d 个问题中得到分数 %d 的方案数: %d\n", n, k, result) // Number of ways to get score %d with %d questions: %d

	// 来自思考过程的测试用例：
	// fmt.Println(solveQuiz(3, 1)) // 期望结果: 4
	// fmt.Println(solveQuiz(3, 0)) // 期望结果: 8
	// fmt.Println(solveQuiz(2, 2)) // 期望结果: 1
	// fmt.Println(solveQuiz(1, 2)) // 期望结果: 0 (k > n)
	// fmt.Println(solveQuiz(0, 0)) // 期望结果: 1
	// fmt.Println(solveQuiz(0, 1)) // 期望结果: 0
}
