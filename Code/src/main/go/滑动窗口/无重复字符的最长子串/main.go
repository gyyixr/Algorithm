package main

import "fmt"

/*
给定一个字符串 s ，请你找出其中不含有重复字符的 最长 子串 的长度。
示例 1:

输入: s = "abcabcbb"
输出: 3
解释: 因为无重复字符的最长子串是 "abc"，所以其长度为 3。
示例 2:

输入: s = "bbbbb"
输出: 1
解释: 因为无重复字符的最长子串是 "b"，所以其长度为 1。
示例 3:

输入: s = "pwwkew"
输出: 3
解释: 因为无重复字符的最长子串是 "wke"，所以其长度为 3。
     请注意，你的答案必须是 子串 的长度，"pwke" 是一个子序列，不是子串。
*/

/*
实现原理
滑动窗口机制：使用两个指针 left 和 right 维护一个动态窗口
哈希表去重：使用 map[rune]bool 记录当前窗口中的字符
动态调整：当遇到重复字符时，收缩左边界；否则扩展右边界
*/
func lengthOfLongestSubstring(s string) int {
	chars := []rune(s)
	left, right, curLen, maxLen := 0, 0, 0, 0
	marked := make(map[rune]bool)

	for right < len(chars) {
		if _, ok := marked[chars[right]]; !ok {
			marked[chars[right]] = true
			curLen++
			maxLen = max(maxLen, curLen)
			right++
		} else {
			delete(marked, chars[left])
			left++
			curLen--
		}
	}

	return maxLen
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	fmt.Println(lengthOfLongestSubstring("abcabcbb"))
	fmt.Println(lengthOfLongestSubstring("bbbbb"))
	fmt.Println(lengthOfLongestSubstring("pwwkew"))
}
