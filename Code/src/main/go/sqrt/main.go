package main

import (
	"fmt"
	"math"
)

func sqrtIterationMethod(x float64) float64 {
	// 初始猜测值，可以选择 x/2 或其他值
	z := x / 2.0

	// 迭代次数，可以根据需要调整

	// 使用迭代法进行逼近
	for {
		z = (z + x/z) / 2
		if math.Abs(z*z-x) < 0.0000000001 {
			break
		}
	}

	return z
}

func main() {
	number := 24.0 // 你想计算平方根的数

	// 使用迭代法计算平方根
	squareRoot := sqrtIterationMethod(number)

	// 手动格式化输出到小数点后10位
	fmt.Printf("Square root of %.2f is %.10f\n", number, squareRoot)
}
