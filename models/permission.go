package models

import (
	"time"

	"gorm.io/gorm"
)

// Permission 权限模型
// Permission 权限模型
type Permission struct {
    ID          uint           `gorm:"primarykey" json:"id"`
    Name        string         `gorm:"size:50;not null" json:"name"`
    Code        string         `gorm:"size:50;not null;unique" json:"code"`
    Type        string         `gorm:"size:20;not null" json:"type"` // menu: 菜单, button: 按钮
    ParentID    *uint          `gorm:"default:null;index" json:"parent_id"`
    Path        string         `gorm:"size:200" json:"path"`         // 前端路由路径
    Component   string         `gorm:"size:200" json:"component"`    // 前端组件路径
    Icon        string         `gorm:"size:50" json:"icon"`          // 图标
    Sort        int            `gorm:"default:0" json:"sort"`        // 排序
    Status      int            `gorm:"default:1" json:"status"`      // 1: 启用, 0: 禁用
    Hidden      int            `gorm:"default:0" json:"hidden"`      // 是否隐藏菜单 0-显示 1-隐藏
    Cache       int            `gorm:"default:1" json:"cache"`       // 是否缓存页面 0-不缓存 1-缓存
    IsExternal  int            `gorm:"default:0" json:"is_external"` // 是否为外部链接 0-否 1-是
    ExternalURL string         `gorm:"size:500" json:"external_url"` // 外部链接地址
    OpenType    int            `gorm:"default:1" json:"open_type"`   // 打开方式 0-内嵌 1-新窗口
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// RolePermission 角色权限关联模型
type RolePermission struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	RoleID       uint           `gorm:"not null;index:idx_role_perm" json:"role_id"`
	PermissionID uint           `gorm:"not null;index:idx_role_perm" json:"permission_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}
