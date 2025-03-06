package v1

import (
	"datax-admin/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterRoleRoutes 注册角色管理相关路由
func RegisterRoleRoutes(authenticated *gin.RouterGroup) {
	roleController := controllers.NewRoleController()
	
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
}