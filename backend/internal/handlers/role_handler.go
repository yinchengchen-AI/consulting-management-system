package handlers

import (
	"strconv"

	"consulting-system/internal/services"
	"consulting-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// RoleHandler 角色处理器
type RoleHandler struct {
	roleService *services.RoleService
}

// NewRoleHandler 创建角色处理器
func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleService: services.NewRoleService(),
	}
}

// CreateRole 创建角色
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req services.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	role, err := h.roleService.CreateRole(&req)
	if err != nil {
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Created(c, role)
}

// UpdateRole 更新角色
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	var req services.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	role, err := h.roleService.UpdateRole(uint(id), &req)
	if err != nil {
		if err.Error() == "角色不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Success(c, role)
}

// DeleteRole 删除角色
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	if err := h.roleService.DeleteRole(uint(id)); err != nil {
		if err.Error() == "角色不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.Error(c, response.CodeConflict, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// GetRole 获取角色详情
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	role, err := h.roleService.GetRole(uint(id))
	if err != nil {
		if err.Error() == "角色不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, role)
}

// ListRoles 获取角色列表
func (h *RoleHandler) ListRoles(c *gin.Context) {
	var req services.ListRolesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	result, err := h.roleService.ListRoles(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.BuildMeta(result.Total, req.Page, req.PageSize)
	response.SuccessWithMeta(c, result.List, meta)
}

// GetAllRoles 获取所有角色
func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.roleService.GetAllRoles()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, roles)
}

// AssignRolesToUser 为用户分配角色
func (h *RoleHandler) AssignRolesToUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req struct {
		RoleIDs []uint `json:"role_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	if err := h.roleService.AssignRolesToUser(uint(userID), req.RoleIDs); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "角色分配成功"})
}

// GetUserRoles 获取用户角色
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	roles, err := h.roleService.GetUserRoles(uint(userID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, roles)
}

// RegisterRoutes 注册路由
func (h *RoleHandler) RegisterRoutes(r *gin.RouterGroup) {
	roles := r.Group("/roles")
	{
		roles.GET("", h.ListRoles)
		roles.POST("", h.CreateRole)
		roles.GET("/all", h.GetAllRoles)
		roles.GET("/:id", h.GetRole)
		roles.PUT("/:id", h.UpdateRole)
		roles.DELETE("/:id", h.DeleteRole)
	}

	// 用户角色分配
	r.POST("/users/:userId/roles", h.AssignRolesToUser)
	r.GET("/users/:userId/roles", h.GetUserRoles)
}
