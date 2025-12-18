package services

import (
	"datax-admin/config"
	"datax-admin/models"
	"datax-admin/types"
	"datax-admin/utils/cache"
	"datax-admin/utils/logger"
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
	dashboardSvc    *DashboardService
	initOnce        sync.Once
)

const (
	dashboardCacheKey = "dashboard_data"
	cacheDuration     = 5 * time.Minute
	updateInterval    = 1 * time.Minute // 数据更新间隔
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
type DashboardService struct {
	stopChan chan struct{}
}

// NewDashboardService 创建仪表盘服务
func NewDashboardService() *DashboardService {
	initOnce.Do(func() {
		dashboardSvc = &DashboardService{
			stopChan: make(chan struct{}),
		}
		go dashboardSvc.startBackgroundUpdate()
	})
	return dashboardSvc
}

// Stop 停止后台更新任务
func (s *DashboardService) Stop() {
	close(s.stopChan)
}

// startBackgroundUpdate 启动后台更新任务
func (s *DashboardService) startBackgroundUpdate() {
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	// 立即执行一次更新
	s.updateDashboardData()

	for {
		select {
		case <-ticker.C:
			s.updateDashboardData()
		case <-s.stopChan:
			return
		}
	}
}

// updateDashboardData 更新仪表盘数据
func (s *DashboardService) updateDashboardData() {
	var (
		stats        types.DashboardStats
		recentLogins []types.RecentLogin
		systemInfo   *types.SystemInfo
		trendData    []types.JobExecutionTrend
		wg           sync.WaitGroup
		errChan      = make(chan error, 4)
	)

	// 并发获取各项数据
	wg.Add(4)

	// 获取统计数据
	go func() {
		defer wg.Done()
		if err := s.getStats(&stats); err != nil {
			errChan <- err
		}
	}()

	// 获取最近登录记录
	go func() {
		defer wg.Done()
		var err error
		recentLogins, err = s.getRecentLogins()
		if err != nil {
			errChan <- err
		}
	}()

	// 获取系统信息
	go func() {
		defer wg.Done()
		var err error
		systemInfo, err = s.GetSystemInfo()
		if err != nil {
			errChan <- err
		}
	}()

	// 获取趋势数据
	go func() {
		defer wg.Done()
		var err error
		trendData, err = s.getJobExecutionTrend()
		if err != nil {
			errChan <- err
		}
	}()

	// 等待所有goroutine完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误
	for err := range errChan {
		// 只记录错误，不中断更新
		logger.Info("Dashboard update error: %v", err)
		return
	}

	// 构建响应数据
	response := &types.DashboardResponse{
		Stats:        stats,
		RecentLogins: recentLogins,
		SystemInfo:   *systemInfo,
		TrendData:    trendData,
	}

	// 更新缓存
	dashboardCache := cache.GetOrCreate[*types.DashboardResponse]("dashboard")
	dashboardCache.Set(dashboardCacheKey, response, cacheDuration)
}

// GetDashboardData 获取仪表盘数据
func (s *DashboardService) GetDashboardData() (*types.DashboardResponse, error) {
	// 从缓存获取数据
	dashboardCache := cache.GetOrCreate[*types.DashboardResponse]("dashboard")
	if data, exists := dashboardCache.Get(dashboardCacheKey); exists {
		return data, nil
	}

	// 如果缓存中没有数据，触发一次更新并等待
	s.updateDashboardData()

	// 再次尝试从缓存获取
	if data, exists := dashboardCache.Get(dashboardCacheKey); exists {
		return data, nil
	}

	return nil, fmt.Errorf("failed to get dashboard data")
}

// GetSystemInfo 获取系统信息
func (s *DashboardService) GetSystemInfo() (*types.SystemInfo, error) {
	// 获取数据库版本
	var version string
	versionQuery := getVersionQuery()
	if err := models.DB.Raw(versionQuery).Scan(&version).Error; err != nil {
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

// getVersionQuery 根据数据库类型生成版本查询语句
func getVersionQuery() string {
	dbType := config.GlobalConfig.Database.Type
	if dbType == "" {
		dbType = "mysql"
	}

	switch dbType {
	case "postgres", "postgresql":
		return "SELECT version()"
	default:
		return "SELECT VERSION()"
	}
}

// getJobExecutionTrend 获取最近7天的任务执行趋势
func (s *DashboardService) getJobExecutionTrend() ([]types.JobExecutionTrend, error) {
	now := time.Now()
	// 获取今天的开始时间
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// 6天前（共7天）
	startDate := today.AddDate(0, 0, -6)

	var histories []models.JobHistory
	// 只查询需要的字段
	if err := models.DB.Model(&models.JobHistory{}).
		Select("status, created_at").
		Where("created_at >= ?", startDate).
		Find(&histories).Error; err != nil {
		return nil, err
	}

	// 内存中聚合数据: 日期字符串 -> {成功数, 失败数}
	statsMap := make(map[string]struct{ Success, Failed int64 })

	for _, h := range histories {
		dateStr := h.CreatedAt.Format("2006-01-02")
		stat := statsMap[dateStr]
		if h.Status == 1 {
			stat.Success++
		} else if h.Status == 0 {
			stat.Failed++
		}
		statsMap[dateStr] = stat
	}

	// 构建最近7天的结果，确保日期连续
	var trends []types.JobExecutionTrend
	for i := 0; i < 7; i++ {
		date := startDate.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		stat := statsMap[dateStr]
		trends = append(trends, types.JobExecutionTrend{
			Date:         dateStr,
			SuccessCount: stat.Success,
			FailedCount:  stat.Failed,
		})
	}

	return trends, nil
}

// getStats 获取统计数据
func (s *DashboardService) getStats(stats *types.DashboardStats) error {
	var wg sync.WaitGroup
	var errChan = make(chan error, 4)

	wg.Add(4)

	// 获取用户总数
	go func() {
		defer wg.Done()
		var userCount int64
		if err := models.DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
			errChan <- err
			return
		}
		stats.UserCount = userCount
	}()

	// 获取角色总数
	go func() {
		defer wg.Done()
		var roleCount int64
		if err := models.DB.Model(&models.Role{}).Count(&roleCount).Error; err != nil {
			errChan <- err
			return
		}
		stats.RoleCount = roleCount
	}()

	// 获取权限总数
	go func() {
		defer wg.Done()
		var permissionCount int64
		if err := models.DB.Model(&models.Permission{}).Count(&permissionCount).Error; err != nil {
			errChan <- err
			return
		}
		stats.PermissionCount = permissionCount
	}()

	// 获取任务执行统计
	go func() {
		defer wg.Done()
		type StatusCount struct {
			Status int
			Count  int64
		}
		var results []StatusCount
		if err := models.DB.Model(&models.JobHistory{}).
			Select("status, count(*) as count").
			Group("status").
			Scan(&results).Error; err != nil {
			errChan <- err
			return
		}

		var successCount, failedCount int64
		for _, r := range results {
			if r.Status == 1 {
				successCount = r.Count
			} else if r.Status == 0 {
				failedCount = r.Count
			}
		}
		stats.SuccessCount = successCount
		stats.FailedCount = failedCount
	}()

	// 获取在线用户数
	stats.OnlineCount = GetOnlineUserCount()

	wg.Wait()
	close(errChan)

	// 检查是否有错误
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// getRecentLogins 获取最近登录记录
func (s *DashboardService) getRecentLogins() ([]types.RecentLogin, error) {
	var loginLogs []models.LoginLog
	if err := models.DB.Model(&models.LoginLog{}).
		Select("username, login_time, ip").
		Order("login_time desc").
		Limit(6).
		Find(&loginLogs).Error; err != nil {
		return nil, err
	}

	recentLogins := make([]types.RecentLogin, len(loginLogs))
	for i, log := range loginLogs {
		recentLogins[i] = types.RecentLogin{
			Username:  log.Username,
			LoginTime: log.LoginTime.Format("2006-01-02 15:04:05"),
			IP:        log.IP,
		}
	}

	return recentLogins, nil
}
