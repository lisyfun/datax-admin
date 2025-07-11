package main

// 这是一个独立的数据重置工具
// 使用方法:
//   go run reset_data.go                    # 使用默认配置文件 config.yaml
//   go run reset_data.go -db mysql          # 指定使用MySQL
//   go run reset_data.go -db postgres       # 指定使用PostgreSQL
//   go run reset_data.go -config config-mysql.yaml    # 指定配置文件

import (
	"datax-admin/config"
	"datax-admin/models"
	"datax-admin/utils/logger"
	"flag"
	"fmt"
	"os"
)

func main() {
	// 解析命令行参数
	var dbType = flag.String("db", "mysql", "数据库类型 (mysql|postgres)")
	var configFile = flag.String("config", "", "配置文件路径")
	var help = flag.Bool("help", false, "显示帮助信息")
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// 根据参数设置配置文件
	if *configFile == "" {
		switch *dbType {
		case "postgres", "postgresql":
			*configFile = "config-postgres.yaml"
		case "mysql":
			*configFile = "config-mysql.yaml"
		default:
			*configFile = "config.yaml"
		}
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(*configFile); os.IsNotExist(err) {
		fmt.Printf("❌ 配置文件不存在: %s\n", *configFile)
		fmt.Println("可用的配置文件:")
		fmt.Println("  - config.yaml (默认)")
		fmt.Println("  - config-mysql.yaml (MySQL)")
		fmt.Println("  - config-postgres.yaml (PostgreSQL)")
		os.Exit(1)
	}

	fmt.Printf("📋 使用配置文件: %s\n", *configFile)
	fmt.Printf("🗄️  数据库类型: %s\n", *dbType)

	// 初始化配置
	config.InitConfigWithFile(*configFile)

	// 初始化日志系统
	if err := logger.Init(); err != nil {
		fmt.Printf("❌ 日志系统初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库
	models.InitDB()

	// 显示当前数据库信息
	dbConfig := config.GlobalConfig.Database
	fmt.Printf("🔗 数据库连接: %s@%s:%d/%s\n", dbConfig.Username, dbConfig.Host, dbConfig.Port, dbConfig.DBName)

	fmt.Println("\n⚠️  警告：此操作将删除所有用户、角色和权限数据！")
	fmt.Print("确认要继续吗？(y/N): ")

	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "y" && confirm != "Y" {
		fmt.Println("❌ 操作已取消")
		return
	}

	fmt.Println("🔄 开始重置数据...")

	// 清空数据
	if err := clearAllData(); err != nil {
		fmt.Printf("❌ 清空数据失败: %v\n", err)
		os.Exit(1)
	}

	// 重新初始化
	if err := models.InitDefaultData(); err != nil {
		fmt.Printf("❌ 重新初始化失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ 数据重置完成！")
	fmt.Println("👤 默认管理员账号：")
	fmt.Println("   用户名: admin")
	fmt.Println("   密码: admin123")
}

func clearAllData() error {
	db := models.DB

	// 按照外键依赖顺序删除数据
	if err := db.Exec("DELETE FROM user_roles").Error; err != nil {
		return fmt.Errorf("删除用户角色关联失败: %v", err)
	}

	if err := db.Exec("DELETE FROM role_permissions").Error; err != nil {
		return fmt.Errorf("删除角色权限关联失败: %v", err)
	}

	if err := db.Exec("DELETE FROM users").Error; err != nil {
		return fmt.Errorf("删除用户失败: %v", err)
	}

	if err := db.Exec("DELETE FROM roles").Error; err != nil {
		return fmt.Errorf("删除角色失败: %v", err)
	}

	if err := db.Exec("DELETE FROM permissions").Error; err != nil {
		return fmt.Errorf("删除权限失败: %v", err)
	}

	// 重置序列/自增ID
	dbType := config.GlobalConfig.Database.Type
	if dbType == "" {
		dbType = "mysql"
	}

	switch dbType {
	case "postgres", "postgresql":
		// PostgreSQL 重置序列
		sequences := []string{"users_id_seq", "roles_id_seq", "permissions_id_seq"}
		for _, seq := range sequences {
			sql := fmt.Sprintf("SELECT setval('%s', 1, false)", seq)
			if err := db.Exec(sql).Error; err != nil {
				fmt.Printf("重置序列 %s 失败: %v\n", seq, err)
			}
		}
	default:
		// MySQL 重置自增ID
		tables := []string{"users", "roles", "permissions"}
		for _, table := range tables {
			sql := fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = 1", table)
			if err := db.Exec(sql).Error; err != nil {
				fmt.Printf("重置 %s 自增ID失败: %v\n", table, err)
			}
		}
	}

	fmt.Println("数据清空完成")
	return nil
}

func showHelp() {
	fmt.Println("📖 DataX Admin 数据重置工具")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  go run reset_data.go [选项]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -db string")
	fmt.Println("        数据库类型 (mysql|postgres) (默认: mysql)")
	fmt.Println("  -config string")
	fmt.Println("        配置文件路径 (如果不指定，会根据 -db 参数自动选择)")
	fmt.Println("  -help")
	fmt.Println("        显示此帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  go run reset_data.go                           # 使用默认MySQL配置")
	fmt.Println("  go run reset_data.go -db mysql                 # 使用MySQL配置")
	fmt.Println("  go run reset_data.go -db postgres              # 使用PostgreSQL配置")
	fmt.Println("  go run reset_data.go -config config.yaml       # 指定配置文件")
	fmt.Println("  go run reset_data.go -config config-mysql.yaml # 指定MySQL配置文件")
	fmt.Println()
	fmt.Println("⚠️  警告: 此工具会删除所有用户、角色和权限数据，请谨慎使用！")
}
