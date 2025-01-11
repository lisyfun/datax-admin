package types

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	UserCount       int64 `json:"userCount"`       // 用户总数
	RoleCount       int64 `json:"roleCount"`       // 角色总数
	PermissionCount int64 `json:"permissionCount"` // 权限总数
	OnlineCount     int64 `json:"onlineCount"`     // 在线用户数
	SuccessCount    int64 `json:"successCount"`    // 成功任务数
	FailedCount     int64 `json:"failedCount"`     // 失败任务数
}

// JobExecutionTrend 任务执行趋势
type JobExecutionTrend struct {
	Date         string `json:"date"`         // 日期
	SuccessCount int64  `json:"successCount"` // 成功数
	FailedCount  int64  `json:"failedCount"`  // 失败数
}

// RecentLogin 最近登录记录
type RecentLogin struct {
	Username  string `json:"username"`  // 用户名
	LoginTime string `json:"loginTime"` // 登录时间
	IP        string `json:"ip"`        // IP地址
}

// SystemInfo 系统信息
type SystemInfo struct {
	SystemName   string `json:"systemName"`   // 系统名称
	Version      string `json:"version"`      // 系统版本
	OS           string `json:"os"`           // 操作系统
	GoVersion    string `json:"goVersion"`    // Go版本
	DBVersion    string `json:"dbVersion"`    // 数据库版本
	Uptime       string `json:"uptime"`       // 系统运行时间
	CPUUsage     string `json:"cpuUsage"`     // CPU使用率
	MemoryTotal  string `json:"memoryTotal"`  // 内存总量
	MemoryUsed   string `json:"memoryUsed"`   // 已用内存
	MemoryUsage  string `json:"memoryUsage"`  // 内存使用率
	DiskTotal    string `json:"diskTotal"`    // 磁盘总量
	DiskUsed     string `json:"diskUsed"`     // 已用磁盘
	DiskUsage    string `json:"diskUsage"`    // 磁盘使用率
	NumGoroutine int    `json:"numGoroutine"` // Goroutine数量
	NumCPU       int    `json:"numCPU"`       // CPU核心数
}

// DashboardResponse 仪表盘响应数据
type DashboardResponse struct {
	Stats        DashboardStats      `json:"stats"`        // 统计数据
	RecentLogins []RecentLogin       `json:"recentLogins"` // 最近登录记录
	SystemInfo   SystemInfo          `json:"systemInfo"`   // 系统信息
	TrendData    []JobExecutionTrend `json:"trendData"`    // 趋势数据
}
