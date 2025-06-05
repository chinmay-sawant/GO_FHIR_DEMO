package middleware

import (
	"time"

	"go-fhir-demo/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RequestTracker middleware tracks request and response times
func RequestTracker() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Custom log format for request tracking
		logger.Infof("Request: %s %s | Status: %d | Duration: %v | IP: %s | UserAgent: %s",
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Request.UserAgent(),
		)
		return ""
	})
}

// RequestTimer middleware adds request timing information to context
func RequestTimer() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Set start time in context
		c.Set("request_start_time", startTime)

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Set duration in context for potential use in handlers
		c.Set("request_duration", duration)

		// Log request completion
		logger.Debugf("Request completed: %s %s in %v",
			c.Request.Method,
			c.Request.URL.Path,
			duration,
		)
	}
}

// CORS middleware for handling Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// ErrorHandler middleware for centralized error handling
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle any errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			logger.Errorf("Request error: %v", err.Err)

			// Don't override status if it's already set
			if c.Writer.Status() == 200 {
				c.JSON(500, gin.H{
					"error":   "Internal server error",
					"message": err.Error(),
				})
			}
		}
	}
}
