package services

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SSHClient SSH客户端结构
type SSHClient struct {
	client     *ssh.Client
	session    *ssh.Session
	sftpClient *sftp.Client
	stdin      io.WriteCloser
	stdout     io.Reader
	stderr     io.Reader
	closeOnce  sync.Once
}

// NewSSHClient 创建新的SSH客户端
func NewSSHClient(host string, port int, username, password string) (*SSHClient, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("SSH连接失败: %v", err)
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("创建SSH会话失败: %v", err)
	}

	// 创建SFTP客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		session.Close()
		client.Close()
		return nil, fmt.Errorf("创建SFTP客户端失败: %v", err)
	}

	// 请求伪终端
	if err := session.RequestPty("xterm", 40, 80, ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}); err != nil {
		session.Close()
		client.Close()
		return nil, fmt.Errorf("请求伪终端失败: %v", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		session.Close()
		client.Close()
		return nil, fmt.Errorf("获取标准输入失败: %v", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		client.Close()
		return nil, fmt.Errorf("获取标准输出失败: %v", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		session.Close()
		client.Close()
		return nil, fmt.Errorf("获取标准错误失败: %v", err)
	}

	// 启动shell
	if err := session.Shell(); err != nil {
		session.Close()
		client.Close()
		return nil, fmt.Errorf("启动Shell失败: %v", err)
	}

	return &SSHClient{
		client:     client,
		session:    session,
		sftpClient: sftpClient,
		stdin:      stdin,
		stdout:     stdout,
		stderr:     stderr,
	}, nil
}

// Write 写入数据到SSH会话
func (s *SSHClient) Write(data []byte) (int, error) {
	return s.stdin.Write(data)
}

// Read 从SSH会话读取数据
func (s *SSHClient) Read(p []byte) (int, error) {
	return s.stdout.Read(p)
}

// Close 关闭SSH连接
func (s *SSHClient) Close() error {
	var err error
	s.closeOnce.Do(func() {
		if s.session != nil {
			s.session.Close()
		}
		if s.client != nil {
			err = s.client.Close()
		}
		if s.sftpClient != nil {
			err = s.sftpClient.Close()
		}
	})
	return err
}

// ResizeTerminal 调整终端大小
func (s *SSHClient) ResizeTerminal(width, height int) error {
	return s.session.WindowChange(height, width)
}

// UploadFile 上传文件
func (c *SSHClient) UploadFile(src io.Reader, destPath string) error {
	if c.sftpClient == nil {
		return fmt.Errorf("SFTP客户端未创建")
	}

	// 创建目标文件
	destFile, err := c.sftpClient.Create(destPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer destFile.Close()

	// 复制文件内容
	_, err = io.Copy(destFile, src)
	if err != nil {
		return fmt.Errorf("复制文件失败: %v", err)
	}

	return nil
}
