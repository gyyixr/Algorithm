//给定1千万个随机数，范围在1到1亿之间，要求找出所有在这个范围内但没有出现在随机数中的数字。
package main

import (
	"fmt"
)

/**
# 深入理解位图中的 `index` 和 `bit` 的作用

让我用最直观的方式来解释这两个关键参数，就像教一个完全不懂的人理解这个概念。

## 1. 把位图想象成一个巨大的开关面板

想象你面前有一面墙，上面有1亿个小开关（每个开关代表一个数字），这些开关排列成许多行：

- **每一行**有64个开关（如果我们使用uint64）
- **总行数** = 总开关数/每行开关数 = 1亿/64 ≈ 1,562,500行

## 2. `index` 的作用 - 找对行

`index` 就是告诉你：
- **你要操作的开关在哪一行**
- 计算方式：`index = (数字-1) / 每行开关数`

### 实际例子：
假设我们要处理数字65：
- (65-1)/64 = 64/64 = 1 → 第1行（从0开始计数）
- 因为：
  - 第0行：1-64
  - 第1行：65-128
  - 以此类推

## 3. `bit` 的作用 - 找对列

`bit` 就是告诉你：
- **你要操作的开关在这一行的哪个位置**
- 计算方式：`bit = (数字-1) % 每行开关数`

### 继续数字65的例子：
- (65-1)%64 = 64%64 = 0 → 这一行的第0个位置（最左边）
- 所以数字65对应：
  - 第1行
  - 该行的第0个开关

## 4. 为什么从`数字-1`开始计算？

因为：
- 数字从**1**开始计数
- 但计算机从**0**开始计数
- 所以需要-1来对齐

## 5. 实际内存中的表示

每一行（每个uint64）在内存中是这样存储的：

```
bit位置： 63 62 61 ... 7 6 5 4 3 2 1 0
值：      0  0  0  ... 0 0 0 0 0 0 0 0
```

当我们说"第0个位置"时，指的是最右边的bit（二进制的最低位）

## 6. 完整操作示例：标记数字65

1. 计算位置：
   - index = (65-1)/64 = 1
   - bit = (65-1)%64 = 0

2. 生成掩码：
   - 1 << 0 = 1 (二进制000...0001)

3. 设置位：
   - bitmap[1] = bitmap[1] | 000...0001
   - 这样就把第1行的第0位设为了1

## 7. 为什么这样设计？

这样设计的好处：
- **极省空间**：1亿个数字只需要 ≈12MB 内存
- **超快查询**：计算简单，只需一次除法和取模
- **CPU友好**：位操作是CPU最基本的操作之一

## 8. 常见误区澄清

误区1："bit是二进制的位数"
- 不是！bit是**在该行中的位置索引**，从0开始

误区2："index是直接的内存地址"
- 不是！index是数组的索引，每个元素是一个uint64

误区3："数字越大bit值越大"
- 不一定！bit值会在0-63之间循环，因为每行只有64个位置

## 9. 可视化理解

让我们用更小的例子（假设每行8个数字）来看：

```
行号(index) | 包含的数字
-----------|-----------
0          | 1-8
1          | 9-16
2          | 17-24
...        | ...
```

数字10的位置：
- index = (10-1)/8 = 1（第1行）
- bit = (10-1)%8 = 1（该行的第1个位置，从0开始）

## 10. 总结

- `index` = 找对行（哪个uint64）
- `bit` = 找对列（该uint64的哪一位）
- 组合起来就能精确定位到特定位

这种设计就像：
- 知道书在哪一个书架上（index）
- 知道书在该书架的哪一层（bit）
- 组合起来就能快速找到任何一本书
*/

//bitmap法
func findMissingNumbersBitmap(randomNumbers []int, maxNum int) []int {
	// 计算需要的字节数 (每个字节8位)
	bitmapSize := (maxNum + 7) / 8
	bitmap := make([]byte, bitmapSize)

	// 标记存在的数字
	for _, num := range randomNumbers {
		if num < 1 || num > maxNum {
			continue
		}
		index := (num - 1) / 8
		bit := (num - 1) % 8
		//一步操作 bitmap[index] |= 1 << bit
		// 拆解后的三步操作
		bitMask := byte(1 << bit)         // 第一步：创建位掩码
		current := bitmap[index]          // 第二步：获取当前值
		bitmap[index] = current | bitMask // 第三步：设置新值

	}

	// 收集缺失的数字
	missing := make([]int, 0)
	for num := 1; num <= maxNum; num++ {
		index := (num - 1) / 8
		bit := (num - 1) % 8
		if (bitmap[index] & (1 << bit)) == 0 {
			missing = append(missing, num)
		}
	}

	return missing
}

func main() {
	// 示例使用
	randomNumbers := []int{2, 3, 5, 7, 11} // 假设这是你的1千万个随机数
	maxNum := 20                           // 1亿
	missing := findMissingNumbersBitmap(randomNumbers, maxNum)
	fmt.Println("Missing numbers:", missing)
}

//这个优化版本使用了uint64而不是byte作为位图的基本单位，减少了内存访问次数，提高了性能。对于非常大的数据集，还可以考虑使用goroutine并行处理不同范围的数字检查。
func findMissingNumbersOptimized(randomNumbers []int, maxNum int) []int {
	// 使用更大的单位(如uint64)来减少内存访问次数
	bitmapSize := (maxNum + 63) / 64
	bitmap := make([]uint64, bitmapSize)

	// 并行标记存在的数字(如果数据量很大)
	for _, num := range randomNumbers {
		if num < 1 || num > maxNum {
			continue
		}
		index := (num - 1) / 64
		bit := (num - 1) % 64
		bitmap[index] |= 1 << bit
	}

	// 预分配缺失数组空间(假设缺失数量约为maxNum-len(randomNumbers))
	missing := make([]int, 0, maxNum-len(randomNumbers))

	// 并行检查缺失数字(可以分多个goroutine处理不同范围)
	for num := 1; num <= maxNum; num++ {
		index := (num - 1) / 64
		bit := (num - 1) % 64
		if (bitmap[index] & (1 << bit)) == 0 {
			missing = append(missing, num)
		}
	}

	return missing
}
