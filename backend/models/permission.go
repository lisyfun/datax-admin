package models

import (
	"time"

	"gorm.io/gorm"
)

// Permission 权限模型
type Permission struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Code      string         `gorm:"size:50;not null;unique" json:"code"`
	Type      string         `gorm:"size:20;not null" json:"type"` // menu: 菜单, button: 按钮
	ParentID  *uint          `gorm:"default:null;index" json:"parent_id"`
	Path      string         `gorm:"size:200" json:"path"`      // 前端路由路径
	Component string         `gorm:"size:200" json:"component"` // 前端组件路径
	Icon      string         `gorm:"size:50" json:"icon"`       // 图标
	Sort      int            `gorm:"default:0" json:"sort"`     // 排序
	Status    int            `gorm:"default:1" json:"status"`   // 1: 启用, 0: 禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// RolePermission 角色权限关联模型
type RolePermission struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	RoleID       uint           `gorm:"not null;index" json:"role_id"`
	PermissionID uint           `gorm:"not null;index" json:"permission_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}
