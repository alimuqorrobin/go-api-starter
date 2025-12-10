package deps

import (
    "database/sql"

    "gorm.io/gorm"
    "go-api-starter/config"
    "go-api-starter/internal/pkg/logger"
    usersvc "go-api-starter/internal/user/service"
)

type Dependencies struct {
    DB          *sql.DB
    GormDB      *gorm.DB
    UserService *usersvc.UserService
    Logger      *logger.Logger
    Config      *config.Config
}
