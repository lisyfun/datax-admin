package v1

import (
	"datax-admin/controllers"
	"datax-admin/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户管理相关路由
func RegisterUserRoutes(authenticated *gin.RouterGroup) {
	userController := controllers.NewUserController()

	// 用户个人相关（仅需登录）
	authenticated.GET("/user/info", userController.GetUserInfo)
	authenticated.PUT("/user/password", userController.UpdatePassword)
	authenticated.PUT("/user/profile", userController.UpdateProfile)
	authenticated.POST("/user/logout", userController.Logout)

	// 用户管理相关（需要相应权限）
	authenticated.GET("/users", middleware.RequirePermission("system.users.query"), userController.GetUserList)
	authenticated.POST("/users", middleware.RequirePermission("system.users.create"), userController.CreateUser)
	authenticated.PUT("/users/:id", middleware.RequirePermission("system.users.update"), userController.UpdateUser)
	authenticated.PUT("/users/:id/status", middleware.RequirePermission("system.users.update"), userController.UpdateUserStatus)
	authenticated.PUT("/users/:id/password/reset", middleware.RequirePermission("system.users.update"), userController.ResetPassword)
	authenticated.DELETE("/users/:id", middleware.RequirePermission("system.users.delete"), userController.DeleteUser)
}
