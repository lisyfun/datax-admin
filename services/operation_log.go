package services

import (
	"datax-admin/models"
	"datax-admin/types"
	"encoding/json"
	"fmt"
	"time"
)

// OperationLogService 操作日志服务
type OperationLogService struct{}

// NewOperationLogService 创建操作日志服务
func NewOperationLogService() *OperationLogService {
	return &OperationLogService{}
}

// CreateLog 创建操作日志
func (s *OperationLogService) CreateLog(log *models.OperationLog) error {
	return models.DB.Create(log).Error
}

// GetLogList 获取操作日志列表
func (s *OperationLogService) GetLogList(req *types.OperationLogListRequest) (*types.OperationLogListResponse, error) {
	var logs []models.OperationLog
	var total int64

	query := models.DB.Model(&models.OperationLog{})

	// 搜索条件
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Module != "" {
		query = query.Where("module = ?", req.Module)
	}
	if req.Action != "" {
		query = query.Where("action = ?", req.Action)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.StartTime != "" {
		query = query.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		query = query.Where("created_at <= ?", req.EndTime)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&logs).Error; err != nil {
		return nil, err
	}

	return &types.OperationLogListResponse{
		List:     logs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// DeleteLog 删除操作日志
func (s *OperationLogService) DeleteLog(id uint) error {
	return models.DB.Delete(&models.OperationLog{}, id).Error
}

// BatchDeleteLogs 批量删除操作日志
func (s *OperationLogService) BatchDeleteLogs(ids []uint) error {
	return models.DB.Delete(&models.OperationLog{}, ids).Error
}

// ClearLogs 清空操作日志
func (s *OperationLogService) ClearLogs(beforeDays int) error {
	if beforeDays <= 0 {
		return fmt.Errorf("天数必须大于0")
	}

	beforeTime := time.Now().AddDate(0, 0, -beforeDays)
	return models.DB.Where("created_at < ?", beforeTime).Delete(&models.OperationLog{}).Error
}

// LogOperation 记录操作日志的辅助函数
func LogOperation(userID uint, username, module, action, description, ip, userAgent string, requestData interface{}, status int, errorMsg string) {
	var requestDataStr string
	if requestData != nil {
		if data, err := json.Marshal(requestData); err == nil {
			requestDataStr = string(data)
		}
	}

	log := &models.OperationLog{
		UserID:      userID,
		Username:    username,
		Module:      module,
		Action:      action,
		Description: description,
		IP:          ip,
		UserAgent:   userAgent,
		RequestData: requestDataStr,
		Status:      status,
		ErrorMsg:    errorMsg,
	}

	service := NewOperationLogService()
	if err := service.CreateLog(log); err != nil {
		// 记录日志失败不应该影响主要业务逻辑，静默处理
		// 可以考虑使用专门的日志系统记录此类错误
	}
}
