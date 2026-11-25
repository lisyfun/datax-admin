package middleware

import (
	"crypto/rand"
	"datax-admin/config"
	"datax-admin/models"
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
			if b[i] == 0 {
				b[i] = 1
			}
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
		Secure:   authConfig.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})

	sessionMiddleware := sessions.Sessions("datax_session", store)

	return func(c *gin.Context) {
		// 服务端密钥轮换（例如未配置 AUTH_SECRET 时每次启动随机生成）后，
		// 旧 Cookie 无法通过签名校验。gin-contrib/sessions 此时会返回 nil 会话，
		// 后续 Get/Set 会触发空指针 panic，导致登录也返回 500。
		// 这里提前校验并丢弃无效 Cookie，让请求以全新会话继续。
		origCtx := c.Request.Context()
		if _, err := store.Get(c.Request, "datax_session"); err != nil {
			discardCookie(c.Request, "datax_session")
			// gorilla/sessions 会把（带错误的）解析结果缓存在 request context 中，
			// 必须丢弃该 context，否则后续读取仍会拿到 nil 会话。
			c.Request = c.Request.WithContext(origCtx)
		}
		sessionMiddleware(c)
	}
}

// discardCookie 从请求中移除指定 Cookie，同时保留其它 Cookie。
func discardCookie(r *http.Request, name string) {
	cookies := r.Cookies()
	r.Header.Del("Cookie")
	for _, cookie := range cookies {
		if cookie.Name == name {
			continue
		}
		r.AddCookie(cookie)
	}
}

// SessionAuth Session认证中间件
func SessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 安全地获取session
		session := sessions.Default(c)
		if session == nil {
			// cookie 损坏或签名不匹配（例如服务端密钥轮换后），按未登录处理
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息或会话已过期"})
			c.Abort()
			return
		}

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
	if c.Request.TLS != nil {
		return true
	}
	proto := c.GetHeader("X-Forwarded-Proto")
	return strings.EqualFold(proto, "https")
}
