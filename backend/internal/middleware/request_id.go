package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

// RequestID 请求 ID 中间件
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取或生成请求 ID
		requestID := c.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 将请求 ID 存入上下文和响应头
		c.Set("request_id", requestID)
		c.Header(requestIDHeader, requestID)

		c.Next()
	}
}

// GetRequestID 从上下文中获取请求 ID
func GetRequestID(c *gin.Context) string {
	requestID, exists := c.Get("request_id")
	if !exists {
		return ""
	}
	return requestID.(string)
}
