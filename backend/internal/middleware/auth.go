package middleware

import (
	"context"
	"strings"
	"time"

	"consulting-system/config"
	"consulting-system/internal/database"
	"consulting-system/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

// JWTClaims JWT声明
type JWTClaims struct {
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID uint, username string, roles []string) (string, string, error) {
	cfg := config.GlobalConfig.JWT

	// 访问令牌
	accessClaims := JWTClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}

	// 刷新令牌
	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.RefreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   string(rune(userID)),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}

	// 存储刷新令牌到Redis
	ctx := context.Background()
	key := "refresh_token:" + string(rune(userID))
	database.RedisClient.Set(ctx, key, refreshTokenString, cfg.RefreshTokenTTL)

	return accessTokenString, refreshTokenString, nil
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证信息")
			c.Abort()
			return
		}

		// 提取Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		claims, err := ParseToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "Token无效或已过期")
			c.Abort()
			return
		}

		// 检查token是否在黑名单中
		ctx := context.Background()
		key := "blacklist:" + parts[1]
		exists, _ := database.RedisClient.Exists(ctx, key).Result()
		if exists > 0 {
			response.Unauthorized(c, "Token已失效")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) uint {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	return userID.(uint)
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}

// GetRoles 从上下文获取用户角色
func GetRoles(c *gin.Context) []string {
	roles, exists := c.Get("roles")
	if !exists {
		return []string{}
	}
	return roles.([]string)
}

// HasRole 检查是否有指定角色
func HasRole(c *gin.Context, role string) bool {
	roles := GetRoles(c)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// IsSuperAdmin 检查是否是超级管理员
func IsSuperAdmin(c *gin.Context) bool {
	return HasRole(c, "super_admin")
}

// RefreshToken 刷新Token
func RefreshToken(refreshToken string) (string, string, error) {
	// 解析刷新令牌
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return "", "", err
	}

	// 从Redis验证刷新令牌
	claims, _ := token.Claims.(jwt.MapClaims)
	userID := uint(claims["sub"].(float64))

	ctx := context.Background()
	key := "refresh_token:" + string(rune(userID))
	storedToken, err := database.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil || storedToken != refreshToken {
		return "", "", err
	}

	// 获取用户信息和角色
	// 这里简化处理，实际应该从数据库查询
	return GenerateToken(userID, "", []string{})
}

// BlacklistToken 将Token加入黑名单
func BlacklistToken(tokenString string, expiration time.Duration) error {
	ctx := context.Background()
	key := "blacklist:" + tokenString
	return database.RedisClient.Set(ctx, key, "1", expiration).Err()
}
