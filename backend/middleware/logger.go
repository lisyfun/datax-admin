package middleware

import (
	"datax-admin/utils/logger"
	"fmt"

	"github.com/gin-gonic/gin"
)

// CustomLogger 自定义日志中间件
func CustomLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 自定义日志格式
		return fmt.Sprintf("%-7s %s %s %d %s %s %s %s %s \n",
			"[Go]",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
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
		logger.Error("服务器内部错误: %v", err)
		c.AbortWithStatus(500)
	})
}
