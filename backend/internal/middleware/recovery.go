package middleware

import (
	"consulting-system/backend/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery 错误恢复中间件
func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取请求信息
				fields := []zap.Field{
					zap.Any("error", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("ip", c.ClientIP()),
				}

				// 添加请求 ID
				if requestID := c.GetString("request_id"); requestID != "" {
					fields = append(fields, zap.String("request_id", requestID))
				}

				log.Error("Panic recovered", fields...)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "Internal Server Error",
					"error":   "服务器内部错误",
				})
			}
		}()

		c.Next()
	}
}
