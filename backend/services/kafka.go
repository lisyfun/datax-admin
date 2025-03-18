package services

import (
	"context"
	"crypto/tls"
	"datax-admin/models"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
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

// TopicMetadata 主题元数据
type TopicMetadata struct {
	Topic      string
	Partitions []kafka.Partition
}

// TopicConfig 主题配置
type TopicConfig struct {
	Topic             string
	NumPartitions     int32
	ReplicationFactor int16
}

// ValidateClusterConnection 验证集群连接
func (s *KafkaService) ValidateClusterConnection(cluster *models.KafkaCluster) error {
	// 使用 kafka-go 验证连接
	conn, err := s.getKafkaConn(cluster, "")
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}
	defer conn.Close()

	// 尝试获取 broker 控制器来验证连接
	_, err = conn.Controller()
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

	// 检查集群状态
	if err := s.CheckClusterStatus(&cluster); err != nil {
		// 减少日志输出
		return &cluster, nil
	}

	// 如果集群不可用，直接返回
	if !cluster.Status {
		return &cluster, nil
	}

	// 连接到 Kafka 集群
	conn, err := s.getKafkaConn(&cluster, "")
	if err != nil {
		// 减少日志输出
		return &cluster, nil
	}
	defer conn.Close()

	// 获取 Topic 数量
	topics, err := conn.ReadPartitions()
	if err == nil {
		// 统计唯一的 topic 名称
		topicMap := make(map[string]bool)
		for _, p := range topics {
			topicMap[p.Topic] = true
		}
		cluster.TopicCount = len(topicMap)
	}

	// 获取消费者组数量
	// 由于 kafka-go 没有直接获取消费者组的 API，我们设置一个默认值
	cluster.ConsumerGroupCount = 1

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

	// 直接获取集群信息，不包括状态和统计数据
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&clusters).Error; err != nil {
		return nil, err
	}

	// 计算 Broker 数量
	for i := range clusters {
		clusters[i].BrokerCount = len(strings.Split(clusters[i].BrokerServers, ","))
		// 不在列表页面检查集群状态，避免频繁连接
		// 只设置基本信息
	}

	return &KafkaClusterListResponse{
		Total: total,
		Items: clusters,
	}, nil
}

// CheckClusterStatus 检查集群状态
func (s *KafkaService) CheckClusterStatus(cluster *models.KafkaCluster) error {
	// 更新最后检查时间
	cluster.LastCheckTime = time.Now()

	// 尝试连接集群
	conn, err := s.getKafkaConn(cluster, "")
	if err != nil {
		cluster.Status = false
		return models.DB.Save(cluster).Error
	}
	defer conn.Close()

	// 尝试获取控制器来验证连接
	_, err = conn.Controller()
	if err != nil {
		cluster.Status = false
		return models.DB.Save(cluster).Error
	}

	// 连接成功，更新状态
	cluster.Status = true
	return models.DB.Save(cluster).Error
}

// StartClusterHealthCheck 启动集群健康检查
func (s *KafkaService) StartClusterHealthCheck(checkInterval time.Duration) {
	// 增加检查间隔，减少频繁连接
	if checkInterval < 10*time.Minute {
		checkInterval = 10 * time.Minute // 将最小间隔增加到10分钟
	}

	fmt.Printf("Kafka集群健康检查已启动，检查间隔: %v\n", checkInterval)

	ticker := time.NewTicker(checkInterval)
	go func() {
		for range ticker.C {
			var clusters []models.KafkaCluster
			if err := models.DB.Find(&clusters).Error; err != nil {
				fmt.Printf("获取Kafka集群列表失败: %v\n", err)
				continue
			}

			// 限制并发检查数量
			semaphore := make(chan struct{}, 2) // 最多2个并发检查

			for _, cluster := range clusters {
				clusterCopy := cluster // 创建副本避免闭包问题

				semaphore <- struct{}{} // 获取信号量
				go func() {
					defer func() {
						<-semaphore // 释放信号量
					}()

					// 使用超时上下文控制检查时间
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()

					// 在独立的goroutine中执行检查，避免阻塞主goroutine
					done := make(chan error, 1)
					go func() {
						err := s.CheckClusterStatus(&clusterCopy)
						done <- err
					}()

					// 等待检查完成或超时
					select {
					case err := <-done:
						if err != nil {
							fmt.Printf("检查Kafka集群 [%s] 状态失败: %v\n", clusterCopy.Name, err)
						}
					case <-ctx.Done():
						fmt.Printf("检查Kafka集群 [%s] 状态超时\n", clusterCopy.Name)
					}
				}()
			}
		}
	}()
}

// enrichClusterStats 丰富集群统计信息 - 仅在需要详细信息时使用
func (s *KafkaService) enrichClusterStats(cluster *models.KafkaCluster) {
	// 检查集群状态
	if err := s.CheckClusterStatus(cluster); err != nil {
		// 减少日志输出，只在调试模式下输出
		// fmt.Printf("检查集群状态失败: %v\n", err)
		return
	}

	// 如果集群不可用，直接返回
	if !cluster.Status {
		return
	}

	// 连接到 Kafka 集群
	conn, err := s.getKafkaConn(cluster, "")
	if err != nil {
		// 减少日志输出，只在调试模式下输出
		// fmt.Printf("获取 Kafka 连接失败: %v\n", err)
		return
	}
	defer conn.Close()

	// 获取 Topic 数量
	topics, err := conn.ReadPartitions()
	if err == nil {
		// 统计唯一的 topic 名称
		topicMap := make(map[string]bool)
		for _, p := range topics {
			topicMap[p.Topic] = true
		}
		cluster.TopicCount = len(topicMap)
		// 减少日志输出，只在调试模式下输出
		// fmt.Printf("获取到主题数量: %d\n", cluster.TopicCount)
	} else {
		// 减少日志输出，只在调试模式下输出
		// fmt.Printf("获取主题列表失败: %v\n", err)
	}

	// 获取消费者组数量
	// 由于 kafka-go 没有直接获取消费者组的 API，我们设置一个默认值
	cluster.ConsumerGroupCount = 1
	// 减少日志输出，只在调试模式下输出
	// fmt.Printf("设置默认消费者组数量: %d\n", cluster.ConsumerGroupCount)
}

// getKafkaConn 获取 Kafka 连接
func (s *KafkaService) getKafkaConn(cluster *models.KafkaCluster, topic string) (*kafka.Conn, error) {
	// 解析 broker 地址
	brokers := strings.Split(cluster.BrokerServers, ",")
	if len(brokers) == 0 {
		return nil, fmt.Errorf("未配置 Kafka broker 地址")
	}

	// 使用第一个 broker 建立连接
	broker := brokers[0]

	// 设置连接超时（减少为3秒）
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 创建 TCP 连接
	dialer := &kafka.Dialer{
		Timeout:   2 * time.Second, // 减少连接超时
		DualStack: false,           // 禁用 IPv6
	}

	// 添加认证信息
	if cluster.SecurityProtocol != "" {
		switch cluster.SecurityProtocol {
		case "SASL_PLAINTEXT":
			mechanism := plain.Mechanism{
				Username: cluster.Username,
				Password: cluster.Password,
			}
			dialer.SASLMechanism = mechanism
		case "SASL_SSL":
			mechanism := plain.Mechanism{
				Username: cluster.Username,
				Password: cluster.Password,
			}
			dialer.SASLMechanism = mechanism
			dialer.TLS = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
	}

	// 创建连接
	conn, err := dialer.DialContext(ctx, "tcp", broker)
	if err != nil {
		return nil, fmt.Errorf("连接 Kafka 失败: %v", err)
	}

	return conn, nil
}

// ListTopics 获取 Topic 列表
func (s *KafkaService) ListTopics(clusterID uint, page, pageSize int, search string) (*KafkaTopicListResponse, error) {
	// 获取集群信息
	cluster, err := s.GetKafkaCluster(clusterID)
	if err != nil {
		return nil, err
	}

	// 连接到 Kafka 集群
	conn, err := s.getKafkaConn(cluster, "")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 获取所有分区
	partitions, err := conn.ReadPartitions()
	if err != nil {
		return nil, err
	}

	// 将搜索关键词转为小写，提高搜索效率
	searchLower := strings.ToLower(search)

	// 按主题分组并过滤
	topicMap := make(map[string][]kafka.Partition)
	for _, p := range partitions {
		// 如果有搜索条件，则进行不区分大小写的模糊匹配
		if searchLower != "" && !strings.Contains(strings.ToLower(p.Topic), searchLower) {
			continue
		}
		topicMap[p.Topic] = append(topicMap[p.Topic], p)
	}

	// 构建主题列表
	var topicList []KafkaTopic

	// 创建一个通道来接收处理完成的主题
	resultChan := make(chan KafkaTopic, len(topicMap))
	// 创建一个错误通道
	errChan := make(chan error, 1)
	// 使用计数器跟踪正在进行的goroutine数量
	semaphore := make(chan struct{}, 5) // 最多5个并发请求，避免过多连接

	// 并发处理每个主题
	for name, parts := range topicMap {
		// 限制并发数
		semaphore <- struct{}{}

		go func(name string, parts []kafka.Partition) {
			defer func() { <-semaphore }() // 释放信号量

			// 计算副本数
			replicas := 0
			if len(parts) > 0 {
				replicas = len(parts[0].Replicas)
			}

			// 获取主题的日志大小信息，使用超时控制
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			// 创建一个完成通道
			doneChan := make(chan struct{})
			var logSize int64 = 0
			var avgLogSize int64 = 0

			// 在后台获取日志大小
			go func() {
				defer close(doneChan)
				// 获取主题的日志大小信息
				_, _, size, err := s.GetTopicInfo(clusterID, name)
				if err == nil && size > 0 {
					logSize = size
					// 计算平均日志大小
					if len(parts) > 0 {
						avgLogSize = logSize / int64(len(parts))
					}
				}
			}()

			// 等待日志大小获取完成或超时
			select {
			case <-doneChan:
				// 日志大小获取成功
			case <-ctx.Done():
				// 超时，使用默认值
			}

			// 创建主题对象
			topic := KafkaTopic{
				Name:        name,
				Partitions:  len(parts),
				Replicas:    replicas,
				AvgLogSize:  formatBytes(avgLogSize),
				LogSize:     formatBytes(logSize),
				ClusterID:   cluster.ID,
				ClusterName: cluster.Name,
			}

			resultChan <- topic
		}(name, parts)
	}

	// 收集结果
	for i := 0; i < len(topicMap); i++ {
		select {
		case topic := <-resultChan:
			topicList = append(topicList, topic)
		case err := <-errChan:
			return nil, err
		}
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

	// 创建 Kafka 读取器
	brokers := strings.Split(cluster.BrokerServers, ",")
	// 减少日志输出，只在调试模式下输出
	// fmt.Printf("消费消息，连接 Kafka brokers: %v\n", brokers)

	// 创建 dialer
	dialer := &kafka.Dialer{
		Timeout:   1 * time.Second, // 减少连接超时时间
		DualStack: false,           // 禁用 IPv6
	}

	// 添加认证信息
	if cluster.SecurityProtocol != "" {
		switch cluster.SecurityProtocol {
		case "SASL_PLAINTEXT":
			mechanism := plain.Mechanism{
				Username: cluster.Username,
				Password: cluster.Password,
			}
			dialer.SASLMechanism = mechanism
		case "SASL_SSL":
			mechanism := plain.Mechanism{
				Username: cluster.Username,
				Password: cluster.Password,
			}
			dialer.SASLMechanism = mechanism
			dialer.TLS = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
	}

	// 创建读取器
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		Partition:   partition,
		Dialer:      dialer,
		MinBytes:    1,                      // 最小字节数设为1，不等待数据累积
		MaxBytes:    10e6,                   // 10MB
		MaxWait:     100 * time.Millisecond, // 最大等待时间设为100毫秒
		StartOffset: offset,
	})
	defer reader.Close()

	// 设置上下文超时，减少为2秒
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 创建一个计时器，用于在没有更多消息时提前退出
	noMessageTimer := time.NewTimer(500 * time.Millisecond)
	defer noMessageTimer.Stop()

	// 读取消息
	messages := make([]KafkaMessage, 0, count)
	messageChan := make(chan kafka.Message)
	errorChan := make(chan error)

	// 在后台读取消息
	go func() {
		for i := 0; i < count; i++ {
			// 设置每条消息的读取超时
			msgCtx, msgCancel := context.WithTimeout(ctx, 500*time.Millisecond)
			msg, err := reader.ReadMessage(msgCtx)
			msgCancel()

			if err != nil {
				if err != context.DeadlineExceeded {
					errorChan <- err
				}
				return
			}
			messageChan <- msg
			// 重置无消息计时器
			if !noMessageTimer.Stop() {
				select {
				case <-noMessageTimer.C:
				default:
				}
			}
			noMessageTimer.Reset(500 * time.Millisecond)
		}
	}()

	// 主循环
	for len(messages) < count {
		select {
		case msg := <-messageChan:
			// 应用过滤器
			if keyFilter != "" && !strings.Contains(string(msg.Key), keyFilter) {
				continue
			}
			if valueFilter != "" && !strings.Contains(string(msg.Value), valueFilter) {
				continue
			}

			// 创建消息对象
			message := KafkaMessage{
				Partition: msg.Partition,
				Key:       string(msg.Key),
				Offset:    msg.Offset,
				Timestamp: msg.Time,
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

		case err := <-errorChan:
			// 如果有错误但已经读取了一些消息，则返回已读取的消息
			if len(messages) > 0 {
				return messages, nil
			}
			return nil, err

		case <-noMessageTimer.C:
			// 如果一段时间内没有新消息，提前返回
			if len(messages) > 0 {
				return messages, nil
			}
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

// CreateTopic 创建 Topic
func (s *KafkaService) CreateTopic(clusterID uint, name string, partitions, replicas int) error {
	// 获取集群信息
	cluster, err := s.GetKafkaCluster(clusterID)
	if err != nil {
		return err
	}

	// 连接到 Kafka 集群
	conn, err := s.getKafkaConn(cluster, "")
	if err != nil {
		return err
	}
	defer conn.Close()

	// 创建主题
	return conn.CreateTopics(kafka.TopicConfig{
		Topic:             name,
		NumPartitions:     partitions,
		ReplicationFactor: replicas,
	})
}

// DeleteTopic 删除 Topic
func (s *KafkaService) DeleteTopic(clusterID uint, name string) error {
	// 获取集群信息
	cluster, err := s.GetKafkaCluster(clusterID)
	if err != nil {
		return err
	}

	// 连接到 Kafka 集群
	conn, err := s.getKafkaConn(cluster, "")
	if err != nil {
		return err
	}
	defer conn.Close()

	// 删除主题
	return conn.DeleteTopics(name)
}

// AlterTopic 修改 Topic
func (s *KafkaService) AlterTopic(clusterID uint, name string, partitions int) error {
	// 获取集群信息
	cluster, err := s.GetKafkaCluster(clusterID)
	if err != nil {
		return err
	}

	// 连接到 Kafka 集群
	conn, err := s.getKafkaConn(cluster, "")
	if err != nil {
		return err
	}
	defer conn.Close()

	// 获取控制器
	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	// 连接到控制器
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	// 由于 kafka-go 没有直接的 CreatePartitions 方法，我们使用管理命令
	// 这里我们只能返回一个不支持的错误
	return fmt.Errorf("修改分区功能在 kafka-go 中不直接支持")
}

// GetTopicDetails 获取 Topic 详情
func (s *KafkaService) GetTopicDetails(clusterID uint, name string) (*KafkaTopic, error) {
	// 获取集群信息
	cluster, err := s.GetKafkaCluster(clusterID)
	if err != nil {
		return nil, err
	}

	// 连接到 Kafka 集群
	conn, err := s.getKafkaConn(cluster, "")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 获取主题分区
	partitions, err := conn.ReadPartitions(name)
	if err != nil {
		return nil, err
	}

	if len(partitions) == 0 {
		return nil, fmt.Errorf("topic %s not found", name)
	}

	// 计算副本数
	replicas := len(partitions[0].Replicas)

	// 获取主题的日志大小信息，使用超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 创建一个完成通道
	doneChan := make(chan struct{})
	var logSize int64 = 0
	var avgLogSize int64 = 0

	// 在后台获取日志大小
	go func() {
		defer close(doneChan)
		// 获取主题的日志大小信息
		_, _, size, err := s.GetTopicInfo(clusterID, name)
		if err == nil && size > 0 {
			logSize = size
			// 计算平均日志大小
			if len(partitions) > 0 {
				avgLogSize = logSize / int64(len(partitions))
			}
		}
	}()

	// 等待日志大小获取完成或超时
	select {
	case <-doneChan:
		// 日志大小获取成功
	case <-ctx.Done():
		// 超时，使用默认值
	}

	topic := &KafkaTopic{
		Name:        name,
		Partitions:  len(partitions),
		Replicas:    replicas,
		AvgLogSize:  formatBytes(avgLogSize),
		LogSize:     formatBytes(logSize),
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
	conn, err := s.getKafkaConn(cluster, "")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 获取主题分区
	partitions, err := conn.ReadPartitions(name)
	if err != nil {
		return nil, err
	}

	if len(partitions) == 0 {
		return nil, fmt.Errorf("topic %s not found", name)
	}

	// 创建分区 ID 数组
	result := make([]int32, len(partitions))
	for i, p := range partitions {
		result[i] = int32(p.ID)
	}

	return result, nil
}

// GetPartitionOffsets 获取分区偏移量
func (s *KafkaService) GetPartitionOffsets(clusterID uint, topicName string, partition int32) (int64, int64, error) {
	// 获取集群信息
	cluster, err := s.GetKafkaCluster(clusterID)
	if err != nil {
		return 0, 0, fmt.Errorf("获取集群信息失败: %v", err)
	}

	// 创建 Kafka 读取器配置
	brokers := strings.Split(cluster.BrokerServers, ",")
	if len(brokers) == 0 {
		return 0, 0, fmt.Errorf("未配置 Kafka broker 地址")
	}

	// 确保broker地址包含主机名和端口号
	brokerAddr := brokers[0]
	if !strings.Contains(brokerAddr, ":") {
		return 0, 0, fmt.Errorf("broker地址格式错误，应为host:port格式")
	}

	// 减少日志输出，只在调试模式下输出
	// fmt.Printf("获取分区偏移量，连接 Kafka broker: %s, topic: %s, partition: %d\n",
	//	brokerAddr, topicName, partition)

	// 创建 dialer
	dialer := &kafka.Dialer{
		Timeout:   3 * time.Second, // 增加连接超时时间
		DualStack: false,           // 禁用 IPv6
	}

	// 添加认证信息
	if cluster.SecurityProtocol != "" {
		switch cluster.SecurityProtocol {
		case "SASL_PLAINTEXT":
			mechanism := plain.Mechanism{
				Username: cluster.Username,
				Password: cluster.Password,
			}
			dialer.SASLMechanism = mechanism
		case "SASL_SSL":
			mechanism := plain.Mechanism{
				Username: cluster.Username,
				Password: cluster.Password,
			}
			dialer.SASLMechanism = mechanism
			dialer.TLS = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用LookupPartition获取分区信息
	partitionInfo, err := dialer.LookupPartition(ctx, "tcp", brokerAddr, topicName, int(partition))
	if err != nil {
		return 0, 0, fmt.Errorf("查找分区信息失败: %v", err)
	}

	// 连接到特定分区
	conn, err := dialer.DialPartition(ctx, "tcp", brokerAddr, partitionInfo)
	if err != nil {
		return 0, 0, fmt.Errorf("连接分区失败: %v", err)
	}
	defer conn.Close()

	// 获取最早偏移量
	oldest, err := conn.ReadFirstOffset()
	if err != nil {
		return 0, 0, fmt.Errorf("获取最早偏移量失败: %v", err)
	}

	// 获取最新偏移量
	newest, err := conn.ReadLastOffset()
	if err != nil {
		return 0, 0, fmt.Errorf("获取最新偏移量失败: %v", err)
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

	// 获取broker地址
	brokers := strings.Split(cluster.BrokerServers, ",")
	if len(brokers) == 0 {
		return 0, 0, 0, fmt.Errorf("未配置 Kafka broker 地址")
	}

	// 确保broker地址包含主机名和端口号
	brokerAddr := brokers[0]
	if !strings.Contains(brokerAddr, ":") {
		return 0, 0, 0, fmt.Errorf("broker地址格式错误，应为host:port格式")
	}

	// 创建 dialer
	dialer := &kafka.Dialer{
		Timeout:   1 * time.Second, // 减少连接超时时间
		DualStack: false,           // 禁用 IPv6
	}

	// 添加认证信息
	if cluster.SecurityProtocol != "" {
		switch cluster.SecurityProtocol {
		case "SASL_PLAINTEXT":
			mechanism := plain.Mechanism{
				Username: cluster.Username,
				Password: cluster.Password,
			}
			dialer.SASLMechanism = mechanism
		case "SASL_SSL":
			mechanism := plain.Mechanism{
				Username: cluster.Username,
				Password: cluster.Password,
			}
			dialer.SASLMechanism = mechanism
			dialer.TLS = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
	}

	// 创建上下文，设置较短的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 使用LookupPartitions获取所有分区信息
	partitions, err := dialer.LookupPartitions(ctx, "tcp", brokerAddr, topicName)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("查找分区信息失败: %v", err)
	}

	var beginningOffset, endOffset, size int64

	// 限制处理的分区数量，避免处理太多分区导致超时
	maxPartitions := 5
	if len(partitions) > maxPartitions {
		partitions = partitions[:maxPartitions]
	}

	// 创建一个通道来接收处理完成的结果
	type partitionResult struct {
		oldest int64
		newest int64
	}
	resultChan := make(chan partitionResult, len(partitions))

	// 使用WaitGroup等待所有goroutine完成
	var wg sync.WaitGroup

	// 并发处理每个分区
	for _, p := range partitions {
		wg.Add(1)
		go func(p kafka.Partition) {
			defer wg.Done()

			// 创建分区连接
			connCtx, connCancel := context.WithTimeout(ctx, 1*time.Second)
			defer connCancel()

			conn, err := dialer.DialPartition(connCtx, "tcp", brokerAddr, p)
			if err != nil {
				return
			}
			defer conn.Close()

			// 获取最早偏移量
			oldest, err := conn.ReadFirstOffset()
			if err != nil {
				return
			}

			// 获取最新偏移量
			newest, err := conn.ReadLastOffset()
			if err != nil {
				return
			}

			// 发送结果
			resultChan <- partitionResult{oldest: oldest, newest: newest}
		}(p)
	}

	// 启动一个goroutine等待所有分区处理完成并关闭结果通道
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	for result := range resultChan {
		if result.oldest < beginningOffset || beginningOffset == 0 {
			beginningOffset = result.oldest
		}
		if result.newest > endOffset {
			endOffset = result.newest
		}
		size += result.newest - result.oldest
	}

	return beginningOffset, endOffset, size, nil
}

// GetPartitionOffset 获取特定分区特定类型的偏移量
func (s *KafkaService) GetPartitionOffset(clusterID uint, topicName string, partition int32, offsetType string) (int64, error) {
	// 获取集群信息
	cluster, err := s.GetKafkaCluster(clusterID)
	if err != nil {
		return 0, fmt.Errorf("获取集群信息失败: %v", err)
	}

	// 创建 Kafka 读取器配置
	brokers := strings.Split(cluster.BrokerServers, ",")
	if len(brokers) == 0 {
		return 0, fmt.Errorf("未配置 Kafka broker 地址")
	}

	// 确保broker地址包含主机名和端口号
	brokerAddr := brokers[0]
	if !strings.Contains(brokerAddr, ":") {
		return 0, fmt.Errorf("broker地址格式错误，应为host:port格式")
	}

	// 减少日志输出，只在调试模式下输出
	// fmt.Printf("获取偏移量，连接 Kafka broker: %s, topic: %s, partition: %d, type: %s\n",
	//	brokerAddr, topicName, partition, offsetType)

	// 创建 dialer
	dialer := &kafka.Dialer{
		Timeout:   3 * time.Second, // 增加连接超时时间
		DualStack: false,           // 禁用 IPv6
	}

	// 添加认证信息
	if cluster.SecurityProtocol != "" {
		switch cluster.SecurityProtocol {
		case "SASL_PLAINTEXT":
			mechanism := plain.Mechanism{
				Username: cluster.Username,
				Password: cluster.Password,
			}
			dialer.SASLMechanism = mechanism
		case "SASL_SSL":
			mechanism := plain.Mechanism{
				Username: cluster.Username,
				Password: cluster.Password,
			}
			dialer.SASLMechanism = mechanism
			dialer.TLS = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用LookupPartition获取分区信息
	partitionInfo, err := dialer.LookupPartition(ctx, "tcp", brokerAddr, topicName, int(partition))
	if err != nil {
		return 0, fmt.Errorf("查找分区信息失败: %v", err)
	}

	// 连接到特定分区
	conn, err := dialer.DialPartition(ctx, "tcp", brokerAddr, partitionInfo)
	if err != nil {
		return 0, fmt.Errorf("连接分区失败: %v", err)
	}
	defer conn.Close()

	var offset int64
	var offsetErr error

	// 根据类型获取偏移量
	if offsetType == "earliest" {
		offset, offsetErr = conn.ReadFirstOffset()
		if offsetErr != nil {
			return 0, fmt.Errorf("获取最早偏移量失败: %v", offsetErr)
		}
		// 减少日志输出，只在调试模式下输出
		// fmt.Printf("获取到最早偏移量: %d\n", offset)
	} else if offsetType == "latest" {
		offset, offsetErr = conn.ReadLastOffset()
		if offsetErr != nil {
			return 0, fmt.Errorf("获取最新偏移量失败: %v", offsetErr)
		}
		// 减少日志输出，只在调试模式下输出
		// fmt.Printf("获取到最新偏移量: %d\n", offset)
	} else {
		return 0, fmt.Errorf("无效的偏移量类型: %s", offsetType)
	}

	return offset, nil
}

// formatBytes 将字节数格式化为人类可读的字符串
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
