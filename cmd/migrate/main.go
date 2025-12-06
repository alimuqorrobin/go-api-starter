package main

import (
    "fmt"
    "os"

    "github.com/yourname/go-backend-enterprise/config"
    "github.com/yourname/go-backend-enterprise/internal/pkg/logger"
    "github.com/yourname/go-backend-enterprise/internal/pkg/db"
    "github.com/yourname/go-backend-enterprise/internal/pkg/migrate"
)

func main() {
    cfg := config.LoadFromEnv()
    lg, _ := logger.NewLogger(cfg.LogDir, cfg.LogRotationDays, cfg.LogMaxAgeDays, cfg.LogLevel)
    dbSQL, _, err := db.InitDB(lg, cfg)
    if err != nil {
        lg.Fatal("db init failed", logger.F("err", err.Error()))
    }
    if len(os.Args) > 1 && os.Args[1] == "down" {
        if err := migrate.RollbackLast(dbSQL); err != nil {
            fmt.Println("rollback failed:", err)
        } else {
            fmt.Println("rollback ok")
        }
        return
    }
    if err := migrate.ApplyPending(dbSQL); err != nil {
        fmt.Println("migrate failed:", err)
    } else {
        fmt.Println("migrate ok")
    }
}
