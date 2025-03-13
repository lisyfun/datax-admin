package v1

import (
	"datax-admin/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterKafkaRoutes 注册 Kafka 相关路由
func RegisterKafkaRoutes(r *gin.RouterGroup) {
	kafkaController := controllers.NewKafkaController()

	// Kafka 集群管理
	kafka := r.Group("/kafka")
	{
		clusters := kafka.Group("/clusters")
		{
			// 集群基础操作
			clusters.GET("", kafkaController.ListClusters)
			clusters.POST("", kafkaController.CreateCluster)
			clusters.GET("/:id", kafkaController.GetCluster)
			clusters.PUT("/:id", kafkaController.UpdateCluster)
			clusters.DELETE("/:id", kafkaController.DeleteCluster)

			// Topic 管理
			clusters.GET("/:id/topics", kafkaController.ListTopics)
			clusters.POST("/:id/topics", kafkaController.CreateTopic)
			clusters.GET("/:id/topics/:topicName", kafkaController.GetTopicDetails)
			clusters.PUT("/:id/topics/:topicName", kafkaController.AlterTopic)
			clusters.DELETE("/:id/topics/:topicName", kafkaController.DeleteTopic)

			// 消息消费
			clusters.GET("/:id/topics/:topicName/messages", kafkaController.ConsumeMessages)
			clusters.GET("/:id/topics/:topicName/partitions", kafkaController.GetTopicPartitions)
			clusters.GET("/:id/topics/:topicName/partitions/:partition/offsets", kafkaController.GetPartitionOffsets)
		}
	}
}
