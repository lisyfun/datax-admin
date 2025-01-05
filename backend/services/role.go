package services

import (
	"datax-admin/models"
	"datax-admin/types"
	"errors"

	"gorm.io/gorm"
)

// RoleService 角色服务
type RoleService struct{}

// CreateRole 创建角色
func (s *RoleService) CreateRole(req *types.CreateRoleRequest) error {
	var count int64
	if err := models.DB.Model(&models.Role{}).Where("code = ?", req.Code).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("角色编码已存在")
	}

	role := &models.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      1,
	}

	return models.DB.Create(role).Error
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(id uint, req *types.UpdateRoleRequest) error {
	return models.DB.Model(&models.Role{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
	}).Error
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(id uint) error {
	return models.DB.Transaction(func(tx *gorm.DB) error {
		// 删除角色
		if err := tx.Delete(&models.Role{}, id).Error; err != nil {
			return err
		}

		// 删除角色关联的权限
		if err := tx.Where("role_id = ?", id).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}

		// 删除用户关联的角色
		if err := tx.Where("role_id = ?", id).Delete(&models.UserRole{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetRoleList 获取角色列表
func (s *RoleService) GetRoleList(req *types.RoleListRequest) (*types.RoleListResponse, error) {
	var total int64
	var roles []models.Role

	query := models.DB.Model(&models.Role{})
	if req.Keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR description LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).Find(&roles).Error; err != nil {
		return nil, err
	}

	items := make([]types.RoleResponse, len(roles))
	for i, role := range roles {
		items[i] = types.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
			Status:      role.Status,
		}
	}

	return &types.RoleListResponse{
		Total: total,
		Items: items,
	}, nil
}

// UpdateRolePermissions 更新角色权限
func (s *RoleService) UpdateRolePermissions(roleID uint, req *types.UpdateRolePermissionsRequest) error {
	return models.DB.Transaction(func(tx *gorm.DB) error {
		// 删除原有权限
		if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}

		// 添加新权限
		for _, permID := range req.PermissionIDs {
			rp := &models.RolePermission{
				RoleID:       roleID,
				PermissionID: permID,
			}
			if err := tx.Create(rp).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetRolePermissions 获取角色的权限ID列表
func (s *RoleService) GetRolePermissions(roleID uint) ([]uint, error) {
	var permissions []uint
	err := models.DB.Model(&models.RolePermission{}).
		Where("role_id = ?", roleID).
		Pluck("permission_id", &permissions).Error
	return permissions, err
}

// UpdateUserRoles 更新用户角色
func (s *RoleService) UpdateUserRoles(userID uint, req *types.UpdateUserRolesRequest) error {
	return models.DB.Transaction(func(tx *gorm.DB) error {
		// 删除原有角色
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
			return err
		}

		// 添加新角色
		for _, roleID := range req.RoleIDs {
			ur := &models.UserRole{
				UserID: userID,
				RoleID: roleID,
			}
			if err := tx.Create(ur).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetUserRoles 获取用户的角色ID列表
func (s *RoleService) GetUserRoles(userID uint) ([]uint, error) {
	var roles []uint
	err := models.DB.Model(&models.UserRole{}).
		Where("user_id = ?", userID).
		Pluck("role_id", &roles).Error
	return roles, err
}
