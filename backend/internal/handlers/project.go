package handlers

import (
	"consulting-system/backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProjectHandler 项目处理器
type ProjectHandler struct {
	db *gorm.DB
}

// NewProjectHandler 创建项目处理器
func NewProjectHandler(db *gorm.DB) *ProjectHandler {
	return &ProjectHandler{db: db}
}

// List 获取项目列表
// @Summary 获取项目列表
// @Description 获取所有项目的列表，支持分页和搜索
// @Tags 项目管理
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
// @Router /projects [get]
func (h *ProjectHandler) List(c *gin.Context) {
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
	query := h.db.Model(&models.Project{}).Preload("Customer").Preload("Manager")

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
	var projects []models.Project
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("created_at DESC").Limit(params.PageSize).Offset(offset).Find(&projects).Error; err != nil {
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
			"list":  projects,
			"total": total,
			"page":  params.Page,
			"size":  params.PageSize,
		},
	})
}

// Get 获取项目详情
// @Summary 获取项目详情
// @Description 根据项目 ID 获取项目详细信息
// @Tags 项目管理
// @Accept json
// @Produce json
// @Param id path string true "项目 ID"
// @Success 200 {object} models.Project
// @Security Bearer
// @Router /projects/{id} [get]
func (h *ProjectHandler) Get(c *gin.Context) {
	id := c.Param("id")

	var project models.Project
	if err := h.db.Preload("Customer").Preload("Contract").Preload("Manager").First(&project, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "项目不存在",
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
		"data":    project,
	})
}

// Create 创建项目
// @Summary 创建项目
// @Description 创建新项目
// @Tags 项目管理
// @Accept json
// @Produce json
// @Param request body models.CreateProjectRequest true "项目信息"
// @Success 201 {object} models.Project
// @Security Bearer
// @Router /projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	project := models.Project{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		CustomerID:  req.CustomerID,
		ContractID:  req.ContractID,
		ManagerID:   req.ManagerID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Budget:      req.Budget,
		Priority:    req.Priority,
	}

	if project.Priority == 0 {
		project.Priority = 3
	}

	if err := h.db.Create(&project).Error; err != nil {
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
		"data":    project,
	})
}

// Update 更新项目
// @Summary 更新项目
// @Description 更新项目信息
// @Tags 项目管理
// @Accept json
// @Produce json
// @Param id path string true "项目 ID"
// @Param request body models.UpdateProjectRequest true "项目信息"
// @Success 200 {object} models.Project
// @Security Bearer
// @Router /projects/{id} [put]
func (h *ProjectHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var project models.Project
	if err := h.db.First(&project, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "项目不存在",
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

	var req models.UpdateProjectRequest
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
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.ManagerID != "" {
		updates["manager_id"] = req.ManagerID
	}
	if req.StartDate != nil {
		updates["start_date"] = req.StartDate
	}
	if req.EndDate != nil {
		updates["end_date"] = req.EndDate
	}
	if req.Budget > 0 {
		updates["budget"] = req.Budget
	}
	if req.ActualCost > 0 {
		updates["actual_cost"] = req.ActualCost
	}
	if req.Progress >= 0 {
		updates["progress"] = req.Progress
	}
	if req.Priority > 0 {
		updates["priority"] = req.Priority
	}
	if req.Deliverables != "" {
		updates["deliverables"] = req.Deliverables
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}

	if len(updates) > 0 {
		if err := h.db.Model(&project).Updates(updates).Error; err != nil {
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
		"data":    project,
	})
}

// Delete 删除项目
// @Summary 删除项目
// @Description 删除指定项目
// @Tags 项目管理
// @Accept json
// @Produce json
// @Param id path string true "项目 ID"
// @Success 204
// @Security Bearer
// @Router /projects/{id} [delete]
func (h *ProjectHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	var project models.Project
	if err := h.db.First(&project, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "项目不存在",
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

	if err := h.db.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败",
			"error":   err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
