package types

import "datax-admin/models"

// OperationLogListRequest 操作日志列表请求
type OperationLogListRequest struct {
	Page      int    `form:"page" binding:"required,min=1"`
	PageSize  int    `form:"page_size" binding:"required,min=1,max=100"`
	Username  string `form:"username"`
	Module    string `form:"module"`
	Action    string `form:"action"`
	Status    *int   `form:"status"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
}

// OperationLogListResponse 操作日志列表响应
type OperationLogListResponse struct {
	List     []models.OperationLog `json:"list"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
}

// BatchDeleteLogsRequest 批量删除日志请求
type BatchDeleteLogsRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

// ClearLogsRequest 清空日志请求
type ClearLogsRequest struct {
	BeforeDays int `json:"before_days" binding:"required,min=1"`
}
