// ============================================================
// 包名：middleware
// 功能：RBAC 权限校验中间件
// 路径：backend/internal/middleware/role_check.go
// 说明：检查当前用户是否有权限访问该路由
// ============================================================

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleCheck 权限校验中间件
// minRole: 所需的最低角色（1=员工, 2=管理员, 3=超级管理员）
func RoleCheck(minRole int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从上下文获取用户角色（由 Auth 中间件设置）
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未认证",
			})
			return
		}

		// 2. 检查权限
		userRole := role.(int)
		if userRole < minRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足，需要更高级别权限",
			})
			return
		}

		c.Next()
	}
}

// IsAdmin 检查是否为管理员（角色 >= 2）
func IsAdmin() gin.HandlerFunc {
	return RoleCheck(2)
}

// IsSuperAdmin 检查是否为超级管理员（角色 = 3）
func IsSuperAdmin() gin.HandlerFunc {
	return RoleCheck(3)
}

// CheckPermission 自定义权限检查（支持多角色）
func CheckPermission(allowedRoles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未认证",
			})
			return
		}

		userRole := role.(int)
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "权限不足",
		})
	}
}
