package services

import (
	"errors"

	"consulting-system/internal/database"
	"consulting-system/internal/models"

	"gorm.io/gorm"
)

// RoleService 角色服务
type RoleService struct{}

// NewRoleService 创建角色服务实例
func NewRoleService() *RoleService {
	return &RoleService{}
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required,max=50"`
	Code        string   `json:"code" binding:"required,max=50"`
	Description string   `json:"description" binding:"omitempty,max=255"`
	Permissions []string `json:"permissions" binding:"omitempty"`
	Status      int      `json:"status" binding:"omitempty,oneof=1 2"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string   `json:"name" binding:"omitempty,max=50"`
	Description string   `json:"description" binding:"omitempty,max=255"`
	Permissions []string `json:"permissions" binding:"omitempty"`
	Status      int      `json:"status" binding:"omitempty,oneof=1 2"`
}

// ListRolesRequest 角色列表请求
type ListRolesRequest struct {
	Name     string `form:"name"`
	Code     string `form:"code"`
	Status   int    `form:"status"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

// ListRolesResponse 角色列表响应
type ListRolesResponse struct {
	List  []models.Role `json:"list"`
	Total int64         `json:"total"`
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(req *CreateRoleRequest) (*models.Role, error) {
	// 检查角色代码是否已存在
	var existingRole models.Role
	if err := database.DB.Where("code = ?", req.Code).First(&existingRole).Error; err == nil {
		return nil, errors.New("角色代码已存在")
	}

	// 检查角色名称是否已存在
	if err := database.DB.Where("name = ?", req.Name).First(&existingRole).Error; err == nil {
		return nil, errors.New("角色名称已存在")
	}

	role := models.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Permissions: req.Permissions,
		Status:      req.Status,
	}

	if role.Status == 0 {
		role.Status = 1
	}

	if err := database.DB.Create(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(id uint, req *UpdateRoleRequest) (*models.Role, error) {
	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		return nil, err
	}

	// 不允许修改超级管理员
	if role.Code == "super_admin" {
		return nil, errors.New("不能修改超级管理员角色")
	}

	// 更新字段
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Permissions != nil {
		updates["permissions"] = req.Permissions
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}

	if err := database.DB.Model(&role).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(id uint) error {
	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		return err
	}

	// 不允许删除超级管理员
	if role.Code == "super_admin" {
		return errors.New("不能删除超级管理员角色")
	}

	// 检查是否有用户使用该角色
	var count int64
	database.DB.Model(&models.User{}).Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("该角色下存在用户，无法删除")
	}

	return database.DB.Delete(&role).Error
}

// GetRole 获取角色详情
func (s *RoleService) GetRole(id uint) (*models.Role, error) {
	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		return nil, err
	}
	return &role, nil
}

// GetRoleByCode 根据代码获取角色
func (s *RoleService) GetRoleByCode(code string) (*models.Role, error) {
	var role models.Role
	if err := database.DB.Where("code = ?", code).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// ListRoles 获取角色列表
func (s *RoleService) ListRoles(req *ListRolesRequest) (*ListRolesResponse, error) {
	var roles []models.Role
	var total int64

	query := database.DB.Model(&models.Role{})

	// 应用过滤条件
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		query = query.Where("code LIKE ?", "%"+req.Code+"%")
	}
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&roles).Error; err != nil {
		return nil, err
	}

	return &ListRolesResponse{
		List:  roles,
		Total: total,
	}, nil
}

// GetAllRoles 获取所有角色
func (s *RoleService) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := database.DB.Where("status = ?", 1).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// AssignRolesToUser 为用户分配角色
func (s *RoleService) AssignRolesToUser(userID uint, roleIDs []uint) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	var roles []models.Role
	if err := database.DB.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return err
	}

	// 清除原有角色并分配新角色
	if err := database.DB.Model(&user).Association("Roles").Clear(); err != nil {
		return err
	}

	if len(roles) > 0 {
		return database.DB.Model(&user).Association("Roles").Append(roles)
	}

	return nil
}

// GetUserRoles 获取用户角色
func (s *RoleService) GetUserRoles(userID uint) ([]models.Role, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	var roles []models.Role
	if err := database.DB.Model(&user).Association("Roles").Find(&roles); err != nil {
		return nil, err
	}

	return roles, nil
}
