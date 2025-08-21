package types

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Code        string `json:"code" binding:"required,max=50"`
	Type        string `json:"type" binding:"required,oneof=menu button"`
	ParentID    *uint  `json:"parent_id"`
	Path        string `json:"path" binding:"max=200"`
	Component   string `json:"component" binding:"max=200"`
	Icon        string `json:"icon" binding:"max=50"`
	Sort        int    `json:"sort" binding:"min=0"`
	Hidden      int    `json:"hidden" binding:"oneof=0 1"`
	Cache       int    `json:"cache" binding:"oneof=0 1"`
	IsExternal  int    `json:"is_external" binding:"oneof=0 1"` // 是否为外部链接，0-否，1-是
	ExternalURL string `json:"external_url" binding:"max=500"`   // 外部链接地址
	OpenType    int    `json:"open_type" binding:"oneof=0 1"`   // 打开方式，0-内嵌，1-新窗口
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Type        string `json:"type" binding:"required,oneof=menu button"`
	ParentID    *uint  `json:"parent_id"`
	Path        string `json:"path" binding:"max=200"`
	Component   string `json:"component" binding:"max=200"`
	Icon        string `json:"icon" binding:"max=50"`
	Sort        int    `json:"sort" binding:"min=0"`
	Status      int    `json:"status" binding:"oneof=0 1"`
	Hidden      int    `json:"hidden" binding:"oneof=0 1"`
	Cache       int    `json:"cache" binding:"oneof=0 1"`
	IsExternal  int    `json:"is_external" binding:"oneof=0 1"` // 是否为外部链接，0-否，1-是
	ExternalURL string `json:"external_url" binding:"max=500"`   // 外部链接地址
	OpenType    int    `json:"open_type" binding:"oneof=0 1"`   // 打开方式，0-内嵌，1-新窗口
}

// PermissionResponse 权限响应
type PermissionResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Code        string               `json:"code"`
	Type        string               `json:"type"`
	ParentID    *uint                `json:"parent_id"`
	Path        string               `json:"path"`
	Component   string               `json:"component"`
	Icon        string               `json:"icon"`
	Sort        int                  `json:"sort"`
	Status      int                  `json:"status"`
	Hidden      int                  `json:"hidden"`
	Cache       int                  `json:"cache"`
	IsExternal  int                  `json:"is_external"`  // 是否为外部链接，0-否，1-是
	ExternalURL string               `json:"external_url"` // 外部链接地址
	OpenType    int                  `json:"open_type"`    // 打开方式，0-内嵌，1-新窗口
	Children    []PermissionResponse `json:"children,omitempty"`
}

// PermissionListRequest 权限列表请求
type PermissionListRequest struct {
	Type string `form:"type" binding:"omitempty,oneof=menu button"`
}

// PermissionTreeResponse 权限树响应
type PermissionTreeResponse struct {
	List []PermissionResponse `json:"list"`
}
