package handlers

import (
	"strconv"

	"consulting-system/internal/middleware"
	"consulting-system/internal/services"
	"consulting-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// ==================== NoticeHandler ====================

// NoticeHandler 通知处理器
type NoticeHandler struct {
	noticeService *services.NoticeService
}

// NewNoticeHandler 创建通知处理器
func NewNoticeHandler() *NoticeHandler {
	return &NoticeHandler{
		noticeService: services.NewNoticeService(),
	}
}

// CreateNotice 创建通知
func (h *NoticeHandler) CreateNotice(c *gin.Context) {
	var req services.CreateNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	createdBy := middleware.GetUserID(c)
	notice, err := h.noticeService.CreateNotice(&req, createdBy)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, notice)
}

// UpdateNotice 更新通知
func (h *NoticeHandler) UpdateNotice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的通知ID")
		return
	}

	var req services.UpdateNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	notice, err := h.noticeService.UpdateNotice(uint(id), &req)
	if err != nil {
		if err.Error() == "通知不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, notice)
}

// DeleteNotice 删除通知
func (h *NoticeHandler) DeleteNotice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的通知ID")
		return
	}

	if err := h.noticeService.DeleteNotice(uint(id)); err != nil {
		if err.Error() == "通知不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// GetNotice 获取通知详情
func (h *NoticeHandler) GetNotice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的通知ID")
		return
	}

	userID := middleware.GetUserID(c)
	notice, err := h.noticeService.GetNotice(uint(id), userID)
	if err != nil {
		if err.Error() == "通知不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, notice)
}

// ListNotices 获取通知列表
func (h *NoticeHandler) ListNotices(c *gin.Context) {
	var req services.ListNoticesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	userRoles := middleware.GetRoles(c)
	result, err := h.noticeService.ListNotices(&req, userRoles)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.BuildMeta(result.Total, req.Page, req.PageSize)
	response.SuccessWithMeta(c, result.List, meta)
}

// MarkAsRead 标记通知为已读
func (h *NoticeHandler) MarkAsRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的通知ID")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.noticeService.MarkAsRead(uint(id), userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "标记成功"})
}

// GetUnreadCount 获取未读通知数量
func (h *NoticeHandler) GetUnreadCount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	userRoles := middleware.GetRoles(c)

	count, err := h.noticeService.GetUnreadCount(userID, userRoles)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"unread_count": count})
}

// RegisterRoutes 注册路由
func (h *NoticeHandler) RegisterRoutes(r *gin.RouterGroup) {
	notices := r.Group("/notices")
	{
		notices.GET("", h.ListNotices)
		notices.POST("", h.CreateNotice)
		notices.GET("/unread-count", h.GetUnreadCount)
		notices.GET("/:id", h.GetNotice)
		notices.PUT("/:id", h.UpdateNotice)
		notices.DELETE("/:id", h.DeleteNotice)
		notices.POST("/:id/read", h.MarkAsRead)
	}
}

// ==================== DocumentHandler ====================

// DocumentHandler 文档处理器
type DocumentHandler struct {
	documentService *services.DocumentService
}

// NewDocumentHandler 创建文档处理器
func NewDocumentHandler() *DocumentHandler {
	return &DocumentHandler{
		documentService: services.NewDocumentService(),
	}
}

// CreateDocument 创建文档
func (h *DocumentHandler) CreateDocument(c *gin.Context) {
	var req services.CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	createdBy := middleware.GetUserID(c)
	document, err := h.documentService.CreateDocument(&req, createdBy)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, document)
}

// UpdateDocument 更新文档
func (h *DocumentHandler) UpdateDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	var req services.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	document, err := h.documentService.UpdateDocument(uint(id), &req)
	if err != nil {
		if err.Error() == "文档不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, document)
}

// DeleteDocument 删除文档
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	if err := h.documentService.DeleteDocument(uint(id)); err != nil {
		if err.Error() == "文档不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// GetDocument 获取文档详情
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	document, err := h.documentService.GetDocument(uint(id))
	if err != nil {
		if err.Error() == "文档不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, document)
}

// ListDocuments 获取文档列表
func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	var req services.ListDocumentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	result, err := h.documentService.ListDocuments(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.BuildMeta(result.Total, req.Page, req.PageSize)
	response.SuccessWithMeta(c, result.List, meta)
}

// RegisterRoutes 注册路由
func (h *DocumentHandler) RegisterRoutes(r *gin.RouterGroup) {
	documents := r.Group("/documents")
	{
		documents.GET("", h.ListDocuments)
		documents.POST("", h.CreateDocument)
		documents.GET("/:id", h.GetDocument)
		documents.PUT("/:id", h.UpdateDocument)
		documents.DELETE("/:id", h.DeleteDocument)
	}
}

// ==================== ContractHandler ====================

// ContractHandler 合同处理器
type ContractHandler struct {
	contractService *services.ContractService
}

// NewContractHandler 创建合同处理器
func NewContractHandler() *ContractHandler {
	return &ContractHandler{
		contractService: services.NewContractService(),
	}
}

// CreateContract 创建合同
func (h *ContractHandler) CreateContract(c *gin.Context) {
	var req services.CreateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	createdBy := middleware.GetUserID(c)
	contract, err := h.contractService.CreateContract(&req, createdBy)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Created(c, contract)
}

// UpdateContract 更新合同
func (h *ContractHandler) UpdateContract(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的合同ID")
		return
	}

	var req services.UpdateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	contract, err := h.contractService.UpdateContract(uint(id), &req)
	if err != nil {
		if err.Error() == "合同不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, contract)
}

// DeleteContract 删除合同
func (h *ContractHandler) DeleteContract(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的合同ID")
		return
	}

	if err := h.contractService.DeleteContract(uint(id)); err != nil {
		if err.Error() == "合同不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// GetContract 获取合同详情
func (h *ContractHandler) GetContract(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的合同ID")
		return
	}

	contract, err := h.contractService.GetContract(uint(id))
	if err != nil {
		if err.Error() == "合同不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, contract)
}

// ListContracts 获取合同列表
func (h *ContractHandler) ListContracts(c *gin.Context) {
	var req services.ListContractsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	result, err := h.contractService.ListContracts(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.BuildMeta(result.Total, req.Page, req.PageSize)
	response.SuccessWithMeta(c, result.List, meta)
}

// GetExpiringContracts 获取即将到期的合同
func (h *ContractHandler) GetExpiringContracts(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	contracts, err := h.contractService.GetExpiringContracts(days)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, contracts)
}

// SignContract 签署合同
func (h *ContractHandler) SignContract(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的合同ID")
		return
	}

	contract, err := h.contractService.SignContract(uint(id))
	if err != nil {
		if err.Error() == "合同不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Success(c, contract)
}

// TerminateContract 终止合同
func (h *ContractHandler) TerminateContract(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的合同ID")
		return
	}

	contract, err := h.contractService.TerminateContract(uint(id))
	if err != nil {
		if err.Error() == "合同不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Success(c, contract)
}

// GetContractStats 获取合同统计
func (h *ContractHandler) GetContractStats(c *gin.Context) {
	stats, err := h.contractService.GetContractStats()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

// RegisterRoutes 注册路由
func (h *ContractHandler) RegisterRoutes(r *gin.RouterGroup) {
	contracts := r.Group("/contracts")
	{
		contracts.GET("", h.ListContracts)
		contracts.POST("", h.CreateContract)
		contracts.GET("/expiring", h.GetExpiringContracts)
		contracts.GET("/stats", h.GetContractStats)
		contracts.GET("/:id", h.GetContract)
		contracts.PUT("/:id", h.UpdateContract)
		contracts.POST("/:id/sign", h.SignContract)
		contracts.POST("/:id/terminate", h.TerminateContract)
		contracts.DELETE("/:id", h.DeleteContract)
	}
}

// ==================== SettingHandler ====================

// SettingHandler 设置处理器
type SettingHandler struct {
	settingService *services.SettingService
}

// NewSettingHandler 创建设置处理器
func NewSettingHandler() *SettingHandler {
	return &SettingHandler{
		settingService: services.NewSettingService(),
	}
}

// CreateConfig 创建配置
func (h *SettingHandler) CreateConfig(c *gin.Context) {
	var req services.CreateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	config, err := h.settingService.CreateConfig(&req)
	if err != nil {
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Created(c, config)
}

// UpdateConfig 更新配置
func (h *SettingHandler) UpdateConfig(c *gin.Context) {
	key := c.Param("key")

	var req services.UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	config, err := h.settingService.UpdateConfig(key, &req)
	if err != nil {
		if err.Error() == "配置项不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, config)
}

// DeleteConfig 删除配置
func (h *SettingHandler) DeleteConfig(c *gin.Context) {
	key := c.Param("key")

	if err := h.settingService.DeleteConfig(key); err != nil {
		if err.Error() == "配置项不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// GetConfig 获取配置
func (h *SettingHandler) GetConfig(c *gin.Context) {
	key := c.Param("key")

	config, err := h.settingService.GetConfig(key)
	if err != nil {
		if err.Error() == "配置项不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, config)
}

// ListConfigs 获取所有配置
func (h *SettingHandler) ListConfigs(c *gin.Context) {
	configs, err := h.settingService.ListConfigs()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, configs)
}

// ListLogs 获取操作日志列表
func (h *SettingHandler) ListLogs(c *gin.Context) {
	var req services.ListLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	result, err := h.settingService.ListLogs(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.BuildMeta(result.Total, req.Page, req.PageSize)
	response.SuccessWithMeta(c, result.List, meta)
}

// RegisterRoutes 注册路由
func (h *SettingHandler) RegisterRoutes(r *gin.RouterGroup) {
	// 配置管理
	configs := r.Group("/configs")
	{
		configs.GET("", h.ListConfigs)
		configs.POST("", h.CreateConfig)
		configs.GET("/:key", h.GetConfig)
		configs.PUT("/:key", h.UpdateConfig)
		configs.DELETE("/:key", h.DeleteConfig)
	}

	// 操作日志
	r.GET("/logs", h.ListLogs)
}
