package models

import (
	"time"

	"gorm.io/gorm"
)

// Menu 菜单模型
type Menu struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	ParentID  uint           `gorm:"default:0;index" json:"parent_id"` // 父菜单ID
	Name      string         `gorm:"size:50;not null" json:"name"`     // 菜单名称
	Path      string         `gorm:"size:200" json:"path"`             // 路由路径
	Component string         `gorm:"size:200" json:"component"`        // 组件路径
	Icon      string         `gorm:"size:50" json:"icon"`              // 图标
	Sort      int            `gorm:"default:0" json:"sort"`            // 排序
	Status    int            `gorm:"default:1" json:"status"`          // 状态：1-启用，0-禁用
	Hidden    bool           `gorm:"default:false" json:"hidden"`      // 是否隐藏
	Cache     bool           `gorm:"default:true" json:"cache"`        // 是否缓存
	Type      int            `gorm:"default:1" json:"type"`            // 类型：1-菜单，2-按钮
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}
