package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext 从 gin context 中获取当前登录用户的 ID
func GetUserIDFromContext(c *gin.Context) (uint, error) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("未找到用户ID")
	}

	// 类型断言
	id, ok := userID.(uint)
	if !ok {
		return 0, errors.New("用户ID类型错误")
	}

	return id, nil
}
