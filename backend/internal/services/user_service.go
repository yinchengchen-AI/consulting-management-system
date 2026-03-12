package services

import (
	"errors"

	"consulting-system/internal/database"
	"consulting-system/internal/models"
	"consulting-system/internal/utils"

	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct{}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username   string `json:"username" binding:"required,min=3,max=50"`
	Password   string `json:"password" binding:"required,min=6"`
	Email      string `json:"email" binding:"omitempty,email"`
	Phone      string `json:"phone" binding:"omitempty"`
	RealName   string `json:"real_name" binding:"omitempty,max=50"`
	Department string `json:"department" binding:"omitempty,max=50"`
	RoleIDs    []uint `json:"role_ids" binding:"omitempty"`
	Status     int    `json:"status" binding:"omitempty,oneof=1 2"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email      string `json:"email" binding:"omitempty,email"`
	Phone      string `json:"phone" binding:"omitempty"`
	RealName   string `json:"real_name" binding:"omitempty,max=50"`
	Department string `json:"department" binding:"omitempty,max=50"`
	RoleIDs    []uint `json:"role_ids" binding:"omitempty"`
	Status     int    `json:"status" binding:"omitempty,oneof=1 2"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ListUsersRequest 用户列表请求
type ListUsersRequest struct {
	Username   string `form:"username"`
	RealName   string `form:"real_name"`
	Department string `form:"department"`
	Status     int    `form:"status"`
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=10"`
}

// ListUsersResponse 用户列表响应
type ListUsersResponse struct {
	List  []models.User `json:"list"`
	Total int64         `json:"total"`
}

// CreateUser 创建用户
func (s *UserService) CreateUser(req *CreateUserRequest, createdBy uint) (*models.User, error) {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			return nil, errors.New("邮箱已存在")
		}
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := models.User{
		Username:   req.Username,
		Password:   hashedPassword,
		Email:      req.Email,
		Phone:      req.Phone,
		RealName:   req.RealName,
		Department: req.Department,
		Status:     req.Status,
	}

	if user.Status == 0 {
		user.Status = 1
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	// 关联角色
	if len(req.RoleIDs) > 0 {
		var roles []models.Role
		if err := database.DB.Where("id IN ?", req.RoleIDs).Find(&roles).Error; err == nil {
			database.DB.Model(&user).Association("Roles").Append(roles)
		}
	}

	// 重新加载用户（包含角色）
	database.DB.Preload("Roles").First(&user, user.ID)

	return &user, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(id uint, req *UpdateUserRequest) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 检查邮箱是否被其他用户使用
	if req.Email != "" && req.Email != user.Email {
		var existingUser models.User
		if err := database.DB.Where("email = ? AND id != ?", req.Email, id).First(&existingUser).Error; err == nil {
			return nil, errors.New("邮箱已被其他用户使用")
		}
	}

	// 更新字段
	updates := map[string]interface{}{}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Department != "" {
		updates["department"] = req.Department
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 更新角色
	if len(req.RoleIDs) > 0 {
		var roles []models.Role
		if err := database.DB.Where("id IN ?", req.RoleIDs).Find(&roles).Error; err == nil {
			database.DB.Model(&user).Association("Roles").Clear()
			database.DB.Model(&user).Association("Roles").Append(roles)
		}
	}

	// 重新加载用户
	database.DB.Preload("Roles").First(&user, id)

	return &user, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint) error {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 软删除
	return database.DB.Delete(&user).Error
}

// GetUser 获取用户详情
func (s *UserService) GetUser(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.Preload("Roles").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(req *ListUsersRequest) (*ListUsersResponse, error) {
	var users []models.User
	var total int64

	query := database.DB.Model(&models.User{})

	// 应用过滤条件
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.RealName != "" {
		query = query.Where("real_name LIKE ?", "%"+req.RealName+"%")
	}
	if req.Department != "" {
		query = query.Where("department LIKE ?", "%"+req.Department+"%")
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
	if err := query.Preload("Roles").Offset(offset).Limit(req.PageSize).Find(&users).Error; err != nil {
		return nil, err
	}

	return &ListUsersResponse{
		List:  users,
		Total: total,
	}, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	return database.DB.Model(&user).Update("password", hashedPassword).Error
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(userID uint, newPassword string) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	return database.DB.Model(&user).Update("password", hashedPassword).Error
}
