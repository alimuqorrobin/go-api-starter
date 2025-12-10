package server

import (
    "database/sql"

    "go-api-starter/config"
    "go-api-starter/internal/pkg/logger"
    "go-api-starter/internal/user/service"
    "gorm.io/gorm"
)

type Dependencies struct {
    DB *sql.DB
    GormDB *gorm.DB
    UserService *service.UserService
    Logger *logger.Logger
    Config *config.Config
}
