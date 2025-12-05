package routes

import (
	"go-api-starter/config"
	"go-api-starter/internal/domain/user"
	"go-api-starter/internal/infrastructure/database"
	"go-api-starter/internal/interfaces/controllers"
	"go-api-starter/internal/interfaces/middleware"
	"go-api-starter/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *database.MySQLDB, cfg *config.Config) {
	mysql := db.Conn

	// Repositories
	var userRepo user.Repository = user.NewGormRepository(mysql)

	// Usecases
	userUC := usecase.NewUserUsecase(userRepo, mysql)

	// Controllers
	userController := controllers.NewUserController(userUC)

	// Public routes
	public := r.Group("/api")
	{
		public.GET("/user/:id", userController.GetUser)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth(cfg))
	{
		protected.POST("/user", userController.CreateUser)
	}
}
