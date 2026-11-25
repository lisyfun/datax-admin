// Package secret 提供敏感数据（密码、私钥等）的静态加密能力。
//
// 使用 AES-256-GCM，密文格式为 "enc:v1:<base64(nonce|ciphertext)>"。
// 对于没有该前缀的历史数据，Decrypt 会原样返回，从而兼容升级前写入的明文。
package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const prefix = "enc:v1:"

var (
	mu       sync.RWMutex
	aead     cipher.AEAD
	initOnce sync.Once
	initErr  error
)

// Init 初始化加密密钥。
//
// key 为空时会依次尝试：
//  1. 环境变量 SECURITY_ENCRYPTION_KEY
//  2. 工作目录下的 data/encryption.key（不存在则自动生成，权限 0600）
func Init(key string) error {
	initOnce.Do(func() {
		if key == "" {
			key = os.Getenv("SECURITY_ENCRYPTION_KEY")
		}
		if key == "" {
			var err error
			key, err = loadOrCreateKeyFile("data/encryption.key")
			if err != nil {
				initErr = err
				return
			}
		}

		sum := sha256.Sum256([]byte(key))
		block, err := aes.NewCipher(sum[:])
		if err != nil {
			initErr = err
			return
		}
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			initErr = err
			return
		}
		aead = gcm
	})
	return initErr
}

func loadOrCreateKeyFile(path string) (string, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return "", fmt.Errorf("创建密钥目录失败: %w", err)
	}

	// #nosec G304 -- 密钥文件路径由程序固定生成（data/encryption.key 或配置项）
	if data, err := os.ReadFile(path); err == nil {
		return strings.TrimSpace(string(data)), nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("读取加密密钥失败: %w", err)
	}

	buf := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", fmt.Errorf("生成加密密钥失败: %w", err)
	}
	key := base64.StdEncoding.EncodeToString(buf)

	// #nosec G304 -- 密钥文件路径由程序固定生成（data/encryption.key 或配置项）
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		// 并发启动时可能被其他进程抢先创建
		if errors.Is(err, os.ErrExist) {
			// #nosec G304 -- 同上，路径为程序自身的密钥文件
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return "", readErr
			}
			return strings.TrimSpace(string(data)), nil
		}
		return "", fmt.Errorf("写入加密密钥失败: %w", err)
	}
	defer f.Close()
	if _, err := f.WriteString(key); err != nil {
		return "", fmt.Errorf("写入加密密钥失败: %w", err)
	}
	return key, nil
}

// IsEncrypted 判断字符串是否为当前格式的密文。
func IsEncrypted(value string) bool {
	return strings.HasPrefix(value, prefix)
}

// Encrypt 加密明文。空字符串与已加密字符串会原样返回。
func Encrypt(plain string) (string, error) {
	if plain == "" || IsEncrypted(plain) {
		return plain, nil
	}

	mu.RLock()
	gcm := aead
	mu.RUnlock()
	if gcm == nil {
		return "", errors.New("加密模块未初始化")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return prefix + base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密密文。非密文（历史明文）会原样返回。
func Decrypt(value string) (string, error) {
	if value == "" || !IsEncrypted(value) {
		return value, nil
	}

	mu.RLock()
	gcm := aead
	mu.RUnlock()
	if gcm == nil {
		return "", errors.New("加密模块未初始化")
	}

	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, prefix))
	if err != nil {
		return "", fmt.Errorf("密文格式错误: %w", err)
	}
	if len(raw) < gcm.NonceSize() {
		return "", errors.New("密文长度错误")
	}

	nonce, ciphertext := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("解密失败（密钥可能已更换）: %w", err)
	}
	return string(plain), nil
}
