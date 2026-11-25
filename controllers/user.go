package controllers

import (
	"datax-admin/config"
	"datax-admin/services"
	"datax-admin/types"
	"datax-admin/utils/loginlimit"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 同一 IP + 用户名在 15 分钟内最多失败 5 次；
// 同时限制单个用户名的失败总数，避免通过伪造 X-Forwarded-For 绕过。
var (
	loginLimiter     = loginlimit.New(5, 15*time.Minute)
	loginUserLimiter = loginlimit.New(20, 15*time.Minute)
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
	// 默认关闭公开注册：数据集成平台涉及服务器与数据库凭据，
	// 用户应由管理员通过「用户管理」创建。
	if !config.GlobalConfig.Auth.AllowRegister {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "系统未开放注册，请联系管理员创建账号"})
		return
	}

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

	// 获取客户端 IP：统一使用 gin 的 ClientIP，它只会按 server.trusted_proxies
	// 配置解析 X-Forwarded-For，避免直接信任可被伪造的请求头。
	req.IP = ctx.ClientIP()

	usernameKey := strings.ToLower(strings.TrimSpace(req.Username))
	limiterKey := req.IP + "|" + usernameKey
	if !loginLimiter.Allowed(limiterKey) || !loginUserLimiter.Allowed(usernameKey) {
		// 两个限流器可能各自触发，取剩余时间较长的那个展示
		remaining := loginLimiter.Remaining(limiterKey)
		if userRemaining := loginUserLimiter.Remaining(usernameKey); userRemaining > remaining {
			remaining = userRemaining
		}
		minutes := int(remaining.Minutes()) + 1
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"error": fmt.Sprintf("登录失败次数过多，请 %d 分钟后再试", minutes),
		})
		return
	}

	resp, err := c.userService.Login(&req)
	if err != nil {
		loginLimiter.Fail(limiterKey)
		loginUserLimiter.Fail(usernameKey)
		// 使用 400 而非 401：401 会被前端全局拦截为「登录已过期」并跳转
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loginLimiter.Reset(limiterKey)
	loginUserLimiter.Reset(usernameKey)

	// 设置Session
	session := sessions.Default(ctx)
	if session == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Session初始化失败"})
		return
	}

	// 检查session是否有效
	defer func() {
		if r := recover(); r != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Session操作失败"})
			return
		}
	}()

	session.Set("userID", resp.User.ID)
	session.Set("username", resp.User.Username)
	// 为当前会话设置更安全的 Cookie 选项
	session.Options(sessions.Options{
		Path:     "/",
		MaxAge:   config.GlobalConfig.Auth.Expiration,
		HttpOnly: true,
		Secure:   isHTTPS(ctx),
		SameSite: http.SameSiteLaxMode,
	})
	if err := session.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Session保存失败"})
		return
	}

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

	if uint(userID) == operatorID && req.Status != nil && *req.Status == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "不能禁用自己的账号"})
		return
	}
	if req.Status != nil && *req.Status == 0 {
		if last, checkErr := c.userService.IsLastEnabledAdmin(uint(userID)); checkErr == nil && last {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "系统至少需要保留一个启用的管理员"})
			return
		}
	}

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

// CreateUser 管理员创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req types.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.userService.CreateUser(&req); err != nil {
		services.LogOperation(
			ctx.GetUint("userID"), ctx.GetString("username"),
			"user", "create", fmt.Sprintf("创建用户失败: %s", req.Username),
			ctx.ClientIP(), ctx.GetHeader("User-Agent"),
			map[string]interface{}{"username": req.Username}, 0, err.Error(),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	services.LogOperation(
		ctx.GetUint("userID"), ctx.GetString("username"),
		"user", "create", fmt.Sprintf("创建用户: %s", req.Username),
		ctx.ClientIP(), ctx.GetHeader("User-Agent"),
		map[string]interface{}{"username": req.Username}, 1, "",
	)
	ctx.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

// UpdateUser 管理员更新用户资料
func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req types.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	operatorID := ctx.GetUint("userID")
	if uint(userID) == operatorID && req.Status != nil && *req.Status == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "不能禁用自己的账号"})
		return
	}
	if req.Status != nil && *req.Status == 0 {
		if last, checkErr := c.userService.IsLastEnabledAdmin(uint(userID)); checkErr == nil && last {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "系统至少需要保留一个启用的管理员"})
			return
		}
	}

	if err := c.userService.UpdateUser(uint(userID), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	services.LogOperation(
		operatorID, ctx.GetString("username"),
		"user", "update", fmt.Sprintf("更新用户资料，用户ID: %d", userID),
		ctx.ClientIP(), ctx.GetHeader("User-Agent"),
		map[string]interface{}{"user_id": userID}, 1, "",
	)
	ctx.JSON(http.StatusOK, gin.H{"message": "更新成功"})
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

	if uint(userID) == operatorID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "不能删除自己的账号"})
		return
	}
	if last, checkErr := c.userService.IsLastEnabledAdmin(uint(userID)); checkErr == nil && last {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "系统至少需要保留一个启用的管理员"})
		return
	}

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
	if session != nil {
		session.Clear()
		session.Options(sessions.Options{
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   isHTTPS(ctx),
			SameSite: http.SameSiteLaxMode,
		})
		_ = session.Save()
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// isHTTPS 判断当前请求是否 HTTPS 或经由代理为 HTTPS
func isHTTPS(c *gin.Context) bool {
	if c.Request.TLS != nil {
		return true
	}
	proto := c.GetHeader("X-Forwarded-Proto")
	return strings.EqualFold(proto, "https")
}
