package services

import (
	"datax-admin/models"
	"datax-admin/types"
	"errors"

	"gorm.io/gorm"
)

// PermissionService 权限服务
type PermissionService struct{}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(req *types.CreatePermissionRequest) error {
	var count int64
	if err := models.DB.Model(&models.Permission{}).Where("code = ?", req.Code).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("权限编码已存在")
	}

	permission := &models.Permission{
		Name:      req.Name,
		Code:      req.Code,
		Type:      req.Type,
		ParentID:  req.ParentID,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Sort:      req.Sort,
		Status:    1,
	}

	return models.DB.Create(permission).Error
}

// UpdatePermission 更新权限
func (s *PermissionService) UpdatePermission(id uint, req *types.UpdatePermissionRequest) error {
	updates := map[string]interface{}{
		"name":      req.Name,
		"type":      req.Type,
		"parent_id": req.ParentID,
		"path":      req.Path,
		"component": req.Component,
		"icon":      req.Icon,
		"sort":      req.Sort,
		"status":    req.Status,
	}

	return models.DB.Model(&models.Permission{}).Where("id = ?", id).Updates(updates).Error
}

// DeletePermission 删除权限
func (s *PermissionService) DeletePermission(id uint) error {
	return models.DB.Transaction(func(tx *gorm.DB) error {
		// 检查是否有子权限
		var count int64
		if err := tx.Model(&models.Permission{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("请先删除子权限")
		}

		// 删除权限
		if err := tx.Delete(&models.Permission{}, id).Error; err != nil {
			return err
		}

		// 删除角色关联的权限
		if err := tx.Where("permission_id = ?", id).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetPermissionTree 获取权限树
func (s *PermissionService) GetPermissionTree(req *types.PermissionListRequest) (*types.PermissionTreeResponse, error) {
	var permissions []models.Permission

	query := models.DB.Model(&models.Permission{}).Order("sort")
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	if err := query.Find(&permissions).Error; err != nil {
		return nil, err
	}

	// 构建权限树
	permMap := make(map[uint]*types.PermissionResponse)
	var rootPerms []types.PermissionResponse

	// 第一次遍历，创建所有节点
	for _, p := range permissions {
		perm := types.PermissionResponse{
			ID:        p.ID,
			Name:      p.Name,
			Code:      p.Code,
			Type:      p.Type,
			ParentID:  p.ParentID,
			Path:      p.Path,
			Component: p.Component,
			Icon:      p.Icon,
			Sort:      p.Sort,
			Status:    p.Status,
			Children:  make([]types.PermissionResponse, 0),
		}
		permMap[p.ID] = &perm
		if p.ParentID == nil {
			rootPerms = append(rootPerms, perm)
		}
	}

	// 第二次遍历，构建树形结构
	for _, p := range permissions {
		if p.ParentID != nil {
			if parent, ok := permMap[*p.ParentID]; ok {
				parent.Children = append(parent.Children, *permMap[p.ID])
			}
		}
	}

	return &types.PermissionTreeResponse{List: rootPerms}, nil
}

// GetUserPermissions 获取用户的所有权限
func (s *PermissionService) GetUserPermissions(userID uint) ([]types.PermissionResponse, error) {
	var permissions []models.Permission

	err := models.DB.Distinct("permissions.*").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.status = 1", userID).
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	result := make([]types.PermissionResponse, len(permissions))
	for i, p := range permissions {
		result[i] = types.PermissionResponse{
			ID:        p.ID,
			Name:      p.Name,
			Code:      p.Code,
			Type:      p.Type,
			ParentID:  p.ParentID,
			Path:      p.Path,
			Component: p.Component,
			Icon:      p.Icon,
			Sort:      p.Sort,
			Status:    p.Status,
		}
	}

	return result, nil
}
