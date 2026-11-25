package v1

import (
	"datax-admin/controllers"
	"datax-admin/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterTerminalRoutes 注册终端管理相关路由
func RegisterTerminalRoutes(authenticated *gin.RouterGroup) {
	terminalController := controllers.NewTerminalController()

	terminals := authenticated.Group("/terminals")
	{
		terminals.POST("", middleware.RequirePermission("terminal.list.create"), terminalController.CreateTerminal)
		terminals.PUT("/:id", middleware.RequirePermission("terminal.list.update"), terminalController.UpdateTerminal)
		terminals.DELETE("/:id", middleware.RequirePermission("terminal.list.delete"), terminalController.DeleteTerminal)
		terminals.GET("", middleware.RequirePermission("terminal.list.query", "terminal.list"), terminalController.GetTerminalList)
		terminals.GET("/:id", middleware.RequirePermission("terminal.list.query", "terminal.list"), terminalController.GetTerminalByID)
		terminals.PUT("/:id/status", middleware.RequirePermission("terminal.list.update"), terminalController.UpdateTerminalStatus)
		terminals.POST("/:id/connect", middleware.RequirePermission("terminal.list.query", "terminal.list"), middleware.RequirePermission("terminal.list.connect"), terminalController.ConnectTerminal)
		terminals.POST("/:id/disconnect", middleware.RequirePermission("terminal.list.update"), terminalController.DisconnectTerminal)
		terminals.POST("/:id/upload", middleware.RequirePermission("terminal.list.query", "terminal.list"), middleware.RequirePermission("terminal.list.upload"), terminalController.UploadFiles)
		terminals.GET("/:id/download", middleware.RequirePermission("terminal.list.query", "terminal.list"), middleware.RequirePermission("terminal.list.download"), terminalController.DownloadFile)
		terminals.GET("/:id/files", middleware.RequirePermission("terminal.list.query", "terminal.list"), terminalController.GetFileList)
	}
}

// RegisterWebSocketRoutes 注册WebSocket相关路由
func RegisterWebSocketRoutes(wsGroup *gin.RouterGroup) {
	terminalController := controllers.NewTerminalController()

	wsGroup.GET("/terminals/:id", middleware.RequirePermission("terminal.list.query", "terminal.list"), middleware.RequirePermission("terminal.list.connect"), terminalController.ConnectTerminal)
}
