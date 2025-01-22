package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// CustomLogger 自定义日志中间件
func CustomLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 自定义日志格式
		return fmt.Sprintf("[Go] %s | %s | %d | %s | %s | %s | %s | %s\n",
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			param.Method,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Path,
			param.Request.Proto,
			param.ErrorMessage,
		)
	})
}

// CustomRecovery 自定义恢复中间件
func CustomRecovery() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, err interface{}) {
		// 添加 [Go] 前缀到错误日志
		fmt.Printf("[Go] ERROR: %v\n", err)
		c.AbortWithStatus(500)
	})
}
