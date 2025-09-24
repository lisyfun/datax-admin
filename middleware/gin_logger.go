package middleware

import (
	"datax-admin/utils/logger"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// GinLogger 自定义GIN日志中间件
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		// 根据状态码选择日志级别
		logMsg := fmt.Sprintf("%3d | %8v | %s | %s | %s",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)

		// 添加响应大小信息
		if bodySize >= 0 {
			logMsg += fmt.Sprintf(" | %d bytes", bodySize)
		}

		// 根据状态码记录不同级别的日志
		switch {
		case statusCode >= 400 && statusCode < 500:
			logger.Warn(logMsg)
		case statusCode >= 500:
			logger.Error(logMsg)
		default:
			logger.Info(logMsg)
		}
	}
}

// GinLoggerWithSkipper 带跳过条件的GIN日志中间件
func GinLoggerWithSkipper(skipper func(*gin.Context) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果满足跳过条件，则不记录日志
		if skipper != nil && skipper(c) {
			c.Next()
			return
		}

		// 使用默认的日志中间件
		GinLogger()(c)
	}
}

// SkipHealthCheck 跳过健康检查接口的日志记录
func SkipHealthCheck(c *gin.Context) bool {
	return c.Request.URL.Path == "/health" || c.Request.URL.Path == "/ping"
}
