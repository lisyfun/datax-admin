package types

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Code        string `json:"code" binding:"required,max=50"`
	Description string `json:"description" binding:"max=200"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description" binding:"max=200"`
}

// RoleResponse 角色响应
type RoleResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

// RoleListRequest 角色列表请求
type RoleListRequest struct {
	Page     int    `form:"page,default=1" binding:"required,min=1"`
	PageSize int    `form:"page_size,default=10" binding:"required,min=5,max=100"`
	Keyword  string `form:"keyword" binding:"omitempty,max=50"`
}

// RoleListResponse 角色列表响应
type RoleListResponse struct {
	Total int64          `json:"total"`
	Items []RoleResponse `json:"items"`
}

// UpdateRolePermissionsRequest 更新角色权限请求
type UpdateRolePermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}

// UpdateUserRolesRequest 更新用户角色请求
type UpdateUserRolesRequest struct {
	RoleIDs []uint `json:"role_ids" binding:"required"`
}
