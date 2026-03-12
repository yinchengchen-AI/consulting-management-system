package services

import (
	"errors"
	"time"

	"consulting-system/internal/database"
	"consulting-system/internal/models"

	"gorm.io/gorm"
)

// PaymentService 收款服务
type PaymentService struct{}

// NewPaymentService 创建收款服务实例
func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// CreatePaymentPlanRequest 创建收款计划请求
type CreatePaymentPlanRequest struct {
	CustomerID  uint       `json:"customer_id" binding:"required"`
	ServiceID   *uint      `json:"service_id" binding:"omitempty"`
	InvoiceID   *uint      `json:"invoice_id" binding:"omitempty"`
	Amount      float64    `json:"amount" binding:"required,min=0"`
	PlannedDate *time.Time `json:"planned_date" binding:"omitempty"`
	Remark      string     `json:"remark" binding:"omitempty"`
}

// UpdatePaymentPlanRequest 更新收款计划请求
type UpdatePaymentPlanRequest struct {
	Amount      float64    `json:"amount" binding:"omitempty,min=0"`
	PlannedDate *time.Time `json:"planned_date" binding:"omitempty"`
	Remark      string     `json:"remark" binding:"omitempty"`
}

// CreateReceiptRequest 创建收款记录请求
type CreateReceiptRequest struct {
	PlanID        uint      `json:"plan_id" binding:"required"`
	Amount        float64   `json:"amount" binding:"required,min=0"`
	ReceivedDate  time.Time `json:"received_date" binding:"required"`
	PaymentMethod int       `json:"payment_method" binding:"required,oneof=1 2 3 4"`
	Account       string    `json:"account" binding:"omitempty,max=100"`
	Remark        string    `json:"remark" binding:"omitempty"`
}

// CreateRefundRequest 创建退款请求
type CreateRefundRequest struct {
	ReceiptID uint    `json:"receipt_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,min=0"`
	Reason    string  `json:"reason" binding:"required"`
}

// AuditRefundRequest 审核退款请求
type AuditRefundRequest struct {
	Status int `json:"status" binding:"required,oneof=2 3"` // 2-批准, 3-拒绝
}

// ListPaymentPlansRequest 收款计划列表请求
type ListPaymentPlansRequest struct {
	CustomerID uint   `form:"customer_id"`
	ServiceID  uint   `form:"service_id"`
	Status     int    `form:"status"`
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=10"`
}

// ListPaymentPlansResponse 收款计划列表响应
type ListPaymentPlansResponse struct {
	List  []models.PaymentPlan `json:"list"`
	Total int64                `json:"total"`
}

// CreatePaymentPlan 创建收款计划
func (s *PaymentService) CreatePaymentPlan(req *CreatePaymentPlanRequest, createdBy uint) (*models.PaymentPlan, error) {
	// 检查客户是否存在
	var customer models.Customer
	if err := database.DB.First(&customer, req.CustomerID).Error; err != nil {
		return nil, errors.New("客户不存在")
	}

	plan := models.PaymentPlan{
		CustomerID:  req.CustomerID,
		ServiceID:   req.ServiceID,
		InvoiceID:   req.InvoiceID,
		Amount:      req.Amount,
		PlannedDate: req.PlannedDate,
		Status:      1, // 待收款
		Remark:      req.Remark,
		CreatedBy:   createdBy,
	}

	if err := database.DB.Create(&plan).Error; err != nil {
		return nil, err
	}

	database.DB.Preload("Customer").First(&plan, plan.ID)

	return &plan, nil
}

// UpdatePaymentPlan 更新收款计划
func (s *PaymentService) UpdatePaymentPlan(id uint, req *UpdatePaymentPlanRequest) (*models.PaymentPlan, error) {
	var plan models.PaymentPlan
	if err := database.DB.First(&plan, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("收款计划不存在")
		}
		return nil, err
	}

	// 已收完的计划不能修改
	if plan.Status == 3 {
		return nil, errors.New("已收完的计划不能修改")
	}

	updates := map[string]interface{}{}
	if req.Amount > 0 {
		updates["amount"] = req.Amount
	}
	if req.PlannedDate != nil {
		updates["planned_date"] = req.PlannedDate
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	if err := database.DB.Model(&plan).Updates(updates).Error; err != nil {
		return nil, err
	}

	database.DB.Preload("Customer").First(&plan, id)

	return &plan, nil
}

// DeletePaymentPlan 删除收款计划
func (s *PaymentService) DeletePaymentPlan(id uint) error {
	var plan models.PaymentPlan
	if err := database.DB.First(&plan, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("收款计划不存在")
		}
		return err
	}

	// 已收完的计划不能删除
	if plan.Status == 3 {
		return errors.New("已收完的计划不能删除")
	}

	return database.DB.Delete(&plan).Error
}

// GetPaymentPlan 获取收款计划详情
func (s *PaymentService) GetPaymentPlan(id uint) (*models.PaymentPlan, error) {
	var plan models.PaymentPlan
	if err := database.DB.Preload("Customer").Preload("Receipts").First(&plan, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("收款计划不存在")
		}
		return nil, err
	}
	return &plan, nil
}

// ListPaymentPlans 获取收款计划列表
func (s *PaymentService) ListPaymentPlans(req *ListPaymentPlansRequest) (*ListPaymentPlansResponse, error) {
	var plans []models.PaymentPlan
	var total int64

	query := database.DB.Model(&models.PaymentPlan{})

	if req.CustomerID != 0 {
		query = query.Where("customer_id = ?", req.CustomerID)
	}
	if req.ServiceID != 0 {
		query = query.Where("service_id = ?", req.ServiceID)
	}
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("Customer").Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&plans).Error; err != nil {
		return nil, err
	}

	return &ListPaymentPlansResponse{
		List:  plans,
		Total: total,
	}, nil
}

// CreateReceipt 创建收款记录
func (s *PaymentService) CreateReceipt(req *CreateReceiptRequest, createdBy uint) (*models.Receipt, error) {
	// 检查收款计划是否存在
	var plan models.PaymentPlan
	if err := database.DB.First(&plan, req.PlanID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("收款计划不存在")
		}
		return nil, err
	}

	// 计算已收款金额
	var receivedAmount float64
	database.DB.Model(&models.Receipt{}).Where("plan_id = ?", req.PlanID).Select("COALESCE(SUM(amount), 0)").Scan(&receivedAmount)

	// 检查是否超额收款
	if receivedAmount+req.Amount > plan.Amount {
		return nil, errors.New("收款金额超过计划金额")
	}

	receipt := models.Receipt{
		PlanID:        req.PlanID,
		Amount:        req.Amount,
		ReceivedDate:  req.ReceivedDate,
		PaymentMethod: req.PaymentMethod,
		Account:       req.Account,
		Remark:        req.Remark,
		CreatedBy:     createdBy,
	}

	if err := database.DB.Create(&receipt).Error; err != nil {
		return nil, err
	}

	// 更新收款计划状态
	newReceivedAmount := receivedAmount + req.Amount
	if newReceivedAmount >= plan.Amount {
		database.DB.Model(&plan).Update("status", 3) // 已收完
	} else if newReceivedAmount > 0 {
		database.DB.Model(&plan).Update("status", 2) // 部分收款
	}

	return &receipt, nil
}

// DeleteReceipt 删除收款记录
func (s *PaymentService) DeleteReceipt(id uint) error {
	var receipt models.Receipt
	if err := database.DB.First(&receipt, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("收款记录不存在")
		}
		return err
	}

	// 检查是否有退款记录
	var refundCount int64
	database.DB.Model(&models.Refund{}).Where("receipt_id = ?", id).Count(&refundCount)
	if refundCount > 0 {
		return errors.New("该收款记录有退款，无法删除")
	}

	// 更新收款计划状态
	var plan models.PaymentPlan
	if err := database.DB.First(&plan, receipt.PlanID).Error; err == nil {
		var receivedAmount float64
		database.DB.Model(&models.Receipt{}).Where("plan_id = ? AND id != ?", receipt.PlanID, id).
			Select("COALESCE(SUM(amount), 0)").Scan(&receivedAmount)

		if receivedAmount == 0 {
			database.DB.Model(&plan).Update("status", 1) // 待收款
		} else if receivedAmount < plan.Amount {
			database.DB.Model(&plan).Update("status", 2) // 部分收款
		}
	}

	return database.DB.Delete(&receipt).Error
}

// GetReceipt 获取收款记录详情
func (s *PaymentService) GetReceipt(id uint) (*models.Receipt, error) {
	var receipt models.Receipt
	if err := database.DB.Preload("Refunds").First(&receipt, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("收款记录不存在")
		}
		return nil, err
	}
	return &receipt, nil
}

// CreateRefund 创建退款申请
func (s *PaymentService) CreateRefund(req *CreateRefundRequest, createdBy uint) (*models.Refund, error) {
	// 检查收款记录是否存在
	var receipt models.Receipt
	if err := database.DB.First(&receipt, req.ReceiptID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("收款记录不存在")
		}
		return nil, err
	}

	// 检查退款金额
	var refundedAmount float64
	database.DB.Model(&models.Refund{}).Where("receipt_id = ? AND status IN (1, 2, 4)", req.ReceiptID).
		Select("COALESCE(SUM(amount), 0)").Scan(&refundedAmount)

	if refundedAmount+req.Amount > receipt.Amount {
		return nil, errors.New("退款金额超过收款金额")
	}

	refund := models.Refund{
		ReceiptID: req.ReceiptID,
		Amount:    req.Amount,
		Reason:    req.Reason,
		Status:    1, // 待审核
		CreatedBy: createdBy,
	}

	if err := database.DB.Create(&refund).Error; err != nil {
		return nil, err
	}

	return &refund, nil
}

// AuditRefund 审核退款
func (s *PaymentService) AuditRefund(id uint, req *AuditRefundRequest) (*models.Refund, error) {
	var refund models.Refund
	if err := database.DB.First(&refund, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("退款记录不存在")
		}
		return nil, err
	}

	// 只能审核待审核的退款
	if refund.Status != 1 {
		return nil, errors.New("只能审核待审核的退款")
	}

	if err := database.DB.Model(&refund).Update("status", req.Status).Error; err != nil {
		return nil, err
	}

	return &refund, nil
}

// CompleteRefund 完成退款
func (s *PaymentService) CompleteRefund(id uint) (*models.Refund, error) {
	var refund models.Refund
	if err := database.DB.First(&refund, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("退款记录不存在")
		}
		return nil, err
	}

	// 只能完成已批准的退款
	if refund.Status != 2 {
		return nil, errors.New("只能完成已批准的退款")
	}

	if err := database.DB.Model(&refund).Update("status", 4).Error; err != nil {
		return nil, err
	}

	return &refund, nil
}

// GetPaymentStats 获取收款统计
func (s *PaymentService) GetPaymentStats() (map[string]interface{}, error) {
	var totalPlan int64
	var totalReceipt float64
	var totalPlanAmount float64

	database.DB.Model(&models.PaymentPlan{}).Count(&totalPlan)
	database.DB.Model(&models.Receipt{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalReceipt)
	database.DB.Model(&models.PaymentPlan{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalPlanAmount)

	return map[string]interface{}{
		"total_plan":        totalPlan,
		"total_receipt":     totalReceipt,
		"total_plan_amount": totalPlanAmount,
		"receipt_rate":      0.0,
	}, nil
}
