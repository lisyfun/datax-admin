package types

// Page 通用分页响应
type Page[T any] struct {
	Total int64 `json:"total"`
	Items []T   `json:"items"`
}
