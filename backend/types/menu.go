package types

// CreateMenuRequest 创建菜单请求
type CreateMenuRequest struct {
	ParentID  uint   `json:"parent_id"`
	Name      string `json:"name" binding:"required,max=50"`
	Path      string `json:"path" binding:"max=200"`
	Component string `json:"component" binding:"max=200"`
	Icon      string `json:"icon" binding:"max=50"`
	Sort      int    `json:"sort"`
	Hidden    bool   `json:"hidden"`
	Cache     bool   `json:"cache"`
	Type      int    `json:"type" binding:"required,oneof=1 2"`
}

// UpdateMenuRequest 更新菜单请求
type UpdateMenuRequest struct {
	ParentID  uint   `json:"parent_id"`
	Name      string `json:"name" binding:"required,max=50"`
	Path      string `json:"path" binding:"max=200"`
	Component string `json:"component" binding:"max=200"`
	Icon      string `json:"icon" binding:"max=50"`
	Sort      int    `json:"sort"`
	Status    int    `json:"status" binding:"oneof=0 1"`
	Hidden    bool   `json:"hidden"`
	Cache     bool   `json:"cache"`
	Type      int    `json:"type" binding:"required,oneof=1 2"`
}

// MenuResponse 菜单响应
type MenuResponse struct {
	ID        uint           `json:"id"`
	ParentID  uint           `json:"parent_id"`
	Name      string         `json:"name"`
	Path      string         `json:"path"`
	Component string         `json:"component"`
	Icon      string         `json:"icon"`
	Sort      int            `json:"sort"`
	Status    int            `json:"status"`
	Hidden    bool           `json:"hidden"`
	Cache     bool           `json:"cache"`
	Type      int            `json:"type"`
	Children  []MenuResponse `json:"children,omitempty"`
}

// MenuListRequest 菜单列表请求
type MenuListRequest struct {
	Keyword string `form:"keyword"`
	Type    int    `form:"type"`
}

// MenuListResponse 菜单列表响应
type MenuListResponse struct {
	List []MenuResponse `json:"list"`
}
