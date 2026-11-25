package v1

import (
	"datax-admin/controllers"
	"datax-admin/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterKafkaRoutes 注册 Kafka 相关路由
func RegisterKafkaRoutes(r *gin.RouterGroup) {
	kafkaController := controllers.NewKafkaController()

	readPerms := []string{"tools.kafka.cluster.query", "tools.kafka.cluster", "tools.kafka.message"}
	writePerms := []string{"tools.kafka.cluster.update", "tools.kafka.cluster.create"}
	read := middleware.RequirePermission(readPerms...)
	write := middleware.RequirePermission(writePerms...)

	kafka := r.Group("/kafka")
	{
		clusters := kafka.Group("/clusters")
		{
			// 集群基础操作
			clusters.GET("", read, kafkaController.ListClusters)
			clusters.POST("", middleware.RequirePermission("tools.kafka.cluster.create"), kafkaController.CreateCluster)
			clusters.GET("/:id", read, kafkaController.GetCluster)
			clusters.PUT("/:id", middleware.RequirePermission("tools.kafka.cluster.update"), kafkaController.UpdateCluster)
			clusters.DELETE("/:id", middleware.RequirePermission("tools.kafka.cluster.delete"), kafkaController.DeleteCluster)

			// Topic 管理
			clusters.GET("/:id/topics", read, kafkaController.ListTopics)
			clusters.POST("/:id/topics", write, kafkaController.CreateTopic)
			clusters.GET("/:id/topics/:topicName", read, kafkaController.GetTopicDetails)
			clusters.PUT("/:id/topics/:topicName", write, kafkaController.AlterTopic)
			clusters.DELETE("/:id/topics/:topicName", write, kafkaController.DeleteTopic)

			// 消息消费
			clusters.GET("/:id/topics/:topicName/messages", read, kafkaController.ConsumeMessages)
			clusters.GET("/:id/topics/:topicName/partitions", read, kafkaController.GetTopicPartitions)
			clusters.GET("/:id/topics/:topicName/partitions/:partition/offsets", read, kafkaController.GetPartitionOffsets)
			clusters.GET("/:id/topics/:topicName/partitions/:partition/offset", read, kafkaController.GetPartitionOffset)
			clusters.GET("/:id/topics/:topicName/info", read, kafkaController.GetTopicInfo)
		}
	}
}
