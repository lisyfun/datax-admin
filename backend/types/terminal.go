package types

import "time"

// CreateTerminalRequest 创建终端请求
type CreateTerminalRequest struct {
	Name     string `json:"name" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
}

// UpdateTerminalRequest 更新终端请求
type UpdateTerminalRequest struct {
	Name     string `json:"name" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password,omitempty"`
}

// TerminalListRequest 终端列表请求
type TerminalListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"pageSize" binding:"required,min=1,max=100"`
	Name     string `form:"name"`
	Host     string `form:"host"`
	Status   string `form:"status"`
}

// TerminalInfo 终端信息
type TerminalInfo struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Username  string    `json:"username"`
	Status    string    `json:"status"`
	LastSeen  time.Time `json:"lastSeen"`
	CreatedAt time.Time `json:"createdAt"`
}

// TerminalListResponse 终端列表响应
type TerminalListResponse struct {
	Total int64          `json:"total"`
	List  []TerminalInfo `json:"list"`
}
