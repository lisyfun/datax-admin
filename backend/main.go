package main

import (
	"datax-admin/config"
	"datax-admin/models"
	"datax-admin/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化数据库
	models.InitDB()

	// 设置gin模式
	gin.SetMode(config.GlobalConfig.Server.Mode)

	// 创建路由
	r := gin.Default()

	// 注册路由
	routes.SetupRoutes(r)

	// 启动服务器
	log.Printf("Server starting on port %s", config.GlobalConfig.Server.Port)
	if err := r.Run(config.GlobalConfig.Server.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
