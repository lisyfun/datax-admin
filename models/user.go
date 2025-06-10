package models

import (
	"datax-admin/utils"
	"time"
)

// User 用户模型
type User struct {
	ID        uint   `gorm:"primarykey"`
	Username  string `gorm:"size:50;not null;uniqueIndex;comment:用户名"`
	Password  string `gorm:"size:100;not null;comment:密码"`
	Nickname  string `gorm:"size:50;comment:昵称"`
	Email     string `gorm:"size:100;comment:邮箱"`
	Avatar    string `gorm:"size:255;comment:头像"`
	Status    int    `gorm:"default:1;comment:状态 0:禁用 1:启用"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

// SetPassword 设置密码
func (u *User) SetPassword(password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

// CheckPassword 检查密码
func (u *User) CheckPassword(password string) bool {
	return utils.ComparePasswords(u.Password, password)
}
