package v1

import (
	"datax-admin/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterMenuRoutes 注册菜单管理相关路由
func RegisterMenuRoutes(authenticated *gin.RouterGroup) {
	menuController := controllers.NewMenuController()

	authenticated.POST("/menus", menuController.CreateMenu)
	authenticated.PUT("/menus/:id", menuController.UpdateMenu)
	authenticated.DELETE("/menus/:id", menuController.DeleteMenu)
	authenticated.GET("/menus", menuController.GetMenuList)
	// 注意：/user/menus 路由已移至权限路由中，由权限管理统一处理菜单功能
}
