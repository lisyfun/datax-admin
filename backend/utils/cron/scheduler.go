package cron

import (
	"sync"

	"github.com/robfig/cron/v3"
)

// Scheduler 任务调度器
type Scheduler struct {
	cron     *cron.Cron
	jobIDs   map[string]cron.EntryID
	jobMutex sync.RWMutex
}

// NewScheduler 创建调度器
func NewScheduler() *Scheduler {
	return &Scheduler{
		cron:   cron.New(cron.WithSeconds()),
		jobIDs: make(map[string]cron.EntryID),
	}
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.cron.Start()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cron.Stop()
}

// AddJob 添加任务
func (s *Scheduler) AddJob(name string, spec string, cmd func()) error {
	s.jobMutex.Lock()
	defer s.jobMutex.Unlock()

	// 如果任务已存在，先移除
	if _, exists := s.jobIDs[name]; exists {
		s.cron.Remove(s.jobIDs[name])
		delete(s.jobIDs, name)
	}

	// 添加新任务
	id, err := s.cron.AddFunc(spec, cmd)
	if err != nil {
		return err
	}

	s.jobIDs[name] = id
	return nil
}

// Remove 移除任务
func (s *Scheduler) Remove(name string) {
	s.jobMutex.Lock()
	defer s.jobMutex.Unlock()

	if id, exists := s.jobIDs[name]; exists {
		s.cron.Remove(id)
		delete(s.jobIDs, name)
	}
}

// IsValidCronExpr 验证Cron表达式
func IsValidCronExpr(expr string) bool {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(expr)
	return err == nil
}
