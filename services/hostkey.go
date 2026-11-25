package services

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"

	"datax-admin/config"
	"datax-admin/utils/logger"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

var knownHostsMu sync.Mutex

// hostKeyCallback 返回 SSH 主机密钥校验回调。
//
// 默认使用 known_hosts 校验；首次连接未知主机时采用 TOFU（首次信任并记录）。
// 只有显式配置 ssh.insecure_ignore_host_key=true 时才跳过校验。
func hostKeyCallback() (ssh.HostKeyCallback, error) {
	if config.GlobalConfig.SSH.InsecureIgnoreHostKey {
		logger.Info("警告：SSH 主机密钥校验已禁用（ssh.insecure_ignore_host_key=true）")
		// #nosec G106 -- 由 ssh.insecure_ignore_host_key 显式开启，默认关闭，仅限测试环境
		return ssh.InsecureIgnoreHostKey(), nil
	}

	path := config.GlobalConfig.SSH.KnownHostsPath
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("无法确定 known_hosts 路径: %w", err)
		}
		path = filepath.Join(home, ".ssh", "known_hosts")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, fmt.Errorf("创建 known_hosts 目录失败: %w", err)
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		// #nosec G304 -- 路径来自 ssh.known_hosts_path 配置或默认的 ~/.ssh/known_hosts
		f, createErr := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
		if createErr != nil {
			return nil, fmt.Errorf("创建 known_hosts 失败: %w", createErr)
		}
		_ = f.Close()
	}

	verify, err := knownhosts.New(path)
	if err != nil {
		return nil, fmt.Errorf("加载 known_hosts 失败: %w", err)
	}

	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		verifyErr := verify(hostname, remote, key)
		if verifyErr == nil {
			return nil
		}

		var keyErr *knownhosts.KeyError
		if errors.As(verifyErr, &keyErr) && len(keyErr.Want) == 0 {
			// 首次连接该主机：记录密钥后放行
			knownHostsMu.Lock()
			defer knownHostsMu.Unlock()

			// #nosec G304 -- 同上，known_hosts 路径由服务端配置决定
			f, openErr := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if openErr != nil {
				return fmt.Errorf("记录主机密钥失败: %w", openErr)
			}
			defer f.Close()

			line := knownhosts.Line([]string{knownhosts.Normalize(hostname)}, key)
			if _, writeErr := f.WriteString(line + "\n"); writeErr != nil {
				return fmt.Errorf("记录主机密钥失败: %w", writeErr)
			}
			logger.Info("首次连接主机 %s，已记录主机密钥", hostname)
			return nil
		}

		return fmt.Errorf("SSH 主机密钥校验失败（可能存在中间人攻击，请检查 %s）: %w", path, verifyErr)
	}, nil
}
