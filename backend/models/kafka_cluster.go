package models

import (
	"time"
)

// KafkaCluster 表示 Kafka 集群信息
type KafkaCluster struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"size:100;not null;unique" json:"name"`
	BrokerServers    string    `gorm:"size:1000;not null" json:"broker_servers"`
	DelayMessage     bool      `gorm:"default:false" json:"delay_message"`
	SecurityProtocol string    `gorm:"size:50" json:"security_protocol"`
	SaslMechanism    string    `gorm:"size:50" json:"sasl_mechanism"`
	Username         string    `gorm:"size:100" json:"username"`
	Password         string    `gorm:"size:255" json:"password"`
	TopicCount       int       `gorm:"-" json:"topic_count"`
	BrokerCount      int       `gorm:"-" json:"broker_count"`
	ConsumerCount    int       `gorm:"-" json:"consumer_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TableName 指定表名
func (KafkaCluster) TableName() string {
	return "kafka_clusters"
}
