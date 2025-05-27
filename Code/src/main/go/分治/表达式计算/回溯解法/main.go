package main

import (
	"fmt"
	"strconv"
)

// Token 可以是数字或运算符
type Token struct {
	IsOperator bool
	Value      int    // 如果是数字
	Operator   string // 如果是运算符
}

// parseExpression 将字符串表达式转换为 Token 列表
func parseExpression(expression string) []Token {
	tokens := []Token{}
	numStr := ""
	for _, char := range expression {
		if char >= '0' && char <= '9' {
			numStr += string(char)
		} else {
			if numStr != "" {
				val, _ := strconv.Atoi(numStr)
				tokens = append(tokens, Token{IsOperator: false, Value: val})
				numStr = ""
			}
			tokens = append(tokens, Token{IsOperator: true, Operator: string(char)})
		}
	}
	if numStr != "" {
		val, _ := strconv.Atoi(numStr)
		tokens = append(tokens, Token{IsOperator: false, Value: val})
	}
	return tokens
}

func calculate(a, b int, op string) int {
	switch op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	}
	return 0 // Should not happen for valid ops
}

//func backtrack(currentTokens []Token, results map[int]bool) {
//	// 基本情况：如果只剩下一个数字 Token，则将其添加到结果中
//	if len(currentTokens) == 1 && !currentTokens[0].IsOperator {
//		results[currentTokens[0].Value] = true
//		return
//	}
//
//	// 遍历所有可能的运算符进行计算
//	// 我们需要找到一个 数字-运算符-数字 的模式
//	for i := 0; i < len(currentTokens); i++ {
//		// 寻找一个运算符
//		if currentTokens[i].IsOperator {
//			// 确保运算符左边和右边是数字，并且在列表范围内
//			if i > 0 && i < len(currentTokens)-1 &&
//				!currentTokens[i-1].IsOperator && !currentTokens[i+1].IsOperator {
//
//				leftVal := currentTokens[i-1].Value
//				op := currentTokens[i].Operator
//				rightVal := currentTokens[i+1].Value
//
//				// 执行计算
//				resVal := calculate(leftVal, rightVal, op)
//
//				// 构建新的 Token 列表 (状态转换)
//				newTokens := []Token{}
//				newTokens = append(newTokens, currentTokens[:i-1]...)                  // 运算符左边的部分
//				newTokens = append(newTokens, Token{IsOperator: false, Value: resVal}) // 计算结果
//				if i+2 < len(currentTokens) {
//					newTokens = append(newTokens, currentTokens[i+2:]...) // 运算符右边的部分
//				}
//
//				// 递归调用
//				backtrack(newTokens, results)
//				// "撤销选择" 是隐式的，因为 newTokens 是新创建的，
//				// currentTokens 在当前循环的下一次迭代中保持不变。
//			}
//		}
//	}
//}

// computeWithBacktracking 是外部调用的入口函数
func computeWithBacktracking(expression string) []int {
	initialTokens := parseExpression(expression)
	// 使用 map 来自动处理重复结果，符合题目示例的输出（例如 -10, -10）
	// 如果严格要求不重复，map[int]bool 就很好。但题目允许重复，所以直接用 list。
	// 为了方便，我们还是先用 map 去重，然后转回 list。
	// 或者，如果题目允许输出包含重复，那么 backtrack 的第二个参数应该是 *[]int，直接 append。
	// 示例输出 [-34,-14,-10,-10,10] 包含重复，所以我们直接收集。

	var allResults []int
	// 为了让回溯能直接添加结果，我们改用一个辅助函数和闭包，或者传递结果列表的指针。
	// 这里我们定义一个内部回溯函数来直接修改 allResults。

	var actualBacktrack func(tokens []Token)
	actualBacktrack = func(currentTokens []Token) {
		if len(currentTokens) == 1 && !currentTokens[0].IsOperator {
			allResults = append(allResults, currentTokens[0].Value)
			return
		}

		for i := 0; i < len(currentTokens); i++ {
			if currentTokens[i].IsOperator {
				if i > 0 && i < len(currentTokens)-1 &&
					!currentTokens[i-1].IsOperator && !currentTokens[i+1].IsOperator {
					leftVal := currentTokens[i-1].Value
					op := currentTokens[i].Operator
					rightVal := currentTokens[i+1].Value
					resVal := calculate(leftVal, rightVal, op)

					newTokens := make([]Token, 0, len(currentTokens)-2)
					newTokens = append(newTokens, currentTokens[:i-1]...)
					newTokens = append(newTokens, Token{IsOperator: false, Value: resVal})
					if i+2 < len(currentTokens) {
						newTokens = append(newTokens, currentTokens[i+2:]...)
					}
					actualBacktrack(newTokens)
				}
			}
		}
	}

	actualBacktrack(initialTokens)

	return allResults
}

func main() {
	expression1 := "2-1-1"
	fmt.Printf("输入: %s, 回溯法输出: %v\n", expression1, computeWithBacktracking(expression1))
	// 预期: [0, 2] 或 [2, 0] (可能包含重复，取决于去重策略)

	expression2 := "2*3-4*5"
	fmt.Printf("输入: %s, 回溯法输出: %v\n", expression2, computeWithBacktracking(expression2))
	// 预期: [-34, -14, -10, -10, 10] (允许重复)
}
