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
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
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

	user.Password = req.NewPassword
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
	return models.DB.Model(&models.User{}).Where("id = ?", userID).Update("status", req.Status).Error
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
