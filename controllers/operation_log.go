package controllers

import (
	"datax-admin/services"
	"datax-admin/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OperationLogController 操作日志控制器
type OperationLogController struct {
	logService *services.OperationLogService
}

// NewOperationLogController 创建操作日志控制器
func NewOperationLogController() *OperationLogController {
	return &OperationLogController{
		logService: services.NewOperationLogService(),
	}
}

// GetLogList 获取操作日志列表
func (c *OperationLogController) GetLogList(ctx *gin.Context) {
	var req types.OperationLogListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.logService.GetLogList(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteLog 删除操作日志
func (c *OperationLogController) DeleteLog(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的日志ID"})
		return
	}

	// 记录操作日志
	userID := ctx.GetUint("userID")
	username := ctx.GetString("username")
	services.LogOperation(
		userID,
		username,
		"log",
		"delete",
		"删除操作日志",
		ctx.ClientIP(),
		ctx.GetHeader("User-Agent"),
		map[string]interface{}{"log_id": id},
		1,
		"",
	)

	if err := c.logService.DeleteLog(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// BatchDeleteLogs 批量删除操作日志
func (c *OperationLogController) BatchDeleteLogs(ctx *gin.Context) {
	var req types.BatchDeleteLogsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录操作日志
	userID := ctx.GetUint("userID")
	username := ctx.GetString("username")
	services.LogOperation(
		userID,
		username,
		"log",
		"delete",
		"批量删除操作日志",
		ctx.ClientIP(),
		ctx.GetHeader("User-Agent"),
		req,
		1,
		"",
	)

	if err := c.logService.BatchDeleteLogs(req.IDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "批量删除成功"})
}

// ClearLogs 清空操作日志
func (c *OperationLogController) ClearLogs(ctx *gin.Context) {
	var req types.ClearLogsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录操作日志
	userID := ctx.GetUint("userID")
	username := ctx.GetString("username")
	services.LogOperation(
		userID,
		username,
		"log",
		"clear",
		"清空操作日志",
		ctx.ClientIP(),
		ctx.GetHeader("User-Agent"),
		req,
		1,
		"",
	)

	if err := c.logService.ClearLogs(req.BeforeDays); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "清空成功"})
}
