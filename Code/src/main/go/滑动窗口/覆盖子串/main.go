package main

import (
	"fmt"
	"math"
)

func minWindow(s string, t string) string {
	// 基本情况处理：如果 s 或 t 为空，或者 s 的长度小于 t，则不可能找到覆盖子串
	if len(s) == 0 || len(t) == 0 || len(s) < len(t) {
		return ""
	}

	// targetFreq 用于存储字符串 t 中字符的频率
	// 使用 [128]int 假设字符集为 ASCII。如果需要支持所有 byte 值 (0-255)，应使用 [256]int
	var targetFreq [128]int // 例如：targetFreq['A'] = 1, targetFreq['B'] = 1
	// windowFreq 用于存储当前滑动窗口中字符的频率
	var windowFreq [128]int

	// 统计字符串 t 中每个字符出现的次数
	for i := 0; i < len(t); i++ {
		targetFreq[t[i]]++
	}

	left, right := 0, 0 // 滑动窗口的左右指针
	formed := 0         // 当前窗口中已经满足 t 中字符种类及数量要求的字符种类数
	required := 0       // 字符串 t 中不同字符的种类数 (即 targetFreq 中频率大于0的字符个数)

	// 计算 t 中有多少种不同的字符需要匹配
	for i := 0; i < 128; i++ { // 如果使用 [256]int，这里应为 256
		if targetFreq[i] > 0 {
			required++
		}
	}

	// 如果 t 中没有需要匹配的字符 (例如 t 是空字符串，虽然前面有 len(t) == 0 的判断)
	if required == 0 {
		return ""
	}

	minLength := math.MaxInt32 // 记录最小覆盖子串的长度，初始化为最大整数
	minLeft, minRight := 0, 0  // 记录最小覆盖子串的起始和结束索引（左闭右闭）
	found := false             // 标记是否找到了有效的覆盖子串

	// 开始滑动窗口
	for right < len(s) {
		// 获取右指针指向的字符
		charRight := s[right]

		// 将右侧字符加入窗口，并更新窗口内该字符的频率
		// 注意: 如果 s 中的字符超出了 targetFreq/windowFreq 数组的索引范围 (例如用了[128]int但s中有非ASCII字符)，这里会 panic。
		// 若使用 [256]int 则 byte 类型的值 (0-255) 都在范围内。
		windowFreq[charRight]++

		// 检查新加入的字符是否是 t 中的目标字符，并且其在窗口中的数量是否达到了 t 中的要求
		// targetFreq[charRight] > 0 表示 charRight 是 t 中的一个目标字符
		// windowFreq[charRight] == targetFreq[charRight] 表示窗口中 charRight 的数量已经满足 t 的要求
		if targetFreq[charRight] > 0 && windowFreq[charRight] == targetFreq[charRight] {
			formed++ // 满足条件的字符种类数增加
		}

		// 当窗口中已形成的字符种类数等于 t 中要求的字符种类数时，尝试收缩窗口左边界
		for left <= right && formed == required {
			found = true                      // 找到了一个有效的覆盖子串
			currentLength := right - left + 1 // 当前窗口的长度

			// 如果当前窗口长度小于之前记录的最小长度，则更新最小长度和对应的窗口边界
			if currentLength < minLength {
				minLength = currentLength
				minLeft = left
				minRight = right
			}

			// 准备将左指针指向的字符移出窗口
			charLeft := s[left]
			windowFreq[charLeft]-- // 窗口中该字符的频率减一

			// 如果移出的字符是 t 中的目标字符，并且移出后其在窗口中的数量不再满足 t 的要求
			if targetFreq[charLeft] > 0 && windowFreq[charLeft] < targetFreq[charLeft] {
				formed-- // 满足条件的字符种类数减少，表示当前窗口不再是有效覆盖子串
			}

			// 左指针右移，缩小窗口
			left++
		}

		// 右指针右移，扩大窗口，继续寻找
		right++
	}

	// 如果没有找到任何有效的覆盖子串
	if !found {
		return ""
	}
	// 返回找到的最小覆盖子串 (注意切片是左闭右开，所以是 minLeft 到 minRight+1)
	return s[minLeft : minRight+1]
}

func main() {
	s1 := "ADOBECODEBANC"
	t1 := "ABC"
	fmt.Printf("输入: s = \"%s\", t = \"%s\"\n输出: \"%s\"\n", s1, t1, minWindow(s1, t1)) // 预期: "BANC"

	s2 := "a"
	t2 := "a"
	fmt.Printf("输入: s = \"%s\", t = \"%s\"\n输出: \"%s\"\n", s2, t2, minWindow(s2, t2)) // 预期: "a"

	s3 := "a"
	t3 := "aa"
	fmt.Printf("输入: s = \"%s\", t = \"%s\"\n输出: \"%s\"\n", s3, t3, minWindow(s3, t3)) // 预期: ""

	s4 := "AB"
	t4 := "A"
	fmt.Printf("输入: s = \"%s\", t = \"%s\"\n输出: \"%s\"\n", s4, t4, minWindow(s4, t4)) // 预期: "A"

	s5 := "BBAAC"
	t5 := "BAC"
	fmt.Printf("输入: s = \"%s\", t = \"%s\"\n输出: \"%s\"\n", s5, t5, minWindow(s5, t5)) // 预期: "BAC"

	s6 := "cabwefgewcwaefgcf"
	t6 := "cae"
	fmt.Printf("输入: s = \"%s\", t = \"%s\"\n输出: \"%s\"\n", s6, t6, minWindow(s6, t6)) // 预期: "cwae"
}
