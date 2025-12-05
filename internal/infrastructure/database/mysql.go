package database

import (
	"fmt"

	"go-api-starter/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLDB struct {
	Conn *gorm.DB
}

func NewMySQL(cfg *config.Config) (*MySQLDB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &MySQLDB{Conn: db}, nil
}
