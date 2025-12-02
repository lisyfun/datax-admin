package middleware

import (
    "datax-admin/config"
    "datax-admin/models"
    "crypto/rand"
    "net/http"
    "strings"

    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie" // 添加cookie存储作为备选
    "github.com/gin-gonic/gin"
)

// InitSession 初始化Session中间件
func InitSession() gin.HandlerFunc {
    authConfig := config.GlobalConfig.Auth

    // 使用配置密钥；为空则随机生成 32 字节
    secret := authConfig.Secret
    if secret == "" {
        b := make([]byte, 32)
        _, _ = rand.Read(b)
        for i := range b {
            if b[i] == 0 { b[i] = 1 }
        }
        secret = string(b)
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
        Secure:   false,
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
                Secure:   isHTTPS(c),
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
                Secure:   isHTTPS(c),
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
                Secure:   isHTTPS(c),
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

// isHTTPS 判断当前请求是否 HTTPS 或经由代理为 HTTPS
func isHTTPS(c *gin.Context) bool {
    if c.Request.TLS != nil { return true }
    proto := c.GetHeader("X-Forwarded-Proto")
    return strings.EqualFold(proto, "https")
}
