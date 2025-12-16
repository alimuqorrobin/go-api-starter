package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-api-starter/config"
	"go-api-starter/internal/database"
	"go-api-starter/internal/http/middleware"
	"go-api-starter/internal/http/router"
	"go-api-starter/internal/logger"
	"go-api-starter/internal/migration"
	"go-api-starter/internal/repository"
	"go-api-starter/internal/scheduler"
	"go.uber.org/zap"
)

func main() {
	// ============ LOAD CONFIG ============
	cfg := config.LoadConfig()

	// ============ INITIALIZE LOGGER ============
	log := logger.InitLogger(cfg.LogPath)
	defer log.Sync()

	log.Infof("Starting application in %s environment", cfg.Environment)
	log.Infof("Database: %s://%s:%s/%s", cfg.DBDriver, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// ============ INITIALIZE DATABASE WITH RECOVERY ============
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
		return
	}

	// Check if database connection successful
	if !db.IsConnected() {
		log.Fatalf("Database is not connected properly")
		return
	}

	defer db.Close()

	log.Info("Database connected successfully")

	// ============ RUN MIGRATIONS ============
	migrator := migration.NewMigrator(db.GetDB())
	if err := migrator.Migrate(); err != nil {
		log.Warnf("Migration failed: %v", err)
	} else {
		log.Info("Migrations completed successfully")
	}

	// ============ INITIALIZE RATE LIMITER ============
	limiter := middleware.NewRateLimiter(100, time.Minute)

	// ============ SETUP ROUTER ============
	mux := router.SetupRouter(db, log, limiter, cfg)

	// ============ INITIALIZE SCHEDULER ============
	sch := scheduler.NewScheduler(log)
	sch.Start()

	// Register scheduler tasks
	setupSchedulerTasks(sch, db, log)

	defer sch.Stop()

	// ============ CREATE HTTP SERVER ============
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ============ START SERVER IN GOROUTINE ============
	go func() {
		log.Infof("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// ============ GRACEFUL SHUTDOWN ============
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Infof("Received signal: %v", sig)

	// Graceful shutdown dengan timeout 30 detik
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("Server shutdown error: %v", err)
	}

	log.Info("Server shut down successfully")
	log.Info("Application stopped")
}

// setupSchedulerTasks registers all scheduled tasks
func setupSchedulerTasks(sch *scheduler.Scheduler, db database.Database, logger *zap.SugaredLogger) {
	// Example: Cleanup expired tokens setiap 24 jam
	tokenRepo := repository.NewJWTTokenRepository(db)
	cleanupTask := scheduler.NewCleanupExpiredTokensTask(tokenRepo, logger)
	sch.AddTask(24*time.Hour, cleanupTask)

	logger.Info("Scheduler tasks registered successfully")
}