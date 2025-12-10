package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "go-api-starter/internal/pkg/logger"
)

func RequestLogger(lg *logger.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        lat := time.Since(start)
        lg.Logger.WithFields(map[string]interface{}{
            "path": c.Request.URL.Path,
            "method": c.Request.Method,
            "status": c.Writer.Status(),
            "latency_ms": lat.Milliseconds(),
        }).Info("request")
    }
}
