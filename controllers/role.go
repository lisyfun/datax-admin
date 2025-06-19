package controllers

import (
	"datax-admin/services"
	"datax-admin/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RoleController 角色控制器
type RoleController struct {
	roleService *services.RoleService
}

// NewRoleController 创建角色控制器
func NewRoleController() *RoleController {
	return &RoleController{
		roleService: &services.RoleService{},
	}
}

// CreateRole 创建角色
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var req types.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录操作日志
	operatorID := ctx.GetUint("userID")
	operatorName := ctx.GetString("username")

	if err := c.roleService.CreateRole(&req); err != nil {
		// 记录失败日志
		services.LogOperation(
			operatorID,
			operatorName,
			"role",
			"create",
			"创建角色失败: "+req.Name,
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
		"role",
		"create",
		"创建角色成功: "+req.Name,
		ctx.ClientIP(),
		ctx.GetHeader("User-Agent"),
		req,
		1,
		"",
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

// UpdateRole 更新角色
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	var req types.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	// 记录操作日志
	operatorID := ctx.GetUint("userID")
	operatorName := ctx.GetString("username")

	if err := c.roleService.UpdateRole(uint(id), &req); err != nil {
		// 记录失败日志
		services.LogOperation(
			operatorID,
			operatorName,
			"role",
			"update",
			"更新角色失败: "+req.Name,
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
		"role",
		"update",
		"更新角色成功: "+req.Name,
		ctx.ClientIP(),
		ctx.GetHeader("User-Agent"),
		req,
		1,
		"",
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteRole 删除角色
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	// 记录操作日志
	operatorID := ctx.GetUint("userID")
	operatorName := ctx.GetString("username")

	if err := c.roleService.DeleteRole(uint(id)); err != nil {
		// 记录失败日志
		services.LogOperation(
			operatorID,
			operatorName,
			"role",
			"delete",
			"删除角色失败",
			ctx.ClientIP(),
			ctx.GetHeader("User-Agent"),
			map[string]interface{}{"role_id": id},
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
		"role",
		"delete",
		"删除角色成功",
		ctx.ClientIP(),
		ctx.GetHeader("User-Agent"),
		map[string]interface{}{"role_id": id},
		1,
		"",
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetRoleList 获取角色列表
func (c *RoleController) GetRoleList(ctx *gin.Context) {
	var req types.RoleListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.roleService.GetRoleList(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateRolePermissions 更新角色权限
func (c *RoleController) UpdateRolePermissions(ctx *gin.Context) {
	var req types.UpdateRolePermissionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	// 记录操作日志
	operatorID := ctx.GetUint("userID")
	operatorName := ctx.GetString("username")

	if err := c.roleService.UpdateRolePermissions(uint(id), &req); err != nil {
		// 记录失败日志
		services.LogOperation(
			operatorID,
			operatorName,
			"role",
			"update",
			"更新角色权限失败",
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
		"role",
		"update",
		"更新角色权限成功",
		ctx.ClientIP(),
		ctx.GetHeader("User-Agent"),
		req,
		1,
		"",
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// GetRolePermissions 获取角色权限
func (c *RoleController) GetRolePermissions(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	permissions, err := c.roleService.GetRolePermissions(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

// UpdateUserRoles 更新用户角色
func (c *RoleController) UpdateUserRoles(ctx *gin.Context) {
	var req types.UpdateUserRolesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 记录操作日志
	operatorID := ctx.GetUint("userID")
	operatorName := ctx.GetString("username")

	if err := c.roleService.UpdateUserRoles(uint(id), &req); err != nil {
		// 记录失败日志
		services.LogOperation(
			operatorID,
			operatorName,
			"user",
			"update",
			"更新用户角色失败",
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
		"user",
		"update",
		"更新用户角色成功",
		ctx.ClientIP(),
		ctx.GetHeader("User-Agent"),
		req,
		1,
		"",
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// GetUserRoles 获取用户角色
func (c *RoleController) GetUserRoles(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	roles, err := c.roleService.GetUserRoles(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"roles": roles})
}
