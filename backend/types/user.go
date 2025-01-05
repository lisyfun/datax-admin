package types

// UpdatePasswordRequest 修改密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// UpdateProfileRequest 更新个人信息请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	Avatar   string `json:"avatar" binding:"omitempty,max=255"`
}

// UpdateUserStatusRequest 更新用户状态请求
type UpdateUserStatusRequest struct {
	Status int `json:"status" binding:"required,oneof=0 1"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page     int    `form:"page,default=1" binding:"required,min=1"`
	PageSize int    `form:"page_size,default=10" binding:"required,min=5,max=100"`
	Keyword  string `form:"keyword" binding:"omitempty,max=50"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Total int64          `json:"total"`
	Items []UserResponse `json:"items"`
}
