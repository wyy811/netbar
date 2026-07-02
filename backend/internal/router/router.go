// ============================================================
// 包名：router
// 功能：路由注册
// 路径：backend/internal/router/router.go
// ============================================================

package router

import (
	"github.com/gin-gonic/gin"

	"netbar-management/internal/middleware"
	"netbar-management/internal/pkg/jwt"
)

// SetupRouter 设置路由
func SetupRouter(
	jwtManager *jwt.Manager,
	authHandler *v1.AuthHandler,
	computerHandler *v1.ComputerHandler,
	sessionHandler *v1.SessionHandler,
	memberHandler *v1.MemberHandler,
	auditHandler *v1.AuditHandler,
	statisticsHandler *v1.StatisticsHandler,
	rateHandler *v1.RateHandler,
	userHandler *v1.UserHandler,
) *gin.Engine {

	r := gin.New()

	// ============================================================
	// 1. 全局中间件（按顺序注册）
	// ============================================================

	r.Use(middleware.Recovery())  // 第1个：捕获 panic
	r.Use(middleware.Trace())     // 第2个：链路追踪
	r.Use(middleware.Logger())    // 第3个：请求日志
	r.Use(middleware.CORS())      // 第4个：跨域处理
	r.Use(middleware.RateLimit()) // 第5个：限流

	// ============================================================
	// 2. 健康检查
	// ============================================================

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// ============================================================
	// 3. API v1
	// ============================================================

	v1 := r.Group("/api/v1")
	{
		// 3.1 认证组（不需要认证）
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// 3.2 需要认证的接口
		authorized := v1.Group("/")
		authorized.Use(middleware.Auth(jwtManager)) // 第6个：JWT 验证
		{
			authorized.GET("/profile", authHandler.Profile)
			authorized.POST("/logout", authHandler.Logout)
			authorized.PUT("/profile/password", authHandler.ChangePassword)

			// 机位管理
			computers := authorized.Group("/computers")
			{
				computers.GET("", computerHandler.List)
				computers.GET("/available", computerHandler.GetAvailable)
				computers.GET("/:id", computerHandler.GetByID)
				computers.POST("", computerHandler.Create)
				computers.PUT("/:id", computerHandler.Update)
				computers.DELETE("/:id", computerHandler.Delete)
				computers.PATCH("/:id/status", computerHandler.ChangeStatus)
			}

			// 上机记录
			sessions := authorized.Group("/sessions")
			{
				sessions.POST("/start", sessionHandler.Start)
				sessions.POST("/end", sessionHandler.End)
				sessions.GET("/active", sessionHandler.GetActive)
				sessions.GET("/history", sessionHandler.GetHistory)
				sessions.GET("/revenue", sessionHandler.GetRevenue)
				sessions.POST("/:id/force-end", sessionHandler.ForceEnd)
			}

			// 会员管理
			members := authorized.Group("/members")
			{
				members.GET("", memberHandler.List)
				members.GET("/:id", memberHandler.GetByID)
				members.POST("", memberHandler.Create)
				members.PUT("/:id", memberHandler.Update)
				members.DELETE("/:id", memberHandler.Delete)
				members.POST("/:id/recharge", memberHandler.Recharge)
				members.GET("/:id/balance", memberHandler.GetBalance)
			}

			// 计费规则（需要管理员权限）
			rates := authorized.Group("/rates")
			rates.Use(middleware.RoleCheck(2)) // 第7个：权限校验
			{
				rates.GET("", rateHandler.List)
				rates.POST("", rateHandler.Create)
				rates.PUT("/:id", rateHandler.Update)
				rates.DELETE("/:id", rateHandler.Delete)
			}

			// 管理员功能
			admin := authorized.Group("/admin")
			admin.Use(middleware.RoleCheck(2))
			{
				admin.GET("/audit-logs", auditHandler.List)
				admin.GET("/users", userHandler.List)
				admin.PUT("/users/:id/status", userHandler.UpdateStatus)
			}

			// 超级管理员功能
			superAdmin := authorized.Group("/system")
			superAdmin.Use(middleware.RoleCheck(3))
			{
				superAdmin.GET("/configs", configHandler.List)
				superAdmin.PUT("/configs/:key", configHandler.Update)
			}

			// 统计数据
			statistics := authorized.Group("/statistics")
			{
				statistics.GET("/dashboard", statisticsHandler.GetDashboard)
				statistics.GET("/revenue-trend", statisticsHandler.GetRevenueTrend)
				statistics.GET("/peak-hours", statisticsHandler.GetPeakHours)
			}
		}
	}

	return r
}
