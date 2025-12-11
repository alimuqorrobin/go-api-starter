package router

import (
    "golang-api-starter/config"
    "golang-api-starter/internal/database"
    "golang-api-starter/internal/handler"
    "golang-api-starter/internal/middleware"
    "golang-api-starter/internal/repository"
    "golang-api-starter/internal/service"
    "golang-api-starter/pkg/jwt"
    "golang-api-starter/pkg/logger"
    "golang-api-starter/pkg/response"

    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    "golang.org/x/time/rate"
)

func SetupRouter(db database.Database, cfg *config.Config, log *logger.Logger) *gin.Engine {
    r := gin.New()

    // Global middleware
    r.Use(middleware.RecoveryMiddleware(log))
    r.Use(middleware.LoggerMiddleware(log))
    r.Use(middleware.CORSMiddleware())

    // Rate limiter
    rateLimiter := middleware.NewRateLimiterMiddleware(
        rate.Limit(cfg.RateLimitRequests)/rate.Limit(cfg.RateLimitDuration.Minutes()),
        cfg.RateLimitRequests,
    )
    rateLimiter.CleanupRoutine(cfg.RateLimitDuration)
    r.Use(rateLimiter.Limit())

    // Swagger documentation
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Health check
    r.GET("/health", func(c *gin.Context) {
        response.Success(c, 200, "Server is healthy", gin.H{
            "status": "ok",
            "app":    cfg.AppName,
            "env":    cfg.AppEnv,
        })
    })

    // Initialize dependencies
    userRepo := repository.NewUserRepository(db.GetDB())
    workerPool := service.NewWorkerPool(cfg.WorkerPoolSize)
    userService := service.NewUserService(userRepo, workerPool)
    
    tokenService := jwt.NewTokenService(
        cfg.JWTSecret,
        cfg.JWTExpirationHours,
        cfg.JWTRefreshExpirationHours,
    )
    authService := service.NewAuthService(userService, tokenService)

    // Handlers
    authHandler := handler.NewAuthHandler(authService)
    userHandler := handler.NewUserHandler(userService)

    // API v1 routes
    v1 := r.Group("/api/v1")
    {
        // Public routes (no auth required)
        auth := v1.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
            auth.POST("/refresh", authHandler.RefreshToken)
        }

        // Protected routes (auth required)
        protected := v1.Group("")
        protected.Use(middleware.AuthMiddleware(tokenService))
        {
            // Auth routes
            protected.POST("/auth/logout", authHandler.Logout)

            // User routes
            users := protected.Group("/users")
            {
                users.GET("", userHandler.ListUsers)
                users.POST("", userHandler.CreateUser)
                users.GET("/profile", userHandler.GetProfile)
                users.GET("/:id", userHandler.GetUser)
                users.PUT("/:id", userHandler.UpdateUser)
                users.DELETE("/:id", userHandler.DeleteUser)
                users.POST("/bulk", userHandler.BulkCreateUsers)
            }
        }
    }

    return r
}
```

**Penjelasan:**
- Setup semua routes & middleware
- Dependency injection pattern
- Public routes (no auth)
- Protected routes (require JWT)
- Swagger documentation endpoint
- Health check endpoint

**Middleware Order:**
```
Request
  ↓
1. Recovery (catch panics)
  ↓
2. Logger (log requests)
  ↓
3. CORS (handle cross-origin)
  ↓
4. Rate Limiter (check limits)
  ↓
5. Auth (validate JWT) - for protected routes only
  ↓
Handler
```

**Route Groups:**
```
/swagger/*any          → Swagger UI
/health                → Health check

/api/v1/auth/register  → Public
/api/v1/auth/login     → Public
/api/v1/auth/refresh   → Public
/api/v1/auth/logout    → Protected

/api/v1/users          → Protected (list, create)
/api/v1/users/:id      → Protected (get, update, delete)
/api/v1/users/profile  → Protected (current user)
/api/v1/users/bulk     → Protected (bulk create)
```

**Dependency Injection:**
```
Database
  ↓
Repository
  ↓
Service
  ↓
Handler
  ↓
Router
```

---

### ✅ Part 5 Selesai!

Sudah selesai semua file di Part 5:
- ✅ internal/middleware/auth.go (JWT validation)
- ✅ internal/middleware/rate_limit.go (Rate limiting)
- ✅ internal/middleware/logger.go (Request logging)
- ✅ internal/middleware/recovery.go (Panic recovery)
- ✅ internal/middleware/cors.go (CORS headers)
- ✅ internal/handler/auth_handler.go (Auth endpoints)
- ✅ internal/handler/user_handler.go (User endpoints)
- ✅ internal/router/router.go (Route setup)

**Struktur sekarang:**
```
golang-api-starter/
├── cmd/server/main.go
├── config/config.go
├── internal/
│   ├── database/ (4 files)
│   ├── domain/ (3 files)
│   ├── repository/ (2 files)
│   ├── service/ (4 files)
│   ├── middleware/ (5 files) ✅ BARU
│   ├── handler/ (2 files) ✅ BARU
│   └── router/ (1 file) ✅ BARU
├── pkg/
│   ├── response/ (1 file)
│   ├── logger/ (1 file)
│   ├── jwt/ (1 file)
│   └── validator/ (1 file)
├── go.mod
└── .env.example