package models

import (
	"datax-admin/utils/secret"
	"time"

	"gorm.io/gorm"
)

// KafkaCluster 表示 Kafka 集群信息
type KafkaCluster struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	Name               string    `gorm:"size:100;not null;unique" json:"name"`
	BrokerServers      string    `gorm:"size:1000;not null" json:"brokerServers"`
	DelayMessage       bool      `gorm:"default:false" json:"delayMessage"`
	SecurityProtocol   string    `gorm:"size:50" json:"securityProtocol"`
	SaslMechanism      string    `gorm:"size:50" json:"saslMechanism"`
	Username           string    `gorm:"size:100" json:"username"`
	Password           string    `gorm:"size:512" json:"-"`
	Description        string    `gorm:"size:500" json:"description"`
	Status             bool      `gorm:"default:true;index:idx_kafka_status" json:"status"`
	LastCheckTime      time.Time `gorm:"index:idx_kafka_last_check" json:"lastCheckTime"`
	TopicCount         int       `gorm:"-" json:"topicCount"`
	BrokerCount        int       `gorm:"-" json:"brokerCount"`
	ConsumerGroupCount int       `gorm:"-" json:"consumerGroupCount"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

// TableName 指定表名
func (KafkaCluster) TableName() string {
	return "kafka_clusters"
}

// BeforeSave 加密密码。
func (k *KafkaCluster) BeforeSave(tx *gorm.DB) error {
	var err error
	k.Password, err = secret.Encrypt(k.Password)
	return err
}

// AfterFind 解密密码。
func (k *KafkaCluster) AfterFind(tx *gorm.DB) error {
	var err error
	k.Password, err = secret.Decrypt(k.Password)
	return err
}
