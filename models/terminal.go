package models

import (
	"time"

	"gorm.io/gorm"
)

// Terminal 终端模型
type Terminal struct {
	gorm.Model
	Name          string    `gorm:"type:varchar(100);uniqueIndex;not null;comment:终端名称"`
	Host          string    `gorm:"type:varchar(255);not null;comment:主机地址"`
	Port          int       `gorm:"type:int;not null;default:22;comment:SSH端口"`
	Username      string    `gorm:"type:varchar(50);not null;comment:用户名"`
	AuthType      string    `gorm:"type:varchar(20);not null;default:'password';comment:认证类型(password/key)"`
	Password      string    `gorm:"type:varchar(255);comment:密码"`
	KeyFile       string    `gorm:"type:text;comment:密钥文件内容"`
	KeyPassphrase string    `gorm:"type:varchar(255);comment:密钥文件密码"`
	Status        string    `gorm:"type:varchar(20);not null;default:'offline';comment:状态(online/offline)"`
	LastSeen      time.Time `gorm:"not null;comment:最后在线时间"`
}

// TableName 指定表名
func (Terminal) TableName() string {
	return "terminals"
}

// BeforeCreate GORM 创建钩子
func (t *Terminal) BeforeCreate(tx *gorm.DB) error {
	if t.LastSeen.IsZero() {
		t.LastSeen = time.Now()
	}
	return nil
}
