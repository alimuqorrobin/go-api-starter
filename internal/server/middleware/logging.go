package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    "github.com/yourname/go-backend-enterprise/internal/pkg/logger"
)

func RequestLogger(lg *logger.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        lat := time.Since(start)
        lg.Info("request",
            logger.F("path", c.Request.URL.Path),
            logger.F("method", c.Request.Method),
            logger.F("status", c.Writer.Status()),
            logger.F("latency_ms", lat.Milliseconds()),
        )
    }
}
