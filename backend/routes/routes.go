package routes

import (
	"datax-admin/config"
	"datax-admin/controllers"
	"datax-admin/middleware"
	"datax-admin/utils/logger"
	"net/http"

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

	// 创建控制器
	userController := controllers.NewUserController()
	roleController := controllers.NewRoleController()
	permissionController := controllers.NewPermissionController()
	dashboardController := controllers.NewDashboardController()
	menuController := controllers.NewMenuController()
	jobController := controllers.NewJobController()
	terminalController := controllers.NewTerminalController()

	// WebSocket 路由 - 不需要认证
	wsGroup := r.Group(config.GlobalConfig.Server.BasePath + "/ws")
	{
		wsGroup.GET("/terminals/:id", func(c *gin.Context) {
			logger.Info("收到WebSocket请求: %s", c.Request.URL.Path)
			terminalController.ConnectTerminal(c)
		})
	}

	// API 路由组
	api := r.Group(config.GlobalConfig.Server.BasePath + "/api")
	{
		// API v1 路由组
		v1 := api.Group("/v1")
		{
			// 公开路由
			public := v1.Group("")
			{
				public.GET("/ping", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"message": "pong",
					})
				})

				// 用户认证相关
				public.POST("/register", userController.Register)
				public.POST("/login", userController.Login)
			}

			// 需要认证的路由
			authenticated := v1.Group("")
			authenticated.Use(middleware.JWTAuth())
			{
				// 仪表盘相关
				authenticated.GET("/dashboard", dashboardController.GetDashboardData)

				// 用户个人相关
				authenticated.GET("/user/info", userController.GetUserInfo)
				authenticated.PUT("/user/password", userController.UpdatePassword)
				authenticated.PUT("/user/profile", userController.UpdateProfile)
				authenticated.POST("/user/logout", userController.Logout)

				// 用户管理相关
				authenticated.GET("/users", userController.GetUserList)
				authenticated.PUT("/users/:id/status", userController.UpdateUserStatus)
				authenticated.PUT("/users/:id/password/reset", userController.ResetPassword)
				authenticated.DELETE("/users/:id", userController.DeleteUser)

				// 角色管理相关
				authenticated.POST("/roles", roleController.CreateRole)
				authenticated.PUT("/roles/:id", roleController.UpdateRole)
				authenticated.DELETE("/roles/:id", roleController.DeleteRole)
				authenticated.GET("/roles", roleController.GetRoleList)
				authenticated.GET("/roles/:id/permissions", roleController.GetRolePermissions)
				authenticated.PUT("/roles/:id/permissions", roleController.UpdateRolePermissions)

				// 用户角色管理
				authenticated.GET("/users/:id/roles", roleController.GetUserRoles)
				authenticated.PUT("/users/:id/roles", roleController.UpdateUserRoles)

				// 权限管理相关
				authenticated.POST("/permissions", permissionController.CreatePermission)
				authenticated.PUT("/permissions/:id", permissionController.UpdatePermission)
				authenticated.DELETE("/permissions/:id", permissionController.DeletePermission)
				authenticated.GET("/permissions", permissionController.GetPermissionTree)
				authenticated.GET("/user/permissions", permissionController.GetUserPermissions)

				// 菜单管理相关
				authenticated.POST("/menus", menuController.CreateMenu)
				authenticated.PUT("/menus/:id", menuController.UpdateMenu)
				authenticated.DELETE("/menus/:id", menuController.DeleteMenu)
				authenticated.GET("/menus", menuController.GetMenuList)

				// 任务管理相关
				authenticated.POST("/jobs", jobController.CreateJob)
				authenticated.PUT("/jobs/:id", jobController.UpdateJob)
				authenticated.DELETE("/jobs/:id", jobController.DeleteJob)
				authenticated.POST("/jobs/:id/start", jobController.StartJob)
				authenticated.POST("/jobs/:id/stop", jobController.StopJob)
				authenticated.POST("/jobs/:id/execute", jobController.ExecuteJob)
				authenticated.POST("/jobs/execute", jobController.ExecuteJobs)
				authenticated.GET("/jobs", jobController.GetJobList)
				authenticated.GET("/jobs/history", jobController.GetJobHistoryList)
				authenticated.POST("/jobs/history/clean", jobController.CleanJobHistory)

				// 终端管理相关
				terminals := authenticated.Group("/terminals")
				{
					terminals.POST("", terminalController.CreateTerminal)
					terminals.PUT("/:id", terminalController.UpdateTerminal)
					terminals.DELETE("/:id", terminalController.DeleteTerminal)
					terminals.GET("", terminalController.GetTerminalList)
					terminals.GET("/:id", terminalController.GetTerminalByID)
					terminals.PUT("/:id/status", terminalController.UpdateTerminalStatus)
					terminals.POST("/:id/connect", terminalController.ConnectTerminal)
					terminals.POST("/:id/disconnect", terminalController.DisconnectTerminal)
					terminals.POST("/:id/upload", terminalController.UploadFiles)
				}
			}
		}
	}
}
