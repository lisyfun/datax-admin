package models

import (
	"datax-admin/utils"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"size:50;not null;unique" json:"username"`
	Password  string         `gorm:"size:100;not null" json:"-"`
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Email     string         `gorm:"size:100" json:"email"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Status    int            `gorm:"default:1" json:"status"` // 1: 正常, 0: 禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 在创建记录前加密密码
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.Password = utils.HashPassword(u.Password)
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) bool {
	return utils.CheckPassword(password, u.Password)
}
