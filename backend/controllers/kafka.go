package controllers

import (
	"datax-admin/models"
	"datax-admin/services"
	"net/http"
	"strconv"

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
	// 移除未使用的变量
	// startTime := time.Now()

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

	partition, _ := strconv.Atoi(ctx.DefaultQuery("partition", "0"))
	offset, _ := strconv.ParseInt(ctx.DefaultQuery("offset", "-1"), 10, 64)
	count, _ := strconv.Atoi(ctx.DefaultQuery("count", "100"))
	keyFilter := ctx.Query("keyFilter")
	valueFilter := ctx.Query("valueFilter")
	// 移除未使用的变量
	// groupId := ctx.Query("groupId") // 获取消费者组ID

	// 减少日志输出，只在调试模式下输出
	// fmt.Printf("消费消息请求参数: clusterId=%d, topic=%s, partition=%d, offset=%d, count=%d, keyFilter=%s, valueFilter=%s, groupId=%s\n",
	//	id, topicName, partition, offset, count, keyFilter, valueFilter, groupId)

	messages, err := c.kafkaService.ConsumeMessages(uint(id), topicName, partition, offset, count, keyFilter, valueFilter)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	// 记录性能日志，只在调试模式下输出
	// elapsed := time.Since(startTime)
	// fmt.Printf("消费消息完成: clusterId=%d, topic=%s, partition=%d, offset=%d, count=%d, 获取到 %d 条消息, 耗时: %v\n",
	//	id, topicName, partition, offset, count, len(messages), elapsed)

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
