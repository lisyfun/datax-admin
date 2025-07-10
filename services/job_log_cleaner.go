package services

import (
	"datax-admin/models"
	"datax-admin/utils/logger"
	"os"
	"path/filepath"
	"time"
)

// JobLogCleaner 任务日志清理器
type JobLogCleaner struct {
	maxAge time.Duration // 日志保留时间
}

// NewJobLogCleaner 创建日志清理器
func NewJobLogCleaner(maxAge time.Duration) *JobLogCleaner {
	return &JobLogCleaner{
		maxAge: maxAge,
	}
}

// CleanOldLogs 清理过期的日志文件
func (c *JobLogCleaner) CleanOldLogs() error {
	logger.Info("开始清理过期的任务日志文件...")
	
	// 计算过期时间
	expireTime := time.Now().Add(-c.maxAge)
	
	// 查询过期的历史记录
	var histories []models.JobHistory
	if err := models.DB.Where("created_at < ?", expireTime).Find(&histories).Error; err != nil {
		return err
	}
	
	deletedCount := 0
	errorCount := 0
	
	// 删除对应的日志文件
	for _, history := range histories {
		if history.LogPath != "" {
			if err := os.Remove(history.LogPath); err != nil {
				if !os.IsNotExist(err) {
					logger.Info("删除日志文件失败: %s, 错误: %v", history.LogPath, err)
					errorCount++
				}
			} else {
				deletedCount++
			}
		}
	}
	
	// 删除数据库记录
	result := models.DB.Where("created_at < ?", expireTime).Delete(&models.JobHistory{})
	if result.Error != nil {
		return result.Error
	}
	
	logger.Info("日志清理完成: 删除了 %d 个日志文件, %d 条数据库记录, %d 个错误", 
		deletedCount, result.RowsAffected, errorCount)
	
	return nil
}

// CleanOrphanedLogs 清理孤立的日志文件（数据库中没有对应记录的文件）
func (c *JobLogCleaner) CleanOrphanedLogs() error {
	logger.Info("开始清理孤立的任务日志文件...")
	
	logDir := filepath.Join("logs", "jobs")
	
	// 检查日志目录是否存在
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		logger.Info("日志目录不存在，跳过清理")
		return nil
	}
	
	// 获取所有日志文件路径
	var dbLogPaths []string
	if err := models.DB.Model(&models.JobHistory{}).
		Where("log_path != ''").
		Pluck("log_path", &dbLogPaths).Error; err != nil {
		return err
	}
	
	// 创建数据库中存在的文件路径映射
	dbPathMap := make(map[string]bool)
	for _, path := range dbLogPaths {
		dbPathMap[path] = true
	}
	
	deletedCount := 0
	errorCount := 0
	
	// 遍历日志目录中的所有文件
	err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// 跳过目录
		if info.IsDir() {
			return nil
		}
		
		// 检查文件是否在数据库中存在
		if !dbPathMap[path] {
			// 检查文件是否过期（超过7天的孤立文件）
			if time.Since(info.ModTime()) > 7*24*time.Hour {
				if err := os.Remove(path); err != nil {
					logger.Info("删除孤立日志文件失败: %s, 错误: %v", path, err)
					errorCount++
				} else {
					deletedCount++
				}
			}
		}
		
		return nil
	})
	
	if err != nil {
		return err
	}
	
	logger.Info("孤立日志清理完成: 删除了 %d 个孤立文件, %d 个错误", deletedCount, errorCount)
	
	return nil
}

// StartAutoCleanup 启动自动清理（每天凌晨2点执行）
func (c *JobLogCleaner) StartAutoCleanup() {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		
		// 立即执行一次清理
		if err := c.CleanOldLogs(); err != nil {
			logger.Info("自动清理任务日志失败: %v", err)
		}
		
		if err := c.CleanOrphanedLogs(); err != nil {
			logger.Info("自动清理孤立日志失败: %v", err)
		}
		
		for {
			select {
			case <-ticker.C:
				// 检查是否是凌晨2点
				now := time.Now()
				if now.Hour() == 2 {
					if err := c.CleanOldLogs(); err != nil {
						logger.Info("自动清理任务日志失败: %v", err)
					}
					
					if err := c.CleanOrphanedLogs(); err != nil {
						logger.Info("自动清理孤立日志失败: %v", err)
					}
				}
			}
		}
	}()
}
