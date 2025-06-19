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

	// 记录操作日志
	operatorID := ctx.GetUint("userID")
	operatorName := ctx.GetString("username")

	if err := c.permissionService.CreatePermission(&req); err != nil {
		// 记录失败日志
		services.LogOperation(
			operatorID,
			operatorName,
			"permission",
			"create",
			"创建权限失败: "+req.Name,
			ctx.ClientIP(),
			ctx.GetHeader("User-Agent"),
			req,
			0,
			err.Error(),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录成功日志
	services.LogOperation(
		operatorID,
		operatorName,
		"permission",
		"create",
		"创建权限成功: "+req.Name,
		ctx.ClientIP(),
		ctx.GetHeader("User-Agent"),
		req,
		1,
		"",
	)

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

	// 记录操作日志
	operatorID := ctx.GetUint("userID")
	operatorName := ctx.GetString("username")

	if err := c.permissionService.UpdatePermission(uint(id), &req); err != nil {
		// 记录失败日志
		services.LogOperation(
			operatorID,
			operatorName,
			"permission",
			"update",
			"更新权限失败: "+req.Name,
			ctx.ClientIP(),
			ctx.GetHeader("User-Agent"),
			req,
			0,
			err.Error(),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录成功日志
	services.LogOperation(
		operatorID,
		operatorName,
		"permission",
		"update",
		"更新权限成功: "+req.Name,
		ctx.ClientIP(),
		ctx.GetHeader("User-Agent"),
		req,
		1,
		"",
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeletePermission 删除权限
func (c *PermissionController) DeletePermission(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的权限ID"})
		return
	}

	// 记录操作日志
	operatorID := ctx.GetUint("userID")
	operatorName := ctx.GetString("username")

	if err := c.permissionService.DeletePermission(uint(id)); err != nil {
		// 记录失败日志
		services.LogOperation(
			operatorID,
			operatorName,
			"permission",
			"delete",
			"删除权限失败",
			ctx.ClientIP(),
			ctx.GetHeader("User-Agent"),
			map[string]interface{}{"permission_id": id},
			0,
			err.Error(),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录成功日志
	services.LogOperation(
		operatorID,
		operatorName,
		"permission",
		"delete",
		"删除权限成功",
		ctx.ClientIP(),
		ctx.GetHeader("User-Agent"),
		map[string]interface{}{"permission_id": id},
		1,
		"",
	)

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

// GetUserMenus 获取用户菜单
func (c *PermissionController) GetUserMenus(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户ID无效"})
		return
	}

	menus, err := c.permissionService.GetUserMenus(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, menus)
}
