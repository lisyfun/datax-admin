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

	// 自动迁移数据库表
	err = DB.AutoMigrate(
		&User{},
		&Role{},
		&UserRole{},
		&Permission{},
		&RolePermission{},
		&Menu{},
		&Job{},
		&JobHistory{},
		&LoginLog{},
	)
	if err != nil {
		customLogger.Fatal("数据库自动迁移失败: %v", err)
	}

	// 初始化基础数据
	initBaseData()

	customLogger.Info("数据库连接成功")
}

// initBaseData 初始化基础数据
func initBaseData() {
	// 创建超级管理员角色（如果不存在）
	var adminRole Role
	if err := DB.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			adminRole = Role{
				Name:        "超级管理员",
				Code:        "admin",
				Description: "系统超级管理员",
				Status:      1,
			}
			DB.Create(&adminRole)
		}
	}

	// 为 ID 为 1 的用户分配超级管理员角色（如果存在）
	var adminUser User
	if err := DB.First(&adminUser, 1).Error; err == nil {
		var userRole UserRole
		if err := DB.Where("user_id = ? AND role_id = ?", adminUser.ID, adminRole.ID).First(&userRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				userRole = UserRole{
					UserID: adminUser.ID,
					RoleID: adminRole.ID,
				}
				DB.Create(&userRole)
			}
		}
	}
}
