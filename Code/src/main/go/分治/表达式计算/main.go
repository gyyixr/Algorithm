/**
* 给你一个由数字和运算符组成的字符串 expression ，按不同优先级组合数字和运算符，计算并返回所有可能组合的结果。你可以 按任意顺序 返回答案。
* 示例 1：
* 输入：expression = "2-1-1"
* 输出：[0,2]
* 解释：
* ((2-1)-1) = 0
* (2-(1-1)) = 2
* 示例 2：
* 输入：expression = "2*3-4*5"
* 输出：[-34,-14,-10,-10,10]
* 解释：
* (2*(3-(4*5))) = -34
* ((2*3)-(4*5)) = -14
* ((2*(3-4))*5) = -10
* (2*((3-4)*5)) = -10
* (((2*3)-4)*5) = 10
 */
package main

import (
	"fmt"
)

// memo 用于存储已计算过的子表达式的结果，以避免重复计算。
// 键是子表达式字符串，值是该子表达式所有可能计算结果的切片。
var memo map[string][]int

// diffWaysToCompute 是递归函数，用于计算表达式的所有可能结果。
func diffWaysToCompute(expression string) []int {
	// 检查备忘录中是否已有当前表达式的结果
	if res, ok := memo[expression]; ok {
		return res
	}

	var results []int
	isNumber := true // 标记当前表达式是否只是一个数字

	for i, char := range expression {
		// 如果字符是运算符
		if char == '+' || char == '-' || char == '*' {
			isNumber = false // 发现运算符，则表达式不是纯数字

			// 分割表达式为左右两部分
			leftPart := expression[:i]
			rightPart := expression[i+1:]

			// 递归计算左右两部分的结果
			leftResults := diffWaysToCompute(leftPart)
			rightResults := diffWaysToCompute(rightPart)

			// 合并左右两部分的结果
			for _, l := range leftResults {
				for _, r := range rightResults {
					switch char {
					case '+':
						results = append(results, l+r)
					case '-':
						results = append(results, l-r)
					case '*':
						results = append(results, l*r)
					}
				}
			}
		}
	}

	// 基本情况：如果表达式中没有运算符，它就是一个数字
	if isNumber {
		results = append(results, fromStringToInt(expression)) // 将字符串转换为整数
	}

	// 将当前表达式的结果存入备忘录
	memo[expression] = results
	return results
}

func fromStringToInt(s string) int {
	res := 0
	for _, char := range s {
		res = res*10 + int(char-'0') // 将字符转换为整数
	}
	return res
}

// compute 是外部调用的入口函数，它会初始化备忘录。
func compute(expression string) []int {
	memo = make(map[string][]int) // 为每次顶级调用初始化备忘录
	return diffWaysToCompute(expression)
}

func main() {
	expression1 := "2-1-1"
	// 对结果进行排序是为了方便与预期输出比较，题目要求任意顺序即可
	// sort.Ints(result1)
	fmt.Printf("输入: %s, 输出: %v\n", expression1, compute(expression1))
	// 预期: [0, 2] (顺序可能不同)

	expression2 := "2*3-4*5"
	// sort.Ints(result2)
	fmt.Printf("输入: %s, 输出: %v\n", expression2, compute(expression2))
	// 预期: [-34, -14, -10, -10, 10] (顺序和重复均可能不同)
}
