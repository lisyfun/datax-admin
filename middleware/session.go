package middleware

import (
	"datax-admin/config"
	"datax-admin/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie" // 添加cookie存储作为备选
	"github.com/gin-gonic/gin"
)

// InitSession 初始化Session中间件
func InitSession() gin.HandlerFunc {
	authConfig := config.GlobalConfig.Auth

	// 使用固定的32字节密钥
	secret := "a266d59b88e44d618aaab77bac0f7a50"
	if authConfig.Secret != "" {
		secret = authConfig.Secret
	}

	// 确保密钥长度为32字节
	if len(secret) != 32 {
		if len(secret) > 32 {
			secret = secret[:32]
		} else {
			// 填充到32字节
			for len(secret) < 32 {
				secret += "0"
			}
		}
	}

	defaultMaxAge := authConfig.Expiration
	if defaultMaxAge == 0 {
		defaultMaxAge = 86400 // 默认1天
	}

	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   defaultMaxAge,
		HttpOnly: true,
		Secure:   false, // 开发环境设为false
		SameSite: http.SameSiteLaxMode,
	})

	return sessions.Sessions("datax_session", store)
}

// SessionAuth Session认证中间件
func SessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 安全地获取session
		session := sessions.Default(c)

		userID := session.Get("userID")

		if userID == nil {
			// 主动清除无效 session
			session.Clear()
			session.Options(sessions.Options{
				Path:     "/",
				MaxAge:   -1, // 立即过期
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})
			_ = session.Save()

			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息或会话已过期"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中，与JWT方式保持一致
		// 确保userID是uint类型
		var uid uint
		switch v := userID.(type) {
		case uint:
			uid = v
		case int:
			uid = uint(v)
		case float64:
			uid = uint(v)
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID类型错误"})
			c.Abort()
			return
		}

		// 检查用户状态是否有效
		var user models.User
		if err := models.DB.First(&user, uid).Error; err != nil {
			// 用户不存在，清除session
			session.Clear()
			session.Options(sessions.Options{
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})
			_ = session.Save()

			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			c.Abort()
			return
		}

		// 检查用户是否被禁用
		if user.Status == 0 {
			// 用户被禁用，清除session
			session.Clear()
			session.Options(sessions.Options{
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})
			_ = session.Save()

			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户已被禁用"})
			c.Abort()
			return
		}

		c.Set("userID", uid)
		c.Set("username", session.Get("username"))
		c.Next()
	}
}
