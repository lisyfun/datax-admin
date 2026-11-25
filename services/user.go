package services

import (
	"datax-admin/models"
	"datax-admin/types"
	"datax-admin/utils"
	"datax-admin/utils/logger"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct{}

// Register 用户注册
func (s *UserService) Register(req *types.RegisterRequest) error {
	var count int64
	// 使用参数化查询，防止SQL注入
	if err := models.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 对用户输入进行净化处理，防止XSLT和其他注入攻击
	safeUsername := strings.TrimSpace(req.Username)
	safeNickname := strings.TrimSpace(req.Nickname)
	safeEmail := strings.TrimSpace(req.Email)

	// 开启事务
	return models.DB.Transaction(func(tx *gorm.DB) error {
		user := &models.User{
			Username: safeUsername,
			Nickname: safeNickname,
			Email:    safeEmail,
		}
		if err := user.SetPassword(req.Password); err != nil {
			return err
		}

		// 创建用户
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 查找默认角色（普通用户角色）
		var defaultRole models.Role
		if err := tx.Where("code = ? AND status = 1", "user").First(&defaultRole).Error; err != nil {
			// 如果没有找到默认角色，尝试查找第一个可用角色
			if err := tx.Where("status = 1").First(&defaultRole).Error; err != nil {
				// 如果没有任何角色，跳过角色分配
				return nil
			}
		}

		// 为用户分配默认角色
		userRole := &models.UserRole{
			UserID: user.ID,
			RoleID: defaultRole.ID,
		}
		return tx.Create(userRole).Error
	})
}

// Login 用户登录
func (s *UserService) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
	// 防止XSLT注入攻击
	safeUsername := strings.TrimSpace(req.Username)
	safePassword := req.Password // 不修改密码，因为它会被hash比较

	// 验证用户输入是否包含恶意XML标签
	if strings.Contains(safeUsername, "<") {
		return nil, errors.New("用户名包含非法字符")
	}

	var user models.User
	// 使用参数化查询，防止SQL注入
	if err := models.DB.Where("username = ?", safeUsername).First(&user).Error; err != nil {
		// 不区分「用户不存在」和「密码错误」，避免用户名枚举
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status == 0 {
		return nil, errors.New("用户已被禁用")
	}

	if !utils.ComparePasswords(user.Password, safePassword) {
		return nil, errors.New("用户名或密码错误")
	}

	// 不再生成JWT token，改用session
	// token, err := utils.GenerateToken(user.ID, user.Username)
	// if err != nil {
	//	return nil, err
	// }

	// 记录登录日志
	loginLog := models.LoginLog{
		UserID:    user.ID,
		Username:  user.Username,
		IP:        req.IP,
		LoginTime: time.Now(),
	}
	if err := models.DB.Create(&loginLog).Error; err != nil {
		logger.Error("记录登录日志失败: %v", err)
	}

	// 更新用户在线状态
	UpdateUserOnlineStatus(user.ID)

	return &types.LoginResponse{
		// Token: token, // 不再返回token
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
	updates := map[string]any{}
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

	// 基础查询 - 使用数据库兼容的聚合函数
	var selectSQL string
	// GaussDB 与 PostgreSQL 语法一致（Dialector.Name() 返回 "gaussdb"）
	dbType := models.DB.Dialector.Name()
	if dbType == "postgres" || dbType == "gaussdb" {
		selectSQL = "users.*, STRING_AGG(DISTINCT roles.name, ',') as role_names"
	} else {
		selectSQL = "users.*, GROUP_CONCAT(DISTINCT roles.name SEPARATOR ',') as role_names"
	}

	query := models.DB.Model(&models.User{}).
		Select(selectSQL).
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

// CreateUser 管理员创建用户
func (s *UserService) CreateUser(req *types.CreateUserRequest) error {
	var count int64
	if err := models.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}

	return models.DB.Transaction(func(tx *gorm.DB) error {
		user := &models.User{
			Username: strings.TrimSpace(req.Username),
			Nickname: strings.TrimSpace(req.Nickname),
			Email:    strings.TrimSpace(req.Email),
			Status:   1,
		}
		if req.Status != nil {
			user.Status = *req.Status
		}
		if err := user.SetPassword(req.Password); err != nil {
			return err
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		roleIDs := req.RoleIDs
		if len(roleIDs) == 0 {
			var defaultRole models.Role
			if err := tx.Where("code = ? AND status = 1", "user").First(&defaultRole).Error; err == nil {
				roleIDs = []uint{defaultRole.ID}
			}
		}

		for _, roleID := range roleIDs {
			var roleCount int64
			if err := tx.Model(&models.Role{}).Where("id = ? AND status = 1", roleID).Count(&roleCount).Error; err != nil {
				return err
			}
			if roleCount == 0 {
				return fmt.Errorf("角色不存在或已禁用: %d", roleID)
			}
			if err := tx.Create(&models.UserRole{UserID: user.ID, RoleID: roleID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateUser 管理员更新用户资料/状态
func (s *UserService) UpdateUser(userID uint, req *types.UpdateUserRequest) error {
	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	updates := map[string]any{}
	if req.Nickname != nil {
		updates["nickname"] = strings.TrimSpace(*req.Nickname)
	}
	if req.Email != nil {
		updates["email"] = strings.TrimSpace(*req.Email)
	}
	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) == 0 {
		return nil
	}
	return models.DB.Model(&user).Updates(updates).Error
}

// IsLastEnabledAdmin 判断用户是否为最后一个启用状态的管理员
func (s *UserService) IsLastEnabledAdmin(userID uint) (bool, error) {
	var isAdmin int64
	err := models.DB.Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND user_roles.deleted_at IS NULL AND roles.deleted_at IS NULL AND roles.status = 1 AND roles.code = ?", userID, "admin").
		Count(&isAdmin).Error
	if err != nil || isAdmin == 0 {
		return false, err
	}

	var adminCount int64
	err = models.DB.Table("users").
		Joins("JOIN user_roles ON user_roles.user_id = users.id AND user_roles.deleted_at IS NULL").
		Joins("JOIN roles ON roles.id = user_roles.role_id AND roles.deleted_at IS NULL AND roles.status = 1").
		Where("roles.code = ? AND users.status = 1 AND users.deleted_at IS NULL", "admin").
		Distinct("users.id").
		Count(&adminCount).Error
	if err != nil {
		return false, err
	}
	return adminCount <= 1, nil
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
	// 开启事务
	tx := models.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查用户是否存在
	user := &models.User{}
	if err := tx.First(user, userID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 删除用户角色关联
	if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 软删除用户
	if err := tx.Delete(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
