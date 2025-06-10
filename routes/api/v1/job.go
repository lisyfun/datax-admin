package v1

import (
	"datax-admin/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterJobRoutes 注册任务管理相关路由
func RegisterJobRoutes(authenticated *gin.RouterGroup) {
	jobController := controllers.NewJobController()
	
	authenticated.POST("/jobs", jobController.CreateJob)
	authenticated.PUT("/jobs/:id", jobController.UpdateJob)
	authenticated.DELETE("/jobs/:id", jobController.DeleteJob)
	authenticated.POST("/jobs/:id/start", jobController.StartJob)
	authenticated.POST("/jobs/:id/stop", jobController.StopJob)
	authenticated.POST("/jobs/:id/execute", jobController.ExecuteJob)
	authenticated.POST("/jobs/execute", jobController.ExecuteJobs)
	authenticated.GET("/jobs", jobController.GetJobList)
	authenticated.GET("/jobs/history", jobController.GetJobHistoryList)
	authenticated.POST("/jobs/history/clean", jobController.CleanJobHistory)
}