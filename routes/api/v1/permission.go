package v1

import (
	"datax-admin/controllers"
	"datax-admin/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterPermissionRoutes 注册权限管理相关路由
func RegisterPermissionRoutes(authenticated *gin.RouterGroup) {
	permissionController := controllers.NewPermissionController()

	authenticated.POST("/permissions", middleware.RequirePermission("system.permissions.create"), permissionController.CreatePermission)
	authenticated.PUT("/permissions/:id", middleware.RequirePermission("system.permissions.update"), permissionController.UpdatePermission)
	authenticated.DELETE("/permissions/:id", middleware.RequirePermission("system.permissions.delete"), permissionController.DeletePermission)
	authenticated.GET("/permissions", middleware.RequirePermission("system.permissions.query", "system.roles.permission"), permissionController.GetPermissionTree)

	// 当前登录用户自身的菜单/权限，仅需登录
	authenticated.GET("/user/permissions", permissionController.GetUserPermissions)
	authenticated.GET("/user/menus", permissionController.GetUserMenus)
}
