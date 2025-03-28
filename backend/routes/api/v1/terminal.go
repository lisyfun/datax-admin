package v1

import (
	"datax-admin/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterTerminalRoutes 注册终端管理相关路由
func RegisterTerminalRoutes(authenticated *gin.RouterGroup) {
	terminalController := controllers.NewTerminalController()

	terminals := authenticated.Group("/terminals")
	{
		terminals.POST("", terminalController.CreateTerminal)
		terminals.PUT("/:id", terminalController.UpdateTerminal)
		terminals.DELETE("/:id", terminalController.DeleteTerminal)
		terminals.GET("", terminalController.GetTerminalList)
		terminals.GET("/:id", terminalController.GetTerminalByID)
		terminals.PUT("/:id/status", terminalController.UpdateTerminalStatus)
		terminals.POST("/:id/connect", terminalController.ConnectTerminal)
		terminals.POST("/:id/disconnect", terminalController.DisconnectTerminal)
		terminals.POST("/:id/upload", terminalController.UploadFiles)
	}
}

// RegisterWebSocketRoutes 注册WebSocket相关路由
func RegisterWebSocketRoutes(wsGroup *gin.RouterGroup) {
	terminalController := controllers.NewTerminalController()

	wsGroup.GET("/terminals/:id", terminalController.ConnectTerminal)
}
