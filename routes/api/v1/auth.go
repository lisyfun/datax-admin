package v1

import (
	"datax-admin/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes 注册认证相关路由
func RegisterAuthRoutes(public *gin.RouterGroup) {
	userController := controllers.NewUserController()
	
	public.POST("/register", userController.Register)
	public.POST("/login", userController.Login)
}