//最长回文字串
package main

import "fmt"

// 中心拓展法
func longestPalindromeDP(s string) string {
	if len(s) < 1 {
		return ""
	}

	start, end := 0, 0
	for i := 0; i < len(s); i++ {
		// 检查奇数长度的回文
		len1 := expandAroundCenter(s, i, i)
		// 检查偶数长度的回文
		len2 := expandAroundCenter(s, i, i+1)
		// 获取两者中较大的长度
		maxLen := max(len1, len2)
		if maxLen > end-start {
			start = i - (maxLen-1)/2
			end = i + maxLen/2
		}
	}
	return s[start : end+1]
}

func expandAroundCenter(s string, left, right int) int {
	for left >= 0 && right < len(s) && s[left] == s[right] {
		left--
		right++
	}
	return right - left - 1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	fmt.Println(longestPalindromeDP("babad")) // 输出 "bab" 或 "aba"
	fmt.Println(longestPalindromeDP("cbbcd")) // 输出 "bb"
}
