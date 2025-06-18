package controllers

import (
	"datax-admin/models"
	"datax-admin/services"
	"datax-admin/utils/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// KafkaController Kafka 控制器
type KafkaController struct {
	kafkaService *services.KafkaService
}

// NewKafkaController 创建 Kafka 控制器
func NewKafkaController() *KafkaController {
	return &KafkaController{
		kafkaService: &services.KafkaService{},
	}
}

// ListClusters 获取集群列表
func (c *KafkaController) ListClusters(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	search := ctx.Query("search")

	// 设置缓存控制头，允许客户端缓存结果1分钟
	ctx.Header("Cache-Control", "public, max-age=60")

	clusters, err := c.kafkaService.ListKafkaClusters(page, pageSize, search)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
		"data":    clusters,
	})
}

// GetCluster 获取集群详情
func (c *KafkaController) GetCluster(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的集群ID",
		})
		return
	}

	// 设置缓存控制头，允许客户端缓存结果1分钟
	ctx.Header("Cache-Control", "public, max-age=60")

	cluster, err := c.kafkaService.GetKafkaCluster(uint(id))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
		"data":    cluster,
	})
}

// CreateCluster 创建集群
func (c *KafkaController) CreateCluster(ctx *gin.Context) {
	var cluster models.KafkaCluster
	if err := ctx.ShouldBindJSON(&cluster); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "参数错误：" + err.Error(),
		})
		return
	}

	if err := c.kafkaService.CreateKafkaCluster(&cluster); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
		"data":    cluster,
	})
}

// UpdateCluster 更新集群
func (c *KafkaController) UpdateCluster(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的集群ID",
		})
		return
	}

	var cluster models.KafkaCluster
	if err := ctx.ShouldBindJSON(&cluster); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "参数错误：" + err.Error(),
		})
		return
	}

	cluster.ID = uint(id)
	if err := c.kafkaService.UpdateKafkaCluster(&cluster); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
		"data":    cluster,
	})
}

// DeleteCluster 删除集群
func (c *KafkaController) DeleteCluster(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的集群ID",
		})
		return
	}

	if err := c.kafkaService.DeleteKafkaCluster(uint(id)); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
	})
}

// ListTopics 获取主题列表
func (c *KafkaController) ListTopics(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的集群ID",
		})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	search := ctx.Query("search")

	// 设置缓存控制头，允许客户端缓存结果2分钟
	// 如果有搜索条件，则不缓存结果
	if search == "" {
		ctx.Header("Cache-Control", "public, max-age=120")
	} else {
		ctx.Header("Cache-Control", "no-cache")
	}

	topics, err := c.kafkaService.ListTopics(uint(id), page, pageSize, search)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
		"data":    topics,
	})
}

// GetTopicDetails 获取主题详情
func (c *KafkaController) GetTopicDetails(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的集群ID",
		})
		return
	}

	topicName := ctx.Param("topicName")
	if topicName == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "主题名称不能为空",
		})
		return
	}

	// 设置缓存控制头，允许客户端缓存结果5分钟
	ctx.Header("Cache-Control", "public, max-age=300")

	topic, err := c.kafkaService.GetTopicDetails(uint(id), topicName)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
		"data":    topic,
	})
}

// CreateTopic 创建主题
func (c *KafkaController) CreateTopic(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的集群ID",
		})
		return
	}

	var params struct {
		Name       string `json:"name" binding:"required"`
		Partitions int    `json:"partitions" binding:"required,min=1"`
		Replicas   int    `json:"replicas" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "参数错误：" + err.Error(),
		})
		return
	}

	if err := c.kafkaService.CreateTopic(uint(id), params.Name, params.Partitions, params.Replicas); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
	})
}

// DeleteTopic 删除主题
func (c *KafkaController) DeleteTopic(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的集群ID",
		})
		return
	}

	topicName := ctx.Param("topicName")
	if topicName == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "主题名称不能为空",
		})
		return
	}

	if err := c.kafkaService.DeleteTopic(uint(id), topicName); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
	})
}

// AlterTopic 修改主题
func (c *KafkaController) AlterTopic(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的集群ID",
		})
		return
	}

	topicName := ctx.Param("topicName")
	if topicName == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "主题名称不能为空",
		})
		return
	}

	var params struct {
		Partitions int `json:"partitions" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "参数错误：" + err.Error(),
		})
		return
	}

	if err := c.kafkaService.AlterTopic(uint(id), topicName, params.Partitions); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
	})
}

// ConsumeMessages 消费消息
func (c *KafkaController) ConsumeMessages(ctx *gin.Context) {
	startTime := time.Now()

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的集群ID",
		})
		return
	}

	topicName := ctx.Param("topicName")
	if topicName == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "主题名称不能为空",
		})
		return
	}

	// 获取并解析请求参数
	partition, _ := strconv.Atoi(ctx.DefaultQuery("partition", "0"))

	// 明确处理offset参数
	offsetParam := ctx.DefaultQuery("offset", "-1")
	var offset int64
	if offsetParam == "" || offsetParam == "-1" {
		// 如果没有提供offset或为-1，则表示使用默认值
		// 在这种情况下，我们将在服务层确定是使用最早还是最新偏移量
		offset = -1
		logger.Info("未指定具体偏移量，使用默认值: %d", offset)
	} else {
		// 解析指定的偏移量
		offset, err = strconv.ParseInt(offsetParam, 10, 64)
		if err != nil {
			logger.Info("偏移量参数格式错误: %s, 将使用默认值 -1", offsetParam)
			offset = -1
		} else {
			logger.Info("使用指定的偏移量: %d", offset)
		}
	}

	// 解析其他参数
	count, _ := strconv.Atoi(ctx.DefaultQuery("count", "10"))
	if count > 100 {
		count = 100 // 限制最大消息数量
		logger.Info("请求的消息数量超过限制，已调整为最大值: %d", count)
	}

	keyFilter := ctx.Query("keyFilter")
	valueFilter := ctx.Query("valueFilter")
	groupId := ctx.Query("groupId") // 获取消费者组ID，目前未使用

	logger.Info("接收到消费消息请求: 集群=%d, 主题=%s, 分区=%d, 偏移量=%d, 数量=%d, 键过滤=%s, 值过滤=%s, 消费组=%s",
		id, topicName, partition, offset, count, keyFilter, valueFilter, groupId)

	// 调用服务层方法获取消息
	messages, err := c.kafkaService.ConsumeMessages(uint(id), topicName, partition, offset, count, keyFilter, valueFilter)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	// 记录性能日志
	elapsed := time.Since(startTime)
	logger.Info("消费消息完成: 集群=%d, 主题=%s, 分区=%d, 偏移量=%d, 数量=%d, 获取到 %d 条消息, 耗时: %v",
		id, topicName, partition, offset, count, len(messages), elapsed)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
		"data":    messages,
	})
}

// GetTopicPartitions 获取主题分区信息
func (c *KafkaController) GetTopicPartitions(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的集群ID",
		})
		return
	}

	topicName := ctx.Param("topicName")
	if topicName == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "主题名称不能为空",
		})
		return
	}

	partitions, err := c.kafkaService.GetTopicPartitions(uint(id), topicName)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
		"data":    partitions,
	})
}

// GetPartitionOffsets 获取分区偏移量
func (c *KafkaController) GetPartitionOffsets(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的集群ID",
		})
		return
	}

	topicName := ctx.Param("topicName")
	if topicName == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "主题名称不能为空",
		})
		return
	}

	partition, err := strconv.ParseInt(ctx.Param("partition"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的分区ID",
		})
		return
	}

	oldest, newest, err := c.kafkaService.GetPartitionOffsets(uint(id), topicName, int32(partition))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
		"data": gin.H{
			"oldest": oldest,
			"newest": newest,
		},
	})
}

// GetTopicInfo 获取主题信息（起始偏移量、结束偏移量、消息数量）
func (c *KafkaController) GetTopicInfo(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的集群ID",
		})
		return
	}

	topicName := ctx.Param("topicName")
	if topicName == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "主题名称不能为空",
		})
		return
	}

	// 获取主题信息
	beginningOffset, endOffset, size, err := c.kafkaService.GetTopicInfo(uint(id), topicName)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
		"data": gin.H{
			"beginningOffset": beginningOffset,
			"endOffset":       endOffset,
			"size":            size,
		},
	})
}

// GetPartitionOffset 获取特定分区特定类型的偏移量
func (c *KafkaController) GetPartitionOffset(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的集群ID",
		})
		return
	}

	topicName := ctx.Param("topicName")
	if topicName == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "主题名称不能为空",
		})
		return
	}

	partition, err := strconv.ParseInt(ctx.Param("partition"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的分区ID",
		})
		return
	}

	// 获取偏移量类型（earliest或latest）
	offsetType := ctx.Query("type")
	if offsetType != "earliest" && offsetType != "latest" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "无效的偏移量类型，必须是 earliest 或 latest",
		})
		return
	}

	// 获取偏移量
	offset, err := c.kafkaService.GetPartitionOffset(uint(id), topicName, int32(partition), offsetType)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "操作成功",
		"data":    offset,
	})
}
