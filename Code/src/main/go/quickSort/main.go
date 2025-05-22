package main

import (
	"fmt"
	"math/rand"
	"time"
)

// quickSort 函数实现了快速排序算法
func quickSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}

	left, right := 0, len(arr)-1

	// 随机选择基准元 (pivot) 以提高平均性能
	pivotIndex := rand.Intn(len(arr))
	arr[pivotIndex], arr[right] = arr[right], arr[pivotIndex] // 将基准元放到末尾

	// 分区操作：将小于基准元的元素放到左边，大于基准元的元素放到右边
	for i := range arr {
		if arr[i] < arr[right] {
			arr[left], arr[i] = arr[i], arr[left]
			left++
		}
	}

	// 将基准元放到正确的位置
	arr[left], arr[right] = arr[right], arr[left]

	// 递归地对左右子数组进行排序
	quickSort(arr[:left])   // 排序基准元左边的子数组
	quickSort(arr[left+1:]) // 排序基准元右边的子数组

	return arr
}

func main() {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	s := []int{3, 6, 8, 10, 1, 2, 1, 7, 5, 4, 9}
	fmt.Println("排序前:", s)
	sortedArr := quickSort(s)
	fmt.Println("排序后:", sortedArr)

	s2 := []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
	fmt.Println("排序前:", s2)
	sortedArr2 := quickSort(s2)
	fmt.Println("排序后:", sortedArr2)

	s3 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println("排序前:", s3)
	sortedArr3 := quickSort(s3)
	fmt.Println("排序后:", sortedArr3)
}
