package middleware

import (
	"datax-admin/models"
	"datax-admin/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequirePermission 要求当前用户拥有任意一个指定权限码。
//
// 管理员角色（code=admin）直接放行，避免新增权限码后需要同步修正内置角色。
// 必须在 SessionAuth 之后使用。
func RequirePermission(codes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("userID")
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}

		allowed, err := userHasAnyPermission(userID, codes)
		if err != nil {
			logger.Error("权限校验失败: user=%d path=%s err=%v", userID, c.Request.URL.Path, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "权限校验失败"})
			c.Abort()
			return
		}
		if !allowed {
			logger.Info("拒绝越权访问: user=%d method=%s path=%s required=%v",
				userID, c.Request.Method, c.Request.URL.Path, codes)
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限执行此操作"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func userHasAnyPermission(userID uint, codes []string) (bool, error) {
	var isAdmin int64
	if err := models.DB.Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND user_roles.deleted_at IS NULL AND roles.deleted_at IS NULL AND roles.status = 1 AND roles.code = ?", userID, "admin").
		Count(&isAdmin).Error; err != nil {
		return false, err
	}
	if isAdmin > 0 {
		return true, nil
	}

	if len(codes) == 0 {
		return false, nil
	}

	var count int64
	err := models.DB.Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id AND role_permissions.deleted_at IS NULL").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id AND user_roles.deleted_at IS NULL").
		Joins("JOIN roles ON roles.id = user_roles.role_id AND roles.deleted_at IS NULL AND roles.status = 1").
		Where("user_roles.user_id = ? AND permissions.deleted_at IS NULL AND permissions.status = 1 AND permissions.code IN ?", userID, codes).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
