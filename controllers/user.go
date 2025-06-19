package controllers

import (
	"datax-admin/services"
	"datax-admin/types"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService *services.UserService
}

// NewUserController 创建用户控制器
func NewUserController() *UserController {
	return &UserController{
		userService: &services.UserService{},
	}
}

// Register 用户注册
func (c *UserController) Register(ctx *gin.Context) {
	var req types.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.userService.Register(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var req types.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取客户端 IP
	if ip := ctx.GetHeader("X-Real-IP"); ip != "" {
		req.IP = ip
	} else if ip := ctx.GetHeader("X-Forwarded-For"); ip != "" {
		req.IP = ip
	} else {
		req.IP = ctx.ClientIP()
	}

	resp, err := c.userService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置Session
	session := sessions.Default(ctx)
	session.Set("userID", resp.User.ID)
	session.Set("username", resp.User.Username)
	session.Save()

	ctx.JSON(http.StatusOK, resp)
}

// GetUserInfo 获取用户信息
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	resp, err := c.userService.GetUserInfo(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdatePassword 修改密码
func (c *UserController) UpdatePassword(ctx *gin.Context) {
	var req types.UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint("userID")
	if err := c.userService.UpdatePassword(userID, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// UpdateProfile 更新个人信息
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	var req types.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint("userID")
	if err := c.userService.UpdateProfile(userID, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "个人信息更新成功"})
}

// UpdateUserStatus 更新用户状态
func (c *UserController) UpdateUserStatus(ctx *gin.Context) {
	var req types.UpdateUserStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 记录操作日志
	operatorID := ctx.GetUint("userID")
	operatorName := ctx.GetString("username")

	if err := c.userService.UpdateUserStatus(uint(userID), &req); err != nil {
		// 记录失败日志
		services.LogOperation(
			operatorID,
			operatorName,
			"user",
			"update",
			fmt.Sprintf("更新用户状态失败，用户ID: %d", userID),
			ctx.ClientIP(),
			ctx.GetHeader("User-Agent"),
			req,
			0,
			err.Error(),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录成功日志
	services.LogOperation(
		operatorID,
		operatorName,
		"user",
		"update",
		fmt.Sprintf("更新用户状态成功，用户ID: %d，新状态: %d", userID, req.Status),
		ctx.ClientIP(),
		ctx.GetHeader("User-Agent"),
		req,
		1,
		"",
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "用户状态更新成功"})
}

// GetUserList 获取用户列表
func (c *UserController) GetUserList(ctx *gin.Context) {
	var req types.UserListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.userService.GetUserList(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// ResetPassword 重置用户密码
func (c *UserController) ResetPassword(ctx *gin.Context) {
	// 获取用户ID
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 绑定请求参数
	var req types.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录操作日志
	operatorID := ctx.GetUint("userID")
	operatorName := ctx.GetString("username")

	// 调用服务
	if err := c.userService.ResetPassword(uint(userID), &req); err != nil {
		// 记录失败日志
		services.LogOperation(
			operatorID,
			operatorName,
			"user",
			"update",
			fmt.Sprintf("重置用户密码失败，用户ID: %d", userID),
			ctx.ClientIP(),
			ctx.GetHeader("User-Agent"),
			map[string]interface{}{"user_id": userID},
			0,
			err.Error(),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 记录成功日志
	services.LogOperation(
		operatorID,
		operatorName,
		"user",
		"update",
		fmt.Sprintf("重置用户密码成功，用户ID: %d", userID),
		ctx.ClientIP(),
		ctx.GetHeader("User-Agent"),
		map[string]interface{}{"user_id": userID},
		1,
		"",
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	// 获取用户ID
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 记录操作日志
	operatorID := ctx.GetUint("userID")
	operatorName := ctx.GetString("username")

	// 调用服务删除用户
	if err := c.userService.DeleteUser(uint(userID)); err != nil {
		// 记录失败日志
		services.LogOperation(
			operatorID,
			operatorName,
			"user",
			"delete",
			fmt.Sprintf("删除用户失败，用户ID: %d", userID),
			ctx.ClientIP(),
			ctx.GetHeader("User-Agent"),
			map[string]interface{}{"user_id": userID},
			0,
			err.Error(),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 记录成功日志
	services.LogOperation(
		operatorID,
		operatorName,
		"user",
		"delete",
		fmt.Sprintf("删除用户成功，用户ID: %d", userID),
		ctx.ClientIP(),
		ctx.GetHeader("User-Agent"),
		map[string]interface{}{"user_id": userID},
		1,
		"",
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// Logout 用户登出
func (c *UserController) Logout(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	username := ctx.GetString("username")

	if err := c.userService.Logout(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 记录登出日志
	services.LogOperation(
		userID,
		username,
		"user",
		"logout",
		"用户登出",
		ctx.ClientIP(),
		ctx.GetHeader("User-Agent"),
		nil,
		1,
		"",
	)

	// 清除 Session
	session := sessions.Default(ctx)
	session.Clear()
	session.Options(sessions.Options{
		Path:     "/",
		MaxAge:   -1, // 立即过期
		HttpOnly: true,
	})
	_ = session.Save()

	ctx.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// 检查SQL注入的辅助函数
func containsSQLInjection(input string) bool {
	// 检查常见的SQL注入模式
	sqlPatterns := []string{
		"--", ";", "/*", "*/", "@@", "@",
		"char", "nchar", "varchar", "nvarchar",
		"alter", "begin", "cast", "create", "cursor",
		"declare", "delete", "drop", "end", "exec",
		"execute", "fetch", "insert", "kill", "select",
		"sys", "sysobjects", "syscolumns",
		"table", "update", "xp_",
		"or 1=1", "or 1=", "or 1 =",
		"union", "UNION", "HAVING",
		"1=1", "1 = 1",
	}

	for _, pattern := range sqlPatterns {
		if strings.Contains(strings.ToLower(input), strings.ToLower(pattern)) {
			return true
		}
	}

	// 检查XML/XSLT注入
	if strings.Contains(input, "<") && strings.Contains(input, ">") {
		return true
	}

	return false
}
