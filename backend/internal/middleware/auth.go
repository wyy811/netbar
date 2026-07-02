// ============================================================
// 包名：middleware
// 功能：JWT 验证中间件
// 路径：backend/internal/middleware/auth.go
// 说明：验证 JWT，提取用户信息存入上下文
// ============================================================

package middleware

import (
	"net/http"
	"strings"

	"netbar-management/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// Auth JWT 验证中间件
func Auth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 Authorization Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证令牌",
			})
			return
		}

		// 2. 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证格式错误，请使用 Bearer Token",
			})
			return
		}

		tokenString := parts[1]

		// 3. 解析验证 Token
		claims, err := jwtManager.ParseToken(tokenString)
		if err != nil {
			// 判断错误类型
			errMsg := "无效令牌"
			if strings.Contains(err.Error(), "expired") {
				errMsg = "令牌已过期，请刷新"
			} else if strings.Contains(err.Error(), "blacklisted") {
				errMsg = "令牌已被撤销，请重新登录"
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": errMsg,
			})
			return
		}

		// 4. 验证 Token 类型（必须是 Access Token）
		if claims.TokenType != jwt.TokenTypeAccess {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的令牌类型",
			})
			return
		}

		// 5. 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// GetCurrentUserID 从上下文获取当前用户 ID（供 Handler 使用）
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetCurrentRole 从上下文获取当前用户角色
func GetCurrentRole(c *gin.Context) (int, bool) {
	role, exists := c.Get("role")
	if !exists {
		return 0, false
	}
	return role.(int), true
}
