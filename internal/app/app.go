package app

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "go-api-starter/config"
    "go-api-starter/internal/pkg/db"
    "go-api-starter/internal/pkg/logger"
    "go-api-starter/internal/pkg/migrate"
    userrepo "go-api-starter/internal/user/repository"
    usersvc "go-api-starter/internal/user/service"
    depspkg "go-api-starter/internal/deps"
    "go-api-starter/internal/server"
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

    if err := migrate.ApplyPending(sqlDB); err != nil {
        return nil, fmt.Errorf("migrate failed: %w", err)
    }

    userRepo := userrepo.NewUserRepo(sqlDB)
    userService := usersvc.NewUserService(userRepo, lg, cfg)

    deps := &depspkg.Dependencies{
        DB: sqlDB,
        GormDB: gormDB,
        UserService: userService,
        Logger: lg,
        Config: cfg,
    }

    r := server.NewRouter(deps)
    return &App{Router: r, Config: cfg}, nil
}
