package handlers

import (
	"strconv"

	"consulting-system/internal/middleware"
	"consulting-system/internal/services"
	"consulting-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// ServiceTypeHandler 服务类型处理器
type ServiceTypeHandler struct {
	serviceTypeService *services.ServiceTypeService
}

// NewServiceTypeHandler 创建服务类型处理器
func NewServiceTypeHandler() *ServiceTypeHandler {
	return &ServiceTypeHandler{
		serviceTypeService: services.NewServiceTypeService(),
	}
}

// CreateServiceType 创建服务类型
func (h *ServiceTypeHandler) CreateServiceType(c *gin.Context) {
	var req services.CreateServiceTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	serviceType, err := h.serviceTypeService.CreateServiceType(&req)
	if err != nil {
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Created(c, serviceType)
}

// UpdateServiceType 更新服务类型
func (h *ServiceTypeHandler) UpdateServiceType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的服务类型ID")
		return
	}

	var req services.UpdateServiceTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	serviceType, err := h.serviceTypeService.UpdateServiceType(uint(id), &req)
	if err != nil {
		if err.Error() == "服务类型不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, serviceType)
}

// DeleteServiceType 删除服务类型
func (h *ServiceTypeHandler) DeleteServiceType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的服务类型ID")
		return
	}

	if err := h.serviceTypeService.DeleteServiceType(uint(id)); err != nil {
		if err.Error() == "服务类型不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// GetServiceType 获取服务类型详情
func (h *ServiceTypeHandler) GetServiceType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的服务类型ID")
		return
	}

	serviceType, err := h.serviceTypeService.GetServiceType(uint(id))
	if err != nil {
		if err.Error() == "服务类型不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, serviceType)
}

// ListServiceTypes 获取服务类型列表
func (h *ServiceTypeHandler) ListServiceTypes(c *gin.Context) {
	var req services.ListServiceTypesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	result, err := h.serviceTypeService.ListServiceTypes(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.BuildMeta(result.Total, req.Page, req.PageSize)
	response.SuccessWithMeta(c, result.List, meta)
}

// GetServiceTypeTree 获取服务类型树
func (h *ServiceTypeHandler) GetServiceTypeTree(c *gin.Context) {
	tree, err := h.serviceTypeService.GetServiceTypeTree()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, tree)
}

// RegisterRoutes 注册路由
func (h *ServiceTypeHandler) RegisterRoutes(r *gin.RouterGroup) {
	types := r.Group("/service-types")
	{
		types.GET("", h.ListServiceTypes)
		types.POST("", h.CreateServiceType)
		types.GET("/tree", h.GetServiceTypeTree)
		types.GET("/:id", h.GetServiceType)
		types.PUT("/:id", h.UpdateServiceType)
		types.DELETE("/:id", h.DeleteServiceType)
	}
}

// ==================== ServiceOrderHandler ====================

// ServiceOrderHandler 服务订单处理器
type ServiceOrderHandler struct {
	serviceOrderService *services.ServiceOrderService
}

// NewServiceOrderHandler 创建服务订单处理器
func NewServiceOrderHandler() *ServiceOrderHandler {
	return &ServiceOrderHandler{
		serviceOrderService: services.NewServiceOrderService(),
	}
}

// CreateServiceOrder 创建服务订单
func (h *ServiceOrderHandler) CreateServiceOrder(c *gin.Context) {
	var req services.CreateServiceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	createdBy := middleware.GetUserID(c)
	order, err := h.serviceOrderService.CreateServiceOrder(&req, createdBy)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Created(c, order)
}

// UpdateServiceOrder 更新服务订单
func (h *ServiceOrderHandler) UpdateServiceOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的服务订单ID")
		return
	}

	var req services.UpdateServiceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	order, err := h.serviceOrderService.UpdateServiceOrder(uint(id), &req)
	if err != nil {
		if err.Error() == "服务订单不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, order)
}

// UpdateProgress 更新服务进度
func (h *ServiceOrderHandler) UpdateProgress(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的服务订单ID")
		return
	}

	var req services.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	order, err := h.serviceOrderService.UpdateProgress(uint(id), &req)
	if err != nil {
		if err.Error() == "服务订单不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, order)
}

// DeleteServiceOrder 删除服务订单
func (h *ServiceOrderHandler) DeleteServiceOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的服务订单ID")
		return
	}

	if err := h.serviceOrderService.DeleteServiceOrder(uint(id)); err != nil {
		if err.Error() == "服务订单不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// GetServiceOrder 获取服务订单详情
func (h *ServiceOrderHandler) GetServiceOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的服务订单ID")
		return
	}

	order, err := h.serviceOrderService.GetServiceOrder(uint(id))
	if err != nil {
		if err.Error() == "服务订单不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, order)
}

// ListServiceOrders 获取服务订单列表
func (h *ServiceOrderHandler) ListServiceOrders(c *gin.Context) {
	var req services.ListServiceOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	result, err := h.serviceOrderService.ListServiceOrders(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.BuildMeta(result.Total, req.Page, req.PageSize)
	response.SuccessWithMeta(c, result.List, meta)
}

// CreateCommunication 创建沟通纪要
func (h *ServiceOrderHandler) CreateCommunication(c *gin.Context) {
	serviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的服务订单ID")
		return
	}

	var req services.CreateCommunicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	createdBy := middleware.GetUserID(c)
	communication, err := h.serviceOrderService.CreateCommunication(uint(serviceID), &req, createdBy)
	if err != nil {
		if err.Error() == "服务订单不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, communication)
}

// GetCommunications 获取沟通纪要列表
func (h *ServiceOrderHandler) GetCommunications(c *gin.Context) {
	serviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的服务订单ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	communications, total, err := h.serviceOrderService.GetCommunications(uint(serviceID), page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.BuildMeta(total, page, pageSize)
	response.SuccessWithMeta(c, communications, meta)
}

// DeleteCommunication 删除沟通纪要
func (h *ServiceOrderHandler) DeleteCommunication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("communicationId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的沟通纪要ID")
		return
	}

	if err := h.serviceOrderService.DeleteCommunication(uint(id)); err != nil {
		if err.Error() == "沟通纪要不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// GetServiceStats 获取服务统计
func (h *ServiceOrderHandler) GetServiceStats(c *gin.Context) {
	stats, err := h.serviceOrderService.GetServiceStats()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// RegisterRoutes 注册路由
func (h *ServiceOrderHandler) RegisterRoutes(r *gin.RouterGroup) {
	orders := r.Group("/service-orders")
	{
		orders.GET("", h.ListServiceOrders)
		orders.POST("", h.CreateServiceOrder)
		orders.GET("/stats", h.GetServiceStats)
		orders.GET("/:id", h.GetServiceOrder)
		orders.PUT("/:id", h.UpdateServiceOrder)
		orders.PATCH("/:id/progress", h.UpdateProgress)
		orders.DELETE("/:id", h.DeleteServiceOrder)
		orders.POST("/:id/communications", h.CreateCommunication)
		orders.GET("/:id/communications", h.GetCommunications)
	}

	r.DELETE("/communications/:communicationId", h.DeleteCommunication)
}
