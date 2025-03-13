package services

import (
	"context"
	"datax-admin/models"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

// KafkaService 提供 Kafka 集群管理相关的服务
type KafkaService struct{}

// KafkaClusterListResponse 集群列表响应
type KafkaClusterListResponse struct {
	Total int64                 `json:"total"`
	Items []models.KafkaCluster `json:"items"`
}

// KafkaTopicListResponse Topic列表响应
type KafkaTopicListResponse struct {
	Total int64        `json:"total"`
	Items []KafkaTopic `json:"items"`
}

// KafkaTopic Topic信息
type KafkaTopic struct {
	Name        string `json:"name"`
	Partitions  int    `json:"partitions"`
	Replicas    int    `json:"replicas"`
	AvgLogSize  string `json:"avgLogSize"`
	LogSize     string `json:"logSize"`
	ClusterID   uint   `json:"clusterId"`
	ClusterName string `json:"clusterName"`
}

// KafkaMessage 消息信息
type KafkaMessage struct {
	Partition int       `json:"partition"`
	Key       string    `json:"key"`
	Offset    int64     `json:"offset"`
	Timestamp time.Time `json:"timestamp"`
	MsgType   string    `json:"msgType"`
	MsgData   string    `json:"msgData"`
	Value     string    `json:"value"`
}

// ValidateClusterConnection 验证集群连接
func (s *KafkaService) ValidateClusterConnection(cluster *models.KafkaCluster) error {
	admin, err := s.getKafkaAdminClient(cluster)
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}
	defer admin.Close()

	// 尝试列出主题来验证连接
	_, err = admin.ListTopics()
	if err != nil {
		return fmt.Errorf("验证失败: %v", err)
	}

	return nil
}

// CreateKafkaCluster 创建 Kafka 集群
func (s *KafkaService) CreateKafkaCluster(cluster *models.KafkaCluster) error {
	// 先验证连接
	if err := s.ValidateClusterConnection(cluster); err != nil {
		return err
	}

	// 验证成功后创建集群
	return models.DB.Create(cluster).Error
}

// UpdateKafkaCluster 更新 Kafka 集群
func (s *KafkaService) UpdateKafkaCluster(cluster *models.KafkaCluster) error {
	// 先验证连接
	if err := s.ValidateClusterConnection(cluster); err != nil {
		return err
	}

	// 验证成功后更新集群
	return models.DB.Save(cluster).Error
}

// DeleteKafkaCluster 删除 Kafka 集群
func (s *KafkaService) DeleteKafkaCluster(id uint) error {
	return models.DB.Delete(&models.KafkaCluster{}, id).Error
}

// GetKafkaCluster 获取 Kafka 集群详情
func (s *KafkaService) GetKafkaCluster(id uint) (*models.KafkaCluster, error) {
	var cluster models.KafkaCluster
	if err := models.DB.First(&cluster, id).Error; err != nil {
		return nil, err
	}

	// 获取集群统计信息
	s.enrichClusterStats(&cluster)

	return &cluster, nil
}

// ListKafkaClusters 获取 Kafka 集群列表
func (s *KafkaService) ListKafkaClusters(page, pageSize int, search string) (*KafkaClusterListResponse, error) {
	var clusters []models.KafkaCluster
	var total int64

	query := models.DB.Model(&models.KafkaCluster{})

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&clusters).Error; err != nil {
		return nil, err
	}

	// 获取每个集群的统计信息
	for i := range clusters {
		s.enrichClusterStats(&clusters[i])
	}

	return &KafkaClusterListResponse{
		Total: total,
		Items: clusters,
	}, nil
}

// enrichClusterStats 丰富集群统计信息
func (s *KafkaService) enrichClusterStats(cluster *models.KafkaCluster) {
	// 连接到 Kafka 集群
	admin, err := s.getKafkaAdminClient(cluster)
	if err != nil {
		return
	}
	defer admin.Close()

	// 获取 Topic 数量
	topics, err := admin.ListTopics()
	if err == nil {
		cluster.TopicCount = len(topics)
	}

	// 获取 Broker 数量
	// 从 BrokerServers 字段计算 Broker 数量
	cluster.BrokerCount = len(strings.Split(cluster.BrokerServers, ","))

	// 获取消费者组数量
	groups, err := admin.ListConsumerGroups()
	if err == nil {
		cluster.ConsumerGroupCount = len(groups)
	}
}

// getKafkaAdminClient 获取 Kafka 管理客户端
func (s *KafkaService) getKafkaAdminClient(cluster *models.KafkaCluster) (sarama.ClusterAdmin, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0 // 使用较新的版本

	// 设置安全配置
	if cluster.SecurityProtocol != "" {
		switch cluster.SecurityProtocol {
		case "SASL_PLAINTEXT":
			config.Net.SASL.Enable = true
			config.Net.SASL.Mechanism = sarama.SASLMechanism(cluster.SaslMechanism)
			config.Net.SASL.User = cluster.Username
			config.Net.SASL.Password = cluster.Password
		case "SASL_SSL":
			config.Net.SASL.Enable = true
			config.Net.SASL.Mechanism = sarama.SASLMechanism(cluster.SaslMechanism)
			config.Net.SASL.User = cluster.Username
			config.Net.SASL.Password = cluster.Password
			config.Net.TLS.Enable = true
		}
	}

	brokers := strings.Split(cluster.BrokerServers, ",")
	return sarama.NewClusterAdmin(brokers, config)
}

// ListTopics 获取 Topic 列表
func (s *KafkaService) ListTopics(clusterID uint, page, pageSize int, search string) (*KafkaTopicListResponse, error) {
	// 获取集群信息
	cluster, err := s.GetKafkaCluster(clusterID)
	if err != nil {
		return nil, err
	}

	// 连接到 Kafka 集群
	admin, err := s.getKafkaAdminClient(cluster)
	if err != nil {
		return nil, err
	}
	defer admin.Close()

	// 获取所有 Topic
	topics, err := admin.ListTopics()
	if err != nil {
		return nil, err
	}

	// 获取 Topic 详情
	topicNames := make([]string, 0, len(topics))
	for name := range topics {
		if search != "" && !strings.Contains(name, search) {
			continue
		}
		topicNames = append(topicNames, name)
	}

	// 获取 Topic 元数据
	metadata, err := admin.DescribeTopics(topicNames)
	if err != nil {
		return nil, err
	}

	var topicList []KafkaTopic
	for _, topicMetadata := range metadata {
		// 计算平均日志大小和总日志大小
		avgLogSize := "0 KB"
		logSize := "0 KB"

		topic := KafkaTopic{
			Name:        topicMetadata.Name,
			Partitions:  len(topicMetadata.Partitions),
			Replicas:    len(topicMetadata.Partitions[0].Replicas),
			AvgLogSize:  avgLogSize,
			LogSize:     logSize,
			ClusterID:   cluster.ID,
			ClusterName: cluster.Name,
		}

		topicList = append(topicList, topic)
	}

	// 分页处理
	total := int64(len(topicList))
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= int(total) {
		return &KafkaTopicListResponse{
			Total: total,
			Items: []KafkaTopic{},
		}, nil
	}

	if end > int(total) {
		end = int(total)
	}

	return &KafkaTopicListResponse{
		Total: total,
		Items: topicList[start:end],
	}, nil
}

// ConsumeMessages 消费消息
func (s *KafkaService) ConsumeMessages(clusterID uint, topic string, partition int, offset int64, count int, keyFilter, valueFilter string) ([]KafkaMessage, error) {
	// 获取集群信息
	cluster, err := s.GetKafkaCluster(clusterID)
	if err != nil {
		return nil, err
	}

	// 创建消费者配置
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// 设置安全配置
	if cluster.SecurityProtocol != "" {
		switch cluster.SecurityProtocol {
		case "SASL_PLAINTEXT":
			config.Net.SASL.Enable = true
			config.Net.SASL.Mechanism = sarama.SASLMechanism(cluster.SaslMechanism)
			config.Net.SASL.User = cluster.Username
			config.Net.SASL.Password = cluster.Password
		case "SASL_SSL":
			config.Net.SASL.Enable = true
			config.Net.SASL.Mechanism = sarama.SASLMechanism(cluster.SaslMechanism)
			config.Net.SASL.User = cluster.Username
			config.Net.SASL.Password = cluster.Password
			config.Net.TLS.Enable = true
		}
	}

	// 连接到 Kafka 集群
	brokers := strings.Split(cluster.BrokerServers, ",")
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}
	defer consumer.Close()

	// 创建分区消费者
	partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), offset)
	if err != nil {
		return nil, err
	}
	defer partitionConsumer.Close()

	// 消费消息
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var messages []KafkaMessage
	messageCount := 0

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			// 过滤消息
			if keyFilter != "" && string(msg.Key) != keyFilter {
				continue
			}

			if valueFilter != "" && !strings.Contains(string(msg.Value), valueFilter) {
				continue
			}

			// 解析消息
			message := KafkaMessage{
				Partition: int(msg.Partition),
				Key:       string(msg.Key),
				Offset:    msg.Offset,
				Timestamp: msg.Timestamp,
				Value:     string(msg.Value),
			}

			// 尝试解析 JSON 格式的消息
			var jsonData map[string]interface{}
			if err := json.Unmarshal(msg.Value, &jsonData); err == nil {
				if msgType, ok := jsonData["msgType"].(string); ok {
					message.MsgType = msgType
				}

				if msgData, ok := jsonData["msgData"]; ok {
					msgDataBytes, _ := json.Marshal(msgData)
					message.MsgData = string(msgDataBytes)
				}
			}

			messages = append(messages, message)
			messageCount++

			if messageCount >= count {
				return messages, nil
			}

		case err := <-partitionConsumer.Errors():
			return nil, err

		case <-ctx.Done():
			if len(messages) == 0 {
				return nil, errors.New("超时，未能获取到消息")
			}
			return messages, nil
		}
	}
}

// CreateTopic 创建 Topic
func (s *KafkaService) CreateTopic(clusterID uint, name string, partitions, replicas int) error {
	// 获取集群信息
	cluster, err := s.GetKafkaCluster(clusterID)
	if err != nil {
		return err
	}

	// 连接到 Kafka 集群
	admin, err := s.getKafkaAdminClient(cluster)
	if err != nil {
		return err
	}
	defer admin.Close()

	// 创建 Topic
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     int32(partitions),
		ReplicationFactor: int16(replicas),
	}

	return admin.CreateTopic(name, topicDetail, false)
}

// DeleteTopic 删除 Topic
func (s *KafkaService) DeleteTopic(clusterID uint, name string) error {
	// 获取集群信息
	cluster, err := s.GetKafkaCluster(clusterID)
	if err != nil {
		return err
	}

	// 连接到 Kafka 集群
	admin, err := s.getKafkaAdminClient(cluster)
	if err != nil {
		return err
	}
	defer admin.Close()

	// 删除 Topic
	return admin.DeleteTopic(name)
}

// AlterTopic 修改 Topic
func (s *KafkaService) AlterTopic(clusterID uint, name string, partitions int) error {
	// 获取集群信息
	cluster, err := s.GetKafkaCluster(clusterID)
	if err != nil {
		return err
	}

	// 连接到 Kafka 集群
	admin, err := s.getKafkaAdminClient(cluster)
	if err != nil {
		return err
	}
	defer admin.Close()

	// 修改 Topic 分区数
	return admin.CreatePartitions(name, int32(partitions), nil, false)
}

// GetTopicDetails 获取 Topic 详情
func (s *KafkaService) GetTopicDetails(clusterID uint, name string) (*KafkaTopic, error) {
	// 获取集群信息
	cluster, err := s.GetKafkaCluster(clusterID)
	if err != nil {
		return nil, err
	}

	// 连接到 Kafka 集群
	admin, err := s.getKafkaAdminClient(cluster)
	if err != nil {
		return nil, err
	}
	defer admin.Close()

	// 获取 Topic 详情
	metadata, err := admin.DescribeTopics([]string{name})
	if err != nil {
		return nil, err
	}

	if len(metadata) == 0 {
		return nil, fmt.Errorf("topic %s not found", name)
	}

	topicMetadata := metadata[0]

	// 计算平均日志大小和总日志大小
	avgLogSize := "0 KB"
	logSize := "0 KB"

	topic := &KafkaTopic{
		Name:        topicMetadata.Name,
		Partitions:  len(topicMetadata.Partitions),
		Replicas:    len(topicMetadata.Partitions[0].Replicas),
		AvgLogSize:  avgLogSize,
		LogSize:     logSize,
		ClusterID:   cluster.ID,
		ClusterName: cluster.Name,
	}

	return topic, nil
}

// GetTopicPartitions 获取 Topic 分区信息
func (s *KafkaService) GetTopicPartitions(clusterID uint, name string) ([]int32, error) {
	// 获取集群信息
	cluster, err := s.GetKafkaCluster(clusterID)
	if err != nil {
		return nil, err
	}

	// 连接到 Kafka 集群
	admin, err := s.getKafkaAdminClient(cluster)
	if err != nil {
		return nil, err
	}
	defer admin.Close()

	// 获取 Topic 详情
	metadata, err := admin.DescribeTopics([]string{name})
	if err != nil {
		return nil, err
	}

	if len(metadata) == 0 {
		return nil, fmt.Errorf("topic %s not found", name)
	}

	topicMetadata := metadata[0]

	// 创建一个分区数组，大小为分区总数
	numPartitions := len(topicMetadata.Partitions)
	partitions := make([]int32, numPartitions)

	// 填充所有分区ID
	for i := 0; i < numPartitions; i++ {
		partitions[i] = int32(i)
	}

	return partitions, nil
}

// GetPartitionOffsets 获取分区偏移量
func (s *KafkaService) GetPartitionOffsets(clusterID uint, topic string, partition int32) (int64, int64, error) {
	// 获取集群信息
	cluster, err := s.GetKafkaCluster(clusterID)
	if err != nil {
		return 0, 0, err
	}

	// 创建客户端配置
	config := sarama.NewConfig()

	// 设置安全配置
	if cluster.SecurityProtocol != "" {
		switch cluster.SecurityProtocol {
		case "SASL_PLAINTEXT":
			config.Net.SASL.Enable = true
			config.Net.SASL.Mechanism = sarama.SASLMechanism(cluster.SaslMechanism)
			config.Net.SASL.User = cluster.Username
			config.Net.SASL.Password = cluster.Password
		case "SASL_SSL":
			config.Net.SASL.Enable = true
			config.Net.SASL.Mechanism = sarama.SASLMechanism(cluster.SaslMechanism)
			config.Net.SASL.User = cluster.Username
			config.Net.SASL.Password = cluster.Password
			config.Net.TLS.Enable = true
		}
	}

	// 连接到 Kafka 集群
	brokers := strings.Split(cluster.BrokerServers, ",")
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return 0, 0, err
	}
	defer client.Close()

	// 获取最早偏移量
	oldest, err := client.GetOffset(topic, partition, sarama.OffsetOldest)
	if err != nil {
		return 0, 0, err
	}

	// 获取最新偏移量
	newest, err := client.GetOffset(topic, partition, sarama.OffsetNewest)
	if err != nil {
		return 0, 0, err
	}

	return oldest, newest, nil
}
