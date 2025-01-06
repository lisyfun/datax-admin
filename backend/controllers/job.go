package controllers

import (
	"datax-admin/services"
	"datax-admin/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// JobController 任务控制器
type JobController struct {
	jobService *services.JobService
}

// NewJobController 创建任务控制器
func NewJobController() *JobController {
	return &JobController{
		jobService: services.NewJobService(),
	}
}

// CreateJob 创建任务
func (c *JobController) CreateJob(ctx *gin.Context) {
	var req types.CreateJobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint("userID")
	if err := c.jobService.CreateJob(&req, userID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "任务创建成功"})
}

// UpdateJob 更新任务
func (c *JobController) UpdateJob(ctx *gin.Context) {
	jobID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	var req types.UpdateJobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint("userID")
	if err := c.jobService.UpdateJob(uint(jobID), &req, userID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "任务更新成功"})
}

// DeleteJob 删除任务
func (c *JobController) DeleteJob(ctx *gin.Context) {
	jobID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	if err := c.jobService.DeleteJob(uint(jobID)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "任务删除成功"})
}

// StartJob 启动任务
func (c *JobController) StartJob(ctx *gin.Context) {
	jobID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	if err := c.jobService.StartJob(uint(jobID)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "任务启动成功"})
}

// StopJob 停止任务
func (c *JobController) StopJob(ctx *gin.Context) {
	jobID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	if err := c.jobService.StopJob(uint(jobID)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "任务停止成功"})
}

// GetJobList 获取任务列表
func (c *JobController) GetJobList(ctx *gin.Context) {
	var req types.JobListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.jobService.GetJobList(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetJobHistoryList 获取任务执行历史列表
func (c *JobController) GetJobHistoryList(ctx *gin.Context) {
	var req types.JobHistoryListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.jobService.GetJobHistoryList(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
