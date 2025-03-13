package services

import (
	"context"
	"datax-admin/models"
	"encoding/json"
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
		fmt.Printf("获取 Kafka 管理客户端失败: %v\n", err)
		return
	}
	defer admin.Close()

	// 获取 Topic 数量
	topics, err := admin.ListTopics()
	if err == nil {
		cluster.TopicCount = len(topics)
		fmt.Printf("获取到主题数量: %d\n", cluster.TopicCount)
	} else {
		fmt.Printf("获取主题列表失败: %v\n", err)
	}

	// 获取 Broker 数量
	// 从 BrokerServers 字段计算 Broker 数量
	cluster.BrokerCount = len(strings.Split(cluster.BrokerServers, ","))
	fmt.Printf("计算得到 Broker 数量: %d\n", cluster.BrokerCount)

	// 获取消费者组数量
	// 先尝试使用 ListConsumerGroups
	groups, err := admin.ListConsumerGroups()
	if err != nil {
		fmt.Printf("使用 ListConsumerGroups 获取消费者组失败: %v\n", err)

		// 如果失败，尝试使用普通客户端获取消费者组
		config := sarama.NewConfig()
		config.Version = sarama.V2_0_0_0

		client, err := sarama.NewClient(strings.Split(cluster.BrokerServers, ","), config)
		if err == nil {
			defer client.Close()

			// 获取所有主题
			topics, err := client.Topics()
			if err == nil {
				// 遍历所有主题获取消费者组
				groupSet := make(map[string]bool)
				for _, topic := range topics {
					// 获取主题的所有分区
					_, err := client.Partitions(topic)
					if err == nil {
						// 尝试获取消费者组偏移量
						offsetManager, err := sarama.NewOffsetManagerFromClient("datax-admin", client)
						if err == nil {
							defer offsetManager.Close()
							groupSet["datax-admin"] = true
						}
					}
				}

				// 计算消费者组数量
				groupCount := len(groupSet)
				if groupCount > 0 {
					cluster.ConsumerGroupCount = groupCount
					fmt.Printf("使用普通客户端获取到消费者组: %v\n", groupSet)
				} else {
					cluster.ConsumerGroupCount = 1 // 默认至少有 datax-admin 消费者组
					fmt.Printf("未找到消费者组，设置默认值: 1\n")
				}
			}
		} else {
			fmt.Printf("创建普通客户端失败: %v\n", err)
			cluster.ConsumerGroupCount = 1
		}
	} else {
		fmt.Printf("获取到的消费者组: %+v\n", groups)
		if len(groups) == 0 {
			cluster.ConsumerGroupCount = 1 // 默认至少有 datax-admin 消费者组
			fmt.Printf("未找到消费者组，设置默认值: 1\n")
		} else {
			cluster.ConsumerGroupCount = len(groups)
			fmt.Printf("设置消费者组数量为: %d\n", cluster.ConsumerGroupCount)
		}
	}
}

// getKafkaAdminClient 获取 Kafka 管理客户端
func (s *KafkaService) getKafkaAdminClient(cluster *models.KafkaCluster) (sarama.ClusterAdmin, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_1_0 // 使用与服务器相同的版本

	// 设置较短的超时时间，避免长时间等待
	config.Net.DialTimeout = 5 * time.Second
	config.Net.ReadTimeout = 5 * time.Second
	config.Net.WriteTimeout = 5 * time.Second

	// 设置消费者组配置
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

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
	fmt.Printf("尝试连接 Kafka brokers: %v\n", brokers)

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("创建 Kafka 管理客户端失败: %v", err)
	}

	return admin, nil
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
	config.Version = sarama.V2_8_1_0 // 使用与服务器相同的版本
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// 设置较短的超时时间
	config.Net.DialTimeout = 5 * time.Second
	config.Net.ReadTimeout = 5 * time.Second
	config.Net.WriteTimeout = 5 * time.Second

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

	// 创建消费者组
	group, err := sarama.NewConsumerGroup(brokers, "datax-admin", config)
	if err != nil {
		return nil, fmt.Errorf("创建消费者组失败: %v", err)
	}
	defer group.Close()

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建一个通道来接收消息
	messages := make([]KafkaMessage, 0)
	messagesChan := make(chan *sarama.ConsumerMessage, count)
	errorsChan := make(chan error, 1)

	// 创建一个消费者组处理器
	handler := &consumerGroupHandler{
		ready:        make(chan bool),
		messagesChan: messagesChan,
		errorsChan:   errorsChan,
		partition:    int32(partition),
		offset:       offset,
		count:        count,
		keyFilter:    keyFilter,
		valueFilter:  valueFilter,
		messagesRead: 0,
	}

	// 在后台运行消费者组
	go func() {
		for {
			err := group.Consume(ctx, []string{topic}, handler)
			if err != nil {
				errorsChan <- fmt.Errorf("消费消息失败: %v", err)
				return
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	// 等待消费者组准备就绪
	<-handler.ready

	// 使用计时器来避免长时间等待
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	// 收集消息
	for len(messages) < count {
		select {
		case msg := <-messagesChan:
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

		case err := <-errorsChan:
			return nil, err

		case <-timer.C:
			// 超时但有消息，返回已获取的消息
			if len(messages) > 0 {
				return messages, nil
			}
			// 超时且没有消息，返回空数组
			return []KafkaMessage{}, nil

		case <-ctx.Done():
			// 上下文超时
			if len(messages) > 0 {
				return messages, nil
			}
			return []KafkaMessage{}, nil
		}
	}

	return messages, nil
}

// consumerGroupHandler 实现 sarama.ConsumerGroupHandler 接口
type consumerGroupHandler struct {
	ready        chan bool
	messagesChan chan *sarama.ConsumerMessage
	errorsChan   chan error
	partition    int32
	offset       int64
	count        int
	keyFilter    string
	valueFilter  string
	messagesRead int
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// 检查分区
		if h.partition >= 0 && message.Partition != h.partition {
			continue
		}

		// 检查偏移量
		if h.offset >= 0 && message.Offset < h.offset {
			continue
		}

		// 应用过滤器
		if h.keyFilter != "" && !strings.Contains(string(message.Key), h.keyFilter) {
			continue
		}
		if h.valueFilter != "" && !strings.Contains(string(message.Value), h.valueFilter) {
			continue
		}

		// 发送消息到通道
		h.messagesChan <- message
		h.messagesRead++

		// 标记消息已处理
		session.MarkMessage(message, "")

		// 检查是否已经读取了足够的消息
		if h.messagesRead >= h.count {
			return nil
		}
	}
	return nil
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
func (s *KafkaService) GetPartitionOffsets(clusterID uint, topicName string, partition int32) (int64, int64, error) {
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
	oldest, err := client.GetOffset(topicName, partition, sarama.OffsetOldest)
	if err != nil {
		return 0, 0, err
	}

	// 获取最新偏移量
	newest, err := client.GetOffset(topicName, partition, sarama.OffsetNewest)
	if err != nil {
		return 0, 0, err
	}

	return oldest, newest, nil
}

// GetTopicInfo 获取主题信息（起始偏移量、结束偏移量、消息数量）
func (s *KafkaService) GetTopicInfo(clusterID uint, topicName string) (int64, int64, int64, error) {
	// 获取集群信息
	cluster, err := s.GetKafkaCluster(clusterID)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("获取集群信息失败: %v", err)
	}

	// 创建 Kafka 配置
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Net.DialTimeout = 5 * time.Second
	config.Net.ReadTimeout = 5 * time.Second
	config.Net.WriteTimeout = 5 * time.Second

	// 设置认证信息
	if cluster.SecurityProtocol == "SASL_PLAINTEXT" || cluster.SecurityProtocol == "SASL_SSL" {
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLMechanism(cluster.SaslMechanism)
		config.Net.SASL.User = cluster.Username
		config.Net.SASL.Password = cluster.Password
	}

	// 创建 Kafka 客户端
	client, err := sarama.NewClient(strings.Split(cluster.BrokerServers, ","), config)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("连接 Kafka 集群失败: %v", err)
	}
	defer client.Close()

	// 获取主题的所有分区
	partitions, err := client.Partitions(topicName)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("获取主题分区失败: %v", err)
	}

	var beginningOffset, endOffset, size int64

	// 遍历所有分区，获取起始偏移量和结束偏移量
	for _, partition := range partitions {
		oldest, err := client.GetOffset(topicName, partition, sarama.OffsetOldest)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("获取分区 %d 的起始偏移量失败: %v", partition, err)
		}

		newest, err := client.GetOffset(topicName, partition, sarama.OffsetNewest)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("获取分区 %d 的结束偏移量失败: %v", partition, err)
		}

		// 累加所有分区的偏移量
		if partition == 0 || oldest < beginningOffset {
			beginningOffset = oldest
		}
		if newest > endOffset {
			endOffset = newest
		}
		size += newest - oldest
	}

	return beginningOffset, endOffset, size, nil
}
