package v1

import (
	"datax-admin/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterRedisRoutes 注册 Redis 管理相关路由
func RegisterRedisRoutes(authenticated *gin.RouterGroup) {
	c := controllers.NewRedisController()

	// 连接管理
	authenticated.GET("/redis/connections", c.ListConnections)
	authenticated.POST("/redis/connections", c.CreateConnection)
	authenticated.PUT("/redis/connections/:id", c.UpdateConnection)
	authenticated.DELETE("/redis/connections/:id", c.DeleteConnection)
	authenticated.POST("/redis/connections/:id/test", c.TestConnection)
	authenticated.POST("/redis/connections/:id/select", c.SelectDB)

	// 键操作
	authenticated.GET("/redis/keys", c.ListKeys)
	authenticated.GET("/redis/keys/:key", c.GetKey)
	authenticated.GET("/redis/key", c.GetKey)
	authenticated.POST("/redis/key", c.GetKeyPost)
	authenticated.POST("/redis/keys", c.SetKey)
	authenticated.DELETE("/redis/keys/:key", c.DeleteKey)
	authenticated.POST("/redis/key/delete", c.DeleteKeyPost)
	authenticated.GET("/redis/keys/:key/ttl", c.GetTTL)
	authenticated.PUT("/redis/keys/:key/ttl", c.SetTTL)
	authenticated.POST("/redis/key/ttl", c.GetTTLPost)
	authenticated.POST("/redis/key/ttl/set", c.SetTTLPost)
	authenticated.POST("/redis/keys/rename", c.RenameKey)
	authenticated.POST("/redis/keys/copy", c.CopyKey)
	authenticated.POST("/redis/keys/move", c.MoveKey)
	authenticated.GET("/redis/export", c.ExportKeys)
	authenticated.GET("/redis/keys/count", c.CountKeys)
	authenticated.POST("/redis/cli", c.ExecuteCLI)
	authenticated.POST("/redis/cli/bulk", c.ExecuteCLIBulk)
}
