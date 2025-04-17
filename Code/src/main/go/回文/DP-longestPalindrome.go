//最长回文字串-DP
package main

import "fmt"

func longestPalindrome(s string) string {
	maxLen := 0
	left := 0
	length := 0
	dp := make([][]bool, len(s))
	for i := 0; i < len(s); i++ {
		dp[i] = make([]bool, len(s))
		dp[i][i] = true
	}

	for i := len(s) - 1; i >= 0; i-- {
		for j := i; j < len(s); j++ {
			if s[i] == s[j] {
				if j-i <= 1 { // 情况一和情况二
					length = j - i
					dp[i][j] = true
				} else if dp[i+1][j-1] { // 情况三
					length = j - i
					dp[i][j] = true
				}
			}
		}
		if length > maxLen {
			maxLen = length
			left = i
		}
	}
	return s[left : left+maxLen+1]
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
