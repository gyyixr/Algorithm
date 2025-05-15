package main

import "fmt"

func ceilDiv(m, n int) int {
	return (m + n - 1) / n
}

func sgemmNaiveCPU(A, B, C []float32, M, N, K int) {
	for x := 0; x < M; x++ {
		for y := 0; y < N; y++ {
			sum := float32(0.0)
			for i := 0; i < K; i++ {
				sum += A[x*K+i] * B[i*N+y]
			}
			C[x*N+y] = sum
		}
	}
}

func main() {
	M, N, K := 2, 3, 4
	A := make([]float32, M*K)
	B := make([]float32, K*N)
	C := make([]float32, M*N)

	// 初始化 A 和 B 的元素
	for i, j := range A {
		A[i] = float32(i + 1)
		println(j)
	}
	for i := range B {
		B[i] = float32(i + 1)
	}

	sgemmNaiveCPU(A, B, C, M, N, K)

	// 输出结果矩阵 C
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			fmt.Printf("%.2f ", C[i*N+j])
		}
		fmt.Println()
	}
}
