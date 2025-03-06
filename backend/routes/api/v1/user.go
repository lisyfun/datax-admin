package v1

import (
	"datax-admin/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户管理相关路由
func RegisterUserRoutes(authenticated *gin.RouterGroup) {
	userController := controllers.NewUserController()
	
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
}