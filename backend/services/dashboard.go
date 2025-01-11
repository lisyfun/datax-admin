package services

import (
	"datax-admin/models"
	"datax-admin/types"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

var (
	onlineUsers     = make(map[uint]time.Time)
	onlineUsersLock sync.RWMutex
	startTime       = time.Now()
)

// UpdateUserOnlineStatus 更新用户在线状态
func UpdateUserOnlineStatus(userID uint) {
	onlineUsersLock.Lock()
	defer onlineUsersLock.Unlock()
	onlineUsers[userID] = time.Now()
}

// RemoveUserOnlineStatus 移除用户在线状态
func RemoveUserOnlineStatus(userID uint) {
	onlineUsersLock.Lock()
	defer onlineUsersLock.Unlock()
	delete(onlineUsers, userID)
}

// GetOnlineUserCount 获取在线用户数
func GetOnlineUserCount() int64 {
	onlineUsersLock.RLock()
	defer onlineUsersLock.RUnlock()

	now := time.Now()
	count := int64(0)

	// 清理超过30分钟未活动的用户
	for userID, lastActive := range onlineUsers {
		if now.Sub(lastActive) > 30*time.Minute {
			delete(onlineUsers, userID)
		} else {
			count++
		}
	}

	return count
}

// DashboardService 仪表盘服务
type DashboardService struct{}

// NewDashboardService 创建仪表盘服务
func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

// GetSystemInfo 获取系统信息
func (s *DashboardService) GetSystemInfo() (*types.SystemInfo, error) {
	// 获取数据库版本
	var version string
	if err := models.DB.Raw("SELECT VERSION()").Scan(&version).Error; err != nil {
		return nil, err
	}

	// 获取CPU使用率
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		cpuPercent = []float64{0}
	}

	// 获取内存使用情况
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		memInfo = &mem.VirtualMemoryStat{}
	}

	// 获取磁盘使用情况
	diskInfo, err := disk.Usage("/")
	if err != nil {
		diskInfo = &disk.UsageStat{}
	}

	// 计算系统运行时间
	uptime := time.Since(startTime)
	days := int(uptime.Hours() / 24)
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60

	return &types.SystemInfo{
		SystemName:   "DATAX ADMIN",
		Version:      "v1.0.0",
		OS:           runtime.GOOS + " " + runtime.GOARCH,
		GoVersion:    runtime.Version(),
		DBVersion:    version,
		Uptime:       fmt.Sprintf("%d天%d小时%d分钟", days, hours, minutes),
		CPUUsage:     fmt.Sprintf("%.2f%%", cpuPercent[0]),
		MemoryTotal:  fmt.Sprintf("%.2f GB", float64(memInfo.Total)/1024/1024/1024),
		MemoryUsed:   fmt.Sprintf("%.2f GB", float64(memInfo.Used)/1024/1024/1024),
		MemoryUsage:  fmt.Sprintf("%.2f%%", memInfo.UsedPercent),
		DiskTotal:    fmt.Sprintf("%.2f GB", float64(diskInfo.Total)/1024/1024/1024),
		DiskUsed:     fmt.Sprintf("%.2f GB", float64(diskInfo.Used)/1024/1024/1024),
		DiskUsage:    fmt.Sprintf("%.2f%%", diskInfo.UsedPercent),
		NumGoroutine: runtime.NumGoroutine(),
		NumCPU:       runtime.NumCPU(),
	}, nil
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

	// 获取在线用户数
	stats.OnlineCount = GetOnlineUserCount()

	// 获取任务执行统计
	var successCount, failedCount int64
	if err := models.DB.Model(&models.JobHistory{}).Where("status = ?", 1).Count(&successCount).Error; err != nil {
		return nil, err
	}
	if err := models.DB.Model(&models.JobHistory{}).Where("status = ?", 0).Count(&failedCount).Error; err != nil {
		return nil, err
	}
	stats.SuccessCount = successCount
	stats.FailedCount = failedCount

	// 获取最近登录记录
	var loginLogs []models.LoginLog
	if err := models.DB.Model(&models.LoginLog{}).
		Select("username, login_time, ip").
		Order("login_time desc").
		Limit(10).
		Find(&loginLogs).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	recentLogins := make([]types.RecentLogin, len(loginLogs))
	for i, log := range loginLogs {
		recentLogins[i] = types.RecentLogin{
			Username:  log.Username,
			LoginTime: log.LoginTime.Format("2006-01-02 15:04:05"),
			IP:        log.IP,
		}
	}

	// 获取系统信息
	systemInfo, err := s.GetSystemInfo()
	if err != nil {
		return nil, err
	}

	// 获取最近7天的任务执行趋势
	trendData, err := s.getJobExecutionTrend()
	if err != nil {
		return nil, err
	}

	return &types.DashboardResponse{
		Stats:        stats,
		RecentLogins: recentLogins,
		SystemInfo:   *systemInfo,
		TrendData:    trendData,
	}, nil
}

// getJobExecutionTrend 获取最近7天的任务执行趋势
func (s *DashboardService) getJobExecutionTrend() ([]types.JobExecutionTrend, error) {
	var trends []types.JobExecutionTrend

	// 获取最近7天的日期
	now := time.Now()
	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		startTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endTime := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

		// 获取成功任务数
		var successCount int64
		if err := models.DB.Model(&models.JobHistory{}).
			Where("status = ? AND created_at BETWEEN ? AND ?", 1, startTime, endTime).
			Count(&successCount).Error; err != nil {
			return nil, err
		}

		// 获取失败任务数
		var failedCount int64
		if err := models.DB.Model(&models.JobHistory{}).
			Where("status = ? AND created_at BETWEEN ? AND ?", 0, startTime, endTime).
			Count(&failedCount).Error; err != nil {
			return nil, err
		}

		trends = append(trends, types.JobExecutionTrend{
			Date:         date.Format("2006-01-02"),
			SuccessCount: successCount,
			FailedCount:  failedCount,
		})
	}

	return trends, nil
}
