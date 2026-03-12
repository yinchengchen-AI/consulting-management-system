package handlers

import (
	"consulting-system/backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ContractHandler 合同处理器
type ContractHandler struct {
	db *gorm.DB
}

// NewContractHandler 创建合同处理器
func NewContractHandler(db *gorm.DB) *ContractHandler {
	return &ContractHandler{db: db}
}

// List 获取合同列表
// @Summary 获取合同列表
// @Description 获取所有合同的列表，支持分页和搜索
// @Tags 合同管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Param status query string false "状态筛选"
// @Param type query string false "类型筛选"
// @Param customer_id query string false "客户 ID 筛选"
// @Success 200 {object} map[string]interface{}
// @Security Bearer
// @Router /contracts [get]
func (h *ContractHandler) List(c *gin.Context) {
	var params struct {
		Page       int    `form:"page,default=1"`
		PageSize   int    `form:"page_size,default=10"`
		Keyword    string `form:"keyword"`
		Status     string `form:"status"`
		Type       string `form:"type"`
		CustomerID string `form:"customer_id"`
	}

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 构建查询
	query := h.db.Model(&models.Contract{}).Preload("Customer").Preload("SalesOwner")

	if params.Keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}

	if params.CustomerID != "" {
		query = query.Where("customer_id = ?", params.CustomerID)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"error":   err.Error(),
		})
		return
	}

	// 获取列表
	var contracts []models.Contract
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("created_at DESC").Limit(params.PageSize).Offset(offset).Find(&contracts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"list":  contracts,
			"total": total,
			"page":  params.Page,
			"size":  params.PageSize,
		},
	})
}

// Get 获取合同详情
// @Summary 获取合同详情
// @Description 根据合同 ID 获取合同详细信息
// @Tags 合同管理
// @Accept json
// @Produce json
// @Param id path string true "合同 ID"
// @Success 200 {object} models.Contract
// @Security Bearer
// @Router /contracts/{id} [get]
func (h *ContractHandler) Get(c *gin.Context) {
	id := c.Param("id")

	var contract models.Contract
	if err := h.db.Preload("Customer").Preload("SalesOwner").First(&contract, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "合同不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data":    contract,
	})
}

// Create 创建合同
// @Summary 创建合同
// @Description 创建新合同
// @Tags 合同管理
// @Accept json
// @Produce json
// @Param request body models.CreateContractRequest true "合同信息"
// @Success 201 {object} models.Contract
// @Security Bearer
// @Router /contracts [post]
func (h *ContractHandler) Create(c *gin.Context) {
	var req models.CreateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	contract := models.Contract{
		Name:         req.Name,
		Type:         req.Type,
		CustomerID:   req.CustomerID,
		Amount:       req.Amount,
		TaxRate:      req.TaxRate,
		PaymentTerms: req.PaymentTerms,
		SignedDate:   req.SignedDate,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Description:  req.Description,
		Terms:        req.Terms,
		SalesOwnerID: req.SalesOwnerID,
		SignedBy:     req.SignedBy,
	}

	if contract.TaxRate == 0 {
		contract.TaxRate = 6
	}

	if err := h.db.Create(&contract).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "创建成功",
		"data":    contract,
	})
}

// Update 更新合同
// @Summary 更新合同
// @Description 更新合同信息
// @Tags 合同管理
// @Accept json
// @Produce json
// @Param id path string true "合同 ID"
// @Param request body models.UpdateContractRequest true "合同信息"
// @Success 200 {object} models.Contract
// @Security Bearer
// @Router /contracts/{id} [put]
func (h *ContractHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var contract models.Contract
	if err := h.db.First(&contract, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "合同不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"error":   err.Error(),
		})
		return
	}

	var req models.UpdateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Amount > 0 {
		updates["amount"] = req.Amount
	}
	if req.TaxRate >= 0 {
		updates["tax_rate"] = req.TaxRate
	}
	if req.PaymentTerms != "" {
		updates["payment_terms"] = req.PaymentTerms
	}
	if req.SignedDate != nil {
		updates["signed_date"] = req.SignedDate
	}
	if req.StartDate != nil {
		updates["start_date"] = req.StartDate
	}
	if req.EndDate != nil {
		updates["end_date"] = req.EndDate
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Terms != "" {
		updates["terms"] = req.Terms
	}
	if req.SalesOwnerID != "" {
		updates["sales_owner_id"] = req.SalesOwnerID
	}
	if req.SignedBy != "" {
		updates["signed_by"] = req.SignedBy
	}

	if len(updates) > 0 {
		if err := h.db.Model(&contract).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新失败",
				"error":   err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    contract,
	})
}

// Delete 删除合同
// @Summary 删除合同
// @Description 删除指定合同
// @Tags 合同管理
// @Accept json
// @Produce json
// @Param id path string true "合同 ID"
// @Success 204
// @Security Bearer
// @Router /contracts/{id} [delete]
func (h *ContractHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	var contract models.Contract
	if err := h.db.First(&contract, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "合同不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"error":   err.Error(),
		})
		return
	}

	if err := h.db.Delete(&contract).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败",
			"error":   err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
