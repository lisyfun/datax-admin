package main

import (
	"datax-admin/config"
	"datax-admin/middleware"
	"datax-admin/models"
	"datax-admin/routes"
	"datax-admin/utils/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化日志系统
	if err := logger.Init(); err != nil {
		logger.Fatal("初始化日志系统失败: %v", err)
	}

	// 初始化数据库
	models.InitDB()

	// 设置gin模式
	gin.SetMode(config.GlobalConfig.Server.Mode)

	// 创建路由
	r := gin.New()
	r.Use(middleware.CustomLogger())
	r.Use(middleware.CustomRecovery())

	// 注册路由
	routes.SetupRoutes(r)

	// 启动服务器
	logger.Info("服务启动成功，监听端口%s", config.GlobalConfig.Server.Port)
	if err := r.Run(config.GlobalConfig.Server.Port); err != nil {
		logger.Fatal("服务启动失败: %v", err)
	}
}
