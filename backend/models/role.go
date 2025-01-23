package models

import (
	"time"

	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:50;not null;unique" json:"name"`
	Code        string         `gorm:"size:50;not null;unique" json:"code"`
	Description string         `gorm:"size:200" json:"description"`
	Status      int            `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// UserRole 用户角色关联模型
type UserRole struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null;index:idx_user_role" json:"user_id"`
	RoleID    uint           `gorm:"not null;index:idx_user_role" json:"role_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}
