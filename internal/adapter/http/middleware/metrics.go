package middleware

import (
	"cv-platform/pkg/metrics"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware returns a Gin middleware that collects Prometheus metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get request details
		method := c.Request.Method
		path := c.FullPath()
		status := strconv.Itoa(c.Writer.Status())

		// If no route was matched, use the raw path
		if path == "" {
			path = c.Request.URL.Path
		}

		// Record HTTP metrics
		metrics.RecordHTTPRequest(method, path, status, duration)
	}
}
