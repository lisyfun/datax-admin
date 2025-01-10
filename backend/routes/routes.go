package routes

import (
	"datax-admin/controllers"
	"datax-admin/middleware"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置路由
func SetupRoutes(r *gin.Engine) {
	// CORS 中间件配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 创建控制器
	userController := controllers.NewUserController()
	roleController := controllers.NewRoleController()
	permissionController := controllers.NewPermissionController()

	// API v1 路由组
	v1 := r.Group("/api/v1")
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
			// 用户个人相关
			authenticated.GET("/user/info", userController.GetUserInfo)
			authenticated.PUT("/user/password", userController.UpdatePassword)
			authenticated.PUT("/user/profile", userController.UpdateProfile)

			// 用户管理相关
			authenticated.GET("/users", userController.GetUserList)
			authenticated.PUT("/users/:id/status", userController.UpdateUserStatus)
			authenticated.PUT("/users/:id/password/reset", userController.ResetPassword)
			authenticated.DELETE("/users/:id", userController.DeleteUser) // 删除用户

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
			menuController := controllers.NewMenuController()
			authenticated.POST("/menus", menuController.CreateMenu)
			authenticated.PUT("/menus/:id", menuController.UpdateMenu)
			authenticated.DELETE("/menus/:id", menuController.DeleteMenu)
			authenticated.GET("/menus", menuController.GetMenuList)

			// 任务管理相关
			jobController := controllers.NewJobController()
			authenticated.POST("/jobs", jobController.CreateJob)                // 创建任务
			authenticated.PUT("/jobs/:id", jobController.UpdateJob)             // 更新任务
			authenticated.DELETE("/jobs/:id", jobController.DeleteJob)          // 删除任务
			authenticated.POST("/jobs/:id/start", jobController.StartJob)       // 启动任务
			authenticated.POST("/jobs/:id/stop", jobController.StopJob)         // 停止任务
			authenticated.POST("/jobs/:id/execute", jobController.ExecuteJob)   // 执行任务
			authenticated.GET("/jobs", jobController.GetJobList)                // 获取任务列表
			authenticated.GET("/jobs/history", jobController.GetJobHistoryList) // 获取任务历史列表
		}
	}
}
