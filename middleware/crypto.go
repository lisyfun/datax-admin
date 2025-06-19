package middleware

import (
	"bytes"
	"datax-admin/config"
	"datax-admin/utils/crypto"
	"datax-admin/utils/logger"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var cryptoManager *crypto.CryptoManager

// InitCrypto 初始化加密管理器
func InitCrypto() error {
	cryptoConfig := &crypto.CryptoConfig{
		SecretKey: config.GlobalConfig.Auth.Secret, // 使用配置中的密钥
		EnableLog: true,
	}

	var err error
	cryptoManager, err = crypto.NewCryptoManager(cryptoConfig)
	if err != nil {
		return err
	}

	logger.Info("加密中间件初始化成功")
	return nil
}

// CryptoMiddleware 加密中间件
func CryptoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查加密管理器是否已初始化
		if cryptoManager == nil {
			c.Next()
			return
		}

		// 跳过不需要加密的路径
		if shouldSkipEncryption(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 处理请求解密
		if err := handleRequestDecryption(c); err != nil {
			logger.Error("请求解密失败: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误"})
			c.Abort()
			return
		}

		// 创建响应写入器来拦截响应
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
			statusCode:     0,
		}
		c.Writer = writer

		c.Next()

		// 处理响应加密
		if err := handleResponseEncryption(c, writer); err != nil {
			logger.Error("响应加密失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		}
	}
}

// shouldSkipEncryption 判断是否跳过加密
func shouldSkipEncryption(path string) bool {
	// 跳过的路径列表
	skipPaths := []string{
		"/api/v1/health", // 健康检查
		"/api/v1/ping",   // ping接口
		"/static/",       // 静态资源
		"/assets/",       // 资源文件
		"/ws/",           // WebSocket连接
		"/favicon.ico",   // 网站图标
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	return false
}

// handleRequestDecryption 处理请求解密
func handleRequestDecryption(c *gin.Context) error {
	// 只处理POST、PUT、PATCH请求
	if c.Request.Method != "POST" && c.Request.Method != "PUT" && c.Request.Method != "PATCH" {
		return nil
	}

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body.Close()

	// 如果请求体为空，直接返回
	if len(body) == 0 {
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		return nil
	}

	// 尝试解析加密请求
	var encryptedRequest struct {
		Data      string `json:"data"`
		Timestamp int64  `json:"timestamp"`
		Signature string `json:"signature"`
	}

	if err := json.Unmarshal(body, &encryptedRequest); err != nil {
		// 如果不是加密格式，直接使用原始数据
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		return nil
	}

	// 检查是否包含加密字段
	if encryptedRequest.Data == "" || encryptedRequest.Signature == "" {
		// 不是加密格式，直接使用原始数据
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		return nil
	}

	// 检查加密管理器
	if cryptoManager == nil {
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		return nil
	}

	// 解密数据
	decryptedData, err := cryptoManager.DecryptRequest(
		encryptedRequest.Data,
		encryptedRequest.Timestamp,
		encryptedRequest.Signature,
	)
	if err != nil {
		return err
	}

	// 替换请求体
	c.Request.Body = io.NopCloser(bytes.NewBuffer(decryptedData))
	c.Request.ContentLength = int64(len(decryptedData))

	return nil
}

// handleResponseEncryption 处理响应加密
func handleResponseEncryption(c *gin.Context, writer *responseWriter) error {
	// 检查是否需要加密响应
	if writer.body.Len() == 0 {
		// 恢复原始Writer并写入空响应
		c.Writer = writer.ResponseWriter
		return nil
	}

	// 获取原始响应数据
	originalData := writer.body.Bytes()

	// 检查客户端是否支持加密
	acceptEncryption := c.GetHeader("Accept-Encryption")
	if acceptEncryption != "true" {
		// 客户端不支持加密，返回原始数据
		c.Writer = writer.ResponseWriter
		if writer.statusCode != 0 {
			c.Writer.WriteHeader(writer.statusCode)
		}
		c.Writer.Write(originalData)
		return nil
	}

	// 尝试解析为JSON
	var responseData interface{}
	if err := json.Unmarshal(originalData, &responseData); err != nil {
		// 不是JSON格式，直接返回原始数据
		c.Writer = writer.ResponseWriter
		if writer.statusCode != 0 {
			c.Writer.WriteHeader(writer.statusCode)
		}
		c.Writer.Write(originalData)
		return nil
	}

	// 检查加密管理器
	if cryptoManager == nil {
		c.Writer = writer.ResponseWriter
		if writer.statusCode != 0 {
			c.Writer.WriteHeader(writer.statusCode)
		}
		c.Writer.Write(originalData)
		return nil
	}

	// 加密响应数据
	encryptedResponse, err := cryptoManager.EncryptResponse(responseData)
	if err != nil {
		// 加密失败，返回原始数据
		c.Writer = writer.ResponseWriter
		if writer.statusCode != 0 {
			c.Writer.WriteHeader(writer.statusCode)
		}
		c.Writer.Write(originalData)
		return nil
	}

	// 恢复原始Writer
	c.Writer = writer.ResponseWriter

	// 设置响应头
	c.Header("Content-Type", "application/json")
	c.Header("Content-Encryption", "true")

	// 如果有状态码，先写入状态码
	if writer.statusCode != 0 {
		c.Writer.WriteHeader(writer.statusCode)
	}

	// 序列化加密后的数据并写入响应
	encryptedJSON, err := json.Marshal(encryptedResponse)
	if err != nil {
		// 序列化失败，返回原始数据
		if writer.statusCode != 0 {
			c.Writer.WriteHeader(writer.statusCode)
		}
		c.Writer.Write(originalData)
		return nil
	}

	// 写入加密后的数据
	_, err = c.Writer.Write(encryptedJSON)
	return err
}

// responseWriter 自定义响应写入器
type responseWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	// 不立即写入状态码，等到最后处理
}

func (w *responseWriter) Status() int {
	if w.statusCode == 0 {
		return 200 // 默认状态码
	}
	return w.statusCode
}

// GetCryptoManager 获取加密管理器实例
func GetCryptoManager() *crypto.CryptoManager {
	return cryptoManager
}

// EncryptData 加密数据的便捷方法
func EncryptData(data interface{}) (map[string]interface{}, error) {
	if cryptoManager == nil {
		return nil, nil // 如果加密管理器未初始化，返回原始数据
	}
	return cryptoManager.EncryptResponse(data)
}

// DecryptData 解密数据的便捷方法
func DecryptData(encryptedData string, timestamp int64, signature string) ([]byte, error) {
	if cryptoManager == nil {
		return nil, nil
	}
	return cryptoManager.DecryptRequest(encryptedData, timestamp, signature)
}
