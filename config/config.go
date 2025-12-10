package config

import (
    "os"
    "strconv"
)

type Config struct {
    AppPort string
    DBHost string
    DBUser string
    DBPass string
    DBName string
    DBUseGorm bool
    JwtSecret string
    JwtTtlHours int
    LogDir string
    LogRotationDays int
    LogMaxAgeDays int
    LogLevel string
    RateLimitRPS int
    RateBurst int
    WorkerPoolSize int
}

func LoadFromEnv() *Config {
    atoi := func(k string, d int) int {
        if v := os.Getenv(k); v != "" {
            if n, err := strconv.Atoi(v); err == nil {
                return n
            }
        }
        return d
    }
    boolFromEnv := func(k string, d bool) bool {
        if v := os.Getenv(k); v != "" {
            b, _ := strconv.ParseBool(v)
            return b
        }
        return d
    }

    return &Config{
        AppPort: getenv("APP_PORT", "8080"),
        DBHost: getenv("DB_HOST", "127.0.0.1:3306"),
        DBUser: getenv("DB_USER", "root"),
        DBPass: getenv("DB_PASS", ""),
        DBName: getenv("DB_NAME", "test_golang"),
        DBUseGorm: boolFromEnv("DB_USE_GORM", false),
        JwtSecret: getenv("JWT_SECRET", "changeme"),
        JwtTtlHours: atoi("JWT_TTL_HOURS", 24),
        LogDir: getenv("LOG_DIR", "./logs"),
        LogRotationDays: atoi("LOG_ROTATION_DAYS", 1),
        LogMaxAgeDays: atoi("LOG_MAX_AGE_DAYS", 30),
        LogLevel: getenv("LOG_LEVEL", "info"),
        RateLimitRPS: atoi("RATE_LIMIT_RPS", 50),
        RateBurst: atoi("RATE_BURST", 100),
        WorkerPoolSize: atoi("WORKER_POOL_SIZE", 20),
    }
}

func getenv(k, d string) string {
    if v := os.Getenv(k); v != "" {
        return v
    }
    return d
}
