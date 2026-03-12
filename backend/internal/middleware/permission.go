package middleware

import (
	"consulting-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequirePermission 权限检查中间件
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 超级管理员拥有所有权限
		if IsSuperAdmin(c) {
			c.Next()
			return
		}

		roles := GetRoles(c)
		if len(roles) == 0 {
			response.Forbidden(c, "没有操作权限")
			c.Abort()
			return
		}

		// 检查是否有指定权限
		// 实际项目中应该从数据库或缓存查询角色的权限
		// 这里简化处理
		hasPermission := checkPermission(roles, permission)
		if !hasPermission {
			response.Forbidden(c, "没有操作权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission 需要任意一个权限
func RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if IsSuperAdmin(c) {
			c.Next()
			return
		}

		roles := GetRoles(c)
		if len(roles) == 0 {
			response.Forbidden(c, "没有操作权限")
			c.Abort()
			return
		}

		for _, permission := range permissions {
			if checkPermission(roles, permission) {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "没有操作权限")
		c.Abort()
	}
}

// RequireAllPermissions 需要所有权限
func RequireAllPermissions(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if IsSuperAdmin(c) {
			c.Next()
			return
		}

		roles := GetRoles(c)
		if len(roles) == 0 {
			response.Forbidden(c, "没有操作权限")
			c.Abort()
			return
		}

		for _, permission := range permissions {
			if !checkPermission(roles, permission) {
				response.Forbidden(c, "没有操作权限")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// RequireRole 角色检查中间件
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !HasRole(c, role) {
			response.Forbidden(c, "没有操作权限")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireAnyRole 需要任意一个角色
func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, role := range roles {
			if HasRole(c, role) {
				c.Next()
				return
			}
		}
		response.Forbidden(c, "没有操作权限")
		c.Abort()
	}
}

// checkPermission 检查角色是否有指定权限
// 实际项目中应该从数据库或缓存查询
func checkPermission(roles []string, permission string) bool {
	// 简化处理，实际应该查询数据库
	// 超级管理员权限
	for _, role := range roles {
		if role == "super_admin" {
			return true
		}
	}

	// 管理员权限映射
	adminPermissions := map[string]bool{
		"user:view":         true,
		"user:create":       true,
		"user:update":       true,
		"customer:view":     true,
		"customer:create":   true,
		"customer:update":   true,
		"customer:delete":   true,
		"service:view":      true,
		"service:create":    true,
		"service:update":    true,
		"service:delete":    true,
		"invoice:view":      true,
		"invoice:create":    true,
		"invoice:update":    true,
		"receipt:view":      true,
		"receipt:create":    true,
		"receipt:update":    true,
		"statistics:view":   true,
		"notice:view":       true,
		"notice:create":     true,
		"notice:update":     true,
		"notice:delete":     true,
		"document:view":     true,
		"document:create":   true,
		"document:update":   true,
		"document:delete":   true,
		"contract:view":     true,
		"contract:create":   true,
		"contract:update":   true,
		"contract:delete":   true,
		"setting:view":      true,
		"setting:update":    true,
	}

	// 普通用户权限映射
	userPermissions := map[string]bool{
		"user:view":       true,
		"customer:view":   true,
		"service:view":    true,
		"invoice:view":    true,
		"receipt:view":    true,
		"notice:view":     true,
		"document:view":   true,
		"contract:view":   true,
	}

	for _, role := range roles {
		switch role {
		case "admin":
			if adminPermissions[permission] {
				return true
			}
		case "user":
			if userPermissions[permission] {
				return true
			}
		}
	}

	return false
}
