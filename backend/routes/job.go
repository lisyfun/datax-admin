package routes

import (
	"datax-admin/controllers"
	"datax-admin/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterJobRoutes 注册任务相关路由
func RegisterJobRoutes(r *gin.RouterGroup) {
	jobController := controllers.NewJobController()

	jobs := r.Group("/jobs")
	jobs.Use(middleware.JWTAuth())
	{
		// 任务管理
		jobs.POST("", jobController.CreateJob)                // 创建任务
		jobs.PUT("/:id", jobController.UpdateJob)             // 更新任务
		jobs.DELETE("/:id", jobController.DeleteJob)          // 删除任务
		jobs.POST("/:id/start", jobController.StartJob)       // 启动任务
		jobs.POST("/:id/stop", jobController.StopJob)         // 停止任务
		jobs.GET("", jobController.GetJobList)                // 获取任务列表
		jobs.GET("/history", jobController.GetJobHistoryList) // 获取任务历史列表
	}
}
