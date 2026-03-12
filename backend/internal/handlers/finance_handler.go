package handlers

import (
	"strconv"

	"consulting-system/internal/middleware"
	"consulting-system/internal/services"
	"consulting-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// InvoiceHandler 发票处理器
type InvoiceHandler struct {
	invoiceService *services.InvoiceService
}

// NewInvoiceHandler 创建发票处理器
func NewInvoiceHandler() *InvoiceHandler {
	return &InvoiceHandler{
		invoiceService: services.NewInvoiceService(),
	}
}

// CreateInvoice 创建发票
func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	var req services.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	createdBy := middleware.GetUserID(c)
	invoice, err := h.invoiceService.CreateInvoice(&req, createdBy)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Created(c, invoice)
}

// UpdateInvoice 更新发票
func (h *InvoiceHandler) UpdateInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的发票ID")
		return
	}

	var req services.UpdateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	invoice, err := h.invoiceService.UpdateInvoice(uint(id), &req)
	if err != nil {
		if err.Error() == "发票不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Success(c, invoice)
}

// AuditInvoice 审核开票
func (h *InvoiceHandler) AuditInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的发票ID")
		return
	}

	var req services.AuditInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	invoice, err := h.invoiceService.AuditInvoice(uint(id), &req)
	if err != nil {
		if err.Error() == "发票不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Success(c, invoice)
}

// VoidInvoice 作废发票
func (h *InvoiceHandler) VoidInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的发票ID")
		return
	}

	invoice, err := h.invoiceService.VoidInvoice(uint(id))
	if err != nil {
		if err.Error() == "发票不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Success(c, invoice)
}

// DeleteInvoice 删除发票
func (h *InvoiceHandler) DeleteInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的发票ID")
		return
	}

	if err := h.invoiceService.DeleteInvoice(uint(id)); err != nil {
		if err.Error() == "发票不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// GetInvoice 获取发票详情
func (h *InvoiceHandler) GetInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的发票ID")
		return
	}

	invoice, err := h.invoiceService.GetInvoice(uint(id))
	if err != nil {
		if err.Error() == "发票不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, invoice)
}

// ListInvoices 获取发票列表
func (h *InvoiceHandler) ListInvoices(c *gin.Context) {
	var req services.ListInvoicesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	result, err := h.invoiceService.ListInvoices(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.BuildMeta(result.Total, req.Page, req.PageSize)
	response.SuccessWithMeta(c, result.List, meta)
}

// GetInvoiceStats 获取发票统计
func (h *InvoiceHandler) GetInvoiceStats(c *gin.Context) {
	stats, err := h.invoiceService.GetInvoiceStats()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// RegisterRoutes 注册路由
func (h *InvoiceHandler) RegisterRoutes(r *gin.RouterGroup) {
	invoices := r.Group("/invoices")
	{
		invoices.GET("", h.ListInvoices)
		invoices.POST("", h.CreateInvoice)
		invoices.GET("/stats", h.GetInvoiceStats)
		invoices.GET("/:id", h.GetInvoice)
		invoices.PUT("/:id", h.UpdateInvoice)
		invoices.POST("/:id/audit", h.AuditInvoice)
		invoices.POST("/:id/void", h.VoidInvoice)
		invoices.DELETE("/:id", h.DeleteInvoice)
	}
}

// ==================== PaymentHandler ====================

// PaymentHandler 收款处理器
type PaymentHandler struct {
	paymentService *services.PaymentService
}

// NewPaymentHandler 创建收款处理器
func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		paymentService: services.NewPaymentService(),
	}
}

// CreatePaymentPlan 创建收款计划
func (h *PaymentHandler) CreatePaymentPlan(c *gin.Context) {
	var req services.CreatePaymentPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	createdBy := middleware.GetUserID(c)
	plan, err := h.paymentService.CreatePaymentPlan(&req, createdBy)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Created(c, plan)
}

// UpdatePaymentPlan 更新收款计划
func (h *PaymentHandler) UpdatePaymentPlan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的收款计划ID")
		return
	}

	var req services.UpdatePaymentPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	plan, err := h.paymentService.UpdatePaymentPlan(uint(id), &req)
	if err != nil {
		if err.Error() == "收款计划不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Success(c, plan)
}

// DeletePaymentPlan 删除收款计划
func (h *PaymentHandler) DeletePaymentPlan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的收款计划ID")
		return
	}

	if err := h.paymentService.DeletePaymentPlan(uint(id)); err != nil {
		if err.Error() == "收款计划不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// GetPaymentPlan 获取收款计划详情
func (h *PaymentHandler) GetPaymentPlan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的收款计划ID")
		return
	}

	plan, err := h.paymentService.GetPaymentPlan(uint(id))
	if err != nil {
		if err.Error() == "收款计划不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, plan)
}

// ListPaymentPlans 获取收款计划列表
func (h *PaymentHandler) ListPaymentPlans(c *gin.Context) {
	var req services.ListPaymentPlansRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	result, err := h.paymentService.ListPaymentPlans(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.BuildMeta(result.Total, req.Page, req.PageSize)
	response.SuccessWithMeta(c, result.List, meta)
}

// CreateReceipt 创建收款记录
func (h *PaymentHandler) CreateReceipt(c *gin.Context) {
	var req services.CreateReceiptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	createdBy := middleware.GetUserID(c)
	receipt, err := h.paymentService.CreateReceipt(&req, createdBy)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Created(c, receipt)
}

// DeleteReceipt 删除收款记录
func (h *PaymentHandler) DeleteReceipt(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("receiptId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的收款记录ID")
		return
	}

	if err := h.paymentService.DeleteReceipt(uint(id)); err != nil {
		if err.Error() == "收款记录不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// CreateRefund 创建退款申请
func (h *PaymentHandler) CreateRefund(c *gin.Context) {
	var req services.CreateRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	createdBy := middleware.GetUserID(c)
	refund, err := h.paymentService.CreateRefund(&req, createdBy)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Created(c, refund)
}

// AuditRefund 审核退款
func (h *PaymentHandler) AuditRefund(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("refundId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的退款记录ID")
		return
	}

	var req services.AuditRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	refund, err := h.paymentService.AuditRefund(uint(id), &req)
	if err != nil {
		if err.Error() == "退款记录不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Success(c, refund)
}

// CompleteRefund 完成退款
func (h *PaymentHandler) CompleteRefund(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("refundId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的退款记录ID")
		return
	}

	refund, err := h.paymentService.CompleteRefund(uint(id))
	if err != nil {
		if err.Error() == "退款记录不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Success(c, refund)
}

// GetPaymentStats 获取收款统计
func (h *PaymentHandler) GetPaymentStats(c *gin.Context) {
	stats, err := h.paymentService.GetPaymentStats()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// RegisterRoutes 注册路由
func (h *PaymentHandler) RegisterRoutes(r *gin.RouterGroup) {
	// 收款计划
	plans := r.Group("/payment-plans")
	{
		plans.GET("", h.ListPaymentPlans)
		plans.POST("", h.CreatePaymentPlan)
		plans.GET("/stats", h.GetPaymentStats)
		plans.GET("/:id", h.GetPaymentPlan)
		plans.PUT("/:id", h.UpdatePaymentPlan)
		plans.DELETE("/:id", h.DeletePaymentPlan)
	}

	// 收款记录
	r.POST("/receipts", h.CreateReceipt)
	r.DELETE("/receipts/:receiptId", h.DeleteReceipt)

	// 退款
	r.POST("/refunds", h.CreateRefund)
	r.POST("/refunds/:refundId/audit", h.AuditRefund)
	r.POST("/refunds/:refundId/complete", h.CompleteRefund)
}
