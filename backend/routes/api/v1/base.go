package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterBaseRoutes 注册基础路由
func RegisterBaseRoutes(public *gin.RouterGroup) {
	public.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}