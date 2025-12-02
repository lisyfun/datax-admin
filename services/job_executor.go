package services

import (
	"bufio"
	"bytes"
	"context"
	"datax-admin/config"
	"datax-admin/models"
	"datax-admin/utils/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// executeShellJob 执行Shell任务
func (s *JobService) executeShellJob(job *models.Job, params any, history *models.JobHistory) {
	var shellParams models.JobShellParams
	if err := mapToStruct(params, &shellParams); err != nil {
		history.Status = 0
		// 将日志写入文件而不是直接存储到数据库
		if logErr := history.WriteLogToFile("", fmt.Sprintf("解析Shell参数失败: %v", err)); logErr != nil {
			logger.Info("写入日志文件失败: %v", logErr)
		}
		return
	}

	// 创建上下文，处理超时
	ctx := context.Background()
	if job.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(job.Timeout)*time.Second)
		defer cancel()
	}

	// 初始化日志文件
	if logErr := history.WriteLogToFile(fmt.Sprintf("开始执行Shell命令: %s\n", shellParams.Command), ""); logErr != nil {
		logger.Info("初始化日志文件失败: %v", logErr)
	}

	// 准备命令
	cmd := exec.CommandContext(ctx, "sh", "-c", shellParams.Command)
	if shellParams.WorkDir != "" {
		cmd.Dir = shellParams.WorkDir
	}
	if len(shellParams.Environment) > 0 {
		env := os.Environ()
		for k, v := range shellParams.Environment {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = env
	}

	// 创建管道来实时读取输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		history.Status = 0
		if logErr := history.AppendLogToFile("", fmt.Sprintf("创建输出管道失败: %v\n", err)); logErr != nil {
			logger.Info("写入日志文件失败: %v", logErr)
		}
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		history.Status = 0
		if logErr := history.AppendLogToFile("", fmt.Sprintf("创建错误管道失败: %v\n", err)); logErr != nil {
			logger.Info("写入日志文件失败: %v", logErr)
		}
		return
	}

	// 启动命令
	if err = cmd.Start(); err != nil {
		history.Status = 0
		if logErr := history.AppendLogToFile("", fmt.Sprintf("启动Shell命令失败: %v\n", err)); logErr != nil {
			logger.Info("写入日志文件失败: %v", logErr)
		}
		return
	}

	// 实时读取输出
	var outputBuffer, errorBuffer bytes.Buffer
	done := make(chan bool, 2)

	// 启动goroutine读取stdout
	go func() {
		defer func() { done <- true }()
		buffer := make([]byte, 1024)
		for {
			var n int
			n, err = stdout.Read(buffer)
			if n > 0 {
				content := string(buffer[:n])
				outputBuffer.WriteString(content)
				// 实时追加到日志文件
				if logErr := history.AppendLogToFile(content, ""); logErr != nil {
					logger.Info("追加输出日志失败: %v", logErr)
				}
			}
			if err != nil {
				break
			}
		}
	}()

	// 启动goroutine读取stderr
	go func() {
		defer func() { done <- true }()
		buffer := make([]byte, 1024)
		for {
			var n int
			n, err = stderr.Read(buffer)
			if n > 0 {
				content := string(buffer[:n])
				errorBuffer.WriteString(content)
				// 实时追加到日志文件
				if logErr := history.AppendLogToFile("", content); logErr != nil {
					logger.Info("追加错误日志失败: %v", logErr)
				}
			}
			if err != nil {
				break
			}
		}
	}()

	// 等待命令完成
	err = cmd.Wait()

	// 等待所有goroutine完成
	<-done
	<-done

	if err != nil {
		history.Status = 0
		// 如果有额外的错误信息，追加到日志
		if logErr := history.AppendLogToFile("", fmt.Sprintf("\n命令执行失败: %v\n", err)); logErr != nil {
			logger.Info("写入日志文件失败: %v", logErr)
		}
	} else {
		history.Status = 1
		if logErr := history.AppendLogToFile("\n命令执行成功\n", ""); logErr != nil {
			logger.Info("写入日志文件失败: %v", logErr)
		}
	}
}

// executeHTTPJob 执行HTTP任务
func (s *JobService) executeHTTPJob(job *models.Job, params any, history *models.JobHistory) {
	var httpParams models.JobHTTPParams
	if err := mapToStruct(params, &httpParams); err != nil {
		history.Status = 0
		history.Error = fmt.Sprintf("解析HTTP参数失败: %v", err)
		return
	}

	// 创建HTTP请求
	req, err := http.NewRequest(httpParams.Method, httpParams.URL, strings.NewReader(httpParams.Body))
	if err != nil {
		history.Status = 0
		if logErr := history.WriteLogToFile("", fmt.Sprintf("创建HTTP请求失败: %v", err)); logErr != nil {
			logger.Info("写入日志文件失败: %v", logErr)
		}
		return
	}

	// 添加请求头
	for k, v := range httpParams.Headers {
		req.Header.Set(k, v)
	}

	// 创建带超时的客户端
	client := &http.Client{
		Timeout: time.Duration(job.Timeout) * time.Second,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		history.Status = 0
		if logErr := history.WriteLogToFile("", fmt.Sprintf("发送HTTP请求失败: %v", err)); logErr != nil {
			logger.Info("写入日志文件失败: %v", logErr)
		}
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		history.Status = 0
		if logErr := history.WriteLogToFile("", fmt.Sprintf("读取HTTP响应失败: %v", err)); logErr != nil {
			logger.Info("写入日志文件失败: %v", logErr)
		}
		return
	}

	// 检查状态码
	if len(httpParams.SuccessCode) > 0 {
		isSuccess := false
		for _, code := range httpParams.SuccessCode {
			if resp.StatusCode == code {
				isSuccess = true
				break
			}
		}
		if !isSuccess {
			history.Status = 0
			if logErr := history.WriteLogToFile(string(body), fmt.Sprintf("HTTP响应状态码不符合预期: %d", resp.StatusCode)); logErr != nil {
				logger.Info("写入日志文件失败: %v", logErr)
			}
			return
		}
	}

	history.Status = 1
	if logErr := history.WriteLogToFile(string(body), ""); logErr != nil {
		logger.Info("写入日志文件失败: %v", logErr)
	}
}

// executeDataXJob 执行DataX任务
func (s *JobService) executeDataXJob(job *models.Job, params any, history *models.JobHistory) {
	var dataxParams models.JobDataXParams
	if err := mapToStruct(params, &dataxParams); err != nil {
		history.Status = 0
		if logErr := history.WriteLogToFile("", fmt.Sprintf("解析DataX参数失败: %v", err)); logErr != nil {
			logger.Info("写入日志文件失败: %v", logErr)
		}
		return
	}

	// 替换配置中的参数占位符
	jobConfig := dataxParams.JobConfig
	var paramLog strings.Builder
	paramLog.WriteString("开始执行DataX任务\n")
	if len(dataxParams.Parameters) > 0 {
		paramLog.WriteString("参数替换:\n")
		for key, value := range dataxParams.Parameters {
			placeholder := fmt.Sprintf("${%s}", key)
			jobConfig = strings.ReplaceAll(jobConfig, placeholder, value)
			paramLog.WriteString(fmt.Sprintf("  %s -> %s\n", placeholder, value))
		}
	} else {
		paramLog.WriteString("无参数需要替换\n")
	}

	// 初始化日志文件
	if logErr := history.WriteLogToFile(paramLog.String(), ""); logErr != nil {
		logger.Info("初始化日志文件失败: %v", logErr)
	}

	// 格式化配置内容
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(jobConfig), "", "    "); err != nil {
		history.Status = 0
		if logErr := history.AppendLogToFile("", fmt.Sprintf("格式化配置内容失败: %v", err)); logErr != nil {
			logger.Info("追加日志文件失败: %v", logErr)
		}
		return
	}

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		history.Status = 0
		if logErr := history.WriteLogToFile("", fmt.Sprintf("获取当前工作目录失败: %v", err)); logErr != nil {
			logger.Info("写入日志文件失败: %v", logErr)
		}
		return
	}

	// 在当前目录下创建临时文件
	tmpFile, err := os.CreateTemp(currentDir, "datax-*.json")
	if err != nil {
		history.Status = 0
		if logErr := history.WriteLogToFile("", fmt.Sprintf("创建临时文件失败: %v", err)); logErr != nil {
			logger.Info("写入日志文件失败: %v", logErr)
		}
		return
	}
	tmpFileName := tmpFile.Name()

	// 延迟删除临时文件，但要等到命令执行完成后
	defer func() {
		if err = os.Remove(tmpFileName); err != nil {
			logger.Info("删除临时文件失败: %v", err)
		}
	}()

	// 写入JSON配置（使用替换后的配置）
	if _, err = tmpFile.WriteString(jobConfig); err != nil {
		history.Status = 0
		if logErr := history.WriteLogToFile("", fmt.Sprintf("写入配置失败: %v", err)); logErr != nil {
			logger.Info("写入日志文件失败: %v", logErr)
		}
		tmpFile.Close()
		return
	}

	// 确保数据写入磁盘并关闭文件
	if err = tmpFile.Sync(); err != nil {
		history.Status = 0
		if logErr := history.WriteLogToFile("", fmt.Sprintf("同步文件失败: %v", err)); logErr != nil {
			logger.Info("写入日志文件失败: %v", logErr)
		}
		tmpFile.Close()
		return
	}
	tmpFile.Close()

	// 创建上下文，处理超时
	ctx := context.Background()
	if job.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(job.Timeout)*time.Second)
		defer cancel()
	}

	// 准备命令
	cmd := exec.CommandContext(ctx, config.GlobalConfig.DataX.Bin, "-job", tmpFileName)

	// 设置工作目录为 DataX home 目录
	cmd.Dir = config.GlobalConfig.DataX.Home

	// 记录命令信息
	cmdInfo := fmt.Sprintf("执行命令: %s -job %s\n工作目录: %s\n",
		config.GlobalConfig.DataX.Bin,
		tmpFileName,
		cmd.Dir)

	// 追加命令信息（确保LogPath已设置）
	if history.LogPath != "" {
		if logErr := history.AppendLogToFile(cmdInfo, ""); logErr != nil {
			logger.Info("追加命令信息失败: %v", logErr)
		}
	} else {
		// 如果LogPath还没设置，使用WriteLogToFile来追加
		if logErr := history.WriteLogToFile(cmdInfo, ""); logErr != nil {
			logger.Info("写入命令信息失败: %v", logErr)
		}
	}

	// 创建管道来实时捕获输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		history.Status = 0
		if logErr := history.AppendLogToFile("", fmt.Sprintf("创建stdout管道失败: %v", err)); logErr != nil {
			logger.Info("追加日志文件失败: %v", logErr)
		}
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		history.Status = 0
		if logErr := history.AppendLogToFile("", fmt.Sprintf("创建stderr管道失败: %v", err)); logErr != nil {
			logger.Info("追加日志文件失败: %v", logErr)
		}
		return
	}

	// 启动命令
	if err = cmd.Start(); err != nil {
		history.Status = 0
		if logErr := history.AppendLogToFile("", fmt.Sprintf("启动DataX命令失败: %v", err)); logErr != nil {
			logger.Info("追加日志文件失败: %v", logErr)
		}
		return
	}

	// 用于收集所有输出的缓冲区
	var fullOutput strings.Builder

	// 创建goroutine来实时读取stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fullOutput.WriteString(line + "\n")
			// 实时追加日志
			if logErr := history.AppendLogToFile(line+"\n", ""); logErr != nil {
				logger.Info("追加stdout日志失败: %v", logErr)
			}
		}
	}()

	// 创建goroutine来实时读取stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			fullOutput.WriteString(line + "\n")
			// 实时追加日志
			if logErr := history.AppendLogToFile("", line+"\n"); logErr != nil {
				logger.Info("追加stderr日志失败: %v", logErr)
			}
		}
	}()

	// 等待命令完成
	err = cmd.Wait()

	// 获取完整输出用于状态判断
	fullOutputStr := fullOutput.String()

	// 写入任务完成信息
	completionMsg := "\n=== 任务执行完成 ===\n"
	if logErr := history.AppendLogToFile(completionMsg, ""); logErr != nil {
		logger.Info("追加完成信息失败: %v", logErr)
	}

	// 检查所有输出中是否包含成功完成的标志（同时检查stdout和stderr）
	if strings.Contains(fullOutputStr, "数据同步完成") {
		history.Status = 1
		if logErr := history.AppendLogToFile("", "任务执行成功"); logErr != nil {
			logger.Info("追加成功状态失败: %v", logErr)
		}
	} else if err != nil {
		history.Status = 0
		if logErr := history.AppendLogToFile("", fmt.Sprintf("执行DataX任务失败: %v", err)); logErr != nil {
			logger.Info("追加错误状态失败: %v", logErr)
		}
	} else {
		// 如果没有明确的成功标志，也没有执行错误，则检查是否有错误关键字
		fullOutputLower := strings.ToLower(fullOutputStr)
		isError := strings.Contains(fullOutputLower, "error") ||
			strings.Contains(fullOutputLower, "exception") ||
			strings.Contains(fullOutputLower, "失败")

		if isError {
			history.Status = 0
			if logErr := history.AppendLogToFile("", "执行DataX任务失败，输出中包含错误信息"); logErr != nil {
				logger.Info("追加错误状态失败: %v", logErr)
			}
		} else {
			history.Status = 1
			if logErr := history.AppendLogToFile("", "任务执行完成"); logErr != nil {
				logger.Info("追加完成状态失败: %v", logErr)
			}
		}
	}
}
