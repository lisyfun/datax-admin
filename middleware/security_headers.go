package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

// SecurityHeaders 设置通用安全响应头
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 基本安全头
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("Referrer-Policy", "no-referrer")
        c.Header("X-Frame-Options", "DENY")
        c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=(), fullscreen=()")
        c.Header("Cross-Origin-Opener-Policy", "same-origin")
        c.Header("Cross-Origin-Resource-Policy", "same-origin")

        // 为 API 响应设置最小 CSP（静态页 CSP 由前端网关配置）
        if strings.HasPrefix(c.Request.URL.Path, "/api/") || strings.Contains(c.Request.URL.Path, "/api/") {
            c.Header("Content-Security-Policy", "default-src 'none'; frame-ancestors 'self'; base-uri 'self'; form-action 'self'")
        }

        // 继续处理
        c.Next()

        // 确保错误响应也带有安全头
        if c.Writer.Status() >= http.StatusBadRequest {
            c.Header("X-Content-Type-Options", "nosniff")
        }
    }
}
