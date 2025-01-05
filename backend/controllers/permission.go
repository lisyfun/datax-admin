package controllers

import (
	"datax-admin/services"
	"datax-admin/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PermissionController 权限控制器
type PermissionController struct {
	permissionService *services.PermissionService
}

// NewPermissionController 创建权限控制器
func NewPermissionController() *PermissionController {
	return &PermissionController{
		permissionService: &services.PermissionService{},
	}
}

// CreatePermission 创建权限
func (c *PermissionController) CreatePermission(ctx *gin.Context) {
	var req types.CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.permissionService.CreatePermission(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

// UpdatePermission 更新权限
func (c *PermissionController) UpdatePermission(ctx *gin.Context) {
	var req types.UpdatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的权限ID"})
		return
	}

	if err := c.permissionService.UpdatePermission(uint(id), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeletePermission 删除权限
func (c *PermissionController) DeletePermission(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的权限ID"})
		return
	}

	if err := c.permissionService.DeletePermission(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetPermissionTree 获取权限树
func (c *PermissionController) GetPermissionTree(ctx *gin.Context) {
	var req types.PermissionListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.permissionService.GetPermissionTree(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetUserPermissions 获取用户权限
func (c *PermissionController) GetUserPermissions(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	permissions, err := c.permissionService.GetUserPermissions(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"permissions": permissions})
}
