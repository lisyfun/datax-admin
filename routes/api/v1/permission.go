package v1

import (
	"datax-admin/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterPermissionRoutes 注册权限管理相关路由
func RegisterPermissionRoutes(authenticated *gin.RouterGroup) {
	permissionController := controllers.NewPermissionController()
	
	authenticated.POST("/permissions", permissionController.CreatePermission)
	authenticated.PUT("/permissions/:id", permissionController.UpdatePermission)
	authenticated.DELETE("/permissions/:id", permissionController.DeletePermission)
	authenticated.GET("/permissions", permissionController.GetPermissionTree)
	authenticated.GET("/user/permissions", permissionController.GetUserPermissions)
}