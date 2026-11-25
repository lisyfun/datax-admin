package v1

import (
	"datax-admin/controllers"
	"datax-admin/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRedisRoutes 注册 Redis 管理相关路由
func RegisterRedisRoutes(authenticated *gin.RouterGroup) {
	c := controllers.NewRedisController()

	read := middleware.RequirePermission("tools.redis.query", "tools.redis")
	write := middleware.RequirePermission("tools.redis.set")
	del := middleware.RequirePermission("tools.redis.del")
	ttl := middleware.RequirePermission("tools.redis.ttl")
	conn := middleware.RequirePermission("tools.redis.conn.manage")

	// 连接管理
	authenticated.GET("/redis/connections", read, c.ListConnections)
	authenticated.POST("/redis/connections", conn, c.CreateConnection)
	authenticated.PUT("/redis/connections/:id", conn, c.UpdateConnection)
	authenticated.DELETE("/redis/connections/:id", conn, c.DeleteConnection)
	authenticated.POST("/redis/connections/:id/test", conn, c.TestConnection)
	authenticated.POST("/redis/connections/:id/select", read, c.SelectDB)

	// 键操作
	authenticated.GET("/redis/keys", read, c.ListKeys)
	authenticated.GET("/redis/keys/:key", read, c.GetKey)
	authenticated.GET("/redis/key", read, c.GetKey)
	authenticated.POST("/redis/key", read, c.GetKeyPost)
	authenticated.POST("/redis/keys", write, c.SetKey)
	authenticated.DELETE("/redis/keys/:key", del, c.DeleteKey)
	authenticated.POST("/redis/key/delete", del, c.DeleteKeyPost)
	authenticated.GET("/redis/keys/:key/ttl", read, c.GetTTL)
	authenticated.PUT("/redis/keys/:key/ttl", ttl, c.SetTTL)
	authenticated.POST("/redis/key/ttl", read, c.GetTTLPost)
	authenticated.POST("/redis/key/ttl/set", ttl, c.SetTTLPost)
	authenticated.POST("/redis/keys/rename", write, c.RenameKey)
	authenticated.POST("/redis/keys/copy", write, c.CopyKey)
	authenticated.POST("/redis/keys/move", write, c.MoveKey)
	authenticated.GET("/redis/export", read, c.ExportKeys)
	authenticated.GET("/redis/keys/count", read, c.CountKeys)
	authenticated.POST("/redis/cli", write, c.ExecuteCLI)
	authenticated.POST("/redis/cli/bulk", write, c.ExecuteCLIBulk)
}
