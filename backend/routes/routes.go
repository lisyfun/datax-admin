package routes

import (
	v1 "datax-admin/routes/api/v1"
	"datax-admin/config"
	"datax-admin/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置路由
func SetupRoutes(r *gin.Engine) {
	// CORS 中间件配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Connection", "Upgrade", "Sec-WebSocket-Key", "Sec-WebSocket-Version", "Sec-WebSocket-Extensions", "Sec-WebSocket-Protocol"},
		ExposeHeaders:    []string{"Content-Length", "Upgrade", "Connection"},
		AllowCredentials: true,
	}))

	// WebSocket 路由 - 不需要认证
	wsGroup := r.Group(config.GlobalConfig.Server.BasePath + "/ws")
	v1.RegisterWebSocketRoutes(wsGroup)

	// API 路由组
	api := r.Group(config.GlobalConfig.Server.BasePath + "/api")
	{
		// API v1 路由组
		v1Group := api.Group("/v1")
		{
			// 公开路由
			public := v1Group.Group("")
			{
				v1.RegisterBaseRoutes(public)
				v1.RegisterAuthRoutes(public)
			}

			// 需要认证的路由
			authenticated := v1Group.Group("")
			authenticated.Use(middleware.JWTAuth())
			{
				v1.RegisterDashboardRoutes(authenticated)
				v1.RegisterUserRoutes(authenticated)
				v1.RegisterRoleRoutes(authenticated)
				v1.RegisterPermissionRoutes(authenticated)
				v1.RegisterMenuRoutes(authenticated)
				v1.RegisterJobRoutes(authenticated)
				v1.RegisterTerminalRoutes(authenticated)
			}
		}
	}
}
