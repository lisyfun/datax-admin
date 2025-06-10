package models

import (
	"time"
)

// LoginLog 登录日志
type LoginLog struct {
	ID        uint      `gorm:"primarykey"`
	UserID    uint      `gorm:"not null;comment:用户ID"`
	Username  string    `gorm:"size:50;not null;comment:用户名"`
	IP        string    `gorm:"size:50;not null;comment:登录IP"`
	LoginTime time.Time `gorm:"not null;comment:登录时间"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
