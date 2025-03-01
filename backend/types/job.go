package types

import (
	"datax-admin/models"
	"time"
)

// CreateJobRequest 创建任务请求
type CreateJobRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Type        string `json:"type" binding:"required,oneof=shell http datax"`
	Description string `json:"description" binding:"max=500"`
	CronExpr    string `json:"cron_expr" binding:"required,max=100"`
	Timeout     int    `json:"timeout" binding:"min=0"`
	RetryCount  int    `json:"retry_count" binding:"min=0"`
	RetryDelay  int    `json:"retry_delay" binding:"min=0"`
	Params      any    `json:"params" binding:"required"`
}

// UpdateJobRequest 更新任务请求
type UpdateJobRequest struct {
	Name        string `json:"name" binding:"omitempty,max=100"`
	Description string `json:"description" binding:"max=500"`
	CronExpr    string `json:"cron_expr" binding:"omitempty,max=100"`
	Timeout     int    `json:"timeout" binding:"min=0"`
	RetryCount  int    `json:"retry_count" binding:"min=0"`
	RetryDelay  int    `json:"retry_delay" binding:"min=0"`
	Params      any    `json:"params"`
}

// JobListRequest 任务列表请求
type JobListRequest struct {
	Page     int    `form:"page,default=1" binding:"required,min=1"`
	PageSize int    `form:"page_size,default=10" binding:"required,min=5,max=100"`
	Keyword  string `form:"keyword" binding:"omitempty,max=100"`
	Type     string `form:"type" binding:"omitempty,oneof=shell http datax"`
	Status   *int   `form:"status" binding:"omitempty,oneof=0 1 2"`
}

// JobListResponse 任务列表响应
type JobListResponse struct {
	Total int64        `json:"total"`
	Items []models.Job `json:"items"`
}

// JobHistoryListRequest 任务历史列表请求
type JobHistoryListRequest struct {
	Page      int       `form:"page"`
	PageSize  int       `form:"page_size"`
	JobID     *int      `form:"job_id"`
	Status    *int      `form:"status"`
	Keyword   string    `form:"keyword"`
	StartTime time.Time `form:"start_time"`
	EndTime   time.Time `form:"end_time"`
}

// JobHistoryListResponse 任务历史列表响应
type JobHistoryListResponse struct {
	Total int64               `json:"total"`
	Items []models.JobHistory `json:"items"`
}
