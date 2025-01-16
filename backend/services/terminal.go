package services

import (
	"datax-admin/models"
	"datax-admin/types"
	"errors"
	"time"

	"gorm.io/gorm"
)

// TerminalService 终端服务
type TerminalService struct{}

// CreateTerminal 创建终端
func (s *TerminalService) CreateTerminal(req *types.CreateTerminalRequest) error {
	terminal := &models.Terminal{
		Name:     req.Name,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Status:   "offline",
	}

	return models.DB.Create(terminal).Error
}

// UpdateTerminal 更新终端
func (s *TerminalService) UpdateTerminal(id uint, req *types.UpdateTerminalRequest) error {
	updates := map[string]interface{}{
		"name":     req.Name,
		"host":     req.Host,
		"port":     req.Port,
		"username": req.Username,
	}

	if req.Password != "" {
		updates["password"] = req.Password
	}

	result := models.DB.Model(&models.Terminal{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("终端不存在")
	}
	return nil
}

// DeleteTerminal 删除终端
func (s *TerminalService) DeleteTerminal(id uint) error {
	result := models.DB.Delete(&models.Terminal{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("终端不存在")
	}
	return nil
}

// GetTerminalList 获取终端列表
func (s *TerminalService) GetTerminalList(req *types.TerminalListRequest) (*types.TerminalListResponse, error) {
	var total int64
	var terminals []models.Terminal

	query := models.DB.Model(&models.Terminal{})

	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Host != "" {
		query = query.Where("host LIKE ?", "%"+req.Host+"%")
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 获取分页数据
	if err := query.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).Find(&terminals).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]types.TerminalInfo, len(terminals))
	for i, terminal := range terminals {
		list[i] = types.TerminalInfo{
			ID:        terminal.ID,
			Name:      terminal.Name,
			Host:      terminal.Host,
			Port:      terminal.Port,
			Username:  terminal.Username,
			Status:    terminal.Status,
			LastSeen:  terminal.LastSeen,
			CreatedAt: terminal.CreatedAt,
		}
	}

	return &types.TerminalListResponse{
		Total: total,
		List:  list,
	}, nil
}

// UpdateTerminalStatus 更新终端状态
func (s *TerminalService) UpdateTerminalStatus(id uint, status string) error {
	result := models.DB.Model(&models.Terminal{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    status,
		"last_seen": time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("终端不存在")
	}
	return nil
}

// GetTerminalByID 根据ID获取终端信息
func (s *TerminalService) GetTerminalByID(id uint) (*types.TerminalInfo, error) {
	var terminal models.Terminal
	if err := models.DB.First(&terminal, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("终端不存在")
		}
		return nil, err
	}

	return &types.TerminalInfo{
		ID:        terminal.ID,
		Name:      terminal.Name,
		Host:      terminal.Host,
		Port:      terminal.Port,
		Username:  terminal.Username,
		Password:  terminal.Password,
		Status:    terminal.Status,
		LastSeen:  terminal.LastSeen,
		CreatedAt: terminal.CreatedAt,
	}, nil
}
