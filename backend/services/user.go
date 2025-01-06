package services

import (
	"datax-admin/middleware"
	"datax-admin/models"
	"datax-admin/types"
	"errors"

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
func (s *UserService) Login(req *types.LoginRequest) (*types.TokenResponse, error) {
	var user models.User
	if err := models.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 检查密码格式
	if len(user.Password) < 60 {
		return nil, errors.New("密码格式错误")
	}

	if !user.CheckPassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status == 0 {
		return nil, errors.New("账号已被禁用")
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &types.TokenResponse{Token: token}, nil
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
	if req.Status == nil {
		return errors.New("状态值不能为空")
	}
	if *req.Status != 0 && *req.Status != 1 {
		return errors.New("无效的状态值")
	}

	// 更新状态
	return models.DB.Model(&user).Update("status", *req.Status).Error
}

// GetUserList 获取用户列表
func (s *UserService) GetUserList(req *types.UserListRequest) (*types.UserListResponse, error) {
	var total int64
	var users []models.User

	query := models.DB.Model(&models.User{})
	if req.Keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 获取分页数据
	if err := query.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).Find(&users).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	items := make([]types.UserResponse, len(users))
	for i, user := range users {
		items[i] = types.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Avatar:   user.Avatar,
			Status:   user.Status,
		}
	}

	return &types.UserListResponse{
		Total: total,
		Items: items,
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
