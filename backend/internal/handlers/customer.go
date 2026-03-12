package handlers

import (
	"consulting-system/backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CustomerHandler 客户处理器
type CustomerHandler struct {
	db *gorm.DB
}

// NewCustomerHandler 创建客户处理器
func NewCustomerHandler(db *gorm.DB) *CustomerHandler {
	return &CustomerHandler{db: db}
}

// List 获取客户列表
// @Summary 获取客户列表
// @Description 获取所有客户的列表，支持分页和搜索
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Param type query string false "类型筛选"
// @Param level query string false "等级筛选"
// @Param status query string false "状态筛选"
// @Success 200 {object} map[string]interface{}
// @Security Bearer
// @Router /customers [get]
func (h *CustomerHandler) List(c *gin.Context) {
	var params struct {
		Page     int    `form:"page,default=1"`
		PageSize int    `form:"page_size,default=10"`
		Keyword  string `form:"keyword"`
		Type     string `form:"type"`
		Level    string `form:"level"`
		Status   string `form:"status"`
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
	query := h.db.Model(&models.Customer{}).Preload("SalesOwner")

	if params.Keyword != "" {
		query = query.Where("name LIKE ? OR contact_name LIKE ? OR contact_email LIKE ?",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}

	if params.Level != "" {
		query = query.Where("level = ?", params.Level)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
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
	var customers []models.Customer
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("created_at DESC").Limit(params.PageSize).Offset(offset).Find(&customers).Error; err != nil {
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
			"list":  customers,
			"total": total,
			"page":  params.Page,
			"size":  params.PageSize,
		},
	})
}

// Get 获取客户详情
// @Summary 获取客户详情
// @Description 根据客户 ID 获取客户详细信息
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param id path string true "客户 ID"
// @Success 200 {object} models.Customer
// @Security Bearer
// @Router /customers/{id} [get]
func (h *CustomerHandler) Get(c *gin.Context) {
	id := c.Param("id")

	var customer models.Customer
	if err := h.db.Preload("SalesOwner").First(&customer, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "客户不存在",
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
		"data":    customer,
	})
}

// Create 创建客户
// @Summary 创建客户
// @Description 创建新客户
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param request body models.CreateCustomerRequest true "客户信息"
// @Success 201 {object} models.Customer
// @Security Bearer
// @Router /customers [post]
func (h *CustomerHandler) Create(c *gin.Context) {
	var req models.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	customer := models.Customer{
		Name:         req.Name,
		Type:         req.Type,
		Level:        req.Level,
		Industry:     req.Industry,
		Scale:        req.Scale,
		Website:      req.Website,
		Address:      req.Address,
		Description:  req.Description,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		ContactEmail: req.ContactEmail,
		ContactTitle: req.ContactTitle,
		SalesOwnerID: req.SalesOwnerID,
	}

	if customer.Type == "" {
		customer.Type = models.CustomerTypeEnterprise
	}
	if customer.Level == "" {
		customer.Level = models.CustomerLevelC
	}

	if err := h.db.Create(&customer).Error; err != nil {
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
		"data":    customer,
	})
}

// Update 更新客户
// @Summary 更新客户
// @Description 更新客户信息
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param id path string true "客户 ID"
// @Param request body models.UpdateCustomerRequest true "客户信息"
// @Success 200 {object} models.Customer
// @Security Bearer
// @Router /customers/{id} [put]
func (h *CustomerHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var customer models.Customer
	if err := h.db.First(&customer, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "客户不存在",
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

	var req models.UpdateCustomerRequest
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
	if req.Level != "" {
		updates["level"] = req.Level
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Industry != "" {
		updates["industry"] = req.Industry
	}
	if req.Scale != "" {
		updates["scale"] = req.Scale
	}
	if req.Website != "" {
		updates["website"] = req.Website
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.ContactName != "" {
		updates["contact_name"] = req.ContactName
	}
	if req.ContactPhone != "" {
		updates["contact_phone"] = req.ContactPhone
	}
	if req.ContactEmail != "" {
		updates["contact_email"] = req.ContactEmail
	}
	if req.ContactTitle != "" {
		updates["contact_title"] = req.ContactTitle
	}
	if req.SalesOwnerID != "" {
		updates["sales_owner_id"] = req.SalesOwnerID
	}

	if len(updates) > 0 {
		if err := h.db.Model(&customer).Updates(updates).Error; err != nil {
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
		"data":    customer,
	})
}

// Delete 删除客户
// @Summary 删除客户
// @Description 删除指定客户
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param id path string true "客户 ID"
// @Success 204
// @Security Bearer
// @Router /customers/{id} [delete]
func (h *CustomerHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	var customer models.Customer
	if err := h.db.First(&customer, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "客户不存在",
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

	if err := h.db.Delete(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败",
			"error":   err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
