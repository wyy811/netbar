// ============================================================
// 包名：middleware
// 功能：跨域中间件
// 路径：backend/internal/middleware/cors.go
// 说明：处理跨域请求，设置响应头
// ============================================================

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 允许的来源（生产环境应配置具体域名）
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Header("Access-Control-Allow-Origin", origin)

		// 2. 允许的方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		// 3. 允许的请求头
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Trace-ID, X-Requested-With")

		// 4. 允许携带凭证（Cookie）
		c.Header("Access-Control-Allow-Credentials", "true")

		// 5. 预检缓存时间（秒）
		c.Header("Access-Control-Max-Age", "86400")

		// 6. 处理 OPTIONS 预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
