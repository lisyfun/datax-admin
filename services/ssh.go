package services

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SSHClient SSH客户端结构
type SSHClient struct {
	client        *ssh.Client
	session       *ssh.Session
	sftpClient    *sftp.Client
	stdin         io.WriteCloser
	stdout        io.Reader
	stderr        io.Reader
	closeOnce     sync.Once
	password      string // 添加密码字段
	authType      string // 添加认证类型字段
	keyContent    string // 添加密钥内容字段
	keyPassphrase string // 添加密钥密码字段
}

// NewSSHClient 创建新的SSH客户端（密码认证）
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
		client:        client,
		session:       session,
		sftpClient:    sftpClient,
		stdin:         stdin,
		stdout:        stdout,
		stderr:        stderr,
		password:      password,   // 保存密码
		authType:      "password", // 密码认证
		keyContent:    "",
		keyPassphrase: "",
	}, nil
}

// NewSSHClientWithKey 创建新的SSH客户端（密钥认证）
func NewSSHClientWithKey(host string, port int, username, keyContent, keyPassphrase string) (*SSHClient, error) {
	// 解析私钥
	var signer ssh.Signer
	var err error

	if keyPassphrase != "" {
		// 带密码的私钥
		signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(keyContent), []byte(keyPassphrase))
	} else {
		// 不带密码的私钥
		signer, err = ssh.ParsePrivateKey([]byte(keyContent))
	}

	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %v", err)
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
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
		client:        client,
		session:       session,
		sftpClient:    sftpClient,
		stdin:         stdin,
		stdout:        stdout,
		stderr:        stderr,
		password:      "",            // 密钥认证不需要密码
		authType:      "key",         // 密钥认证
		keyContent:    keyContent,    // 保存密钥内容
		keyPassphrase: keyPassphrase, // 保存密钥密码
	}, nil
}

// NewSSHClientWithAuth 根据认证类型创建SSH客户端
func NewSSHClientWithAuth(host string, port int, username, authType, password, keyContent, keyPassphrase string) (*SSHClient, error) {
	if authType == "key" {
		return NewSSHClientWithKey(host, port, username, keyContent, keyPassphrase)
	} else {
		return NewSSHClient(host, port, username, password)
	}
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
	// 注意：这里不再使用defer删除临时文件，因为异步上传还需要使用

	// 保存文件到临时目录
	buf := make([]byte, 1024*1024) // 1MB缓冲区
	var total int64

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

		// 每传输10MB记录一次进度（静默处理，不输出日志）
		if total%(10*1024*1024) == 0 {
			// 进度检查点，但不输出日志
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
			os.Remove(tempPath)
			return
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			os.Remove(tempPath)
			return
		}

		// 获取临时文件的绝对路径
		absTempPath, err := filepath.Abs(tempPath)
		if err != nil {
			os.Remove(tempPath)
			return
		}

		// 根据认证类型构建不同的scp命令
		var cmd *exec.Cmd
		if c.authType == "key" {
			// 密钥认证：创建临时密钥文件
			keyFile, err := os.CreateTemp("", "ssh-key-*.tmp")
			if err != nil {
				os.Remove(tempPath)
				return
			}
			keyPath := keyFile.Name()

			// 写入密钥内容
			if _, err := keyFile.WriteString(c.keyContent); err != nil {
				keyFile.Close()
				os.Remove(tempPath)
				os.Remove(keyPath)
				return
			}
			keyFile.Close()

			// 设置密钥文件权限为600
			if err := os.Chmod(keyPath, 0600); err != nil {
				os.Remove(tempPath)
				os.Remove(keyPath)
				return
			}

			// 使用密钥文件的scp命令
			scpCmd := fmt.Sprintf("scp -i %s -o StrictHostKeyChecking=no -P %d %s \"%s@%s:%s\"", keyPath, port, absTempPath, c.client.User(), host, destPath)
			cmd = exec.Command("bash", "-c", scpCmd)

			// 执行命令后删除密钥文件
			defer func() {
				os.Remove(keyPath)
			}()
		} else {
			// 密码认证：使用sshpass
			scpCmd := fmt.Sprintf("sshpass -p '%s' scp -o StrictHostKeyChecking=no -P %d %s \"%s@%s:%s\"", c.password, port, absTempPath, c.client.User(), host, destPath)
			cmd = exec.Command("bash", "-c", scpCmd)
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			os.Remove(tempPath)
			return
		}

		// 上传成功后删除临时文件
		os.Remove(tempPath)
	}()

	// 本地保存完成后立即返回成功
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

	// 创建缓冲区
	buf := make([]byte, 32*1024) // 32KB缓冲区
	var total int64

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

		// 每传输10MB记录一次进度（静默处理，不输出日志）
		if total%(10*1024*1024) == 0 {
			// 进度检查点，但不输出日志
		}
	}

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
	// 如果目录路径为空，使用当前目录
	if dirPath == "" {
		dirPath = "."
	}

	// 首先尝试使用SSH命令获取完整的文件列表（包括隐藏文件）
	fileList, err := c.getFileListBySSH(dirPath)
	if err == nil && len(fileList) > 0 {
		return fileList, nil
	}

	// 如果SSH命令失败，回退到SFTP方式

	if c.sftpClient == nil {
		return nil, fmt.Errorf("SFTP客户端未创建")
	}

	// 获取目录下的文件列表
	entries, err := c.sftpClient.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %v", err)
	}

	// 转换为FileInfo结构
	var sftpFileList []FileInfo
	for _, entry := range entries {
		fileInfo := FileInfo{
			Name:    entry.Name(),
			Path:    filepath.Join(dirPath, entry.Name()),
			Size:    entry.Size(),
			Mode:    entry.Mode(),
			ModTime: entry.ModTime(),
			IsDir:   entry.IsDir(),
		}
		sftpFileList = append(sftpFileList, fileInfo)
	}

	return sftpFileList, nil
}

// DownloadFileToLocal 先将远程文件下载到本地临时文件，返回本地路径和文件大小
func (c *SSHClient) DownloadFileToLocal(remotePath string) (localPath string, size int64, err error) {
	if c.sftpClient == nil {
		return "", 0, fmt.Errorf("SFTP客户端未创建")
	}

	// 打开远程文件
	srcFile, err := c.sftpClient.Open(remotePath)
	if err != nil {
		return "", 0, fmt.Errorf("打开远程文件失败: %v", err)
	}
	defer srcFile.Close()

	// 获取文件信息
	fileInfo, err := srcFile.Stat()
	if err != nil {
		return "", 0, fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 创建本地临时文件
	tmpFile, err := os.CreateTemp("", "download-*.tmp")
	if err != nil {
		return "", 0, fmt.Errorf("创建本地临时文件失败: %v", err)
	}
	localPath = tmpFile.Name()
	defer func() {
		tmpFile.Close()
	}()

	// 拷贝内容
	buf := make([]byte, 32*1024)
	total := int64(0)
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			os.Remove(localPath)
			return "", 0, fmt.Errorf("读取远程文件失败: %v", err)
		}
		if n == 0 {
			break
		}
		if _, err := tmpFile.Write(buf[:n]); err != nil {
			os.Remove(localPath)
			return "", 0, fmt.Errorf("写入本地临时文件失败: %v", err)
		}
		total += int64(n)
	}

	// 确保写入磁盘
	if err := tmpFile.Sync(); err != nil {
		os.Remove(localPath)
		return "", 0, fmt.Errorf("同步本地临时文件失败: %v", err)
	}

	return localPath, fileInfo.Size(), nil
}

// getFileListBySSH 通过SSH命令获取文件列表（包括隐藏文件）
func (c *SSHClient) getFileListBySSH(dirPath string) ([]FileInfo, error) {
	if c.client == nil {
		return nil, fmt.Errorf("SSH客户端未创建")
	}

	// 创建SSH会话
	session, err := c.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	// 使用ls -la命令获取详细的文件列表（包括隐藏文件）
	cmd := fmt.Sprintf("ls -la '%s' 2>/dev/null", dirPath)
	output, err := session.Output(cmd)
	if err != nil {
		return nil, fmt.Errorf("执行ls命令失败: %v", err)
	}

	// 解析ls命令的输出
	return c.parseListOutput(string(output), dirPath)
}

// parseListOutput 解析ls -la命令的输出
func (c *SSHClient) parseListOutput(output, dirPath string) ([]FileInfo, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var fileList []FileInfo

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 跳过总计行和当前目录/父目录
		if strings.HasPrefix(line, "total ") ||
			strings.HasSuffix(line, " .") ||
			strings.HasSuffix(line, " ..") {
			continue
		}

		// 解析ls -la的输出格式
		// 格式: drwxr-xr-x 2 user group 4096 Jan 1 12:00 filename
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		// 获取文件名（可能包含空格，所以需要特殊处理）
		nameStartIndex := 8
		fileName := strings.Join(fields[nameStartIndex:], " ")

		// 跳过 . 和 .. 目录
		if fileName == "." || fileName == ".." {
			continue
		}

		// 解析文件大小
		size, _ := strconv.ParseInt(fields[4], 10, 64)

		// 解析修改时间
		timeStr := strings.Join(fields[5:8], " ")
		modTime, _ := time.Parse("Jan 2 15:04", timeStr)
		if modTime.Year() == 0 {
			modTime = modTime.AddDate(time.Now().Year(), 0, 0)
		}

		// 判断是否为目录
		isDir := strings.HasPrefix(fields[0], "d")

		// 解析文件权限（暂时使用默认值）
		mode := os.FileMode(0644)
		if isDir {
			mode = os.FileMode(0755) | os.ModeDir
		}

		fileInfo := FileInfo{
			Name:    fileName,
			Path:    filepath.Join(dirPath, fileName),
			Size:    size,
			Mode:    mode,
			ModTime: modTime,
			IsDir:   isDir,
		}
		fileList = append(fileList, fileInfo)
	}

	return fileList, nil
}
