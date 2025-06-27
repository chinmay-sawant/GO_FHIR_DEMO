package middleware

import (
	"time"

	"go-fhir-demo/pkg/logger"
	"go-fhir-demo/pkg/utils/tracer"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

// RequestTracker middleware tracks request and response times.
// Tracing (e.g., Jaeger) should be handled in handlers/services, not middleware.
func RequestTracker() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start a span for the incoming request
		ctx, span := tracer.StartSpan(c.Request.Context(), "HTTP "+c.Request.Method+" "+c.FullPath())
		defer span.End()

		// Replace the request context with the new context containing the span
		c.Request = c.Request.WithContext(ctx)

		start := time.Now()
		c.Next()
		duration := time.Since(start)
		// Optionally, set attributes on the span for better observability
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.path", c.FullPath()),
			attribute.Int("http.status_code", c.Writer.Status()),
			attribute.String("http.client_ip", c.ClientIP()),
			attribute.String("http.user_agent", c.Request.UserAgent()),
			attribute.String("http.request_id", c.GetString("request_id")),
		)

		logger.WithContext(ctx).Infof("Request: %s %s | Status: %d | Duration: %v | IP: %s | UserAgent: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
			c.ClientIP(),
			c.Request.UserAgent(),
		)
	}
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
