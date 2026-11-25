package v1

import (
	"datax-admin/controllers"
	"datax-admin/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoleRoutes 注册角色管理相关路由
func RegisterRoleRoutes(authenticated *gin.RouterGroup) {
	roleController := controllers.NewRoleController()

	// 角色管理相关
	authenticated.POST("/roles", middleware.RequirePermission("system.roles.create"), roleController.CreateRole)
	authenticated.PUT("/roles/:id", middleware.RequirePermission("system.roles.update"), roleController.UpdateRole)
	authenticated.DELETE("/roles/:id", middleware.RequirePermission("system.roles.delete"), roleController.DeleteRole)
	authenticated.GET("/roles", middleware.RequirePermission("system.roles.query", "system.users.query", "system.roles.permission"), roleController.GetRoleList)
	authenticated.GET("/roles/:id/permissions", middleware.RequirePermission("system.roles.query", "system.roles.permission"), roleController.GetRolePermissions)
	authenticated.PUT("/roles/:id/permissions", middleware.RequirePermission("system.roles.permission"), roleController.UpdateRolePermissions)

	// 用户角色管理
	authenticated.GET("/users/:id/roles", middleware.RequirePermission("system.roles.query", "system.users.query"), roleController.GetUserRoles)
	authenticated.PUT("/users/:id/roles", middleware.RequirePermission("system.roles.permission"), roleController.UpdateUserRoles)
}
