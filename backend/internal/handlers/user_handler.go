package handlers

import (
	"strconv"

	"consulting-system/internal/middleware"
	"consulting-system/internal/services"
	"consulting-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req services.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	createdBy := middleware.GetUserID(c)
	user, err := h.userService.CreateUser(&req, createdBy)
	if err != nil {
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Created(c, user)
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req services.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	user, err := h.userService.UpdateUser(uint(id), &req)
	if err != nil {
		if err.Error() == "用户不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Success(c, user)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		if err.Error() == "用户不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// GetUser 获取用户详情
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.userService.GetUser(uint(id))
	if err != nil {
		if err.Error() == "用户不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, user)
}

// ListUsers 获取用户列表
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req services.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	result, err := h.userService.ListUsers(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.BuildMeta(result.Total, req.Page, req.PageSize)
	response.SuccessWithMeta(c, result.List, meta)
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req services.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	if err := h.userService.ChangePassword(userID, &req); err != nil {
		response.Error(c, response.CodeValidationErr, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "密码修改成功"})
}

// ResetPassword 重置密码
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	if err := h.userService.ResetPassword(uint(id), req.NewPassword); err != nil {
		if err.Error() == "用户不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "密码重置成功"})
}

// RegisterRoutes 注册路由
func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.GET("", h.ListUsers)
		users.POST("", h.CreateUser)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
		users.POST("/:id/reset-password", h.ResetPassword)
	}

	// 当前用户相关路由
	r.POST("/change-password", h.ChangePassword)
}
