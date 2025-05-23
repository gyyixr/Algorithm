package main

import (
	"fmt"
	"sync"
	"time"
)

// TokenBucket 定义了令牌桶的结构
type TokenBucket struct {
	capacity      int64      // 桶的容量
	rate          int64      // 令牌放入速率 (每秒放入多少个)
	tokens        int64      // 当前桶内令牌数量
	lastTokenTime time.Time  // 上一次放令牌的时间
	mu            sync.Mutex // 用于并发安全的互斥锁
}

// NewTokenBucket 创建一个新的令牌桶
func NewTokenBucket(capacity, rate int64) *TokenBucket {
	return &TokenBucket{
		capacity:      capacity,
		rate:          rate,
		tokens:        capacity, // 初始时桶是满的
		lastTokenTime: time.Now(),
	}
}

// Allow 尝试获取一个令牌，如果成功返回 true，否则返回 false
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	// 计算从上次放令牌到现在，应该新增多少令牌
	elapsed := now.Sub(tb.lastTokenTime).Seconds()
	newTokens := int64(elapsed * float64(tb.rate))

	if newTokens > 0 {
		tb.tokens += newTokens
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity // 令牌数量不能超过桶的容量
		}
		tb.lastTokenTime = now
	}

	// 检查是否有足够的令牌
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

func main() {
	// 创建一个容量为 5，每秒生成 2 个令牌的令牌桶
	limiter := NewTokenBucket(5, 2)

	// 模拟请求
	for i := 1; i <= 10; i++ {
		if limiter.Allow() {
			fmt.Printf("请求 %d: 允许通过\n", i)
		} else {
			fmt.Printf("请求 %d: 被限制\n", i)
		}
		time.Sleep(200 * time.Millisecond) // 模拟请求间隔
	}

	fmt.Println("\n等待一段时间让令牌桶填充...")
	time.Sleep(3 * time.Second)

	// 再次模拟请求
	for i := 11; i <= 15; i++ {
		if limiter.Allow() {
			fmt.Printf("请求 %d: 允许通过\n", i)
		} else {
			fmt.Printf("请求 %d: 被限制\n", i)
		}
		time.Sleep(300 * time.Millisecond)
	}
}
