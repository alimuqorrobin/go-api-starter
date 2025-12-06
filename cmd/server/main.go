package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "time"

    "github.com/yourname/go-backend-enterprise/config"
    "github.com/yourname/go-backend-enterprise/internal/app"
    "github.com/yourname/go-backend-enterprise/internal/pkg/logger"
)

func main() {
    cfg := config.LoadFromEnv()

    lg, err := logger.NewLogger(cfg.LogDir, cfg.LogRotationDays, cfg.LogMaxAgeDays, cfg.LogLevel)
    if err != nil {
        fmt.Println("logger init failed:", err)
        os.Exit(1)
    }
    defer lg.Sync()

    appInstance, err := app.New(lg, cfg)
    if err != nil {
        lg.Fatal("failed to init app", logger.F("err", err.Error()))
    }

    addr := fmt.Sprintf(":%s", cfg.AppPort)
    srv := &http.Server{
        Addr:         addr,
        Handler:      appInstance.Router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    lg.Info("http server starting", logger.F("addr", addr))
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            lg.Fatal("listen error", logger.F("err", err.Error()))
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    <-quit
    lg.Info("shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        lg.Fatal("server forced to shutdown", logger.F("err", err.Error()))
    }
    lg.Info("server stopped gracefully")
}
