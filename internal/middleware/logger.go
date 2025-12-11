package middleware

import (
    "time"

    "go-api-starter/pkg/logger"
    "github.com/gin-gonic/gin"
)

// LoggerMiddleware logs all requests
func LoggerMiddleware(log *logger.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        // Process request
        c.Next()

        // Calculate latency
        latency := time.Since(start)

        // Get status code
        statusCode := c.Writer.Status()

        // Get client IP
        clientIP := c.ClientIP()

        // Get request method
        method := c.Request.Method

        // Build query string
        if raw != "" {
            path = path + "?" + raw
        }

        // Log request details
        log.Infow("Request processed",
            "method", method,
            "path", path,
            "status", statusCode,
            "latency", latency.String(),
            "client_ip", clientIP,
            "user_agent", c.Request.UserAgent(),
        )

        // Log errors if any
        if len(c.Errors) > 0 {
            log.Errorw("Request errors",
                "method", method,
                "path", path,
                "errors", c.Errors.String(),
            )
        }
    }
}

// {
//   "level": "INFO",
//   "timestamp": "2024-01-15T10:30:45Z",
//   "method": "GET",
//   "path": "/api/v1/users?page=1",
//   "status": 200,
//   "latency": "15ms",
//   "client_ip": "192.168.1.1",
//   "user_agent": "Mozilla/5.0..."
// }