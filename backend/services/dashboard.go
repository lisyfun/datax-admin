package services

import (
	"datax-admin/models"
	"datax-admin/types"
	"runtime"
	"time"
)

// DashboardService 仪表盘服务
type DashboardService struct{}

// NewDashboardService 创建仪表盘服务
func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

// GetDashboardData 获取仪表盘数据
func (s *DashboardService) GetDashboardData() (*types.DashboardResponse, error) {
	var stats types.DashboardStats

	// 获取用户总数
	var userCount int64
	if err := models.DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return nil, err
	}
	stats.UserCount = userCount

	// 获取角色总数
	var roleCount int64
	if err := models.DB.Model(&models.Role{}).Count(&roleCount).Error; err != nil {
		return nil, err
	}
	stats.RoleCount = roleCount

	// 获取权限总数
	var permissionCount int64
	if err := models.DB.Model(&models.Permission{}).Count(&permissionCount).Error; err != nil {
		return nil, err
	}
	stats.PermissionCount = permissionCount

	// TODO: 获取在线用户数（需要实现用户会话管理）
	stats.OnlineCount = 1

	// 获取任务执行统计
	var successCount, failedCount int64
	if err := models.DB.Model(&models.JobHistory{}).Where("status = ?", 1).Count(&successCount).Error; err != nil {
		return nil, err
	}
	if err := models.DB.Model(&models.JobHistory{}).Where("status = ?", 2).Count(&failedCount).Error; err != nil {
		return nil, err
	}
	stats.SuccessCount = successCount
	stats.FailedCount = failedCount

	// 获取最近登录记录
	var recentLogins []types.RecentLogin
	if err := models.DB.Model(&models.LoginLog{}).
		Select("username, login_time as loginTime, ip").
		Order("login_time desc").
		Limit(10).
		Scan(&recentLogins).Error; err != nil {
		return nil, err
	}

	// 获取系统信息
	systemInfo := types.SystemInfo{
		SystemName: "DATAX ADMIN",
		Version:    "v1.0.0",
		OS:         runtime.GOOS + " " + runtime.GOARCH,
		GoVersion:  runtime.Version(),
		DBVersion:  "MySQL 8.0", // TODO: 从数据库获取版本信息
	}

	// 获取最近7天的任务执行趋势
	trendData, err := s.getJobExecutionTrend()
	if err != nil {
		return nil, err
	}

	return &types.DashboardResponse{
		Stats:        stats,
		RecentLogins: recentLogins,
		SystemInfo:   systemInfo,
		TrendData:    trendData,
	}, nil
}

// getJobExecutionTrend 获取任务执行趋势
func (s *DashboardService) getJobExecutionTrend() ([]types.JobExecutionTrend, error) {
	// 获取最近7天的日期
	now := time.Now()
	startDate := now.AddDate(0, 0, -6).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	// 查询最近7天的任务执行记录
	var records []struct {
		Date         string
		SuccessCount int64
		FailedCount  int64
	}

	query := `
		SELECT
			DATE(start_time) as date,
			SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) as success_count,
			SUM(CASE WHEN status = 2 THEN 1 ELSE 0 END) as failed_count
		FROM job_histories
		WHERE start_time >= ? AND start_time < DATE_ADD(?, INTERVAL 1 DAY)
		GROUP BY DATE(start_time)
		ORDER BY date ASC
	`

	if err := models.DB.Raw(query, startDate, endDate).Scan(&records).Error; err != nil {
		return nil, err
	}

	// 构建趋势数据，确保每天都有数据
	trendData := make([]types.JobExecutionTrend, 7)
	recordMap := make(map[string]struct {
		SuccessCount int64
		FailedCount  int64
	})

	// 将查询结果转换为map
	for _, record := range records {
		recordMap[record.Date] = struct {
			SuccessCount int64
			FailedCount  int64
		}{
			SuccessCount: record.SuccessCount,
			FailedCount:  record.FailedCount,
		}
	}

	// 填充趋势数据
	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, -6+i).Format("2006-01-02")
		record, exists := recordMap[date]
		if !exists {
			record = struct {
				SuccessCount int64
				FailedCount  int64
			}{0, 0}
		}

		trendData[i] = types.JobExecutionTrend{
			Date:         date,
			SuccessCount: record.SuccessCount,
			FailedCount:  record.FailedCount,
		}
	}

	return trendData, nil
}
