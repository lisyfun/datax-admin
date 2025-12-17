package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

// CryptoConfig 加密配置
type CryptoConfig struct {
	SecretKey string // 32字节的密钥
	EnableLog bool   // 是否启用日志
}

// CryptoManager 加密管理器
type CryptoManager struct {
	config *CryptoConfig
	block  cipher.Block
}

// NewCryptoManager 创建加密管理器
func NewCryptoManager(config *CryptoConfig) (*CryptoManager, error) {
	if len(config.SecretKey) != 32 {
		return nil, errors.New("密钥长度必须为32字节")
	}

	block, err := aes.NewCipher([]byte(config.SecretKey))
	if err != nil {
		return nil, fmt.Errorf("创建AES密码器失败: %v", err)
	}

	return &CryptoManager{
		config: config,
		block:  block,
	}, nil
}

// Encrypt 加密数据 - 使用AES-CBC模式
func (cm *CryptoManager) Encrypt(plaintext []byte) (string, error) {
	// 生成随机IV (16字节)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("生成IV失败: %v", err)
	}

	// 填充数据到块大小的倍数
	paddedData := pkcs7Pad(plaintext, aes.BlockSize)

	// 创建CBC模式加密器
	mode := cipher.NewCBCEncrypter(cm.block, iv)

	// 加密数据
	ciphertext := make([]byte, len(paddedData))
	mode.CryptBlocks(ciphertext, paddedData)

	// 组合IV和密文
	result := append(iv, ciphertext...)

	// 返回base64编码的结果
	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt 解密数据 - 使用AES-CBC模式
func (cm *CryptoManager) Decrypt(ciphertext string) ([]byte, error) {
	// base64解码
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("base64解码失败: %v", err)
	}

	// 检查数据长度
	if len(data) < aes.BlockSize {
		return nil, errors.New("密文数据太短")
	}

	// 提取IV和密文
	iv := data[:aes.BlockSize]
	cipherData := data[aes.BlockSize:]

	// 检查密文长度是否为块大小的倍数
	if len(cipherData)%aes.BlockSize != 0 {
		return nil, errors.New("密文长度不正确")
	}

	// 创建CBC模式解密器
	mode := cipher.NewCBCDecrypter(cm.block, iv)

	// 解密数据
	plaintext := make([]byte, len(cipherData))
	mode.CryptBlocks(plaintext, cipherData)

	// 去除填充
	unpaddedData, err := pkcs7Unpad(plaintext)
	if err != nil {
		return nil, fmt.Errorf("去除填充失败: %v", err)
	}

	return unpaddedData, nil
}

// EncryptString 加密字符串
func (cm *CryptoManager) EncryptString(plaintext string) (string, error) {
	return cm.Encrypt([]byte(plaintext))
}

// DecryptString 解密字符串
func (cm *CryptoManager) DecryptString(ciphertext string) (string, error) {
	data, err := cm.Decrypt(ciphertext)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GenerateSignature 生成签名
func (cm *CryptoManager) GenerateSignature(data string, timestamp int64) string {
	// 使用 HMAC-SHA256 生成签名
	message := fmt.Sprintf("%s%d", data, timestamp)
	h := hmac.New(sha256.New, []byte(cm.config.SecretKey))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifySignature 验证签名
func (cm *CryptoManager) VerifySignature(data string, timestamp int64, signature string) bool {
	// 检查时间戳是否在有效范围内（5分钟）
	now := time.Now().Unix()
	if now-timestamp > 300 || timestamp-now > 300 {
		return false
	}

	expectedSignature := cm.GenerateSignature(data, timestamp)
	return expectedSignature == signature
}

// EncryptResponse 加密响应数据
func (cm *CryptoManager) EncryptResponse(data interface{}) (map[string]interface{}, error) {
	// 将数据序列化为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("序列化数据失败: %v", err)
	}

	// 加密数据
	encryptedData, err := cm.Encrypt(jsonData)
	if err != nil {
		return nil, fmt.Errorf("加密数据失败: %v", err)
	}

	// 生成时间戳和签名
	timestamp := time.Now().Unix()
	signature := cm.GenerateSignature(encryptedData, timestamp)

	return map[string]interface{}{
		"data":      encryptedData,
		"timestamp": timestamp,
		"signature": signature,
	}, nil
}

// DecryptRequest 解密请求数据
func (cm *CryptoManager) DecryptRequest(encryptedData string, timestamp int64, signature string) ([]byte, error) {
	// 验证签名
	if !cm.VerifySignature(encryptedData, timestamp, signature) {
		return nil, errors.New("签名验证失败")
	}

	// 解密数据
	return cm.Decrypt(encryptedData)
}

// pkcs7Pad 添加PKCS7填充
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7Unpad 去除PKCS7填充
func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("数据为空")
	}

	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, errors.New("填充长度错误")
	}

	return data[:(length - unpadding)], nil
}
