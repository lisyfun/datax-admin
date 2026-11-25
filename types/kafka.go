package types

// KafkaClusterRequest 创建/更新 Kafka 集群的请求体。
// 单独定义请求类型是因为模型上的 Password 使用 json:"-" 避免敏感字段回显。
type KafkaClusterRequest struct {
	Name             string `json:"name" binding:"required,max=100"`
	BrokerServers    string `json:"brokerServers" binding:"required,max=1000"`
	DelayMessage     bool   `json:"delayMessage"`
	SecurityProtocol string `json:"securityProtocol" binding:"max=50"`
	SaslMechanism    string `json:"saslMechanism" binding:"max=50"`
	Username         string `json:"username" binding:"max=100"`
	Password         string `json:"password" binding:"max=255"`
	Description      string `json:"description" binding:"max=500"`
	Status           *bool  `json:"status"`
}
