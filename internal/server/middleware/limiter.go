package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

func NewRateLimiterMiddleware(rps int, burst int) gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Limit(rps), burst)
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{"status": "error", "message": "Too many requests"})
            c.Abort()
            return
        }
        c.Next()
    }
}
