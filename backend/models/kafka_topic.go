package models

import (
	"time"
)

// KafkaTopic Kafka主题信息
type KafkaTopic struct {
	ID              uint      `json:"id" gorm:"primarykey"`
	ClusterID       uint      `json:"clusterId" gorm:"not null;index:idx_cluster_topic"`
	Name            string    `json:"name" gorm:"not null;index:idx_cluster_topic"`
	Partitions      int       `json:"partitions"`
	Replicas        int       `json:"replicas"`
	AvgLogSize      int64     `json:"avgLogSize"`
	TotalLogSize    int64     `json:"totalLogSize"`
	LastRefreshTime time.Time `json:"lastRefreshTime"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// TableName 指定表名
func (KafkaTopic) TableName() string {
	return "kafka_topics"
}
