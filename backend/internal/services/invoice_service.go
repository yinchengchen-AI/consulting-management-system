package services

import (
	"errors"
	"time"

	"consulting-system/internal/database"
	"consulting-system/internal/models"

	"gorm.io/gorm"
)

// InvoiceService 发票服务
type InvoiceService struct{}

// NewInvoiceService 创建发票服务实例
func NewInvoiceService() *InvoiceService {
	return &InvoiceService{}
}

// CreateInvoiceRequest 创建发票请求
type CreateInvoiceRequest struct {
	CustomerID  uint           `json:"customer_id" binding:"required"`
	ServiceID   *uint          `json:"service_id" binding:"omitempty"`
	InvoiceType int            `json:"invoice_type" binding:"required,oneof=1 2"`
	Amount      float64        `json:"amount" binding:"required,min=0"`
	TaxRate     float64        `json:"tax_rate" binding:"omitempty,min=0,max=100"`
	InvoiceInfo models.JSON    `json:"invoice_info" binding:"omitempty"`
}

// UpdateInvoiceRequest 更新发票请求
type UpdateInvoiceRequest struct {
	InvoiceType int         `json:"invoice_type" binding:"omitempty,oneof=1 2"`
	Amount      float64     `json:"amount" binding:"omitempty,min=0"`
	TaxRate     float64     `json:"tax_rate" binding:"omitempty,min=0,max=100"`
	InvoiceInfo models.JSON `json:"invoice_info" binding:"omitempty"`
}

// AuditInvoiceRequest 审核发票请求
type AuditInvoiceRequest struct {
	InvoiceNo   string    `json:"invoice_no" binding:"required"`
	InvoiceCode string    `json:"invoice_code" binding:"omitempty"`
	InvoiceDate time.Time `json:"invoice_date" binding:"required"`
}

// ListInvoicesRequest 发票列表请求
type ListInvoicesRequest struct {
	CustomerID  uint   `form:"customer_id"`
	ServiceID   uint   `form:"service_id"`
	InvoiceType int    `form:"invoice_type"`
	Status      int    `form:"status"`
	InvoiceNo   string `form:"invoice_no"`
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
}

// ListInvoicesResponse 发票列表响应
type ListInvoicesResponse struct {
	List  []models.Invoice `json:"list"`
	Total int64            `json:"total"`
}

// CreateInvoice 创建发票
func (s *InvoiceService) CreateInvoice(req *CreateInvoiceRequest, createdBy uint) (*models.Invoice, error) {
	// 检查客户是否存在
	var customer models.Customer
	if err := database.DB.First(&customer, req.CustomerID).Error; err != nil {
		return nil, errors.New("客户不存在")
	}

	// 检查服务订单是否存在
	if req.ServiceID != nil {
		var service models.ServiceOrder
		if err := database.DB.First(&service, *req.ServiceID).Error; err != nil {
			return nil, errors.New("服务订单不存在")
		}
	}

	// 计算税额和总金额
	taxRate := req.TaxRate
	if taxRate == 0 {
		taxRate = 6 // 默认6%税率
	}
	taxAmount := req.Amount * taxRate / 100
	totalAmount := req.Amount + taxAmount

	invoice := models.Invoice{
		CustomerID:  req.CustomerID,
		ServiceID:   req.ServiceID,
		InvoiceType: req.InvoiceType,
		Amount:      req.Amount,
		TaxRate:     taxRate,
		TaxAmount:   taxAmount,
		TotalAmount: totalAmount,
		Status:      1, // 待开票
		InvoiceInfo: req.InvoiceInfo,
		CreatedBy:   createdBy,
	}

	if err := database.DB.Create(&invoice).Error; err != nil {
		return nil, err
	}

	// 加载关联数据
	database.DB.Preload("Customer").First(&invoice, invoice.ID)

	return &invoice, nil
}

// UpdateInvoice 更新发票
func (s *InvoiceService) UpdateInvoice(id uint, req *UpdateInvoiceRequest) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := database.DB.First(&invoice, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("发票不存在")
		}
		return nil, err
	}

	// 已开票的发票不能修改
	if invoice.Status == 2 {
		return nil, errors.New("已开票的发票不能修改")
	}

	updates := map[string]interface{}{}

	if req.InvoiceType != 0 {
		updates["invoice_type"] = req.InvoiceType
	}
	if req.Amount > 0 {
		updates["amount"] = req.Amount
		// 重新计算税额
		taxRate := invoice.TaxRate
		if req.TaxRate > 0 {
			taxRate = req.TaxRate
			updates["tax_rate"] = taxRate
		}
		updates["tax_amount"] = req.Amount * taxRate / 100
		updates["total_amount"] = req.Amount + (req.Amount * taxRate / 100)
	}
	if req.TaxRate > 0 && req.Amount == 0 {
		updates["tax_rate"] = req.TaxRate
		updates["tax_amount"] = invoice.Amount * req.TaxRate / 100
		updates["total_amount"] = invoice.Amount + (invoice.Amount * req.TaxRate / 100)
	}
	if req.InvoiceInfo != nil {
		updates["invoice_info"] = req.InvoiceInfo
	}

	if err := database.DB.Model(&invoice).Updates(updates).Error; err != nil {
		return nil, err
	}

	database.DB.Preload("Customer").First(&invoice, id)

	return &invoice, nil
}

// AuditInvoice 审核开票
func (s *InvoiceService) AuditInvoice(id uint, req *AuditInvoiceRequest) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := database.DB.First(&invoice, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("发票不存在")
		}
		return nil, err
	}

	// 只能审核待开票的发票
	if invoice.Status != 1 {
		return nil, errors.New("只能审核待开票的发票")
	}

	updates := map[string]interface{}{
		"status":       2, // 已开票
		"invoice_no":   req.InvoiceNo,
		"invoice_code": req.InvoiceCode,
		"invoice_date": req.InvoiceDate,
	}

	if err := database.DB.Model(&invoice).Updates(updates).Error; err != nil {
		return nil, err
	}

	database.DB.Preload("Customer").First(&invoice, id)

	return &invoice, nil
}

// VoidInvoice 作废发票
func (s *InvoiceService) VoidInvoice(id uint) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := database.DB.First(&invoice, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("发票不存在")
		}
		return nil, err
	}

	// 只能作废已开票的发票
	if invoice.Status != 2 {
		return nil, errors.New("只能作废已开票的发票")
	}

	if err := database.DB.Model(&invoice).Update("status", 3).Error; err != nil {
		return nil, err
	}

	return &invoice, nil
}

// DeleteInvoice 删除发票
func (s *InvoiceService) DeleteInvoice(id uint) error {
	var invoice models.Invoice
	if err := database.DB.First(&invoice, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("发票不存在")
		}
		return err
	}

	// 已开票的发票不能删除
	if invoice.Status == 2 {
		return errors.New("已开票的发票不能删除")
	}

	return database.DB.Delete(&invoice).Error
}

// GetInvoice 获取发票详情
func (s *InvoiceService) GetInvoice(id uint) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := database.DB.Preload("Customer").First(&invoice, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("发票不存在")
		}
		return nil, err
	}
	return &invoice, nil
}

// ListInvoices 获取发票列表
func (s *InvoiceService) ListInvoices(req *ListInvoicesRequest) (*ListInvoicesResponse, error) {
	var invoices []models.Invoice
	var total int64

	query := database.DB.Model(&models.Invoice{})

	if req.CustomerID != 0 {
		query = query.Where("customer_id = ?", req.CustomerID)
	}
	if req.ServiceID != 0 {
		query = query.Where("service_id = ?", req.ServiceID)
	}
	if req.InvoiceType != 0 {
		query = query.Where("invoice_type = ?", req.InvoiceType)
	}
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}
	if req.InvoiceNo != "" {
		query = query.Where("invoice_no = ?", req.InvoiceNo)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("Customer").Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&invoices).Error; err != nil {
		return nil, err
	}

	return &ListInvoicesResponse{
		List:  invoices,
		Total: total,
	}, nil
}

// GetInvoiceStats 获取发票统计
func (s *InvoiceService) GetInvoiceStats() (map[string]interface{}, error) {
	var total int64
	var byStatus []struct {
		Status int   `json:"status"`
		Count  int64 `json:"count"`
	}
	var totalAmount float64
	var totalTaxAmount float64

	database.DB.Model(&models.Invoice{}).Count(&total)
	database.DB.Model(&models.Invoice{}).Select("status, COUNT(*) as count").Group("status").Scan(&byStatus)
	database.DB.Model(&models.Invoice{}).Select("COALESCE(SUM(amount), 0) as total_amount").Pluck("total_amount", &totalAmount)
	database.DB.Model(&models.Invoice{}).Select("COALESCE(SUM(tax_amount), 0) as total_tax_amount").Pluck("total_tax_amount", &totalTaxAmount)

	return map[string]interface{}{
		"total":          total,
		"by_status":      byStatus,
		"total_amount":   totalAmount,
		"total_tax_amount": totalTaxAmount,
	}, nil
}
