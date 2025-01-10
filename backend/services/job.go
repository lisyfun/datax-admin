package services

import (
	"datax-admin/models"
	"datax-admin/types"
	"datax-admin/utils/cron"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// JobService 任务服务
type JobService struct {
	scheduler *cron.Scheduler
}

// NewJobService 创建任务服务
func NewJobService() *JobService {
	scheduler := cron.NewScheduler()
	scheduler.Start() // 启动调度器
	return &JobService{
		scheduler: scheduler,
	}
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

	updates := map[string]interface{}{
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

	// 添加到调度器
	if err := s.scheduler.AddJob(fmt.Sprintf("job_%d", job.ID), job.CronExpr, func() {
		s.executeJob(&job)
	}); err != nil {
		return err
	}

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
	// 使用原生SQL查询以便于调试
	query := `
		SELECT
			h.*,
			j.name as job_name
		FROM job_histories h
		LEFT JOIN jobs j ON j.id = h.job_id
		WHERE 1=1
	`
	var params []interface{}

	// 构建查询条件
	if req.JobID != nil {
		query += " AND h.job_id = ?"
		params = append(params, *req.JobID)
	}
	if req.Status != nil {
		query += " AND h.status = ?"
		params = append(params, *req.Status)
	}
	if req.Keyword != "" {
		query += " AND j.name LIKE ?"
		params = append(params, "%"+req.Keyword+"%")
	}
	if !req.StartTime.IsZero() {
		query += " AND h.start_time >= ?"
		params = append(params, req.StartTime)
	}
	if !req.EndTime.IsZero() {
		query += " AND h.end_time <= ?"
		params = append(params, req.EndTime)
	}

	// 获取总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM (" + query + ") as t"
	if err := models.DB.Raw(countQuery, params...).Count(&total).Error; err != nil {
		return nil, err
	}

	// 添加排序和分页
	query += " ORDER BY h.id DESC LIMIT ? OFFSET ?"
	params = append(params, req.PageSize, (req.Page-1)*req.PageSize)

	// 执行查询
	var items []models.JobHistory
	if err := models.DB.Raw(query, params...).Scan(&items).Error; err != nil {
		return nil, err
	}

	return &types.JobHistoryListResponse{
		Total: total,
		Items: items,
	}, nil
}

// validateAndSerializeParams 验证并序列化任务参数
func (s *JobService) validateAndSerializeParams(jobType string, params interface{}) (string, error) {
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
	history := &models.JobHistory{
		JobID:     job.ID,
		StartTime: time.Now(),
	}

	// 执行任务并记录结果
	defer func() {
		history.EndTime = time.Now()
		history.Duration = history.EndTime.Sub(history.StartTime).Milliseconds()

		// 添加错误处理和重试逻辑
		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			if err := models.DB.Create(history).Error; err != nil {
				if i == maxRetries-1 {
					// 如果是最后一次重试，记录错误
					fmt.Printf("Failed to save job history after %d retries: %v\n", maxRetries, err)
				}
				// 等待一小段时间后重试
				time.Sleep(time.Second * time.Duration(i+1))
				continue
			}
			break
		}
	}()

	var params interface{}
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
func mapToStruct(in interface{}, out interface{}) error {
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
