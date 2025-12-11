package config

import (
    "log"
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
)

type Config struct {
    // App
    AppName string
    AppEnv  string
    Port    string
    Debug   bool

    // Database
    DBDriver          string
    DBHost            string
    DBPort            string
    DBUser            string
    DBPassword        string
    DBName            string
    DBSSLMode         string
    DBMaxOpenConns    int
    DBMaxIdleConns    int
    DBConnMaxLifetime time.Duration

    // JWT
    JWTSecret                 string
    JWTExpirationHours        int
    JWTRefreshExpirationHours int

    // Rate Limiting
    RateLimitRequests int
    RateLimitDuration time.Duration

    // Concurrency
    WorkerPoolSize    int
    MaxConcurrentJobs int

    // Logging
    LogLevel      string
    LogFilePath   string
    LogMaxSize    int
    LogMaxBackups int
    LogMaxAge     int

    // Migration
    MigrationPath string
}

func LoadConfig() *Config {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️  No .env file found, using environment variables")
    }

    return &Config{
        // App
        AppName: getEnv("APP_NAME", "go-api-starter"),
        AppEnv:  getEnv("APP_ENV", "development"),
        Port:    getEnv("APP_PORT", "8080"),
        Debug:   getEnvBool("APP_DEBUG", true),

        // Database
        DBDriver:          getEnv("DB_DRIVER", "postgres"),
        DBHost:            getEnv("DB_HOST", "localhost"),
        DBPort:            getEnv("DB_PORT", "5432"),
        DBUser:            getEnv("DB_USER", "postgres"),
        DBPassword:        getEnv("DB_PASSWORD", "postgres"),
        DBName:            getEnv("DB_NAME", "golang_api_db"),
        DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
        DBMaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
        DBMaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
        DBConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),

        // JWT
        JWTSecret:                 getEnv("JWT_SECRET", "your-secret-key-minimum-32-characters-long"),
        JWTExpirationHours:        getEnvInt("JWT_EXPIRATION_HOURS", 24),
        JWTRefreshExpirationHours: getEnvInt("JWT_REFRESH_EXPIRATION_HOURS", 168),

        // Rate Limiting
        RateLimitRequests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
        RateLimitDuration: getEnvDuration("RATE_LIMIT_DURATION", 1*time.Minute),

        // Concurrency
        WorkerPoolSize:    getEnvInt("WORKER_POOL_SIZE", 10),
        MaxConcurrentJobs: getEnvInt("MAX_CONCURRENT_JOBS", 100),

        // Logging
        LogLevel:      getEnv("LOG_LEVEL", "info"),
        LogFilePath:   getEnv("LOG_FILE_PATH", "logs/app.log"),
        LogMaxSize:    getEnvInt("LOG_MAX_SIZE", 100),
        LogMaxBackups: getEnvInt("LOG_MAX_BACKUPS", 30),
        LogMaxAge:     getEnvInt("LOG_MAX_AGE", 90),

        // Migration
        MigrationPath: getEnv("MIGRATION_PATH", "file://migrations"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        if boolValue, err := strconv.ParseBool(value); err == nil {
            return boolValue
        }
    }
    return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}