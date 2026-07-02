// ============================================================
// 包名：middleware
// 功能：Panic 恢复中间件
// 路径：backend/internal/middleware/recovery.go
// 说明：捕获 panic，防止程序崩溃，返回 500 错误
// ============================================================

package middleware

import (
	"net/http"
	"runtime/debug"

	"netbar-management/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Recovery  panic 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 1. 获取堆栈信息
				stack := debug.Stack()

				// 2. 记录错误日志
				logger.WithFields(map[string]interface{}{
					"error":     err,
					"stack":     string(stack),
					"path":      c.Request.URL.Path,
					"method":    c.Request.Method,
					"client_ip": c.ClientIP(),
				}).Error("服务器内部错误（panic）")

				// 3. 返回 500 错误
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "服务器内部错误，请稍后重试",
				})
			}
		}()
		c.Next()
	}
}
