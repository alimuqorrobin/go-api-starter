package server

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/yourname/go-backend-enterprise/internal/server/handlers"
    "github.com/yourname/go-backend-enterprise/internal/server/middleware"
)

func NewRouter(d *Dependencies) *gin.Engine {
    r := gin.New()
    r.Use(gin.Recovery())
    r.Use(middleware.RecoveryWithLogger(d.Logger))
    r.Use(middleware.RequestLogger(d.Logger))
    r.Use(middleware.SecureHeaders())
    r.Use(middleware.NewRateLimiterMiddleware(d.Config.RateLimitRPS, d.Config.RateBurst))

    r.MaxMultipartMemory = 8 << 20

    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now()})
    })

    r.POST("/api/register", handlers.RegisterHandler(d))
    r.POST("/api/login", handlers.LoginHandler(d))

    v1 := r.Group("/api/v1")
    v1.Use(middleware.JWTAuth(d.Config, d.Logger))
    {
        users := v1.Group("/users")
        users.GET("/me", handlers.MeHandler(d))
    }

    r.StaticFile("/swagger.yaml", "./docs/swagger.yaml")
    return r
}
