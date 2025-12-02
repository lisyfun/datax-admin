package models

import (
    "time"
    "gorm.io/gorm"
)

// RedisConnection Redis 连接配置
type RedisConnection struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    Name      string         `gorm:"size:100;not null;uniqueIndex;comment:连接名称" json:"name"`
    Host      string         `gorm:"size:200;not null;comment:主机地址" json:"host"`
    Port      int            `gorm:"default:6379;comment:端口" json:"port"`
    Username  string         `gorm:"size:100;comment:用户名" json:"username"`
    Password  string         `gorm:"size:200;comment:密码(加密存储)" json:"-"`
    DB        int            `gorm:"default:0;comment:数据库索引" json:"db"`
    UseTLS    bool           `gorm:"default:false;comment:是否启用TLS" json:"use_tls"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (RedisConnection) TableName() string {
    return "redis_connections"
}

