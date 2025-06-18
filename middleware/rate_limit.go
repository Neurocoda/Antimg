package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 简单的内存速率限制器
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int           // 请求限制数量
	window   time.Duration // 时间窗口
}

// NewRateLimiter 创建新的速率限制器
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	
	// 获取客户端的请求历史
	requests, exists := rl.requests[clientIP]
	if !exists {
		rl.requests[clientIP] = []time.Time{now}
		return true
	}

	// 清理过期的请求记录
	var validRequests []time.Time
	for _, reqTime := range requests {
		if now.Sub(reqTime) < rl.window {
			validRequests = append(validRequests, reqTime)
		}
	}

	// 检查是否超过限制
	if len(validRequests) >= rl.limit {
		rl.requests[clientIP] = validRequests
		return false
	}

	// 添加当前请求
	validRequests = append(validRequests, now)
	rl.requests[clientIP] = validRequests
	return true
}

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)
	
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		if !limiter.Allow(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, please try again later",
				"code":  429,
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}