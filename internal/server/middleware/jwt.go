package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "go-api-starter/config"
    "go-api-starter/internal/pkg/jwt"
    "go-api-starter/internal/pkg/logger"
)

func JWTAuth(cfg *config.Config, lg *logger.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        if auth == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization"})
            return
        }
        token := strings.TrimPrefix(auth, "Bearer ")
        claims, err := jwt.ParseToken(token, cfg.JwtSecret)
        if err != nil {
            lg.Logger.WithField("err", err.Error()).Info("invalid token")
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }
        if sub, ok := claims["sub"]; ok {
            c.Set("sub", sub)
        }
        c.Next()
    }
}
