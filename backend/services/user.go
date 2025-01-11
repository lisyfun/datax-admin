package services

import (
	"datax-admin/models"
	"datax-admin/types"
	"datax-admin/utils"
	"errors"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct{}

// Register 用户注册
func (s *UserService) Register(req *types.RegisterRequest) error {
	var count int64
	if err := models.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}

	user := &models.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
	}
	if err := user.SetPassword(req.Password); err != nil {
		return err
	}

	return models.DB.Create(user).Error
}

// Login 用户登录
func (s *UserService) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
	var user models.User
	if err := models.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	if !utils.ComparePasswords(user.Password, req.Password) {
		return nil, errors.New("密码错误")
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	// 记录登录日志
	loginLog := models.LoginLog{
		UserID:    user.ID,
		Username:  user.Username,
		IP:        req.IP,
		LoginTime: time.Now(),
	}
	if err := models.DB.Create(&loginLog).Error; err != nil {
		log.Printf("记录登录日志失败: %v", err)
	}

	// 更新用户在线状态
	UpdateUserOnlineStatus(user.ID)

	return &types.LoginResponse{
		Token: token,
		User: types.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Email:    user.Email,
			Status:   user.Status,
		},
	}, nil
}

// Logout 用户登出
func (s *UserService) Logout(userID uint) error {
	// 移除用户在线状态
	RemoveUserOnlineStatus(userID)
	return nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID uint) (*types.UserResponse, error) {
	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &types.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Status:   user.Status,
	}, nil
}

// UpdatePassword 修改密码
func (s *UserService) UpdatePassword(userID uint, req *types.UpdatePasswordRequest) error {
	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		return err
	}

	if !user.CheckPassword(req.OldPassword) {
		return errors.New("原密码错误")
	}

	if err := user.SetPassword(req.NewPassword); err != nil {
		return err
	}

	return models.DB.Save(&user).Error
}

// UpdateProfile 更新个人信息
func (s *UserService) UpdateProfile(userID uint, req *types.UpdateProfileRequest) error {
	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	return models.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

// UpdateUserStatus 更新用户状态
func (s *UserService) UpdateUserStatus(userID uint, req *types.UpdateUserStatusRequest) error {
	// 先检查用户是否存在
	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 检查状态值是否有效
	if req.Status != nil && (*req.Status != 0 && *req.Status != 1) {
		return errors.New("无效的状态值")
	}

	// 更新状态
	return models.DB.Model(&user).Update("status", req.Status).Error
}

// GetUserList 获取用户列表
func (s *UserService) GetUserList(req *types.UserListRequest) (*types.UserListResponse, error) {
	var total int64
	var users []struct {
		models.User
		RoleNames string `gorm:"column:role_names"`
	}

	// 基础查询
	query := models.DB.Model(&models.User{}).
		Select("users.*, GROUP_CONCAT(DISTINCT roles.name SEPARATOR ',') as role_names").
		Joins("LEFT JOIN user_roles ON users.id = user_roles.user_id AND user_roles.deleted_at IS NULL").
		Joins("LEFT JOIN roles ON user_roles.role_id = roles.id AND roles.deleted_at IS NULL AND roles.status = 1").
		Group("users.id, users.username, users.password, users.nickname, users.email, users.avatar, users.status, users.created_at, users.updated_at, users.deleted_at")

	if req.Username != "" {
		query = query.Where("users.username LIKE ?", "%"+req.Username+"%")
	}
	if req.Status != nil {
		query = query.Where("users.status = ?", *req.Status)
	}

	// 获取总数（需要去掉 GROUP BY 和 SELECT）
	countQuery := models.DB.Model(&models.User{})
	if req.Username != "" {
		countQuery = countQuery.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Status != nil {
		countQuery = countQuery.Where("status = ?", *req.Status)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// 获取分页数据
	if err := query.Limit(req.PageSize).
		Offset((req.Page - 1) * req.PageSize).
		Find(&users).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]types.UserInfo, len(users))
	for i, user := range users {
		var roleList []types.RoleInfo
		if user.RoleNames != "" {
			roleNames := strings.Split(user.RoleNames, ",")
			roleList = make([]types.RoleInfo, len(roleNames))
			for j, name := range roleNames {
				roleList[j] = types.RoleInfo{
					Name: name,
				}
			}
		}

		list[i] = types.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Avatar:   user.Avatar,
			Status:   user.Status,
			Roles:    roleList,
		}
	}

	return &types.UserListResponse{
		Total: total,
		List:  list,
	}, nil
}

// ResetPassword 重置用户密码
func (s *UserService) ResetPassword(userID uint, req *types.ResetPasswordRequest) error {
	// 先检查用户是否存在
	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 设置新密码
	if err := user.SetPassword(req.Password); err != nil {
		return err
	}

	// 保存前检查密码是否已加密
	if len(user.Password) < 60 {
		return errors.New("密码加密失败")
	}

	if err := models.DB.Save(&user).Error; err != nil {
		return err
	}

	// 验证新密码是否可以正确验证
	if !user.CheckPassword(req.Password) {
		return errors.New("密码设置失败，无法验证")
	}

	return nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(userID uint) error {
	// 检查用户是否存在
	user := &models.User{}
	if err := models.DB.First(user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 软删除用户
	if err := models.DB.Delete(user).Error; err != nil {
		return err
	}

	return nil
}
