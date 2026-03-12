package middleware

import (
	"consulting-system/backend/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger 日志中间件
func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算延迟
		latency := time.Since(start)

		// 获取客户端 IP
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		// 构建日志字段
		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", raw),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
			zap.String("user-agent", c.Request.UserAgent()),
		}

		// 添加请求 ID
		if requestID := c.GetString("request_id"); requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}

		// 添加用户 ID
		if userID, exists := c.Get("user_id"); exists {
			fields = append(fields, zap.Any("user_id", userID))
		}

		// 添加错误信息
		if len(c.Errors) > 0 {
			fields = append(fields, zap.Strings("errors", c.Errors.Errors()))
		}

		// 根据状态码选择日志级别
		switch {
		case statusCode >= 500:
			log.Error("Server Error", fields...)
		case statusCode >= 400:
			log.Warn("Client Error", fields...)
		default:
			log.Info("Request", fields...)
		}
	}
}
