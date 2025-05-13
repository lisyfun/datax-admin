package middleware

import (
	"datax-admin/config"
	"datax-admin/utils/logger"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie" // 添加cookie存储作为备选
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

// InitSession 初始化Session中间件
func InitSession() gin.HandlerFunc {
	redisConfig := config.GlobalConfig.Redis

	// 默认配置，用于配置不可用时
	defaultSecret := "datax_default_secret_key"
	defaultMaxAge := 86400 // 1天

	var store sessions.Store
	var err error

	// 尝试使用Redis存储
	if redisConfig.Host != "" {
		// 构建Redis连接地址
		redisAddr := fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port)

		// 使用Redis存储，从配置中读取参数
		store, err = redis.NewStore(redisConfig.MaxIdle, "tcp", redisAddr, redisConfig.Password, redisConfig.Secret)
		if err != nil {
			logger.Warn("初始化Redis Session存储失败: %v，将使用Cookie存储替代", err)
			// 如果Redis连接失败，使用Cookie存储作为备选
			store = cookie.NewStore([]byte(defaultSecret))
		} else {
			logger.Info("成功使用Redis存储会话")

			// Session Redis配置
			store.Options(sessions.Options{
				Path:     "/",
				MaxAge:   redisConfig.MaxAge,
				HttpOnly: true,
			})

			return sessions.Sessions("datax_session", store)
		}
	} else {
		// 如果没有配置Redis，使用Cookie存储
		logger.Info("未检测到Redis配置，使用Cookie存储会话")
		store = cookie.NewStore([]byte(defaultSecret))
	}

	// Cookie存储的配置
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   defaultMaxAge,
		HttpOnly: true,
	})

	return sessions.Sessions("datax_session", store)
}

// SessionAuth Session认证中间件
func SessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("userID")

		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息或会话已过期"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中，与JWT方式保持一致
		c.Set("userID", userID)
		c.Set("username", session.Get("username"))
		c.Next()
	}
}
