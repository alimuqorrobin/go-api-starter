package db

import (
    "database/sql"
    "fmt"
    "go-api-starter/config"
    "go-api-starter/internal/pkg/logger"
    _ "github.com/go-sql-driver/mysql"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func InitDB(lg *logger.Logger, cfg *config.Config) (*sql.DB, *gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBName)
    sqlDB, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, nil, err
    }
    sqlDB.SetMaxOpenConns(50)
    sqlDB.SetMaxIdleConns(25)
    if err := sqlDB.Ping(); err != nil {
        return nil, nil, err
    }
    lg.Logger.Info("sql.DB connected")
    var gormDB *gorm.DB
    if cfg.DBUseGorm {
        gormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
        if err != nil {
            return sqlDB, nil, err
        }
        lg.Logger.Info("gorm.DB connected")
    }
    return sqlDB, gormDB, nil
}
