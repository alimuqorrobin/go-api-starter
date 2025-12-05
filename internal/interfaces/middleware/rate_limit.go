package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-api-starter/internal/utils"
)

var limiter = utils.NewRateLimiter(50, 100)

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			utils.Error(c, http.StatusTooManyRequests, "Too many requests")
			c.Abort()
			return
		}
		c.Next()
	}
}
