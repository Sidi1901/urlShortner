package middleware

import (
	"time"

	"github.com/Sidi1901/urlShortner/internal/logger"
	"github.com/gin-gonic/gin"
)

type LoggerMiddleware struct {
	level string
}

func NewLoggerMiddleware(level string) *LoggerMiddleware {
	return &LoggerMiddleware{level: level}
}

func (m *LoggerMiddleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Request Info
		method := c.Request.Method
		client_ip := c.ClientIP()
		path := c.Request.URL.Path

		c.Next()

		// After Request
		statusCode := c.Writer.Status()
		latency := time.Since(start)

		logger.Log.WithFields(map[string]interface{}{
			"method":    method,
			"path":      path,
			"status":    statusCode,
			"latency":   latency.String(),
			"client_ip": client_ip,
			"level":     m.level,
		}).Info("Incomming request")

	}
}
