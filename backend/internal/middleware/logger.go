// ============================================================
// 包名：middleware
// 功能：请求日志中间件
// 路径：backend/internal/middleware/logger.go
// 说明：记录每个请求的详细信息
// ============================================================

package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"netbar-management/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

// responseWriter 包装 gin.ResponseWriter，捕获响应内容
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger 请求日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 记录开始时间
		startTime := time.Now()

		// 2. 获取 TraceID
		traceID, _ := c.Get("trace_id")

		// 3. 读取请求体
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				// 重新设置请求体
				c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}

		// 4. 脱敏敏感信息
		requestBody = sanitizeBody(requestBody)

		// 5. 创建响应捕获器
		blw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		// 6. 记录请求开始
		logger.WithFields(map[string]interface{}{
			"trace_id":     traceID,
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"query":        c.Request.URL.RawQuery,
			"client_ip":    c.ClientIP(),
			"user_agent":   c.Request.UserAgent(),
			"request_body": truncateString(requestBody, 1000),
		}).Info("请求开始")

		// 7. 处理请求
		c.Next()

		// 8. 记录请求结束
		duration := time.Since(startTime)
		statusCode := c.Writer.Status()

		fields := map[string]interface{}{
			"trace_id":      traceID,
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"status_code":   statusCode,
			"duration_ms":   duration.Milliseconds(),
			"client_ip":     c.ClientIP(),
			"response_body": truncateString(blw.body.String(), 1000),
		}

		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}

		// 根据状态码记录不同级别
		if statusCode >= 500 {
			logger.WithFields(fields).Error("请求失败（服务端错误）")
		} else if statusCode >= 400 {
			logger.WithFields(fields).Warn("请求失败（客户端错误）")
		} else {
			logger.WithFields(fields).Info("请求完成")
		}
	}
}

// 辅助函数
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "... (truncated)"
}

func sanitizeBody(body string) string {
	// 简单脱敏：替换敏感字段
	sensitive := []string{"password", "token", "refresh_token", "secret"}
	for _, s := range sensitive {
		// 简单替换，生产环境建议用 JSON 解析
		// 这里只做演示，实际应用中应根据需求调整
		body = strings.ReplaceAll(body, s, "****")
	}
	return body
}
