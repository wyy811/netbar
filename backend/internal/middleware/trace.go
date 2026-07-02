// ============================================================
// 包名：middleware
// 功能：请求链路追踪中间件
// 路径：backend/internal/middleware/trace.go
// 说明：生成 TraceID，用于日志关联和问题追踪
// ============================================================

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Trace 请求链路追踪中间件
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 尝试从请求头获取 TraceID（支持链路传递）
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			// 如果没有，生成新的 TraceID
			traceID = uuid.New().String()
		}

		// 2. 生成 SpanID（当前请求的 ID）
		spanID := uuid.New().String()

		// 3. 存入上下文
		c.Set("trace_id", traceID)
		c.Set("span_id", spanID)

		// 4. 设置响应头，方便前端/下游追踪
		c.Header("X-Trace-ID", traceID)

		c.Next()
	}
}

// GetTraceID 从上下文获取 TraceID（供其他中间件使用）
func GetTraceID(c *gin.Context) string {
	if traceID, exists := c.Get("trace_id"); exists {
		return traceID.(string)
	}
	return "unknown"
}

// GetSpanID 从上下文获取 SpanID
func GetSpanID(c *gin.Context) string {
	if spanID, exists := c.Get("span_id"); exists {
		return spanID.(string)
	}
	return "unknown"
}
