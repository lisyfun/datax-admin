package main

import (
	"datax-admin/config"
	"datax-admin/models"
	"datax-admin/routes"
	"datax-admin/services"
	"datax-admin/utils/logger"
	"log"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/utkusen/baitroute/go/pkg/baitroute"
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

	// 初始化任务调度器并加载运行中的任务
	jobService := services.GetJobService()
	if err := jobService.InitJobScheduler(); err != nil {
		logger.Error("初始化任务调度器失败: %v", err)
	} else {
		logger.Info("任务调度器初始化成功")
	}

	// 初始化服务
	kafkaService := services.NewKafkaService()
	kafkaTopicService := services.NewKafkaTopicService(kafkaService)

	// 启动Kafka集群健康检查
	kafkaService.StartClusterHealthCheck(10 * time.Minute)

	// 启动主题同步任务
	kafkaTopicService.StartTopicSyncTask(5 * time.Minute)

	// 设置gin模式
	gin.SetMode(config.GlobalConfig.Server.Mode)

	// 创建路由
	r := gin.Default()

	// 设置文件上传大小限制为 500MB
	r.MaxMultipartMemory = config.GlobalConfig.Server.MaxFileSize << 20

	b := initBaitRoute()

	// 注册蜜罐路由到 Gin
	if err := b.RegisterWithGin(r); err != nil {
		log.Fatalf("Failed to register bait endpoints: %v", err)
	}

	// 注册路由
	routes.SetupRoutes(r)

	// 启动服务器
	logger.Info("服务启动成功，监听端口%s", config.GlobalConfig.Server.Port)
	if err := r.Run(config.GlobalConfig.Server.Port); err != nil {
		logger.Fatal("服务启动失败: %v", err)
	}
}

// BaitRoute
func initBaitRoute() *baitroute.BaitRoute {
	rulesPath := filepath.Join("rules") // 你的 rules 路径
	b, err := baitroute.NewBaitRoute(rulesPath)
	if err != nil {
		log.Fatalf("Failed to initialize baitroute: %v", err)
	}

	b.OnBaitHit(func(alert baitroute.Alert) {
		log.Printf("Bait Alert: Method=%s Path=%s SourceIP=%s Rule=%s",
			alert.Method, alert.Path, alert.SourceIP, alert.RuleName)
	})
	return b
}
