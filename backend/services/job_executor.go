package services

import (
	"bytes"
	"context"
	"datax-admin/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// executeShellJob 执行Shell任务
func (s *JobService) executeShellJob(job *models.Job, params interface{}, history *models.JobHistory) {
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
func (s *JobService) executeHTTPJob(job *models.Job, params interface{}, history *models.JobHistory) {
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
func (s *JobService) executeDataXJob(job *models.Job, params interface{}, history *models.JobHistory) {
	var dataxParams models.JobDataXParams
	if err := mapToStruct(params, &dataxParams); err != nil {
		history.Status = 0
		history.Error = fmt.Sprintf("解析DataX参数失败: %v", err)
		return
	}

	// 创建临时JSON文件
	tmpFile, err := os.CreateTemp("", "datax-*.json")
	if err != nil {
		history.Status = 0
		history.Error = fmt.Sprintf("创建临时文件失败: %v", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	// 写入JSON配置
	if _, err := tmpFile.WriteString(dataxParams.JobConfig); err != nil {
		history.Status = 0
		history.Error = fmt.Sprintf("写入配置失败: %v", err)
		return
	}
	tmpFile.Close()

	// 构建命令参数
	args := []string{
		"-job",
		tmpFile.Name(),
	}

	// 添加速率限制
	if dataxParams.Speed > 0 {
		args = append(args, "-s", fmt.Sprintf("%d", dataxParams.Speed))
	}

	// 添加自定义参数
	if len(dataxParams.Parameters) > 0 {
		for k, v := range dataxParams.Parameters {
			args = append(args, fmt.Sprintf("-p%s=%s", k, v))
		}
	}

	// 创建上下文，处理超时
	ctx := context.Background()
	if job.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(job.Timeout)*time.Second)
		defer cancel()
	}

	// 准备命令
	cmd := exec.CommandContext(ctx, "./bin/datax", args...)

	// 捕获输出
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令
	err = cmd.Run()
	if err != nil {
		history.Status = 0
		history.Error = fmt.Sprintf("执行DataX任务失败: %v\n%s", err, stderr.String())
	} else {
		history.Status = 1
		history.Output = stdout.String()
	}
}
