package middleware

import (
	"consulting-system/backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS(cfg *config.CORSConfig) gin.HandlerFunc {
	allowOrigins := cfg.GetAllowOrigins()
	allowMethods := cfg.GetAllowMethods()
	allowHeaders := cfg.GetAllowHeaders()

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查是否允许该来源
		allowed := false
		for _, o := range allowOrigins {
			if o == "*" || o == origin {
				allowed = true
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}

		if !allowed && len(allowOrigins) > 0 {
			c.Header("Access-Control-Allow-Origin", allowOrigins[0])
		}

		c.Header("Access-Control-Allow-Methods", joinStrings(allowMethods, ", "))
		c.Header("Access-Control-Allow-Headers", joinStrings(allowHeaders, ", "))
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// 处理预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// joinStrings 连接字符串切片
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
