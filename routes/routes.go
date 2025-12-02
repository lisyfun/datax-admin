package routes

import (
    "datax-admin/config"
    "datax-admin/middleware"
    v1 "datax-admin/routes/api/v1"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

// SetupRoutes 配置路由
func SetupRoutes(r *gin.Engine) {
    // 安全响应头
    r.Use(middleware.SecurityHeaders())

    // CORS 中间件配置
    allowed := config.GlobalConfig.Server.AllowedOrigins
    corsCfg := cors.Config{
        AllowOrigins:  allowed,
        AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:  []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders: []string{"Content-Length"},
    }
    // 若未配置白名单，则仅放行同源请求（避免 *）
    if len(allowed) == 0 {
        corsCfg.AllowOriginFunc = func(origin string) bool { return false }
        corsCfg.AllowCredentials = false
    } else {
        corsCfg.AllowCredentials = true
    }
    r.Use(cors.New(corsCfg))

	// 添加安全中间件
	// r.Use(middleware.SecurityMiddleware())

	// 初始化Session中间件 (新增)
	r.Use(middleware.InitSession())

	// 添加加密中间件
	r.Use(middleware.CryptoMiddleware())

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

			// 需要认证的路由 - 改为使用Session认证
			authenticated := v1Group.Group("")
			// authenticated.Use(middleware.JWTAuth()) // 注释掉JWT认证
            authenticated.Use(middleware.SessionAuth()) // 使用Session认证
            {
                v1.RegisterDashboardRoutes(authenticated)
                v1.RegisterUserRoutes(authenticated)
                v1.RegisterRoleRoutes(authenticated)
                v1.RegisterPermissionRoutes(authenticated)
                v1.RegisterJobRoutes(authenticated)
                v1.RegisterTerminalRoutes(authenticated)
                v1.RegisterKafkaRoutes(authenticated)
                v1.RegisterOperationLogRoutes(authenticated)
                v1.RegisterRedisRoutes(authenticated)
            }
        }
    }
}
