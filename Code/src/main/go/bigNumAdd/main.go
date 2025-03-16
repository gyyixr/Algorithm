package main

import (
	"fmt"
)

// reverseString 反转字符串
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// addLargeNumbers 手动实现大数相加
func addLargeNumbers(num1, num2 string) string {
	// 保证 num1 的长度不小于 num2
	if len(num1) < len(num2) {
		num1, num2 = num2, num1
	}

	// 反转字符串，方便从低位开始相加
	num1 = reverseString(num1)
	num2 = reverseString(num2)

	// 用于存储结果的切片
	result := make([]byte, 0)

	carry := 0 // 进位

	// 遍历 num1 的每一位
	for i := 0; i < len(num1); i++ {
		digit1 := int(num1[i] - '0') // 将字符转换为数字

		// 获取 num2 的对应位的数字，如果越界则设置为 0
		var digit2 int
		if i < len(num2) {
			digit2 = int(num2[i] - '0')
		} else {
			digit2 = 0
		}

		// 相加当前位和进位
		sum := digit1 + digit2 + carry

		// 计算新的进位
		carry = sum / 10

		// 计算当前位的值，并添加到结果中
		result = append(result, byte('0'+sum%10))
	}

	// 处理最高位的进位
	if carry > 0 {
		result = append(result, byte('0'+carry))
	}

	// 反转结果字符串
	resultStr := reverseString(string(result))

	return resultStr
}

func main() {
	// 两个大数相加的例子
	num1 := "123456789012345678901234567890"
	num2 := "987654321098765432109876543212"

	// 调用相加函数
	result := addLargeNumbers(num1, num2)

	// 打印结果
	fmt.Println("相加结果：", result)

	str := "你好" // UTF-8编码的字符串

	// 使用 range 迭代每个 Unicode 字符（rune）
	for _, r := range str {
		fmt.Printf("%c ", r)
	}

}
