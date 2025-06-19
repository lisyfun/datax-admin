package v1

import (
	"datax-admin/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterOperationLogRoutes 注册操作日志相关路由
func RegisterOperationLogRoutes(authenticated *gin.RouterGroup) {
	logController := controllers.NewOperationLogController()

	authenticated.GET("/operation-logs", logController.GetLogList)
	authenticated.DELETE("/operation-logs/:id", logController.DeleteLog)
	authenticated.POST("/operation-logs/batch-delete", logController.BatchDeleteLogs)
	authenticated.POST("/operation-logs/clear", logController.ClearLogs)
}
