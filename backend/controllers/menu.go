package controllers

import (
	"datax-admin/services"
	"datax-admin/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MenuController 菜单控制器
type MenuController struct {
	menuService *services.MenuService
}

// NewMenuController 创建菜单控制器
func NewMenuController() *MenuController {
	return &MenuController{
		menuService: &services.MenuService{},
	}
}

// CreateMenu 创建菜单
func (c *MenuController) CreateMenu(ctx *gin.Context) {
	var req types.CreateMenuRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.menuService.CreateMenu(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

// UpdateMenu 更新菜单
func (c *MenuController) UpdateMenu(ctx *gin.Context) {
	var req types.UpdateMenuRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的菜单ID"})
		return
	}

	if err := c.menuService.UpdateMenu(uint(id), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteMenu 删除菜单
func (c *MenuController) DeleteMenu(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的菜单ID"})
		return
	}

	if err := c.menuService.DeleteMenu(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetMenuList 获取菜单列表
func (c *MenuController) GetMenuList(ctx *gin.Context) {
	var req types.MenuListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.menuService.GetMenuList(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.JSON(http.StatusOK, resp)
}
