package controllers

import (
	"datax-admin/services"
	"datax-admin/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TerminalController 终端控制器
type TerminalController struct {
	terminalService *services.TerminalService
}

// NewTerminalController 创建终端控制器
func NewTerminalController() *TerminalController {
	return &TerminalController{
		terminalService: &services.TerminalService{},
	}
}

// CreateTerminal 创建终端
func (c *TerminalController) CreateTerminal(ctx *gin.Context) {
	var req types.CreateTerminalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.terminalService.CreateTerminal(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

// UpdateTerminal 更新终端
func (c *TerminalController) UpdateTerminal(ctx *gin.Context) {
	var req types.UpdateTerminalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的终端ID"})
		return
	}

	if err := c.terminalService.UpdateTerminal(uint(id), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteTerminal 删除终端
func (c *TerminalController) DeleteTerminal(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的终端ID"})
		return
	}

	if err := c.terminalService.DeleteTerminal(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetTerminalList 获取终端列表
func (c *TerminalController) GetTerminalList(ctx *gin.Context) {
	var req types.TerminalListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.terminalService.GetTerminalList(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateTerminalStatus 更新终端状态
func (c *TerminalController) UpdateTerminalStatus(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的终端ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=online offline"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.terminalService.UpdateTerminalStatus(uint(id), req.Status); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}

// GetTerminalByID 获取终端详情
func (c *TerminalController) GetTerminalByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的终端ID"})
		return
	}

	terminal, err := c.terminalService.GetTerminalByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, terminal)
}
