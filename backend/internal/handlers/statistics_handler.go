package handlers

import (
	"consulting-system/internal/services"
	"consulting-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// StatisticsHandler 统计处理器
type StatisticsHandler struct {
	statisticsService *services.StatisticsService
}

// NewStatisticsHandler 创建统计处理器
func NewStatisticsHandler() *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: services.NewStatisticsService(),
	}
}

// GetCustomerStatistics 获取客户统计
func (h *StatisticsHandler) GetCustomerStatistics(c *gin.Context) {
	stats, err := h.statisticsService.GetCustomerStatistics()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// GetServiceStatistics 获取服务统计
func (h *StatisticsHandler) GetServiceStatistics(c *gin.Context) {
	stats, err := h.statisticsService.GetServiceStatistics()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// GetFinanceStatistics 获取财务统计
func (h *StatisticsHandler) GetFinanceStatistics(c *gin.Context) {
	stats, err := h.statisticsService.GetFinanceStatistics()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// GetPerformanceStatistics 获取绩效统计
func (h *StatisticsHandler) GetPerformanceStatistics(c *gin.Context) {
	stats, err := h.statisticsService.GetPerformanceStatistics()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// GetDashboardData 获取仪表盘数据
func (h *StatisticsHandler) GetDashboardData(c *gin.Context) {
	data, err := h.statisticsService.GetDashboardData()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, data)
}

// RegisterRoutes 注册路由
func (h *StatisticsHandler) RegisterRoutes(r *gin.RouterGroup) {
	statistics := r.Group("/statistics")
	{
		statistics.GET("/dashboard", h.GetDashboardData)
		statistics.GET("/customer", h.GetCustomerStatistics)
		statistics.GET("/service", h.GetServiceStatistics)
		statistics.GET("/finance", h.GetFinanceStatistics)
		statistics.GET("/performance", h.GetPerformanceStatistics)
	}
}
