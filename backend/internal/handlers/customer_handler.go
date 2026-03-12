package handlers

import (
	"strconv"

	"consulting-system/internal/middleware"
	"consulting-system/internal/services"
	"consulting-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// CustomerHandler 客户处理器
type CustomerHandler struct {
	customerService *services.CustomerService
}

// NewCustomerHandler 创建客户处理器
func NewCustomerHandler() *CustomerHandler {
	return &CustomerHandler{
		customerService: services.NewCustomerService(),
	}
}

// CreateCustomer 创建客户
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req services.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	createdBy := middleware.GetUserID(c)
	customer, err := h.customerService.CreateCustomer(&req, createdBy)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, customer)
}

// UpdateCustomer 更新客户
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的客户ID")
		return
	}

	var req services.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	customer, err := h.customerService.UpdateCustomer(uint(id), &req)
	if err != nil {
		if err.Error() == "客户不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, customer)
}

// DeleteCustomer 删除客户
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的客户ID")
		return
	}

	if err := h.customerService.DeleteCustomer(uint(id)); err != nil {
		if err.Error() == "客户不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// GetCustomer 获取客户详情
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的客户ID")
		return
	}

	customer, err := h.customerService.GetCustomer(uint(id))
	if err != nil {
		if err.Error() == "客户不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, customer)
}

// ListCustomers 获取客户列表
func (h *CustomerHandler) ListCustomers(c *gin.Context) {
	var req services.ListCustomersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	result, err := h.customerService.ListCustomers(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.BuildMeta(result.Total, req.Page, req.PageSize)
	response.SuccessWithMeta(c, result.List, meta)
}

// CreateFollowUp 创建跟进记录
func (h *CustomerHandler) CreateFollowUp(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的客户ID")
		return
	}

	var req services.CreateFollowUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	createdBy := middleware.GetUserID(c)
	followUp, err := h.customerService.CreateFollowUp(uint(customerID), &req, createdBy)
	if err != nil {
		if err.Error() == "客户不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, followUp)
}

// GetFollowUpRecords 获取跟进记录列表
func (h *CustomerHandler) GetFollowUpRecords(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的客户ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	records, total, err := h.customerService.GetFollowUpRecords(uint(customerID), page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.BuildMeta(total, page, pageSize)
	response.SuccessWithMeta(c, records, meta)
}

// DeleteFollowUp 删除跟进记录
func (h *CustomerHandler) DeleteFollowUp(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("followUpId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的跟进记录ID")
		return
	}

	if err := h.customerService.DeleteFollowUp(uint(id)); err != nil {
		if err.Error() == "跟进记录不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// GetCustomerStats 获取客户统计
func (h *CustomerHandler) GetCustomerStats(c *gin.Context) {
	stats, err := h.customerService.GetCustomerStats()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// RegisterRoutes 注册路由
func (h *CustomerHandler) RegisterRoutes(r *gin.RouterGroup) {
	customers := r.Group("/customers")
	{
		customers.GET("", h.ListCustomers)
		customers.POST("", h.CreateCustomer)
		customers.GET("/stats", h.GetCustomerStats)
		customers.GET("/:id", h.GetCustomer)
		customers.PUT("/:id", h.UpdateCustomer)
		customers.DELETE("/:id", h.DeleteCustomer)
		customers.POST("/:id/follow-ups", h.CreateFollowUp)
		customers.GET("/:id/follow-ups", h.GetFollowUpRecords)
	}

	r.DELETE("/follow-ups/:followUpId", h.DeleteFollowUp)
}
