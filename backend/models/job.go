package models

import (
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
	ID        uint      `json:"id" gorm:"primaryKey"`
	JobID     uint      `json:"job_id"`
	JobName   string    `json:"job_name" gorm:"-"` // 虚拟字段，不存储到数据库
	Status    int       `json:"status"`            // 0:失败 1:成功
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  int64     `json:"duration"` // 执行时长(毫秒)
	Output    string    `json:"output"`   // 执行输出
	Error     string    `json:"error"`    // 错误信息
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Job) TableName() string {
	return "jobs"
}

// TableName 指定表名
func (JobHistory) TableName() string {
	return "job_histories"
}
