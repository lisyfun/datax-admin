package services

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	password   string // 添加密码字段
}

// NewSSHClient 创建新的SSH客户端
func NewSSHClient(host string, port int, username, password string) (*SSHClient, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
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
		password:   password, // 保存密码
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

// ProgressCallback 定义进度回调函数类型
type ProgressCallback func(current, total int64)

// UploadFileWithTemp 通过临时文件上传
func (c *SSHClient) UploadFileWithTemp(src io.Reader, destPath string, fileSize int64, progressCb ProgressCallback) error {
	if c.client == nil {
		return fmt.Errorf("SSH客户端未创建")
	}

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "upload-*.tmp")
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %v", err)
	}
	tempPath := tempFile.Name()
	fmt.Printf("创建临时文件: %s\n", tempPath)
	// 注意：这里不再使用defer删除临时文件，因为异步上传还需要使用

	// 保存文件到临时目录
	fmt.Println("开始保存文件到临时目录...")
	buf := make([]byte, 1024*1024) // 1MB缓冲区
	var total int64
	startTime := time.Now()

	for {
		n, err := src.Read(buf)
		if err != nil {
			if err != io.EOF {
				tempFile.Close()
				os.Remove(tempPath)
				return fmt.Errorf("读取文件失败: %v", err)
			}
			break
		}
		if n == 0 {
			break
		}

		if _, err := tempFile.Write(buf[:n]); err != nil {
			tempFile.Close()
			os.Remove(tempPath)
			return fmt.Errorf("写入临时文件失败: %v", err)
		}

		total += int64(n)
		if progressCb != nil {
			progressCb(total, fileSize)
		}

		// 每传输10MB记录一次进度
		if total%(10*1024*1024) == 0 {
			elapsed := time.Since(startTime)
			speed := float64(total) / elapsed.Seconds() / 1024 / 1024 // MB/s
			fmt.Printf("已保存到临时文件: %.2f MB, 总大小: %.2f MB, 速度: %.2f MB/s\n",
				float64(total)/1024/1024,
				float64(fileSize)/1024/1024,
				speed)
		}
	}

	// 确保所有数据都写入磁盘
	if err := tempFile.Sync(); err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return fmt.Errorf("同步临时文件失败: %v", err)
	}

	// 关闭临时文件
	tempFile.Close()

	// 启动异步上传
	go func() {
		// 从原始客户端获取主机地址和端口
		addr := c.client.RemoteAddr().String()
		host, portStr, err := net.SplitHostPort(addr)
		if err != nil {
			fmt.Printf("解析地址失败: %v\n", err)
			os.Remove(tempPath)
			return
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			fmt.Printf("解析端口失败: %v\n", err)
			os.Remove(tempPath)
			return
		}

		// 获取临时文件的绝对路径
		absTempPath, err := filepath.Abs(tempPath)
		if err != nil {
			fmt.Printf("获取临时文件绝对路径失败: %v\n", err)
			os.Remove(tempPath)
			return
		}
		fmt.Printf("使用临时文件绝对路径: %s\n", absTempPath)

		// 构建scp命令，使用sshpass自动输入密码
		scpCmd := fmt.Sprintf("sshpass -p '%s' scp -o StrictHostKeyChecking=no -P %d %s \"%s@%s:%s\"", c.password, port, absTempPath, c.client.User(), host, destPath)
		fmt.Printf("执行scp命令: %s\n", scpCmd)

		// 使用bash -c执行命令
		cmd := exec.Command("bash", "-c", scpCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("文件上传失败: %v\n", err)
			os.Remove(tempPath)
			return
		}

		// 上传成功后删除临时文件
		os.Remove(tempPath)
		fmt.Printf("临时文件已删除: %s\n", tempPath)

		elapsed := time.Since(startTime)
		speed := float64(total) / elapsed.Seconds() / 1024 / 1024 // MB/s
		fmt.Printf("文件上传完成: %.2f MB, 耗时: %.2f 秒, 平均速度: %.2f MB/s\n",
			float64(total)/1024/1024,
			elapsed.Seconds(),
			speed)
	}()

	// 本地保存完成后立即返回成功
	elapsed := time.Since(startTime)
	speed := float64(total) / elapsed.Seconds() / 1024 / 1024 // MB/s
	fmt.Printf("文件已保存到临时目录: %.2f MB, 耗时: %.2f 秒, 速度: %.2f MB/s\n",
		float64(total)/1024/1024,
		elapsed.Seconds(),
		speed)

	return nil
}

// UploadFile 上传文件（保持原有方法作为备选）
func (c *SSHClient) UploadFile(src io.Reader, destPath string, fileSize int64, progressCb ProgressCallback) error {
	// 默认使用新的上传方法
	return c.UploadFileWithTemp(src, destPath, fileSize, progressCb)
}

// GetFileInfo 获取文件信息
func (c *SSHClient) GetFileInfo(filePath string) (os.FileInfo, error) {
	if c.sftpClient == nil {
		return nil, fmt.Errorf("SFTP客户端未创建")
	}

	fileInfo, err := c.sftpClient.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	return fileInfo, nil
}

// DownloadFile 下载文件
func (c *SSHClient) DownloadFile(filePath string, writer io.Writer) error {
	if c.sftpClient == nil {
		return fmt.Errorf("SFTP客户端未创建")
	}

	// 打开远程文件
	srcFile, err := c.sftpClient.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开远程文件失败: %v", err)
	}
	defer srcFile.Close()

	// 获取文件大小
	fileInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 创建缓冲区
	buf := make([]byte, 32*1024) // 32KB缓冲区
	var total int64
	startTime := time.Now()

	// 读取并写入文件
	for {
		n, err := srcFile.Read(buf)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("读取文件失败: %v", err)
			}
			break
		}
		if n == 0 {
			break
		}

		if _, err := writer.Write(buf[:n]); err != nil {
			return fmt.Errorf("写入文件失败: %v", err)
		}

		total += int64(n)

		// 每传输10MB记录一次进度
		if total%(10*1024*1024) == 0 {
			elapsed := time.Since(startTime)
			speed := float64(total) / elapsed.Seconds() / 1024 / 1024 // MB/s
			fmt.Printf("已下载: %.2f MB, 总大小: %.2f MB, 速度: %.2f MB/s\n",
				float64(total)/1024/1024,
				float64(fileInfo.Size())/1024/1024,
				speed)
		}
	}

	elapsed := time.Since(startTime)
	speed := float64(total) / elapsed.Seconds() / 1024 / 1024 // MB/s
	fmt.Printf("文件下载完成: %.2f MB, 耗时: %.2f 秒, 平均速度: %.2f MB/s\n",
		float64(total)/1024/1024,
		elapsed.Seconds(),
		speed)

	return nil
}

// FileInfo 文件信息结构
type FileInfo struct {
	Name    string      `json:"name"`
	Path    string      `json:"path"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"mode"`
	ModTime time.Time   `json:"modTime"`
	IsDir   bool        `json:"isDir"`
}

// GetFileList 获取目录下的文件列表
func (c *SSHClient) GetFileList(dirPath string) ([]FileInfo, error) {
	if c.sftpClient == nil {
		return nil, fmt.Errorf("SFTP客户端未创建")
	}

	// 如果目录路径为空，使用当前目录
	if dirPath == "" {
		dirPath = "."
	}

	// 获取目录下的文件列表
	entries, err := c.sftpClient.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %v", err)
	}

	// 转换为FileInfo结构
	var fileList []FileInfo
	for _, entry := range entries {
		fileInfo := FileInfo{
			Name:    entry.Name(),
			Path:    filepath.Join(dirPath, entry.Name()),
			Size:    entry.Size(),
			Mode:    entry.Mode(),
			ModTime: entry.ModTime(),
			IsDir:   entry.IsDir(),
		}
		fileList = append(fileList, fileInfo)
	}

	return fileList, nil
}
