package middleware

import (
	"datax-admin/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie" // 添加cookie存储作为备选
	"github.com/gin-gonic/gin"
)

// InitSession 初始化Session中间件
func InitSession() gin.HandlerFunc {
	authConfig := config.GlobalConfig.Auth

	// 默认配置，用于配置不可用时
	defaultSecret := authConfig.Secret
	defaultMaxAge := authConfig.Expiration // 1天

	store := cookie.NewStore([]byte(defaultSecret))
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
			// 主动清除无效 session
			session.Clear()
			session.Options(sessions.Options{
				Path:     "/",
				MaxAge:   -1, // 立即过期
				HttpOnly: true,
			})
			_ = session.Save()

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
