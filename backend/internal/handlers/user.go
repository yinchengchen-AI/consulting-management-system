package handlers

import (
	"consulting-system/backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserHandler 用户处理器
type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler 创建用户处理器
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// List 获取用户列表
// @Summary 获取用户列表
// @Description 获取所有用户的列表，支持分页和搜索
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Param role query string false "角色筛选"
// @Param status query string false "状态筛选"
// @Success 200 {object} map[string]interface{}
// @Security Bearer
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
	var params struct {
		Page     int    `form:"page,default=1"`
		PageSize int    `form:"page_size,default=10"`
		Keyword  string `form:"keyword"`
		Role     string `form:"role"`
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
	query := h.db.Model(&models.User{})

	if params.Keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR real_name LIKE ?",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.Role != "" {
		query = query.Where("role = ?", params.Role)
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
	var users []models.User
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("created_at DESC").Limit(params.PageSize).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"error":   err.Error(),
		})
		return
	}

	// 转换为响应格式
	var responses []*models.UserResponse
	for i := range users {
		responses = append(responses, users[i].ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"list":  responses,
			"total": total,
			"page":  params.Page,
			"size":  params.PageSize,
		},
	})
}

// Get 获取用户详情
// @Summary 获取用户详情
// @Description 根据用户 ID 获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Success 200 {object} models.UserResponse
// @Security Bearer
// @Router /users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.db.First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "用户不存在",
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
		"data":    user.ToResponse(),
	})
}

// Create 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "用户信息"
// @Success 201 {object} models.UserResponse
// @Security Bearer
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := h.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "创建失败",
			"error":   "用户名已存在",
		})
		return
	}

	// 检查邮箱是否已存在
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "创建失败",
			"error":   "邮箱已被注册",
		})
		return
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		RealName: req.RealName,
		Phone:    req.Phone,
		Role:     req.Role,
		Status:   models.UserStatusActive,
	}

	if user.Role == "" {
		user.Role = models.RoleConsultant
	}

	// 设置密码
	if err := user.SetPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建失败",
			"error":   "密码加密失败",
		})
		return
	}

	// 保存用户
	if err := h.db.Create(&user).Error; err != nil {
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
		"data":    user.ToResponse(),
	})
}

// Update 更新用户
// @Summary 更新用户
// @Description 更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Param request body models.UpdateUserRequest true "用户信息"
// @Success 200 {object} models.UserResponse
// @Security Bearer
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.db.First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "用户不存在",
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

	var req models.UpdateUserRequest
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
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		if err := h.db.Model(&user).Updates(updates).Error; err != nil {
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
		"data":    user.ToResponse(),
	})
}

// Delete 删除用户
// @Summary 删除用户
// @Description 删除指定用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Success 204
// @Security Bearer
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.db.First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "用户不存在",
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

	// 软删除
	if err := h.db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败",
			"error":   err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
