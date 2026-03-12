package services

import (
	"consulting-system/internal/database"
	"consulting-system/internal/models"
	"time"
)

// StatisticsService 统计服务
type StatisticsService struct{}

// NewStatisticsService 创建统计服务实例
func NewStatisticsService() *StatisticsService {
	return &StatisticsService{}
}

// CustomerStatistics 客户统计
type CustomerStatistics struct {
	Total           int64                  `json:"total"`
	ByStatus        []StatusCount          `json:"by_status"`
	ByIndustry      []IndustryCount        `json:"by_industry"`
	MonthlyNew      []MonthlyCount         `json:"monthly_new"`
	RecentCustomers []models.Customer      `json:"recent_customers"`
}

// ServiceStatistics 服务统计
type ServiceStatistics struct {
	Total       int64           `json:"total"`
	ByStatus    []StatusCount   `json:"by_status"`
	ByType      []TypeCount     `json:"by_type"`
	TotalAmount float64         `json:"total_amount"`
	MonthlyData []MonthlyAmount `json:"monthly_data"`
}

// FinanceStatistics 财务统计
type FinanceStatistics struct {
	InvoiceTotal      float64         `json:"invoice_total"`
	InvoiceTaxTotal   float64         `json:"invoice_tax_total"`
	ReceiptTotal      float64         `json:"receipt_total"`
	PendingReceipt    float64         `json:"pending_receipt"`
	MonthlyInvoice    []MonthlyAmount `json:"monthly_invoice"`
	MonthlyReceipt    []MonthlyAmount `json:"monthly_receipt"`
}

// PerformanceStatistics 绩效统计
type PerformanceStatistics struct {
	UserPerformance []UserPerformance `json:"user_performance"`
}

// StatusCount 状态统计
type StatusCount struct {
	Status int   `json:"status"`
	Count  int64 `json:"count"`
}

// IndustryCount 行业统计
type IndustryCount struct {
	Industry string `json:"industry"`
	Count    int64  `json:"count"`
}

// TypeCount 类型统计
type TypeCount struct {
	TypeID uint   `json:"type_id"`
	Name   string `json:"name"`
	Count  int64  `json:"count"`
}

// MonthlyCount 月度统计
type MonthlyCount struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

// MonthlyAmount 月度金额统计
type MonthlyAmount struct {
	Month  string  `json:"month"`
	Amount float64 `json:"amount"`
}

// UserPerformance 用户绩效
type UserPerformance struct {
	UserID     uint    `json:"user_id"`
	Username   string  `json:"username"`
	RealName   string  `json:"real_name"`
	CustomerCount int64 `json:"customer_count"`
	ServiceCount  int64 `json:"service_count"`
	TotalAmount   float64 `json:"total_amount"`
}

// GetCustomerStatistics 获取客户统计
func (s *StatisticsService) GetCustomerStatistics() (*CustomerStatistics, error) {
	var total int64
	database.DB.Model(&models.Customer{}).Count(&total)

	var byStatus []StatusCount
	database.DB.Model(&models.Customer{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&byStatus)

	var byIndustry []IndustryCount
	database.DB.Model(&models.Customer{}).
		Select("industry, COUNT(*) as count").
		Where("industry != ?", "").
		Group("industry").
		Scan(&byIndustry)

	// 获取最近6个月的新增客户
	var monthlyNew []MonthlyCount
	for i := 5; i >= 0; i-- {
		startDate := time.Now().AddDate(0, -i, 0).Format("2006-01")
		var count int64
		database.DB.Model(&models.Customer{}).
			Where("TO_CHAR(created_at, 'YYYY-MM') = ?", startDate).
			Count(&count)
		monthlyNew = append(monthlyNew, MonthlyCount{
			Month: startDate,
			Count: count,
		})
	}

	// 获取最近的客户
	var recentCustomers []models.Customer
	database.DB.Order("created_at DESC").Limit(5).Find(&recentCustomers)

	return &CustomerStatistics{
		Total:           total,
		ByStatus:        byStatus,
		ByIndustry:      byIndustry,
		MonthlyNew:      monthlyNew,
		RecentCustomers: recentCustomers,
	}, nil
}

// GetServiceStatistics 获取服务统计
func (s *StatisticsService) GetServiceStatistics() (*ServiceStatistics, error) {
	var total int64
	database.DB.Model(&models.ServiceOrder{}).Count(&total)

	var byStatus []StatusCount
	database.DB.Model(&models.ServiceOrder{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&byStatus)

	var byType []TypeCount
	database.DB.Model(&models.ServiceOrder{}).
		Joins("JOIN service_types ON service_orders.service_type_id = service_types.id").
		Select("service_type_id as type_id, service_types.name, COUNT(*) as count").
		Group("service_type_id, service_types.name").
		Scan(&byType)

	var totalAmount float64
	database.DB.Model(&models.ServiceOrder{}).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalAmount)

	// 获取最近6个月的服务数据
	var monthlyData []MonthlyAmount
	for i := 5; i >= 0; i-- {
		startDate := time.Now().AddDate(0, -i, 0).Format("2006-01")
		var amount float64
		database.DB.Model(&models.ServiceOrder{}).
			Where("TO_CHAR(created_at, 'YYYY-MM') = ?", startDate).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&amount)
		monthlyData = append(monthlyData, MonthlyAmount{
			Month:  startDate,
			Amount: amount,
		})
	}

	return &ServiceStatistics{
		Total:       total,
		ByStatus:    byStatus,
		ByType:      byType,
		TotalAmount: totalAmount,
		MonthlyData: monthlyData,
	}, nil
}

// GetFinanceStatistics 获取财务统计
func (s *StatisticsService) GetFinanceStatistics() (*FinanceStatistics, error) {
	var invoiceTotal float64
	database.DB.Model(&models.Invoice{}).
		Where("status = ?", 2).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&invoiceTotal)

	var invoiceTaxTotal float64
	database.DB.Model(&models.Invoice{}).
		Where("status = ?", 2).
		Select("COALESCE(SUM(tax_amount), 0)").
		Scan(&invoiceTaxTotal)

	var receiptTotal float64
	database.DB.Model(&models.Receipt{}).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&receiptTotal)

	var totalPlanAmount float64
	database.DB.Model(&models.PaymentPlan{}).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalPlanAmount)

	pendingReceipt := invoiceTotal - receiptTotal
	if pendingReceipt < 0 {
		pendingReceipt = 0
	}

	// 获取最近6个月的发票数据
	var monthlyInvoice []MonthlyAmount
	for i := 5; i >= 0; i-- {
		startDate := time.Now().AddDate(0, -i, 0).Format("2006-01")
		var amount float64
		database.DB.Model(&models.Invoice{}).
			Where("status = ? AND TO_CHAR(invoice_date, 'YYYY-MM') = ?", 2, startDate).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&amount)
		monthlyInvoice = append(monthlyInvoice, MonthlyAmount{
			Month:  startDate,
			Amount: amount,
		})
	}

	// 获取最近6个月的收款数据
	var monthlyReceipt []MonthlyAmount
	for i := 5; i >= 0; i-- {
		startDate := time.Now().AddDate(0, -i, 0).Format("2006-01")
		var amount float64
		database.DB.Model(&models.Receipt{}).
			Where("TO_CHAR(received_date, 'YYYY-MM') = ?", startDate).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&amount)
		monthlyReceipt = append(monthlyReceipt, MonthlyAmount{
			Month:  startDate,
			Amount: amount,
		})
	}

	return &FinanceStatistics{
		InvoiceTotal:    invoiceTotal,
		InvoiceTaxTotal: invoiceTaxTotal,
		ReceiptTotal:    receiptTotal,
		PendingReceipt:  pendingReceipt,
		MonthlyInvoice:  monthlyInvoice,
		MonthlyReceipt:  monthlyReceipt,
	}, nil
}

// GetPerformanceStatistics 获取绩效统计
func (s *StatisticsService) GetPerformanceStatistics() (*PerformanceStatistics, error) {
	var userPerformance []UserPerformance

	// 获取所有用户
	var users []models.User
	database.DB.Where("status = ?", 1).Find(&users)

	for _, user := range users {
		var customerCount int64
		database.DB.Model(&models.Customer{}).Where("created_by = ?", user.ID).Count(&customerCount)

		var serviceCount int64
		database.DB.Model(&models.ServiceOrder{}).Where("created_by = ?", user.ID).Count(&serviceCount)

		var totalAmount float64
		database.DB.Model(&models.ServiceOrder{}).Where("created_by = ?", user.ID).
			Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)

		userPerformance = append(userPerformance, UserPerformance{
			UserID:        user.ID,
			Username:      user.Username,
			RealName:      user.RealName,
			CustomerCount: customerCount,
			ServiceCount:  serviceCount,
			TotalAmount:   totalAmount,
		})
	}

	return &PerformanceStatistics{
		UserPerformance: userPerformance,
	}, nil
}

// GetDashboardData 获取仪表盘数据
func (s *StatisticsService) GetDashboardData() (map[string]interface{}, error) {
	customerStats, _ := s.GetCustomerStatistics()
	serviceStats, _ := s.GetServiceStatistics()
	financeStats, _ := s.GetFinanceStatistics()

	return map[string]interface{}{
		"customer": customerStats,
		"service":  serviceStats,
		"finance":  financeStats,
	}, nil
}
