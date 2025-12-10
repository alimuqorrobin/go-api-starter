package server

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "go-api-starter/internal/server/handlers"
    "go-api-starter/internal/server/middleware"
    deps "go-api-starter/internal/deps"
)

func NewRouter(d *deps.Dependencies) *gin.Engine {
    r := gin.New()
    r.Use(gin.Recovery())
    r.Use(middleware.RecoveryWithLogger(d.Logger))
    r.Use(middleware.RequestLogger(d.Logger))
    r.Use(middleware.SecureHeaders())
    r.Use(middleware.NewRateLimiterMiddleware(d.Config.RateLimitRPS, d.Config.RateBurst))

    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status":"ok","time":time.Now()})
    })

    r.POST("/api/register", handlers.NewRegisterHandler(d).Register)
    r.POST("/api/login", handlers.NewAuthHandler(d).Login)

    v1 := r.Group("/api/v1")
    v1.Use(middleware.JWTAuth(d.Config, d.Logger))
    {
        users := v1.Group("/users")
        users.GET("/me", handlers.NewUserHandler(d).Me)
    }

    r.StaticFile("/swagger.yaml", "./docs/swagger.yaml")
    return r
}
