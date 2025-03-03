package services

import (
	"bytes"
	"context"
	"datax-admin/config"
	"datax-admin/models"
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
		history.Error = fmt.Sprintf("解析Shell参数失败: %v", err)
		return
	}

	// 创建上下文，处理超时
	ctx := context.Background()
	if job.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(job.Timeout)*time.Second)
		defer cancel()
	}

	// 准备命令
	cmd := exec.CommandContext(ctx, "sh", "-c", shellParams.Command+" 2>&1")
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

	// 捕获输出
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令
	err := cmd.Run()
	if err != nil {
		history.Status = 0
		history.Error = fmt.Sprintf("执行Shell命令失败: %v\n%s", err, stderr.String())
	} else {
		history.Status = 1
		history.Output = stdout.String()
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
		history.Error = fmt.Sprintf("创建HTTP请求失败: %v", err)
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
		history.Error = fmt.Sprintf("发送HTTP请求失败: %v", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		history.Status = 0
		history.Error = fmt.Sprintf("读取HTTP响应失败: %v", err)
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
			history.Error = fmt.Sprintf("HTTP响应状态码不符合预期: %d", resp.StatusCode)
			history.Output = string(body)
			return
		}
	}

	history.Status = 1
	history.Output = string(body)
}

// executeDataXJob 执行DataX任务
func (s *JobService) executeDataXJob(job *models.Job, params any, history *models.JobHistory) {
	var dataxParams models.JobDataXParams
	if err := mapToStruct(params, &dataxParams); err != nil {
		history.Status = 0
		history.Error = fmt.Sprintf("解析DataX参数失败: %v", err)
		return
	}

	// 格式化配置内容
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(dataxParams.JobConfig), "", "    "); err != nil {
		history.Status = 0
		history.Error = fmt.Sprintf("格式化配置内容失败: %v", err)
		return
	}

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		history.Status = 0
		history.Error = fmt.Sprintf("获取当前工作目录失败: %v", err)
		return
	}

	// 在当前目录下创建临时文件
	tmpFile, err := os.CreateTemp(currentDir, "datax-*.json")
	if err != nil {
		history.Status = 0
		history.Error = fmt.Sprintf("创建临时文件失败: %v", err)
		return
	}
	tmpFileName := tmpFile.Name()

	// 延迟删除临时文件，但要等到命令执行完成后
	defer func() {
		if err := os.Remove(tmpFileName); err != nil {
			fmt.Printf("删除临时文件失败: %v\n", err)
		}
	}()

	// 写入JSON配置
	if _, err := tmpFile.WriteString(dataxParams.JobConfig); err != nil {
		history.Status = 0
		history.Error = fmt.Sprintf("写入配置失败: %v", err)
		tmpFile.Close()
		return
	}

	// 确保数据写入磁盘并关闭文件
	if err := tmpFile.Sync(); err != nil {
		history.Status = 0
		history.Error = fmt.Sprintf("同步文件失败: %v", err)
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

	// 捕获输出
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令
	err = cmd.Run()

	// 获取输出内容
	outMsg := stdout.String()
	errMsg := stderr.String()

	fullOutput := outMsg + errMsg

	// 合并输出内容
	var combinedOutput = fmt.Sprintf("命令信息:\n%s\n\n执行日志:\n%s", cmdInfo, fullOutput)

	// 检查所有输出中是否包含成功完成的标志（同时检查stdout和stderr）
	if strings.Contains(fullOutput, "数据同步完成") {
		history.Status = 1
		history.Output = combinedOutput
	} else if err != nil {
		history.Status = 0
		history.Error = fmt.Sprintf("执行DataX任务失败: %v", err)
		history.Output = combinedOutput
	} else {
		// 如果没有明确的成功标志，也没有执行错误，则检查是否有错误关键字
		fullOutputLower := strings.ToLower(fullOutput)
		isError := strings.Contains(fullOutputLower, "error") ||
			strings.Contains(fullOutputLower, "exception") ||
			strings.Contains(fullOutputLower, "失败")

		if isError {
			history.Status = 0
			history.Error = "执行DataX任务失败，输出中包含错误信息"
			history.Output = combinedOutput
		} else {
			history.Status = 1
			history.Output = combinedOutput
		}
	}
}
