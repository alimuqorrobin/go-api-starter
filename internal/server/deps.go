package server

import (
    "database/sql"

    "github.com/yourname/go-backend-enterprise/config"
    "github.com/yourname/go-backend-enterprise/internal/pkg/logger"
    "github.com/yourname/go-backend-enterprise/internal/user/service"
    "gorm.io/gorm"
)

type Dependencies struct {
    DB          *sql.DB
    GormDB      *gorm.DB
    UserService *service.UserService
    Logger      *logger.Logger
    Config      *config.Config
}
