package main

import (
	"log"

	"go-api-starter/config"
	"go-api-starter/internal/infrastructure/database"
	"go-api-starter/internal/infrastructure/logger"
	"go-api-starter/internal/interfaces/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	logger.Init()

	db, err := database.NewMySQL(cfg)
	if err != nil {
		log.Fatal("Database error:", err)
	}

	r := gin.Default()
	routes.RegisterRoutes(r, db, cfg)

	r.Run(":" + cfg.AppPort)
}
