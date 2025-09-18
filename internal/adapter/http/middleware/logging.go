package middleware

import (
	logger "cv-platform/pkg/log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDKey is the key used to store request ID in context
const RequestIDKey = "request_id"

// RequestLogging middleware adds request ID and simple logging to each request
func RequestLogging(skipPaths []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for certain paths
		for _, path := range skipPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		start := time.Now()

		// Generate unique request ID
		requestID := uuid.New().String()
		c.Set(RequestIDKey, requestID)

		// Create request-scoped log with request ID
		reqLogger := logger.With("request_id", requestID)

		// Store log in context for handlers to use
		ctx := logger.IntoContext(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(ctx)

		// Create simple log for easier logging
		log := logger.FLogFromContext(ctx)

		// Log incoming request
		query := c.Request.URL.RawQuery
		log.Infof("incoming request: %s %s?%s from %s",
			c.Request.Method, c.Request.URL.Path, query, c.ClientIP())

		// Process request
		c.Next()

		// Log response
		duration := time.Since(start)
		log.Infof("request completed: status=%d, duration=%v, size=%d",
			c.Writer.Status(), duration, c.Writer.Size())

		// Log errors if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Errorf("request error: %v", err.Err)
			}
		}
	}
}

// GetRequestID retrieves the request ID from gin context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// LoggerFromContext retrieves the request-scoped log from gin context
func LoggerFromContext(c *gin.Context) *logger.FLogger {
	return logger.FLogFromContext(c.Request.Context())
}
