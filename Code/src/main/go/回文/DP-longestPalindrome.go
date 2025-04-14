//最长回文字串-DP
package main

import "fmt"

func longestPalindrome(s string) string {
	length := len(s)
	if length < 2 {
		return s
	}

	maxLen := 1
	begin := 0
	// dp[i][j] 表示 s[i..j] 是否是回文串
	dp := make([][]bool, length)
	// 初始化：所有长度为 1 的子串都是回文串
	for i := 0; i < length; i++ {
		dp[i] = make([]bool, length)
		dp[i][i] = true
	}

	// 转换为字符数组便于访问
	charArray := []rune(s)

	// 递推开始
	// 先枚举子串长度
	for L := 2; L <= length; L++ {
		// 枚举左边界
		for i := 0; i < length; i++ {
			// 由 L 和 i 可以确定右边界，即 j - i + 1 = L 得
			j := L + i - 1
			// 如果右边界越界，就可以退出当前循环
			if j >= length {
				break
			}

			if charArray[i] != charArray[j] {
				dp[i][j] = false
			} else {
				if j-i < 3 {
					dp[i][j] = true
				} else {
					dp[i][j] = dp[i+1][j-1]
				}
			}

			// 只要 dp[i][j] == true 成立，就表示子串 s[i..j] 是回文，此时记录回文长度和起始位置
			if dp[i][j] && j-i+1 > maxLen {
				maxLen = j - i + 1
				begin = i
			}
		}
	}
	return s[begin : begin+maxLen]
}

func main() {
	testCases := []struct {
		input  string
		expect string
	}{
		{"babad", "bab"}, // 或 "aba"
		{"cbbd", "bb"},
		{"a", "a"},
		{"ac", "a"}, // 或 "c"
		{"aaaa", "aaaa"},
		{"abcba", "abcba"},
		{"abacdfgdcaba", "aba"},
	}

	for _, tc := range testCases {
		result := longestPalindrome(tc.input)
		fmt.Printf("输入: %-12s 输出: %-12s 预期: %-12s 匹配: %v\n",
			tc.input, result, tc.expect, result == tc.expect)
	}
}
