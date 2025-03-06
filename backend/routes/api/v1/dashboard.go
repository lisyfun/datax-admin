package v1

import (
	"datax-admin/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterDashboardRoutes 注册仪表盘相关路由
func RegisterDashboardRoutes(authenticated *gin.RouterGroup) {
	dashboardController := controllers.NewDashboardController()
	
	authenticated.GET("/dashboard", dashboardController.GetDashboardData)
}