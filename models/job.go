package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

// JobType 任务类型
type JobType string

const (
	JobTypeShell JobType = "shell" // Shell 脚本
	JobTypeHTTP  JobType = "http"  // HTTP 请求
	JobTypeDataX JobType = "datax" // DataX 任务
)

// JobStatus 任务状态
type JobStatus int

const (
	JobStatusStop     JobStatus = 0 // 停止
	JobStatusRunning  JobStatus = 1 // 运行中
	JobStatusDisabled JobStatus = 2 // 禁用
)

// Job 任务模型
type Job struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`      // 任务名称
	Type        JobType        `gorm:"size:20;not null" json:"type"`       // 任务类型
	Description string         `gorm:"size:500" json:"description"`        // 任务描述
	CronExpr    string         `gorm:"size:100;not null" json:"cron_expr"` // Cron 表达式
	Status      JobStatus      `gorm:"default:0" json:"status"`            // 任务状态
	Timeout     int            `gorm:"default:0" json:"timeout"`           // 超时时间(秒)
	RetryCount  int            `gorm:"default:0" json:"retry_count"`       // 重试次数
	RetryDelay  int            `gorm:"default:0" json:"retry_delay"`       // 重试间隔(秒)
	Params      string         `gorm:"type:text" json:"params"`            // 任务参数(JSON)
	Creator     uint           `gorm:"not null" json:"creator"`            // 创建者ID
	Updater     uint           `json:"updater"`                            // 更新者ID
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// JobShellParams Shell任务参数
type JobShellParams struct {
	Command     string            `json:"command"`     // 执行命令
	WorkDir     string            `json:"work_dir"`    // 工作目录
	Environment map[string]string `json:"environment"` // 环境变量
}

// JobHTTPParams HTTP任务参数
type JobHTTPParams struct {
	URL         string            `json:"url"`          // 请求URL
	Method      string            `json:"method"`       // 请求方法
	Headers     map[string]string `json:"headers"`      // 请求头
	Body        string            `json:"body"`         // 请求体
	SuccessCode []int             `json:"success_code"` // 成功状态码
}

// JobDataXParams DataX任务参数
type JobDataXParams struct {
	JobConfig  string            `json:"job_config"` // 任务JSON配置
	Parameters map[string]string `json:"parameters"` // 任务参数
}

// JobHistory 任务执行历史
type JobHistory struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	JobID     uint       `json:"job_id"`
	JobName   string     `json:"job_name" gorm:"-"` // 虚拟字段，不存储到数据库
	Status    int        `json:"status"`            // -1:执行中 0:失败 1:成功
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`                 // 使用指针类型，允许NULL
	Duration  int64      `json:"duration"`                 // 执行时长(毫秒)
	LogPath   string     `json:"log_path" gorm:"size:500"` // 日志文件路径
	Output    string     `json:"output" gorm:"-"`          // 执行输出（虚拟字段，从文件读取）
	Error     string     `json:"error" gorm:"-"`           // 错误信息（虚拟字段，从文件读取）
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (Job) TableName() string {
	return "jobs"
}

// TableName 指定表名
func (JobHistory) TableName() string {
	return "job_histories"
}

// JobLogData 日志数据结构
type JobLogData struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

// WriteLogToFile 将日志写入文件
func (jh *JobHistory) WriteLogToFile(output, error string) error {
	// 确保日志目录存在
	logDir := filepath.Join("logs", "jobs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 生成日志文件名：job_{job_id}_{history_id}_{timestamp}.json
	var filename string
	if jh.ID == 0 {
		filename = fmt.Sprintf("job_%d_%d.json", jh.JobID, jh.StartTime.UnixNano())
	} else {
		filename = fmt.Sprintf("job_%d_%d_%d.json", jh.JobID, jh.ID, jh.StartTime.Unix())
	}
	jh.LogPath = filepath.Join(logDir, filename)

	// 准备日志数据
	logData := JobLogData{
		Output: output,
		Error:  error,
	}

	// 序列化为JSON
	data, err := json.MarshalIndent(logData, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化日志数据失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(jh.LogPath, data, 0644); err != nil {
		return fmt.Errorf("写入日志文件失败: %v", err)
	}

	// 如果有ID（已保存到数据库），立即更新log_path字段
	if jh.ID > 0 {
		if err := DB.Model(jh).Update("log_path", jh.LogPath).Error; err != nil {
			// 记录错误但不返回，避免影响主要流程
			return fmt.Errorf("更新日志路径失败: %v", err)
		}
	}

	return nil
}

// AppendLogToFile 追加日志内容到文件（用于实时日志）
func (jh *JobHistory) AppendLogToFile(newOutput, newError string) error {
	if jh.LogPath == "" {
		return fmt.Errorf("日志文件路径为空")
	}

	// 读取现有日志内容
	var existingData JobLogData
	if data, err := os.ReadFile(jh.LogPath); err == nil {
		// 如果文件存在，解析现有内容
		if err := json.Unmarshal(data, &existingData); err != nil {
			// 如果解析失败，初始化为空
			existingData = JobLogData{}
		}
	}

	// 追加新内容
	if newOutput != "" {
		existingData.Output += newOutput
	}
	if newError != "" {
		existingData.Error += newError
	}

	// 写回文件
	data, err := json.MarshalIndent(existingData, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化日志数据失败: %v", err)
	}

	if err := os.WriteFile(jh.LogPath, data, 0644); err != nil {
		return fmt.Errorf("写入日志文件失败: %v", err)
	}

	return nil
}

// ReadLogFromFile 从文件读取日志
func (jh *JobHistory) ReadLogFromFile() error {
	if jh.LogPath == "" {
		return nil
	}

	// 检查文件是否存在
	if _, err := os.Stat(jh.LogPath); os.IsNotExist(err) {
		jh.Output = ""
		jh.Error = ""
		return nil
	}

	// 读取文件内容
	data, err := os.ReadFile(jh.LogPath)
	if err != nil {
		return fmt.Errorf("读取日志文件失败: %v", err)
	}

	// 反序列化JSON
	var logData JobLogData
	if err := json.Unmarshal(data, &logData); err != nil {
		return fmt.Errorf("解析日志文件失败: %v", err)
	}

	jh.Output = logData.Output
	jh.Error = logData.Error

	return nil
}
