package main

import (
	"datax-admin/config"
	"datax-admin/middleware"
	"datax-admin/models"
	"datax-admin/routes"
	"datax-admin/services"
	"datax-admin/utils/logger"
	"datax-admin/utils/secret"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	// 初始化配置
	config.InitConfig()

	// 初始化日志系统
	if err := logger.Init(); err != nil {
		logger.Fatal("初始化日志系统失败: %v", err)
	}

	// 初始化敏感数据加密（必须在数据库初始化前完成，GORM 钩子会用到）
	if config.GlobalConfig.Security.EncryptionKey == "" {
		logger.Info("未配置 security.encryption_key，将使用 data/encryption.key（请持久化该文件，丢失后无法解密已有凭据）")
	}
	if err := secret.Init(config.GlobalConfig.Security.EncryptionKey); err != nil {
		logger.Fatal("初始化敏感数据加密失败: %v", err)
	}

	// 初始化数据库
	models.InitDB()

	// 初始化加密中间件
	if err := middleware.InitCrypto(); err != nil {
		logger.Fatal("加密中间件初始化失败: %v", err)
	}

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

	// 启动任务日志清理器
	if config.GlobalConfig.JobLog.AutoCleanup {
		logCleaner := services.NewJobLogCleaner(time.Duration(config.GlobalConfig.JobLog.MaxAge) * 24 * time.Hour)
		logCleaner.StartAutoCleanup()
		logger.Info("任务日志自动清理器已启动，保留天数: %d", config.GlobalConfig.JobLog.MaxAge)
	}

	// 设置gin模式
	gin.SetMode(config.GlobalConfig.Server.Mode)

	// 创建路由
	r := gin.New()

	// 配置可信反向代理：默认不信任任何代理，避免客户端伪造 X-Forwarded-For
	// 绕过登录限流等基于来源 IP 的保护。部署在 nginx 之后时请配置
	// server.trusted_proxies（例如 ["127.0.0.1", "::1"]）。
	if len(config.GlobalConfig.Server.TrustedProxies) > 0 {
		if err := r.SetTrustedProxies(config.GlobalConfig.Server.TrustedProxies); err != nil {
			logger.Fatal("配置可信代理失败: %v", err)
		}
	} else if err := r.SetTrustedProxies(nil); err != nil {
		logger.Fatal("设置可信代理失败: %v", err)
	}

	// 使用自定义日志中间件
	// r.Use(middleware.GinLogger())
	r.Use(gin.Recovery())

	// 设置文件上传大小限制为 500MB
	r.MaxMultipartMemory = config.GlobalConfig.Server.MaxFileSize << 20

	// 注册路由
	routes.SetupRoutes(r)

	// 启动服务器
	logger.Info("服务启动成功，监听端口%s", config.GlobalConfig.Server.Port)
	if err := r.Run(config.GlobalConfig.Server.Port); err != nil {
		logger.Fatal("服务启动失败: %v", err)
	}
}
