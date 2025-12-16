package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port        string

	// Database
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// JWT
	JWTSecret             string
	JWTExpirationHours    int
	RefreshExpirationDays int

	// Logger
	LogPath string

	// CORS
	CORSAllowedOrigins string
}

func LoadConfig() *Config {
	// Load .env file (optional if not exists)
	// Try multiple paths
	_ = godotenv.Load()                    // Try current dir
	_ = godotenv.Load(".env.local")        // Try .env.local
	_ = godotenv.Load("../../.env")        // Try parent dir

	// Debug: Print loaded environment
	fmt.Println("=== Environment Loaded ===")
	fmt.Printf("DB_DRIVER: %s\n", getEnv("DB_DRIVER", "mysql"))
	fmt.Printf("DB_HOST: %s\n", getEnv("DB_HOST", "localhost"))
	fmt.Printf("DB_USER: %s\n", getEnv("DB_USER", "root"))
	fmt.Printf("ENVIRONMENT: %s\n", getEnv("ENVIRONMENT", "development"))
	fmt.Println("==========================")

	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),

		DBDriver:   getEnv("DB_DRIVER", "mysql"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "golang_api"),

		JWTSecret:             getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpirationHours:    24,
		RefreshExpirationDays: 7,

		LogPath: getEnv("LOG_PATH", "./logs"),

		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// GetDSN - Build database connection string
func GetDSN(cfg *Config) string {
	switch cfg.DBDriver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		)
	case "postgres":
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
		)
	default:
		return ""
	}
}