package services

import (
	"datax-admin/models"
	"datax-admin/utils/logger"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaTopicService 提供主题管理相关的服务
type KafkaTopicService struct {
	kafkaService *KafkaService
}

// NewKafkaTopicService 创建主题服务实例
func NewKafkaTopicService(kafkaService *KafkaService) *KafkaTopicService {
	return &KafkaTopicService{
		kafkaService: kafkaService,
	}
}

// SyncTopics 同步指定集群的主题信息
func (s *KafkaTopicService) SyncTopics(clusterID uint) error {
	// 获取集群信息
	cluster, err := s.kafkaService.GetKafkaCluster(clusterID)
	if err != nil {
		return fmt.Errorf("获取集群信息失败: %v", err)
	}

	// 连接到 Kafka 集群
	conn, err := s.kafkaService.getKafkaConn(cluster, "")
	if err != nil {
		return fmt.Errorf("连接Kafka失败: %v", err)
	}
	defer conn.Close()

	// 获取所有分区信息
	partitions, err := conn.ReadPartitions()
	if err != nil {
		return fmt.Errorf("获取分区信息失败: %v", err)
	}

	// 按主题分组
	topicMap := make(map[string][]kafka.Partition)
	for _, p := range partitions {
		topicMap[p.Topic] = append(topicMap[p.Topic], p)
	}

	// 并发处理每个主题
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // 最多5个并发
	errChan := make(chan error, len(topicMap))

	for topicName, parts := range topicMap {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量

		go func(name string, parts []kafka.Partition) {
			defer func() {
				<-semaphore // 释放信号量
				wg.Done()
			}()

			// 计算副本数
			replicas := 0
			if len(parts) > 0 {
				replicas = len(parts[0].Replicas)
			}

			// 获取主题的日志大小信息
			_, _, totalSize, err := s.kafkaService.GetTopicInfo(clusterID, name)
			if err != nil {
				errChan <- fmt.Errorf("获取主题[%s]信息失败: %v", name, err)
				return
			}

			// 计算平均日志大小
			avgSize := int64(0)
			if len(parts) > 0 {
				avgSize = totalSize / int64(len(parts))
			}

			// 更新或创建主题记录
			topic := &models.KafkaTopic{
				ClusterID:       clusterID,
				Name:            name,
				Partitions:      len(parts),
				Replicas:        replicas,
				AvgLogSize:      avgSize,
				TotalLogSize:    totalSize,
				LastRefreshTime: time.Now(),
			}

			// 使用 upsert 操作
			result := models.DB.Where("cluster_id = ? AND name = ?", clusterID, name).
				Assign(topic).
				FirstOrCreate(topic)

			if result.Error != nil {
				errChan <- fmt.Errorf("保存主题[%s]信息失败: %v", name, result.Error)
			}
		}(topicName, parts)
	}

	// 等待所有goroutine完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("同步主题信息时发生错误: %v", errors)
	}

	return nil
}

// StartTopicSyncTask 启动主题同步任务
func (s *KafkaTopicService) StartTopicSyncTask(syncInterval time.Duration) {
	if syncInterval < 5*time.Minute {
		syncInterval = 5 * time.Minute // 最小同步间隔5分钟
	}

	logger.Info("主题同步任务已启动，同步间隔: %v", syncInterval)

	ticker := time.NewTicker(syncInterval)
	go func() {
		for range ticker.C {
			// 获取所有集群
			var clusters []models.KafkaCluster
			if err := models.DB.Find(&clusters).Error; err != nil {
				logger.Info("获取集群列表失败: %v\n", err)
				continue
			}

			// 为每个集群同步主题信息
			for _, cluster := range clusters {
				if err := s.SyncTopics(cluster.ID); err != nil {
					logger.Info("同步集群[%s]主题信息失败: %v\n", cluster.Name, err)
				} else {
					logger.Info("成功同步集群[%s]主题信息\n", cluster.Name)
				}
			}
		}
	}()
}

// GetTopics 获取主题列表
func (s *KafkaTopicService) GetTopics(clusterID uint, page, pageSize int, search string) (*KafkaTopicListResponse, error) {
	// 先检查数据库中是否有该集群的主题数据
	var count int64
	if err := models.DB.Model(&models.KafkaTopic{}).Where("cluster_id = ?", clusterID).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("检查主题数据失败: %v", err)
	}

	// 如果没有数据，先同步主题信息
	if count == 0 {
		if err := s.SyncTopics(clusterID); err != nil {
			return nil, fmt.Errorf("同步主题数据失败: %v", err)
		}
	}

	var topics []models.KafkaTopic
	var total int64
	query := models.DB.Model(&models.KafkaTopic{}).Where("cluster_id = ?", clusterID)

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 获取分页数据
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&topics).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	items := make([]KafkaTopic, len(topics))
	for i, t := range topics {
		items[i] = KafkaTopic{
			Name:        t.Name,
			Partitions:  t.Partitions,
			Replicas:    t.Replicas,
			AvgLogSize:  formatBytes(t.AvgLogSize),
			LogSize:     formatBytes(t.TotalLogSize),
			ClusterID:   t.ClusterID,
			ClusterName: "", // 可以根据需要查询集群名称
		}
	}

	return &KafkaTopicListResponse{
		Total: total,
		Items: items,
	}, nil
}
