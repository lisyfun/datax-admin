package services

import (
	"datax-admin/models"
	"datax-admin/types"
	"datax-admin/utils/secret"
	"errors"
	"time"

	"gorm.io/gorm"
)

// TerminalService 终端服务
type TerminalService struct {
}

// NewTerminalService 创建终端服务
func NewTerminalService() *TerminalService {
	return &TerminalService{}
}

// CreateTerminal 创建终端
func (s *TerminalService) CreateTerminal(req *types.CreateTerminalRequest) error {
	terminal := &models.Terminal{
		Name:          req.Name,
		Host:          req.Host,
		Port:          req.Port,
		Username:      req.Username,
		AuthType:      req.AuthType,
		Password:      req.Password,
		KeyFile:       req.KeyFile,
		KeyPassphrase: req.KeyPassphrase,
		Status:        "offline",
	}

	return models.DB.Create(terminal).Error
}

// UpdateTerminal 更新终端
func (s *TerminalService) UpdateTerminal(id uint, req *types.UpdateTerminalRequest) error {
	// 使用结构体更新，确保 GORM 的加密钩子被触发。
	// 敏感字段仅在请求中显式提供时才覆盖，避免编辑时把已有凭据清空。
	updates := models.Terminal{
		Name:     req.Name,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		AuthType: req.AuthType,
	}
	fields := []string{"name", "host", "port", "username", "auth_type"}

	if req.Password != "" {
		v, err := secret.Encrypt(req.Password)
		if err != nil {
			return err
		}
		updates.Password = v
		fields = append(fields, "password")
	}
	if req.KeyFile != "" {
		v, err := secret.Encrypt(req.KeyFile)
		if err != nil {
			return err
		}
		updates.KeyFile = v
		fields = append(fields, "key_file")
	}
	if req.KeyPassphrase != "" {
		v, err := secret.Encrypt(req.KeyPassphrase)
		if err != nil {
			return err
		}
		updates.KeyPassphrase = v
		fields = append(fields, "key_passphrase")
	}

	result := models.DB.Model(&models.Terminal{}).Where("id = ?", id).Select(fields).Updates(&updates)
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
			ID:          terminal.ID,
			Name:        terminal.Name,
			Host:        terminal.Host,
			Port:        terminal.Port,
			Username:    terminal.Username,
			AuthType:    terminal.AuthType,
			HasPassword: terminal.Password != "",
			HasKeyFile:  terminal.KeyFile != "",
			Status:      terminal.Status,
			LastSeen:    terminal.LastSeen,
			CreatedAt:   terminal.CreatedAt,
		}
	}

	return &types.TerminalListResponse{
		Total: total,
		List:  list,
	}, nil
}

// UpdateTerminalStatus 更新终端状态
func (s *TerminalService) UpdateTerminalStatus(id uint, status string) error {
	result := models.DB.Model(&models.Terminal{}).Where("id = ?", id).Updates(map[string]any{
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
		ID:            terminal.ID,
		Name:          terminal.Name,
		Host:          terminal.Host,
		Port:          terminal.Port,
		Username:      terminal.Username,
		AuthType:      terminal.AuthType,
		Password:      terminal.Password,
		KeyFile:       terminal.KeyFile,
		KeyPassphrase: terminal.KeyPassphrase,
		HasPassword:   terminal.Password != "",
		HasKeyFile:    terminal.KeyFile != "",
		Status:        terminal.Status,
		LastSeen:      terminal.LastSeen,
		CreatedAt:     terminal.CreatedAt,
	}, nil
}
