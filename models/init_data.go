package models

import (
	"datax-admin/config"
	"datax-admin/utils/logger"
	"time"

	"gorm.io/gorm"
)

// InitDefaultData 初始化默认数据
func InitDefaultData() error {
	logger.Info("开始初始化默认数据...")

	// 检查是否强制重置
	forceReset := config.GlobalConfig.Init.ForceReset

	// 检查是否已有数据
	var userCount, roleCount, permissionCount int64

	if err := DB.Model(&User{}).Count(&userCount).Error; err != nil {
		return err
	}

	if err := DB.Model(&Role{}).Count(&roleCount).Error; err != nil {
		return err
	}

	if err := DB.Model(&Permission{}).Count(&permissionCount).Error; err != nil {
		return err
	}

	hasData := userCount > 0 || roleCount > 0 || permissionCount > 0

	if hasData && !forceReset {
		logger.Info("检测到已有数据，跳过默认数据初始化")
		logger.Info("如需重新初始化，请在配置文件中设置 init.force_reset: true")
		return nil
	}

	if hasData && forceReset {
		logger.Info("强制重置模式：清空现有数据...")
		if err := clearExistingData(); err != nil {
			return err
		}
	}

	// 开启事务
	return DB.Transaction(func(tx *gorm.DB) error {
		// 1. 创建默认角色
		if err := createDefaultRoles(tx); err != nil {
			return err
		}

		// 2. 创建默认权限菜单
		if err := createDefaultPermissions(tx); err != nil {
			return err
		}

		// 3. 分配角色权限
		if err := assignRolePermissions(tx); err != nil {
			return err
		}

		// 4. 创建默认管理员用户
		if err := createDefaultAdmin(tx); err != nil {
			return err
		}

		logger.Info("默认数据初始化完成")
		return nil
	})
}

// createDefaultRoles 创建默认角色
func createDefaultRoles(tx *gorm.DB) error {
	roles := []Role{
		{
			Name:        "超级管理员",
			Code:        "admin",
			Description: "系统超级管理员，拥有所有权限",
			Status:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "普通用户",
			Code:        "user",
			Description: "普通用户角色，拥有基本权限",
			Status:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, role := range roles {
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
	}

	logger.Info("默认角色创建完成")
	return nil
}

// createDefaultPermissions 创建默认权限菜单（与SQL文件保持一致）
func createDefaultPermissions(tx *gorm.DB) error {
	permissions := []Permission{
		// 根菜单
		{ID: 1, Name: "首页", Code: "root", Type: "menu", ParentID: nil, Path: "/", Component: "layouts/default.vue", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 1},

		// 一级菜单
		{ID: 2, Name: "数据面板", Code: "dashboard", Type: "menu", ParentID: uintPtr(1), Path: "dashboard", Component: "views/dashboard/index.vue", Icon: "icon-dashboard", Sort: 0, Status: 1, Hidden: 0, Cache: 1},
		{ID: 3, Name: "任务管理", Code: "job", Type: "menu", ParentID: uintPtr(1), Path: "job", Component: "views/job/index.vue", Icon: "icon-calendar", Sort: 1, Status: 1, Hidden: 0, Cache: 1},
		{ID: 4, Name: "任务列表", Code: "job.list", Type: "menu", ParentID: uintPtr(3), Path: "list", Component: "views/job/list/index.vue", Icon: "icon-unordered-list", Sort: 0, Status: 1, Hidden: 0, Cache: 1},
		{ID: 5, Name: "执行历史", Code: "job.history", Type: "menu", ParentID: uintPtr(3), Path: "history", Component: "views/job/history/index.vue", Icon: "icon-clock-circle", Sort: 1, Status: 1, Hidden: 0, Cache: 1},
		{ID: 6, Name: "终端管理", Code: "terminal", Type: "menu", ParentID: uintPtr(1), Path: "terminal", Component: "views/terminal/index.vue", Icon: "icon-command", Sort: 2, Status: 1, Hidden: 0, Cache: 1},
		{ID: 7, Name: "终端列表", Code: "terminal.list", Type: "menu", ParentID: uintPtr(6), Path: "list", Component: "views/terminal/list/index.vue", Icon: "icon-desktop", Sort: 0, Status: 1, Hidden: 0, Cache: 1},
		{ID: 8, Name: "终端连接", Code: "terminal.connect", Type: "menu", ParentID: uintPtr(6), Path: "connect/:id", Component: "views/terminal/connect/index.vue", Icon: "", Sort: 1, Status: 1, Hidden: 1, Cache: 1},
		{ID: 9, Name: "工具管理", Code: "tools", Type: "menu", ParentID: uintPtr(1), Path: "tools", Component: "views/tools/index.vue", Icon: "icon-common", Sort: 4, Status: 1, Hidden: 0, Cache: 1},
		{ID: 10, Name: "格式化工具", Code: "tools.json-formatter", Type: "menu", ParentID: uintPtr(9), Path: "json-formatter", Component: "views/tools/JsonFormatter.vue", Icon: "icon-code", Sort: 0, Status: 1, Hidden: 0, Cache: 1},
		{ID: 11, Name: "加解密工具", Code: "tools.crypto", Type: "menu", ParentID: uintPtr(9), Path: "crypto", Component: "views/tools/Crypto.vue", Icon: "icon-lock", Sort: 1, Status: 1, Hidden: 0, Cache: 1},
		{ID: 12, Name: "消息管理", Code: "tools.kafka", Type: "menu", ParentID: uintPtr(1), Path: "kafka", Component: "views/kafka/index.vue", Icon: "icon-apps", Sort: 3, Status: 1, Hidden: 0, Cache: 1},
		{ID: 13, Name: "集群管理", Code: "tools.kafka.cluster", Type: "menu", ParentID: uintPtr(12), Path: "cluster", Component: "views/kafka/cluster/index.vue", Icon: "icon-apps", Sort: 0, Status: 1, Hidden: 0, Cache: 1},
		{ID: 14, Name: "主题管理", Code: "tools.kafka.topic", Type: "menu", ParentID: uintPtr(12), Path: "kafka/cluster/:clusterId/topic", Component: "views/kafka/topic/index.vue", Icon: "", Sort: 1, Status: 1, Hidden: 1, Cache: 1},
		{ID: 15, Name: "消息列表", Code: "tools.kafka.message", Type: "menu", ParentID: uintPtr(12), Path: "kafka/clusters/:clusterId/topics/:topicName/messages", Component: "views/kafka/topic/components/MessageList.vue", Icon: "", Sort: 2, Status: 1, Hidden: 1, Cache: 1},
		{ID: 16, Name: "系统管理", Code: "system", Type: "menu", ParentID: uintPtr(1), Path: "system", Component: "views/system/index.vue", Icon: "icon-settings", Sort: 999, Status: 1, Hidden: 0, Cache: 1},
		{ID: 17, Name: "用户管理", Code: "system.users", Type: "menu", ParentID: uintPtr(16), Path: "users", Component: "views/system/users/index.vue", Icon: "icon-user", Sort: 0, Status: 1, Hidden: 0, Cache: 1},
		{ID: 18, Name: "角色管理", Code: "system.roles", Type: "menu", ParentID: uintPtr(16), Path: "roles", Component: "views/system/roles/index.vue", Icon: "icon-user-group", Sort: 1, Status: 1, Hidden: 0, Cache: 1},
		{ID: 19, Name: "权限管理", Code: "system.permissions", Type: "menu", ParentID: uintPtr(16), Path: "permissions", Component: "views/system/permissions/index.vue", Icon: "icon-safe", Sort: 2, Status: 1, Hidden: 0, Cache: 1},
		// {ID: 20, Name: "登录", Code: "login", Type: "menu", ParentID: nil, Path: "/login", Component: "views/login/index.vue", Icon: "", Sort: 0, Status: 1, Hidden: 1, Cache: 0},
		// {ID: 21, Name: "注册", Code: "register", Type: "menu", ParentID: nil, Path: "/register", Component: "views/register/index.vue", Icon: "", Sort: 0, Status: 0, Hidden: 1, Cache: 0},

		// 按钮权限 - 用户管理
		{ID: 22, Name: "用户查询", Code: "system.users.query", Type: "button", ParentID: uintPtr(17), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 23, Name: "用户创建", Code: "system.users.create", Type: "button", ParentID: uintPtr(17), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 24, Name: "用户编辑", Code: "system.users.update", Type: "button", ParentID: uintPtr(17), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 25, Name: "用户删除", Code: "system.users.delete", Type: "button", ParentID: uintPtr(17), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},

		// 按钮权限 - 角色管理
		{ID: 26, Name: "角色查询", Code: "system.roles.query", Type: "button", ParentID: uintPtr(18), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 27, Name: "角色创建", Code: "system.roles.create", Type: "button", ParentID: uintPtr(18), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 28, Name: "角色编辑", Code: "system.roles.update", Type: "button", ParentID: uintPtr(18), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 29, Name: "角色删除", Code: "system.roles.delete", Type: "button", ParentID: uintPtr(18), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 30, Name: "角色权限设置", Code: "system.roles.permission", Type: "button", ParentID: uintPtr(18), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},

		// 按钮权限 - 权限管理
		{ID: 31, Name: "权限查询", Code: "system.permissions.query", Type: "button", ParentID: uintPtr(19), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 32, Name: "权限创建", Code: "system.permissions.create", Type: "button", ParentID: uintPtr(19), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 33, Name: "权限编辑", Code: "system.permissions.update", Type: "button", ParentID: uintPtr(19), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 34, Name: "权限删除", Code: "system.permissions.delete", Type: "button", ParentID: uintPtr(19), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},

		// 按钮权限 - 任务管理
		{ID: 35, Name: "任务查询", Code: "job.list.query", Type: "button", ParentID: uintPtr(4), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 36, Name: "任务创建", Code: "job.list.create", Type: "button", ParentID: uintPtr(4), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 37, Name: "任务编辑", Code: "job.list.update", Type: "button", ParentID: uintPtr(4), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 38, Name: "任务删除", Code: "job.list.delete", Type: "button", ParentID: uintPtr(4), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 39, Name: "任务执行", Code: "job.list.execute", Type: "button", ParentID: uintPtr(4), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},

		// 按钮权限 - 执行历史
		{ID: 40, Name: "历史查询", Code: "job.history.query", Type: "button", ParentID: uintPtr(5), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 41, Name: "历史详情", Code: "job.history.detail", Type: "button", ParentID: uintPtr(5), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},

		// 按钮权限 - 终端管理
		{ID: 42, Name: "终端查询", Code: "terminal.list.query", Type: "button", ParentID: uintPtr(7), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 43, Name: "终端创建", Code: "terminal.list.create", Type: "button", ParentID: uintPtr(7), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 44, Name: "终端编辑", Code: "terminal.list.update", Type: "button", ParentID: uintPtr(7), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 45, Name: "终端删除", Code: "terminal.list.delete", Type: "button", ParentID: uintPtr(7), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 46, Name: "终端连接", Code: "terminal.list.connect", Type: "button", ParentID: uintPtr(7), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},

		// 按钮权限 - Kafka集群管理
		{ID: 47, Name: "集群查询", Code: "tools.kafka.cluster.query", Type: "button", ParentID: uintPtr(13), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 48, Name: "集群创建", Code: "tools.kafka.cluster.create", Type: "button", ParentID: uintPtr(13), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 49, Name: "集群编辑", Code: "tools.kafka.cluster.update", Type: "button", ParentID: uintPtr(13), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 50, Name: "集群删除", Code: "tools.kafka.cluster.delete", Type: "button", ParentID: uintPtr(13), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},

		// 操作日志管理
		{ID: 51, Name: "操作管理", Code: "system.logs", Type: "menu", ParentID: uintPtr(16), Path: "logs", Component: "views/system/logs/index.vue", Icon: "icon-file", Sort: 2, Status: 1, Hidden: 0, Cache: 1},

		// 终端上传下载权限
		{ID: 56, Name: "终端上传", Code: "terminal.list.upload", Type: "button", ParentID: uintPtr(7), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
		{ID: 57, Name: "终端下载", Code: "terminal.list.download", Type: "button", ParentID: uintPtr(7), Path: "", Component: "", Icon: "", Sort: 0, Status: 1, Hidden: 0, Cache: 0},
	}

	for _, permission := range permissions {
		permission.CreatedAt = time.Now()
		permission.UpdatedAt = time.Now()
		if err := tx.Create(&permission).Error; err != nil {
			return err
		}
	}

	logger.Info("默认权限菜单创建完成")
	return nil
}

// assignRolePermissions 分配角色权限（与SQL文件保持一致）
func assignRolePermissions(tx *gorm.DB) error {
	// 获取管理员角色
	var adminRole Role
	if err := tx.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	// 获取普通用户角色
	var userRole Role
	if err := tx.Where("code = ?", "user").First(&userRole).Error; err != nil {
		return err
	}

	// 为管理员分配所有权限（ID 1-57）
	adminPermissionIDs := []uint{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
		22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
		41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 56, 57,
	}

	for _, permissionID := range adminPermissionIDs {
		rolePermission := RolePermission{
			RoleID:       adminRole.ID,
			PermissionID: permissionID,
		}
		if err := tx.Create(&rolePermission).Error; err != nil {
			return err
		}
	}

	// 为普通用户分配基本权限（根据SQL文件中的最新配置）
	userPermissionIDs := []uint{
		2,  // 数据面板
		3,  // 任务管理
		9,  // 工具管理
		4,  // 任务列表
		10, // Json格式化
		11, // 加解密工具
		12, // 消息管理
		13, // 集群管理
		14, // 主题管理
		15, // 消息列表
		56, // 终端上传
		57, // 终端下载
	}

	for _, permissionID := range userPermissionIDs {
		rolePermission := RolePermission{
			RoleID:       userRole.ID,
			PermissionID: permissionID,
		}
		if err := tx.Create(&rolePermission).Error; err != nil {
			return err
		}
	}

	logger.Info("角色权限分配完成")
	return nil
}

// createDefaultAdmin 创建默认管理员用户（与SQL文件保持一致）
func createDefaultAdmin(tx *gorm.DB) error {
	// 创建管理员用户（与SQL文件中的数据一致）
	admin := User{
		ID:        1, // 指定ID为1，与SQL文件一致
		Username:  "admin",
		Nickname:  "超级管理员",
		Email:     "admin@example.com", // 与SQL文件一致
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 从配置中获取默认密码
	defaultPassword := config.GlobalConfig.Init.DefaultPassword
	if defaultPassword == "" {
		defaultPassword = "admin123" // 默认值
	}
	if err := admin.SetPassword(defaultPassword); err != nil {
		return err
	}

	if err := tx.Create(&admin).Error; err != nil {
		return err
	}

	// 获取管理员角色
	var adminRole Role
	if err := tx.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	// 分配管理员角色
	userRole := UserRole{
		UserID:    admin.ID,
		RoleID:    adminRole.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := tx.Create(&userRole).Error; err != nil {
		return err
	}

	// 记录日志时不显示实际密码
	logger.Info("默认管理员用户创建完成 - 用户名: admin, 密码已从配置中读取")
	return nil
}

// clearExistingData 清空现有数据
func clearExistingData() error {
	logger.Info("清空现有数据...")

	// 使用事务确保数据一致性
	return DB.Transaction(func(tx *gorm.DB) error {
		// 按照外键依赖顺序删除数据
		if err := tx.Exec("DELETE FROM user_roles").Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM role_permissions").Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM users").Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM roles").Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM permissions").Error; err != nil {
			return err
		}

		// 重置自增ID/序列
		dbType := config.GlobalConfig.Database.Type
		if dbType == "" {
			dbType = "mysql"
		}

		switch dbType {
		case "postgres", "postgresql":
			// PostgreSQL 使用 TRUNCATE 重置序列（更可靠的方法）
			tables := []string{"permissions", "roles", "users", "role_permissions", "user_roles"}
			for _, table := range tables {
				// 使用 TRUNCATE RESTART IDENTITY 重置序列
				sql := "TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE"
				if err := tx.Exec(sql).Error; err != nil {
					logger.Error("重置表 %s 失败: %v", table, err)
					// 如果 TRUNCATE 失败，尝试手动重置序列
					seqName := table + "_id_seq"
					resetSQL := "SELECT setval('" + seqName + "', 1, false)"
					if err2 := tx.Exec(resetSQL).Error; err2 != nil {
						logger.Error("手动重置序列 %s 也失败: %v", seqName, err2)
					} else {
						logger.Info("手动重置序列成功: %s", seqName)
					}
				} else {
					logger.Info("成功重置表和序列: %s", table)
				}
			}
		default:
			// MySQL 重置自增ID
			tables := []string{"users", "roles", "permissions"}
			for _, table := range tables {
				sql := "ALTER TABLE " + table + " AUTO_INCREMENT = 1"
				if err := tx.Exec(sql).Error; err != nil {
					logger.Error("重置 %s 自增ID失败: %v", table, err)
				} else {
					logger.Info("成功重置 %s 自增ID", table)
				}
			}
		}

		logger.Info("现有数据清空完成")
		return nil
	})
}

// uintPtr 辅助函数，返回uint指针
func uintPtr(u uint) *uint {
	return &u
}
