package controllers

import (
	"datax-admin/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DashboardController 仪表盘控制器
type DashboardController struct {
	dashboardService *services.DashboardService
}

// NewDashboardController 创建仪表盘控制器
func NewDashboardController() *DashboardController {
	return &DashboardController{
		dashboardService: services.NewDashboardService(),
	}
}

// GetDashboardData 获取仪表盘数据
func (c *DashboardController) GetDashboardData(ctx *gin.Context) {
	resp, err := c.dashboardService.GetDashboardData()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
