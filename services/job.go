package services

import (
	"datax-admin/models"
	"datax-admin/types"
	"datax-admin/utils/cron"
	"datax-admin/utils/logger"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

// JobService 任务服务
type JobService struct {
	scheduler *cron.Scheduler
}

var (
	jobServiceInstance *JobService
	jobServiceOnce     sync.Once
)

// GetJobService 获取任务服务的单例实例
func GetJobService() *JobService {
	jobServiceOnce.Do(func() {
		jobServiceInstance = &JobService{
			scheduler: cron.NewScheduler(),
		}
		jobServiceInstance.scheduler.Start() // 启动调度器

		// 记录调度器启动
		logger.Info("任务调度器已启动")
	})
	return jobServiceInstance
}

// NewJobService 创建任务服务 (为保持兼容性保留)
func NewJobService() *JobService {
	return GetJobService()
}

// InitJobScheduler 初始化任务调度器，加载所有运行中的任务
func (s *JobService) InitJobScheduler() error {
	var jobs []models.Job

	// 获取所有状态为"运行中"的任务
	if err := models.DB.Where("status = ?", models.JobStatusRunning).Find(&jobs).Error; err != nil {
		return fmt.Errorf("加载运行中的任务失败: %v", err)
	}

	logger.Info("正在加载 %d 个运行中的任务...", len(jobs))

	// 将所有运行中的任务添加到调度器
	for _, job := range jobs {
		jobCopy := job // 创建副本避免闭包问题

		// 添加到调度器
		if err := s.scheduler.AddJob(
			fmt.Sprintf("job_%d", jobCopy.ID),
			jobCopy.CronExpr,
			func() {
				s.executeJob(&jobCopy)
			},
		); err != nil {
			logger.Info("添加任务 [%s] (ID: %d) 到调度器失败: %v", job.Name, job.ID, err)
			// 如果添加失败，将任务状态更新为停止
			if updateErr := models.DB.Model(&jobCopy).Update("status", models.JobStatusStop).Error; updateErr != nil {
				logger.Info("更新任务状态失败: %v", updateErr)
			}
		} else {
			logger.Info("成功加载任务 [%s] (ID: %d), Cron表达式: %s", job.Name, job.ID, job.CronExpr)
		}
	}

	return nil
}

// CreateJob 创建任务
func (s *JobService) CreateJob(req *types.CreateJobRequest, userID uint) error {
	// 验证Cron表达式
	if !cron.IsValidCronExpr(req.CronExpr) {
		return errors.New("无效的Cron表达式")
	}

	// 验证并序列化任务参数
	params, err := s.validateAndSerializeParams(req.Type, req.Params)
	if err != nil {
		return err
	}

	job := &models.Job{
		Name:        req.Name,
		Type:        models.JobType(req.Type),
		Description: req.Description,
		CronExpr:    req.CronExpr,
		Timeout:     req.Timeout,
		RetryCount:  req.RetryCount,
		RetryDelay:  req.RetryDelay,
		Params:      params,
		Creator:     userID,
		Status:      models.JobStatusStop,
	}

	return models.DB.Create(job).Error
}

// UpdateJob 更新任务
func (s *JobService) UpdateJob(jobID uint, req *types.UpdateJobRequest, userID uint) error {
	var job models.Job
	if err := models.DB.First(&job, jobID).Error; err != nil {
		return err
	}

	// 如果任务正在运行，不允许更新
	if job.Status == models.JobStatusRunning {
		return errors.New("任务正在运行，无法更新")
	}

	// 验证Cron表达式
	if req.CronExpr != "" && !cron.IsValidCronExpr(req.CronExpr) {
		return errors.New("无效的Cron表达式")
	}

	// 验证并序列化任务参数
	var params string
	if req.Params != nil {
		var err error
		params, err = s.validateAndSerializeParams(string(job.Type), req.Params)
		if err != nil {
			return err
		}
	}

	updates := map[string]any{
		"name":        req.Name,
		"description": req.Description,
		"cron_expr":   req.CronExpr,
		"timeout":     req.Timeout,
		"retry_count": req.RetryCount,
		"retry_delay": req.RetryDelay,
		"updater":     userID,
		"updated_at":  time.Now(),
	}

	if req.Params != nil {
		updates["params"] = params
	}

	return models.DB.Model(&job).Updates(updates).Error
}

// DeleteJob 删除任务
func (s *JobService) DeleteJob(jobID uint) error {
	var job models.Job
	if err := models.DB.First(&job, jobID).Error; err != nil {
		return err
	}

	// 如果任务正在运行，不允许删除
	if job.Status == models.JobStatusRunning {
		return errors.New("任务正在运行，无法删除")
	}

	// 停止定时任务
	s.scheduler.Remove(fmt.Sprintf("job_%d", jobID))

	return models.DB.Delete(&job).Error
}

// StartJob 启动任务
func (s *JobService) StartJob(jobID uint) error {
	var job models.Job
	if err := models.DB.First(&job, jobID).Error; err != nil {
		return err
	}

	if job.Status == models.JobStatusRunning {
		return errors.New("任务已在运行")
	}

	if job.Status == models.JobStatusDisabled {
		return errors.New("任务已禁用")
	}

	// 创建任务副本，避免闭包引用问题
	jobCopy := job

	// 添加到调度器
	if err := s.scheduler.AddJob(fmt.Sprintf("job_%d", jobCopy.ID), jobCopy.CronExpr, func() {
		s.executeJob(&jobCopy)
	}); err != nil {
		return err
	}

	logger.Info("任务 [%s] (ID: %d) 已添加到调度器, Cron表达式: %s\n", job.Name, job.ID, job.CronExpr)

	// 更新状态
	return models.DB.Model(&job).Update("status", models.JobStatusRunning).Error
}

// StopJob 停止任务
func (s *JobService) StopJob(jobID uint) error {
	var job models.Job
	if err := models.DB.First(&job, jobID).Error; err != nil {
		return err
	}

	if job.Status != models.JobStatusRunning {
		return errors.New("任务未在运行")
	}

	// 从调度器中移除
	s.scheduler.Remove(fmt.Sprintf("job_%d", job.ID))

	// 更新状态
	return models.DB.Model(&job).Update("status", models.JobStatusStop).Error
}

// GetJobList 获取任务列表
func (s *JobService) GetJobList(req *types.JobListRequest) (*types.JobListResponse, error) {
	var total int64
	var jobs []models.Job

	query := models.DB.Model(&models.Job{})

	// 添加查询条件
	if req.Keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 获取分页数据
	if err := query.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).
		Order("id DESC").Find(&jobs).Error; err != nil {
		return nil, err
	}

	return &types.JobListResponse{
		Total: total,
		Items: jobs,
	}, nil
}

// GetJobHistoryList 获取任务执行历史列表
func (s *JobService) GetJobHistoryList(req *types.JobHistoryListRequest) (*types.JobHistoryListResponse, error) {
	var items []struct {
		models.JobHistory
		JobName string `json:"job_name"`
	}

	query := models.DB.Model(&models.JobHistory{}).
		Select("job_histories.*, jobs.name as job_name").
		Joins("LEFT JOIN jobs ON jobs.id = job_histories.job_id")

	// 构建查询条件
	if req.JobID != nil {
		query = query.Where("job_histories.job_id = ?", *req.JobID)
	}
	if req.Status != nil {
		query = query.Where("job_histories.status = ?", *req.Status)
	}
	if req.Keyword != "" {
		query = query.Where("jobs.name LIKE ?", "%"+req.Keyword+"%")
	}
	if !req.StartTime.IsZero() {
		query = query.Where("start_time >= ?", req.StartTime)
	}
	if !req.EndTime.IsZero() {
		query = query.Where("end_time <= ?", req.EndTime)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 添加排序和分页
	if err := query.Order("job_histories.id DESC").
		Limit(req.PageSize).
		Offset((req.Page - 1) * req.PageSize).
		Find(&items).Error; err != nil {
		return nil, err
	}

	// 转换为 JobHistory 切片
	histories := make([]models.JobHistory, len(items))
	for i, item := range items {
		histories[i] = item.JobHistory
		histories[i].JobName = item.JobName
	}

	return &types.JobHistoryListResponse{
		Total: total,
		Items: histories,
	}, nil
}

// validateAndSerializeParams 验证并序列化任务参数
func (s *JobService) validateAndSerializeParams(jobType string, params any) (string, error) {
	var err error
	switch models.JobType(jobType) {
	case models.JobTypeShell:
		var shellParams models.JobShellParams
		if err = mapToStruct(params, &shellParams); err != nil {
			return "", err
		}
		if shellParams.Command == "" {
			return "", errors.New("shell命令不能为空")
		}
		data, err := json.Marshal(shellParams)
		if err != nil {
			return "", err
		}
		return string(data), nil
	case models.JobTypeHTTP:
		var httpParams models.JobHTTPParams
		if err = mapToStruct(params, &httpParams); err != nil {
			return "", err
		}
		if httpParams.URL == "" {
			return "", errors.New("HTTP URL不能为空")
		}
		if httpParams.Method == "" {
			httpParams.Method = "GET"
		}
		data, err := json.Marshal(httpParams)
		if err != nil {
			return "", err
		}
		return string(data), nil
	case models.JobTypeDataX:
		var dataxParams models.JobDataXParams
		if err = mapToStruct(params, &dataxParams); err != nil {
			return "", err
		}
		if dataxParams.JobConfig == "" {
			return "", errors.New("DataX任务配置不能为空")
		}
		data, err := json.Marshal(dataxParams)
		if err != nil {
			return "", err
		}
		return string(data), nil
	default:
		return "", errors.New("不支持的任务类型")
	}
}

// executeJob 执行任务
func (s *JobService) executeJob(job *models.Job) {
	logger.Info("开始执行任务: [%s] (ID: %d)", job.Name, job.ID)

	history := &models.JobHistory{
		JobID:     job.ID,
		StartTime: time.Now(),
	}

	// 执行任务并记录结果
	defer func() {
		// 添加panic恢复
		if r := recover(); r != nil {
			history.Status = 0
			history.Error = fmt.Sprintf("任务执行panic: %v", r)
			logger.Info("任务执行panic: [%s] (ID: %d): %v", job.Name, job.ID, r)
		}

		history.EndTime = time.Now()
		history.Duration = history.EndTime.Sub(history.StartTime).Milliseconds()

		logger.Info("任务执行完成: [%s] (ID: %d), 状态: %d, 耗时: %dms",
			job.Name, job.ID, history.Status, history.Duration)

		// 添加错误处理和重试逻辑
		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			if err := models.DB.Create(history).Error; err != nil {
				if i == maxRetries-1 {
					// 如果是最后一次重试，记录错误
					logger.Info("保存任务历史记录失败 [%s] (ID: %d): %v", job.Name, job.ID, err)
				}
				// 等待一小段时间后重试
				time.Sleep(time.Second * time.Duration(i+1))
				continue
			}
			break
		}
	}()

	var params any
	if err := json.Unmarshal([]byte(job.Params), &params); err != nil {
		history.Status = 0
		history.Error = fmt.Sprintf("解析任务参数失败: %v", err)
		return
	}

	// 根据任务类型执行不同的逻辑
	switch job.Type {
	case models.JobTypeShell:
		s.executeShellJob(job, params, history)
	case models.JobTypeHTTP:
		s.executeHTTPJob(job, params, history)
	case models.JobTypeDataX:
		s.executeDataXJob(job, params, history)
	}
}

// 辅助函数：将map转换为结构体
func mapToStruct(in any, out any) error {
	bytes, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, out)
}

// ExecuteJob 手动执行任务
func (s *JobService) ExecuteJob(jobID uint) error {
	var job models.Job
	if err := models.DB.First(&job, jobID).Error; err != nil {
		return err
	}

	// 检查任务状态
	if job.Status == models.JobStatusDisabled {
		return errors.New("任务已禁用")
	}

	// 异步执行任务
	go s.executeJob(&job)

	return nil
}

// CleanJobHistory 清理任务历史
func (s *JobService) CleanJobHistory(beforeTime time.Time) error {
	if beforeTime.IsZero() {
		// 清理全部历史记录
		return models.DB.Where("1 = 1").Delete(&models.JobHistory{}).Error
	}
	return models.DB.Where("created_at < ?", beforeTime).Delete(&models.JobHistory{}).Error
}
