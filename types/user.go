package types

// UserInfo 用户信息
type UserInfo struct {
	ID       uint       `json:"id"`
	Username string     `json:"username"`
	Nickname string     `json:"nickname"`
	Avatar   string     `json:"avatar"`
	Email    string     `json:"email"`
	Status   int        `json:"status"`
	Roles    []RoleInfo `json:"roles"`
}

// RoleInfo 角色信息
type RoleInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	// Token string   `json:"token"` // 改用session认证，不再需要返回token
	User UserInfo `json:"user"`
}

// UpdatePasswordRequest 更新密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// UpdateProfileRequest 更新个人信息请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

// UpdateUserStatusRequest 更新用户状态请求
type UpdateUserStatusRequest struct {
	Status *int `json:"status" binding:"required,oneof=0 1"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"pageSize" binding:"required,min=1,max=100"`
	Username string `form:"username"`
	Status   *int   `form:"status"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Total int64      `json:"total"`
	List  []UserInfo `json:"list"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6,max=100"`
}

// CreateUserRequest 管理员创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50,excludesall=<>"`
	Password string `json:"password" binding:"required,min=6,max=100,excludesall=<>"`
	Nickname string `json:"nickname" binding:"omitempty,max=50,excludesall=<>"`
	Email    string `json:"email" binding:"omitempty,email,max=100,excludesall=<>"`
	Status   *int   `json:"status" binding:"omitempty,oneof=0 1"`
	RoleIDs  []uint `json:"roleIds"`
}

// UpdateUserRequest 管理员更新用户请求。
// 使用指针以区分「不修改」和「清空」。
type UpdateUserRequest struct {
	Nickname *string `json:"nickname" binding:"omitempty,max=50"`
	Email    *string `json:"email" binding:"omitempty,max=100"`
	Avatar   *string `json:"avatar" binding:"omitempty,max=255"`
	Status   *int    `json:"status" binding:"omitempty,oneof=0 1"`
}
