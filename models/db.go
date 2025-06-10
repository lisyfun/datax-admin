package models

import (
	"datax-admin/config"
	customLogger "datax-admin/utils/logger"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.GlobalConfig.Database.Username,
		config.GlobalConfig.Database.Password,
		config.GlobalConfig.Database.Host,
		config.GlobalConfig.Database.Port,
		config.GlobalConfig.Database.DBName,
		config.GlobalConfig.Database.Charset,
	)

	logMode := logger.Warn
	switch config.GlobalConfig.Database.LogMode {
	case "info":
		logMode = logger.Info
	case "warn":
		logMode = logger.Warn
	case "error":
		logMode = logger.Error
	case "debug":
		logMode = logger.Info
	}

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
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
}
