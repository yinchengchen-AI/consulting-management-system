package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRole 用户角色
type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleManager   UserRole = "manager"
	RoleConsultant UserRole = "consultant"
	RoleFinance   UserRole = "finance"
	RoleViewer    UserRole = "viewer"
)

// UserStatus 用户状态
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// User 用户模型
type User struct {
	ID        uint       `gorm:"primary_key;autoIncrement" json:"id"`
	UUID      string     `gorm:"type:uuid;uniqueIndex;not null" json:"uuid"`
	Username  string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email     string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string     `gorm:"type:varchar(255);not null" json:"-"`
	RealName  string     `gorm:"type:varchar(50)" json:"real_name"`
	Phone     string     `gorm:"type:varchar(20)" json:"phone"`
	Avatar    string     `gorm:"type:varchar(255)" json:"avatar"`
	Role      UserRole   `gorm:"type:varchar(20);default:'consultant'" json:"role"`
	Status    UserStatus `gorm:"type:varchar(20);default:'active'" json:"status"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UUID == "" {
		u.UUID = uuid.New().String()
	}
	return nil
}

// SetPassword 设置密码（自动哈希）
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsAdmin 检查是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// UserResponse 用户响应结构（不包含敏感信息）
type UserResponse struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	RealName  string     `json:"real_name"`
	Phone     string     `json:"phone"`
	Avatar    string     `json:"avatar"`
	Role      UserRole   `json:"role"`
	Status    UserStatus `json:"status"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
}

// ToResponse 转换为响应结构
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		RealName:  u.RealName,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		Role:      u.Role,
		Status:    u.Status,
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
	}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string   `json:"username" binding:"required,min=3,max=50"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	RealName string   `json:"real_name"`
	Phone    string   `json:"phone"`
	Role     UserRole `json:"role" binding:"omitempty,oneof=admin manager consultant finance viewer"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	RealName string     `json:"real_name"`
	Phone    string     `json:"phone"`
	Avatar   string     `json:"avatar"`
	Role     UserRole   `json:"role" binding:"omitempty,oneof=admin manager consultant finance viewer"`
	Status   UserStatus `json:"status" binding:"omitempty,oneof=active inactive suspended"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
	User         *UserResponse `json:"user"`
}
