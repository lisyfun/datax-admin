package models

import (
	"time"
)

// OperationLog 操作日志
type OperationLog struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"not null;comment:操作用户ID" json:"user_id"`
	Username    string    `gorm:"size:50;not null;comment:操作用户名" json:"username"`
	Module      string    `gorm:"size:50;not null;comment:操作模块" json:"module"`
	Action      string    `gorm:"size:50;not null;comment:操作动作" json:"action"`
	Description string    `gorm:"size:500;comment:操作描述" json:"description"`
	IP          string    `gorm:"size:50;comment:操作IP" json:"ip"`
	UserAgent   string    `gorm:"size:500;comment:用户代理" json:"user_agent"`
	RequestData string    `gorm:"type:text;comment:请求数据" json:"request_data"`
	Status      int       `gorm:"default:1;comment:操作状态 1:成功 0:失败" json:"status"`
	ErrorMsg    string    `gorm:"size:500;comment:错误信息" json:"error_msg"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (OperationLog) TableName() string {
	return "operation_logs"
}

// 操作模块常量
const (
	ModuleUser       = "user"       // 用户管理
	ModuleRole       = "role"       // 角色管理
	ModulePermission = "permission" // 权限管理
	ModuleMenu       = "menu"       // 菜单管理
	ModuleJob        = "job"        // 任务管理
	ModuleTerminal   = "terminal"   // 终端管理
	ModuleKafka      = "kafka"      // Kafka管理
	ModuleSystem     = "system"     // 系统管理
)

// 操作动作常量
const (
	ActionCreate = "create" // 创建
	ActionUpdate = "update" // 更新
	ActionDelete = "delete" // 删除
	ActionView   = "view"   // 查看
	ActionLogin  = "login"  // 登录
	ActionLogout = "logout" // 登出
	ActionExport = "export" // 导出
	ActionImport = "import" // 导入
)

// 操作状态常量
const (
	StatusSuccess = 1 // 成功
	StatusFailed  = 0 // 失败
)
