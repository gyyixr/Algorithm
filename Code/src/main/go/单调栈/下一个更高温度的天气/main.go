package main

func dailyTemperatures(temperatures []int) []int {
	// 1. 初始化结果数组，长度与输入温度数组相同，默认值为 0
	res := make([]int, len(temperatures))

	// 2. 初始化一个栈，并将第一个日期的索引 0 放入栈中
	// 这个栈将用来存储那些我们还没有找到比它更暖和的日期的索引
	stack := []int{0}

	// 3. 从第二个日期开始遍历 (索引 i 从 1 开始)
	for i := 1; i < len(temperatures); i++ {
		// 获取栈顶元素的索引
		top := stack[len(stack)-1]

		// 4. 比较当前温度 temperatures[i] 与栈顶索引对应的温度 temperatures[top]

		// 情况一：当前温度 < 栈顶温度
		// 如果当前温度比栈顶元素对应的温度还要低，
		// 说明当前温度不是栈顶元素的“下一个更暖和的温度”，
		// 并且当前温度也需要等待一个更暖和的未来，所以将当前索引 i 入栈。
		// 栈内元素对应的温度仍然保持递减（或非严格递增）。
		if temperatures[i] < temperatures[top] {
			stack = append(stack, i)
		} else if temperatures[i] == temperatures[top] {
			// 情况二：当前温度 == 栈顶温度
			// 如果当前温度与栈顶元素对应的温度相同，
			// 同样地，它不是栈顶元素严格意义上“更暖和”的温度。
			// 对于题目要求，我们需要的是“更暖和”，即严格大于。
			// 所以也将当前索引 i 入栈，继续等待一个真正更暖和的。
			stack = append(stack, i)
		} else {
			// 情况三：当前温度 > 栈顶温度
			// 这时，当前温度 temperatures[i] 比栈顶索引 top 对应的温度 temperatures[top] 要高。
			// 这意味着对于栈顶的 top 这一天，我们找到了它之后第一个更暖和的日期 i。

			// 循环处理：只要栈不为空，并且当前温度持续大于新的栈顶温度
			for len(stack) != 0 && temperatures[i] > temperatures[top] {
				// 计算等待天数：当前日期索引 i - 栈顶日期索引 top
				res[top] = i - top
				// 将已经找到更暖和日期的栈顶元素弹出
				stack = stack[:len(stack)-1]
				// 如果栈还不为空，更新 top 为新的栈顶元素，继续比较
				if len(stack) != 0 {
					top = stack[len(stack)-1]
				}
			}
			// 当循环结束时，要么栈为空，要么当前温度不再大于新的栈顶温度。
			// 此时，将当前日期索引 i 压入栈中，因为它也需要等待未来一个更暖和的温度，
			// 或者它将成为新的“标杆”（因为它比之前的栈内元素都暖和，或者栈已经空了）。
			stack = append(stack, i)
		}
	}
	// 5. 遍历结束后，res 数组中就包含了每一天需要等待的天数。
	// 对于那些在栈中仍然存在的索引，意味着它们之后没有更暖和的日期，
	// 它们在 res 中对应的结果将保持初始值 0，这符合题目要求。
	return res
}

func main() {
	// 示例用法
	temperatures := []int{73, 74, 75, 71, 69, 72, 76, 73}
	result := dailyTemperatures(temperatures)
	for i, days := range result {
		println("Day", i+1, "waits", days, "days for a warmer temperature.")
	}
}
