package services

import (
	"errors"

	"consulting-system/internal/database"
	"consulting-system/internal/middleware"
	"consulting-system/internal/models"
	"consulting-system/internal/utils"
)

// AuthService 认证服务
type AuthService struct{}

// NewAuthService 创建认证服务实例
func NewAuthService() *AuthService {
	return &AuthService{}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"`
	User         UserInfo    `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID         uint     `json:"id"`
	Username   string   `json:"username"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	RealName   string   `json:"real_name"`
	Department string   `json:"department"`
	Avatar     string   `json:"avatar"`
	Roles      []string `json:"roles"`
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 查找用户
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).Preload("Roles").First(&user).Error; err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 获取用户角色
	roleCodes := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		roleCodes = append(roleCodes, role.Code)
	}

	// 生成Token
	accessToken, refreshToken, err := middleware.GenerateToken(user.ID, user.Username, roleCodes)
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	// 更新最后登录时间
	now := utils.Now()
	database.DB.Model(&user).Update("last_login_at", &now)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    7200, // 2小时
		User: UserInfo{
			ID:         user.ID,
			Username:   user.Username,
			Email:      user.Email,
			Phone:      user.Phone,
			RealName:   user.RealName,
			Department: user.Department,
			Avatar:     user.Avatar,
			Roles:      roleCodes,
		},
	}, nil
}

// Logout 用户登出
func (s *AuthService) Logout(userID uint, token string) error {
	// 将token加入黑名单
	return middleware.BlacklistToken(token, 7200)
}

// RefreshToken 刷新Token
func (s *AuthService) RefreshToken(refreshToken string) (*LoginResponse, error) {
	accessToken, newRefreshToken, err := middleware.RefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("刷新Token失败")
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    7200,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(userID uint) (*UserInfo, error) {
	var user models.User
	if err := database.DB.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	roleCodes := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		roleCodes = append(roleCodes, role.Code)
	}

	return &UserInfo{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Phone:      user.Phone,
		RealName:   user.RealName,
		Department: user.Department,
		Avatar:     user.Avatar,
		Roles:      roleCodes,
	}, nil
}
