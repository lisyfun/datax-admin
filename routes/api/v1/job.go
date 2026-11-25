package v1

import (
	"datax-admin/controllers"
	"datax-admin/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterJobRoutes 注册任务管理相关路由
func RegisterJobRoutes(authenticated *gin.RouterGroup) {
	jobController := controllers.NewJobController()

	authenticated.POST("/jobs", middleware.RequirePermission("job.list.create"), jobController.CreateJob)
	authenticated.PUT("/jobs/:id", middleware.RequirePermission("job.list.update"), jobController.UpdateJob)
	authenticated.DELETE("/jobs/:id", middleware.RequirePermission("job.list.delete"), jobController.DeleteJob)
	authenticated.POST("/jobs/:id/start", middleware.RequirePermission("job.list.execute"), jobController.StartJob)
	authenticated.POST("/jobs/:id/stop", middleware.RequirePermission("job.list.execute"), jobController.StopJob)
	authenticated.POST("/jobs/:id/execute", middleware.RequirePermission("job.list.execute"), jobController.ExecuteJob)
	authenticated.POST("/jobs/execute", middleware.RequirePermission("job.list.execute"), jobController.ExecuteJobs)
	authenticated.GET("/jobs", middleware.RequirePermission("job.list.query", "job.list", "job.list.create", "job.list.update", "job.list.delete", "job.list.execute"), jobController.GetJobList)
	authenticated.GET("/jobs/history", middleware.RequirePermission("job.history.query", "job.history"), jobController.GetJobHistoryList)
	authenticated.GET("/jobs/history/:id", middleware.RequirePermission("job.history.detail", "job.history.query", "job.history"), jobController.GetJobHistoryDetail)
	authenticated.POST("/jobs/history/clean", middleware.RequirePermission("job.history.query"), jobController.CleanJobHistory)
}
