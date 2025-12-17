package router

import (
	"net/http"

	"go-api-starter/config"
	"go-api-starter/internal/database"
	"go-api-starter/internal/http/handlers"
	"go-api-starter/internal/http/middleware"
	"go-api-starter/internal/repository"
	"go-api-starter/internal/service"
	"go-api-starter/pkg/jwt"
	"go.uber.org/zap"
)

func SetupRouter(
	db database.Database,
	logger *zap.SugaredLogger,
	limiter *middleware.RateLimiter,
	cfg *config.Config,
) http.Handler {
	mux := http.NewServeMux()

	// Get GORM DB instance
	gormDB := db.GetDB()

	// ============ INITIALIZE REPOSITORIES ============
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	tokenRepo := repository.NewJWTTokenRepository(db)

	// ============ INITIALIZE SERVICES ============
	// FIXED: Pass gormDB to NewProductService
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(gormDB, productRepo) // ← ADD gormDB parameter
	jwtMgr := jwt.NewManager()
	authService := service.NewAuthService(userRepo, tokenRepo, jwtMgr)

	// ============ INITIALIZE HANDLERS ============
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	authHandler := handlers.NewAuthHandler(authService)

	// ============ HEALTH CHECK ENDPOINT ============
	mux.HandleFunc("GET /health", handleHealth)

	// ============ PUBLIC ROUTES (NO AUTH) ============

	// Authentication
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/refresh", authHandler.RefreshToken)

	// User (registration adalah public)
	mux.HandleFunc("POST /api/users", userHandler.Create)
	mux.HandleFunc("GET /api/users/{id}", userHandler.GetByID)
	mux.HandleFunc("GET /api/users", userHandler.GetAll)

	// Product (GET adalah public)
	mux.HandleFunc("GET /api/products/{id}", productHandler.GetByID)
	mux.HandleFunc("GET /api/products", productHandler.GetAll)

	// ============ PROTECTED ROUTES (REQUIRES AUTH) ============

	// Auth
	mux.Handle("POST /api/auth/logout", 
		middleware.AuthMiddleware(authService)(
			http.HandlerFunc(authHandler.Logout),
		),
	)

	// User (update & delete)
	mux.Handle("PUT /api/users/{id}", 
		middleware.AuthMiddleware(authService)(
			http.HandlerFunc(userHandler.Update),
		),
	)
	mux.Handle("DELETE /api/users/{id}", 
		middleware.AuthMiddleware(authService)(
			http.HandlerFunc(userHandler.Delete),
		),
	)

	// Product (create, update, delete)
	mux.Handle("POST /api/products", 
		middleware.AuthMiddleware(authService)(
			http.HandlerFunc(productHandler.Create),
		),
	)
	mux.Handle("PUT /api/products/{id}", 
		middleware.AuthMiddleware(authService)(
			http.HandlerFunc(productHandler.Update),
		),
	)
	mux.Handle("DELETE /api/products/{id}", 
		middleware.AuthMiddleware(authService)(
			http.HandlerFunc(productHandler.Delete),
		),
	)

	// ============ APPLY GLOBAL MIDDLEWARE ============
	var base http.Handler = mux
	base = middleware.CORS(cfg.CORSAllowedOrigins)(base)
	base = middleware.Logging(logger)(base)
	base = limiter.Limit(base)

	return base
}

// handleHealth returns application health status
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}