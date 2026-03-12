package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorLogger 错误日志中间件
func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 记录错误
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logrus.WithFields(logrus.Fields{
					"error":  err.Error(),
					"method": c.Request.Method,
					"path":   c.Request.URL.Path,
				}).Error("Request error")
			}
		}
	}
}
