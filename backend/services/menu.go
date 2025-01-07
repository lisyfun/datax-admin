package services

import (
	"datax-admin/models"
	"datax-admin/types"
	"errors"
	"sort"
)

// MenuService 菜单服务
type MenuService struct{}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(req *types.CreateMenuRequest) error {
	menu := &models.Menu{
		ParentID:  req.ParentID,
		Name:      req.Name,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Sort:      req.Sort,
		Hidden:    req.Hidden,
		Cache:     req.Cache,
		Type:      req.Type,
		Status:    1,
	}

	return models.DB.Create(menu).Error
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(id uint, req *types.UpdateMenuRequest) error {
	return models.DB.Model(&models.Menu{}).Where("id = ?", id).Updates(map[string]interface{}{
		"parent_id": req.ParentID,
		"name":      req.Name,
		"path":      req.Path,
		"component": req.Component,
		"icon":      req.Icon,
		"sort":      req.Sort,
		"status":    req.Status,
		"hidden":    req.Hidden,
		"cache":     req.Cache,
		"type":      req.Type,
	}).Error
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(id uint) error {
	// 检查是否有子菜单
	var count int64
	if err := models.DB.Model(&models.Menu{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("存在子菜单，无法删除")
	}

	return models.DB.Delete(&models.Menu{}, id).Error
}

// GetMenuList 获取菜单列表
func (s *MenuService) GetMenuList(req *types.MenuListRequest) (*types.MenuListResponse, error) {
	var menus []models.Menu

	query := models.DB.Model(&models.Menu{})
	if req.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+req.Keyword+"%")
	}
	if req.Type > 0 {
		query = query.Where("type = ?", req.Type)
	}

	if err := query.Order("sort").Find(&menus).Error; err != nil {
		return nil, err
	}

	// 转换为树形结构
	menuTree := s.buildMenuTree(menus, 0)

	return &types.MenuListResponse{
		List: menuTree,
	}, nil
}

// buildMenuTree 构建菜单树
func (s *MenuService) buildMenuTree(menus []models.Menu, parentID uint) []types.MenuResponse {
	var tree []types.MenuResponse

	for _, menu := range menus {
		if menu.ParentID == parentID {
			node := types.MenuResponse{
				ID:        menu.ID,
				ParentID:  menu.ParentID,
				Name:      menu.Name,
				Path:      menu.Path,
				Component: menu.Component,
				Icon:      menu.Icon,
				Sort:      menu.Sort,
				Status:    menu.Status,
				Hidden:    menu.Hidden,
				Cache:     menu.Cache,
				Type:      menu.Type,
				Children:  s.buildMenuTree(menus, menu.ID),
			}
			tree = append(tree, node)
		}
	}

	// 按照Sort字段排序
	sort.Slice(tree, func(i, j int) bool {
		return tree[i].Sort < tree[j].Sort
	})

	return tree
}
