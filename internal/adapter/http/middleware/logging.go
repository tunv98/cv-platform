package middleware

import (
	logger "cv-platform/internal/log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestIDKey is the key used to store request ID in context
const RequestIDKey = "request_id"

// RequestLogging middleware adds request ID and structured logging to each request
func RequestLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate unique request ID
		requestID := uuid.New().String()
		c.Set(RequestIDKey, requestID)

		// Create request-scoped logger with request ID
		reqLogger := logger.With("request_id", requestID)

		// Store logger in context for handlers to use
		ctx := logger.IntoContext(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(ctx)

		// Log incoming request
		reqLogger.Info("incoming request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("remote_addr", c.ClientIP()),
		)

		// Process request
		c.Next()

		// Log response
		duration := time.Since(start)
		reqLogger.Info("request completed",
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.Int("response_size", c.Writer.Size()),
		)

		// Log errors if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				reqLogger.Error("request error", zap.Error(err.Err))
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

// LoggerFromContext retrieves the request-scoped logger from gin context
func LoggerFromContext(c *gin.Context) *zap.Logger {
	return logger.FromContext(c.Request.Context())
}

// SimpleLoggerFromContext retrieves the request-scoped simple logger from gin context
func SimpleLoggerFromContext(c *gin.Context) *logger.SimpleLogger {
	return logger.SimpleFromContext(c.Request.Context())
}
