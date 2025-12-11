package database

import (
    "database/sql"
    "fmt"

    "go-api-starter/config"
    _ "github.com/lib/pq"
)

func NewPostgresConnection(cfg *config.Config) (*sql.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.DBHost,
        cfg.DBPort,
        cfg.DBUser,
        cfg.DBPassword,
        cfg.DBName,
        cfg.DBSSLMode,
    )

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open postgres connection: %w", err)
    }

    return db, nil
}