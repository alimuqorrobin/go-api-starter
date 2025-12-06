package app

import (
    "database/sql"
    "fmt"

    "github.com/gin-gonic/gin"
    "github.com/yourname/go-backend-enterprise/config"
    "github.com/yourname/go-backend-enterprise/internal/pkg/db"
    "github.com/yourname/go-backend-enterprise/internal/pkg/logger"
    userrepo "github.com/yourname/go-backend-enterprise/internal/user/repository"
    usersvc "github.com/yourname/go-backend-enterprise/internal/user/service"
    "github.com/yourname/go-backend-enterprise/internal/server"
    "github.com/yourname/go-backend-enterprise/internal/pkg/migrate"
)

type App struct {
    Router *gin.Engine
    Config *config.Config
}

func New(lg *logger.Logger, cfg *config.Config) (*App, error) {
    sqlDB, gormDB, err := db.InitDB(lg, cfg)
    if err != nil {
        return nil, err
    }

    // run migrations (blocking). In production consider running migrations as separate job.
    if err := migrate.ApplyPending(sqlDB); err != nil {
        return nil, fmt.Errorf("migrate failed: %w", err)
    }

    userRepo := userrepo.NewUserRepo(sqlDB)
    userService := usersvc.NewUserService(userRepo, lg, cfg)

    deps := &server.Dependencies{
        DB:          sqlDB,
        GormDB:      gormDB,
        UserService: userService,
        Logger:      lg,
        Config:      cfg,
    }

    r := server.NewRouter(deps)
    return &App{Router: r, Config: cfg}, nil
}
