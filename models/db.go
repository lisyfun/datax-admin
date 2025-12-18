package models

import (
	"datax-admin/config"
	customLogger "datax-admin/utils/logger"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/gaussdb"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	var dialector gorm.Dialector

	// 根据数据库类型构建连接字符串和dialector
	dbType := config.GlobalConfig.Database.Type
	if dbType == "" {
		dbType = "mysql" // 默认使用MySQL
	}

	logMode := logger.Warn
	switch config.GlobalConfig.Database.LogMode {
	case "info":
		logMode = logger.Info
	case "warn":
		logMode = logger.Warn
	case "error":
		logMode = logger.Error
	}

	switch dbType {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			config.GlobalConfig.Database.Username,
			config.GlobalConfig.Database.Password,
			config.GlobalConfig.Database.Host,
			config.GlobalConfig.Database.Port,
			config.GlobalConfig.Database.DBName,
			config.GlobalConfig.Database.Charset,
		)
		dialector = mysql.Open(dsn)
	case "postgres", "postgresql":
		sslMode := config.GlobalConfig.Database.SSLMode
		if sslMode == "" {
			sslMode = "disable"
		}
		timeZone := config.GlobalConfig.Database.TimeZone
		if timeZone == "" {
			timeZone = "Asia/Shanghai"
		}
		schema := config.GlobalConfig.Database.Schema
		if schema == "" {
			schema = "public"
		}
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s search_path=%s",
			config.GlobalConfig.Database.Host,
			config.GlobalConfig.Database.Username,
			config.GlobalConfig.Database.Password,
			config.GlobalConfig.Database.DBName,
			config.GlobalConfig.Database.Port,
			sslMode,
			timeZone,
			schema,
		)
		dialector = postgres.Open(dsn)
	case "gaussdb":
		sslMode := config.GlobalConfig.Database.SSLMode
		if sslMode == "" {
			sslMode = "disable"
		}
		timeZone := config.GlobalConfig.Database.TimeZone
		if timeZone == "" {
			timeZone = "Asia/Shanghai"
		}
		schema := config.GlobalConfig.Database.Schema
		if schema == "" {
			schema = "public"
		}
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s search_path=%s",
			config.GlobalConfig.Database.Host,
			config.GlobalConfig.Database.Username,
			config.GlobalConfig.Database.Password,
			config.GlobalConfig.Database.DBName,
			config.GlobalConfig.Database.Port,
			sslMode,
			timeZone,
			schema,
		)
		dialector = gaussdb.Open(dsn)
	default:
		customLogger.Fatal("不支持的数据库类型: %s", dbType)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger:                                   logger.Default.LogMode(logMode),
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		customLogger.Fatal("数据库连接失败: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		customLogger.Fatal("获取数据库实例失败: %v", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(config.GlobalConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.GlobalConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(config.GlobalConfig.Database.MaxLifetime))

	customLogger.Info("数据库连接成功")

	// 自动迁移数据库表结构
	if err := AutoMigrate(); err != nil {
		customLogger.Fatal("数据库表结构迁移失败: %v", err)
	}

	// 初始化默认数据
	if err := InitDefaultData(); err != nil {
		customLogger.Fatal("初始化默认数据失败: %v", err)
	}
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	// 逐个迁移表，避免外键约束问题
	models := []any{
		&User{},
		&Role{},
		&UserRole{},
		&Permission{},
		&RolePermission{},
		&Job{},
		&JobHistory{},
		&Terminal{},
		&LoginLog{},
		&KafkaCluster{},
		&KafkaTopic{},
		&OperationLog{},
		&RedisConnection{},
	}

	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			// 忽略外键约束删除错误，这些通常是无害的
			if !isIgnorableError(err) {
				customLogger.Error("迁移表失败: %v, 模型: %T", err, model)
			}
			// 继续迁移其他表，不中断整个过程
			continue
		}
	}

	customLogger.Info("数据库表结构迁移完成")

	return nil
}

// isIgnorableError 检查是否是可以忽略的错误
func isIgnorableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	// 忽略外键约束删除错误
	ignorablePatterns := []string{
		"Can't DROP",
		"check that column/key exists",
		"uni_users_username",
		"uni_terminals_name",
	}

	for _, pattern := range ignorablePatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	return false
}
