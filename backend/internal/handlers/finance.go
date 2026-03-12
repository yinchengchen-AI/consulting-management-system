package handlers

import (
	"consulting-system/backend/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FinanceHandler 财务处理器
type FinanceHandler struct {
	db *gorm.DB
}

// NewFinanceHandler 创建财务处理器
func NewFinanceHandler(db *gorm.DB) *FinanceHandler {
	return &FinanceHandler{db: db}
}

// ListIncome 获取收入列表
// @Summary 获取收入列表
// @Description 获取所有收入记录，支持分页和筛选
// @Tags 财务管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Param customer_id query string false "客户 ID"
// @Param status query string false "状态筛选"
// @Success 200 {object} map[string]interface{}
// @Security Bearer
// @Router /finance/income [get]
func (h *FinanceHandler) ListIncome(c *gin.Context) {
	var params struct {
		Page       int    `form:"page,default=1"`
		PageSize   int    `form:"page_size,default=10"`
		StartDate  string `form:"start_date"`
		EndDate    string `form:"end_date"`
		CustomerID string `form:"customer_id"`
		Status     string `form:"status"`
	}

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	query := h.db.Model(&models.Income{}).Preload("Customer").Preload("Contract").Preload("Project")

	if params.StartDate != "" {
		query = query.Where("created_at >= ?", params.StartDate)
	}
	if params.EndDate != "" {
		query = query.Where("created_at <= ?", params.EndDate+" 23:59:59")
	}
	if params.CustomerID != "" {
		query = query.Where("customer_id = ?", params.CustomerID)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	var total int64
	query.Count(&total)

	var incomes []models.Income
	offset := (params.Page - 1) * params.PageSize
	query.Order("created_at DESC").Limit(params.PageSize).Offset(offset).Find(&incomes)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"list":  incomes,
			"total": total,
			"page":  params.Page,
			"size":  params.PageSize,
		},
	})
}

// ListExpense 获取支出列表
// @Summary 获取支出列表
// @Description 获取所有支出记录，支持分页和筛选
// @Tags 财务管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Param project_id query string false "项目 ID"
// @Param type query string false "类型筛选"
// @Param status query string false "状态筛选"
// @Success 200 {object} map[string]interface{}
// @Security Bearer
// @Router /finance/expense [get]
func (h *FinanceHandler) ListExpense(c *gin.Context) {
	var params struct {
		Page      int    `form:"page,default=1"`
		PageSize  int    `form:"page_size,default=10"`
		StartDate string `form:"start_date"`
		EndDate   string `form:"end_date"`
		ProjectID string `form:"project_id"`
		Type      string `form:"type"`
		Status    string `form:"status"`
	}

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	query := h.db.Model(&models.Expense{}).Preload("Project")

	if params.StartDate != "" {
		query = query.Where("created_at >= ?", params.StartDate)
	}
	if params.EndDate != "" {
		query = query.Where("created_at <= ?", params.EndDate+" 23:59:59")
	}
	if params.ProjectID != "" {
		query = query.Where("project_id = ?", params.ProjectID)
	}
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	var total int64
	query.Count(&total)

	var expenses []models.Expense
	offset := (params.Page - 1) * params.PageSize
	query.Order("created_at DESC").Limit(params.PageSize).Offset(offset).Find(&expenses)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"list":  expenses,
			"total": total,
			"page":  params.Page,
			"size":  params.PageSize,
		},
	})
}

// GetSummary 获取财务汇总
// @Summary 获取财务汇总
// @Description 获取财务汇总数据，包括总收入、总支出、净利润等
// @Tags 财务管理
// @Accept json
// @Produce json
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} models.FinanceSummary
// @Security Bearer
// @Router /finance/summary [get]
func (h *FinanceHandler) GetSummary(c *gin.Context) {
	var params struct {
		StartDate string `form:"start_date"`
		EndDate   string `form:"end_date"`
	}

	c.ShouldBindQuery(&params)

	// 默认查询当月
	if params.StartDate == "" {
		params.StartDate = time.Now().Format("2006-01") + "-01"
	}
	if params.EndDate == "" {
		params.EndDate = time.Now().Format("2006-01-02")
	}

	var summary models.FinanceSummary

	// 统计收入
	var incomeResult struct {
		Total   float64
		Pending float64
	}
	h.db.Model(&models.Income{}).
		Select("COALESCE(SUM(amount), 0) as total, COALESCE(SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END), 0) as pending").
		Where("created_at >= ? AND created_at <= ?", params.StartDate, params.EndDate+" 23:59:59").
		Scan(&incomeResult)

	summary.TotalIncome = incomeResult.Total
	summary.PendingIncome = incomeResult.Pending

	// 统计支出
	var expenseResult struct {
		Total   float64
		Pending float64
	}
	h.db.Model(&models.Expense{}).
		Select("COALESCE(SUM(amount), 0) as total, COALESCE(SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END), 0) as pending").
		Where("created_at >= ? AND created_at <= ?", params.StartDate, params.EndDate+" 23:59:59").
		Scan(&expenseResult)

	summary.TotalExpense = expenseResult.Total
	summary.PendingExpense = expenseResult.Pending

	// 计算净利润
	summary.NetProfit = summary.TotalIncome - summary.TotalExpense

	// 计算利润率
	if summary.TotalIncome > 0 {
		summary.ProfitMargin = (summary.NetProfit / summary.TotalIncome) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data":    summary,
	})
}
