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
		Name:        req.Name,
		Code:        req.Code,
		Type:        req.Type,
		ParentID:    req.ParentID,
		Path:        req.Path,
		Component:   req.Component,
		Icon:        req.Icon,
		Sort:        req.Sort,
		Status:      1,
		Hidden:      req.Hidden,
		Cache:       req.Cache,
		IsExternal:  req.IsExternal,
		ExternalURL: req.ExternalURL,
		OpenType:    req.OpenType,
	}

	return models.DB.Create(permission).Error
}

// UpdatePermission 更新权限
func (s *PermissionService) UpdatePermission(id uint, req *types.UpdatePermissionRequest) error {
	updates := map[string]any{
		"name":         req.Name,
		"type":         req.Type,
		"parent_id":    req.ParentID,
		"path":         req.Path,
		"component":    req.Component,
		"icon":         req.Icon,
		"sort":         req.Sort,
		"status":       req.Status,
		"hidden":       req.Hidden,
		"cache":        req.Cache,
		"is_external":  req.IsExternal,
		"external_url": req.ExternalURL,
		"open_type":    req.OpenType,
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

	// 构建权限树 - 使用临时结构来正确处理指针引用
	type tempPermission struct {
		*types.PermissionResponse
		Children []*tempPermission
	}

	permMap := make(map[uint]*tempPermission)
	var rootPermIDs []uint

	// 第一次遍历，创建所有节点
	for _, p := range permissions {
		perm := &tempPermission{
			PermissionResponse: &types.PermissionResponse{
				ID:          p.ID,
				Name:        p.Name,
				Code:        p.Code,
				Type:        p.Type,
				ParentID:    p.ParentID,
				Path:        p.Path,
				Component:   p.Component,
				Icon:        p.Icon,
				Sort:        p.Sort,
				Status:      p.Status,
				Hidden:      p.Hidden,
				Cache:       p.Cache,
				IsExternal:  p.IsExternal,
				ExternalURL: p.ExternalURL,
				OpenType:    p.OpenType,
				Children:    make([]types.PermissionResponse, 0),
			},
			Children: make([]*tempPermission, 0),
		}
		permMap[p.ID] = perm
		if p.ParentID == nil {
			rootPermIDs = append(rootPermIDs, p.ID)
		}
	}

	// 第二次遍历，构建树形结构
	for _, p := range permissions {
		if p.ParentID != nil {
			if parent, ok := permMap[*p.ParentID]; ok {
				if child, childOk := permMap[p.ID]; childOk {
					parent.Children = append(parent.Children, child)
				}
			}
		}
	}

	// 递归函数将临时结构转换为最终结构
	var convertToFinal func(*tempPermission) types.PermissionResponse
	convertToFinal = func(temp *tempPermission) types.PermissionResponse {
		result := *temp.PermissionResponse
		result.Children = make([]types.PermissionResponse, len(temp.Children))
		for i, child := range temp.Children {
			result.Children[i] = convertToFinal(child)
		}
		return result
	}

	// 构建根节点列表
	var rootPerms []types.PermissionResponse
	for _, rootID := range rootPermIDs {
		if rootPerm, ok := permMap[rootID]; ok {
			finalPerm := convertToFinal(rootPerm)
			rootPerms = append(rootPerms, finalPerm)
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
		Order("permissions.sort").
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	result := make([]types.PermissionResponse, len(permissions))
	for i, p := range permissions {
		result[i] = types.PermissionResponse{
			ID:          p.ID,
			Name:        p.Name,
			Code:        p.Code,
			Type:        p.Type,
			ParentID:    p.ParentID,
			Path:        p.Path,
			Component:   p.Component,
			Icon:        p.Icon,
			Sort:        p.Sort,
			Status:      p.Status,
			Hidden:      p.Hidden,
			Cache:       p.Cache,
			IsExternal:  p.IsExternal,
			ExternalURL: p.ExternalURL,
			OpenType:    p.OpenType,
		}
	}

	return result, nil
}

// GetUserMenus 获取用户的菜单权限（树形结构）
func (s *PermissionService) GetUserMenus(userID uint) (*types.PermissionTreeResponse, error) {
	var permissions []models.Permission

	// 分两步查询：
	// 1. 先获取用户有权限的所有菜单（包括隐藏的，用于构建完整的树结构）
	var allUserPermissions []models.Permission
	err := models.DB.Distinct("permissions.*").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id and role_permissions.deleted_at IS NULL").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.status = 1 AND permissions.type = 'menu'", userID).
		Order("permissions.sort").
		Find(&allUserPermissions).Error

	if err != nil {
		return nil, err
	}

	// 2. 过滤出非隐藏的权限用于显示
	permissions = make([]models.Permission, 0)
	hiddenParentIDs := make(map[uint]bool) // 记录隐藏的父级权限ID

	for _, p := range allUserPermissions {
		if p.Hidden == 0 {
			permissions = append(permissions, p)
		} else {
			hiddenParentIDs[p.ID] = true
		}
	}

	// 如果没有找到任何权限，返回空列表
	if len(permissions) == 0 {
		return &types.PermissionTreeResponse{List: []types.PermissionResponse{}}, nil
	}

	// 构建权限树 - 使用临时结构来正确处理指针引用
	type tempPermission struct {
		*types.PermissionResponse
		Children []*tempPermission
	}

	permMap := make(map[uint]*tempPermission)
	var rootPermIDs []uint

	// 第一次遍历，创建所有节点
	for _, p := range permissions {
		perm := &tempPermission{
			PermissionResponse: &types.PermissionResponse{
				ID:          p.ID,
				Name:        p.Name,
				Code:        p.Code,
				Type:        p.Type,
				ParentID:    p.ParentID,
				Path:        p.Path,
				Component:   p.Component,
				Icon:        p.Icon,
				Sort:        p.Sort,
				Status:      p.Status,
				Hidden:      p.Hidden,
				Cache:       p.Cache,
				IsExternal:  p.IsExternal,
				ExternalURL: p.ExternalURL,
				OpenType:    p.OpenType,
				Children:    make([]types.PermissionResponse, 0),
			},
			Children: make([]*tempPermission, 0),
		}
		permMap[p.ID] = perm

		// 判断是否为根权限：没有父级ID
		if p.ParentID == nil {
			rootPermIDs = append(rootPermIDs, p.ID)
		}
	}

	// 检查是否所有权限都有父级但父级权限不在当前权限列表中
	if len(rootPermIDs) == 0 {
		for _, p := range permissions {
			if p.ParentID != nil {
				// 检查父级权限是否在当前权限列表中
				parentExists := false
				for _, parent := range permissions {
					if parent.ID == *p.ParentID {
						parentExists = true
						break
					}
				}
				// 如果父级权限不存在，将当前权限作为根权限
				if !parentExists {
					rootPermIDs = append(rootPermIDs, p.ID)
				}
			}
		}
	}

	// 第二次遍历，构建树形结构
	for _, p := range permissions {
		if p.ParentID != nil {
			if parent, ok := permMap[*p.ParentID]; ok {
				if child, childOk := permMap[p.ID]; childOk {
					parent.Children = append(parent.Children, child)
				}
			}
		}
	}

	// 递归函数将临时结构转换为最终结构
	var convertToFinal func(*tempPermission) types.PermissionResponse
	convertToFinal = func(temp *tempPermission) types.PermissionResponse {
		result := *temp.PermissionResponse
		result.Children = make([]types.PermissionResponse, len(temp.Children))
		for i, child := range temp.Children {
			result.Children[i] = convertToFinal(child)
		}
		return result
	}

	// 构建根节点列表
	var rootPerms []types.PermissionResponse
	for _, rootID := range rootPermIDs {
		if rootPerm, ok := permMap[rootID]; ok {
			rootPerms = append(rootPerms, convertToFinal(rootPerm))
		}
	}

	return &types.PermissionTreeResponse{List: rootPerms}, nil
}
