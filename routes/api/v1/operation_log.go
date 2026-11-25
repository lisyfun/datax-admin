package v1

import (
	"datax-admin/controllers"
	"datax-admin/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterOperationLogRoutes 注册操作日志相关路由
func RegisterOperationLogRoutes(authenticated *gin.RouterGroup) {
	logController := controllers.NewOperationLogController()

	logPerm := middleware.RequirePermission("system.logs")
	authenticated.GET("/operation-logs", logPerm, logController.GetLogList)
	authenticated.DELETE("/operation-logs/:id", logPerm, logController.DeleteLog)
	authenticated.POST("/operation-logs/batch-delete", logPerm, logController.BatchDeleteLogs)
	authenticated.POST("/operation-logs/clear", logPerm, logController.ClearLogs)
}
