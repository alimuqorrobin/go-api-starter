package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "golang-api-starter/config"
    "golang-api-starter/internal/database"
    "golang-api-starter/internal/router"
    "golang-api-starter/pkg/logger"
    _ "golang-api-starter/docs"

    "github.com/gin-gonic/gin"
)

// @title Golang API Starter Enterprise
// @version 2.0
// @description Enterprise-grade REST API with JWT Authentication, Rate Limiting, Daily Logging, Recovery, and Concurrency
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
    // Load configuration
    cfg := config.LoadConfig()

    // Initialize logger
    appLogger := logger.NewLogger(cfg)
    defer appLogger.Sync()

    appLogger.Info("🚀 Starting Golang API Starter...")

    // Initialize database
    db, err := database.NewConnection(cfg)
    if err != nil {
        appLogger.Fatal("❌ Failed to connect to database", "error", err)
    }
    defer db.Close()

    appLogger.Info("✅ Database connected successfully")

    // Run migrations
    migrationManager := database.NewMigrationManager(db.GetDB(), cfg.MigrationPath)
    if err := migrationManager.RunMigrations(); err != nil {
        appLogger.Fatal("❌ Migration failed", "error", err)
    }

    appLogger.Info("✅ Migrations completed successfully")

    // Set Gin mode
    if !cfg.Debug {
        gin.SetMode(gin.ReleaseMode)
    }

    // Setup router
    r := router.SetupRouter(db, cfg, appLogger)

    // Create server
    srv := &http.Server{
        Addr:         fmt.Sprintf(":%s", cfg.Port),
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Start server in goroutine
    go func() {
        appLogger.Info(fmt.Sprintf("🚀 Server running on port %s", cfg.Port))
        appLogger.Info(fmt.Sprintf("📚 Swagger UI: http://localhost:%s/swagger/index.html", cfg.Port))
        appLogger.Info(fmt.Sprintf("🏥 Health check: http://localhost:%s/health", cfg.Port))
        
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            appLogger.Fatal("❌ Failed to start server", "error", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    appLogger.Info("⏳ Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        appLogger.Fatal("❌ Server forced to shutdown", "error", err)
    }

    appLogger.Info("✅ Server exited properly")
}