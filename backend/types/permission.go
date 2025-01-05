package types

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Name      string `json:"name" binding:"required,max=50"`
	Code      string `json:"code" binding:"required,max=50"`
	Type      string `json:"type" binding:"required,oneof=menu button"`
	ParentID  *uint  `json:"parent_id"`
	Path      string `json:"path" binding:"max=200"`
	Component string `json:"component" binding:"max=200"`
	Icon      string `json:"icon" binding:"max=50"`
	Sort      int    `json:"sort" binding:"min=0"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	Name      string `json:"name" binding:"required,max=50"`
	Type      string `json:"type" binding:"required,oneof=menu button"`
	ParentID  *uint  `json:"parent_id"`
	Path      string `json:"path" binding:"max=200"`
	Component string `json:"component" binding:"max=200"`
	Icon      string `json:"icon" binding:"max=50"`
	Sort      int    `json:"sort" binding:"min=0"`
}

// PermissionResponse 权限响应
type PermissionResponse struct {
	ID        uint                 `json:"id"`
	Name      string               `json:"name"`
	Code      string               `json:"code"`
	Type      string               `json:"type"`
	ParentID  *uint                `json:"parent_id"`
	Path      string               `json:"path"`
	Component string               `json:"component"`
	Icon      string               `json:"icon"`
	Sort      int                  `json:"sort"`
	Status    int                  `json:"status"`
	Children  []PermissionResponse `json:"children,omitempty"`
}

// PermissionListRequest 权限列表请求
type PermissionListRequest struct {
	Type string `form:"type" binding:"omitempty,oneof=menu button"`
}

// PermissionTreeResponse 权限树响应
type PermissionTreeResponse struct {
	List []PermissionResponse `json:"list"`
}
