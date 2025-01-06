package models

import (
	"datax-admin/utils"
	"errors"
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

// SetPassword 设置加密后的密码
func (u *User) SetPassword(password string) error {
	if password == "" {
		return errors.New("密码不能为空")
	}
	hashedPassword := utils.HashPassword(password)
	if hashedPassword == "" {
		return errors.New("密码加密失败")
	}
	u.Password = hashedPassword
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) bool {
	if password == "" || u.Password == "" {
		return false
	}
	return utils.CheckPassword(password, u.Password)
}
