// ============================================================
// 包名：middleware
// 功能：限流中间件
// 路径：backend/internal/middleware/rate_limit.go
// 说明：基于 IP 或用户 ID 限制请求频率
// 技术栈：golang.org/x/time/rate
// ============================================================

package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter IP 限流器
type IPRateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*rate.Limiter
	rate     rate.Limit // 每秒允许的请求数
	burst    int        // 突发容量
}

// NewIPRateLimiter 创建 IP 限流器
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

// getLimiter 获取或创建限流器
func (i *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(i.rate, i.burst)
		i.limiters[ip] = limiter
	}

	return limiter
}

// 全局限流器（默认：每秒 10 个请求，突发 20 个）
var defaultRateLimiter = NewIPRateLimiter(10, 20)

// RateLimit 限流中间件（基于 IP）
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := defaultRateLimiter.getLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			return
		}

		c.Next()
	}
}

// RateLimitWithConfig 带配置的限流中间件
func RateLimitWithConfig(r rate.Limit, burst int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(r, burst)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.getLimiter(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			return
		}
		c.Next()
	}
}

// RateLimitByUser 基于用户 ID 的限流（需要 Auth 中间件先执行）
func RateLimitByUser(r rate.Limit, burst int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(r, burst)
	return func(c *gin.Context) {
		// 尝试获取用户 ID
		var key string
		if userID, exists := c.Get("user_id"); exists {
			key = "user:" + string(userID.(uint))
		} else {
			key = "ip:" + c.ClientIP()
		}

		// 这里简化处理，实际应用应使用 Redis 实现分布式限流
		if !limiter.getLimiter(key).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			return
		}
		c.Next()
	}
}
