JWT Auth (Login & Protected Routes)
Modular Migration Engine (bisa add migration cukup dengan file baru)
Repository Pattern
Service Layer
Centralized Error Handling + Panic Recovery
Daily Rotated Logging
Worker Pool (untuk traffic tinggi)
Secure Middleware (Helmet-like)
Rate Limiter
Built-in Transaction Manager
Swagger API Docs
Gin HTTP Router (enterprise ready)
go-backend-enterprise/
│
├── cmd/
│   ├── server/        → Entry point HTTP server
│   └── migrate/       → CLI migration runner
│
├── internal/
│   ├── app/           → Dependency injection / service builder
│   ├── core/
│   │   ├── config/    → Load .env
│   │   ├── logger/    → Logger (daily file logs)
│   │   ├── db/        → DB connection & transaction handler
│   │   └── concurrency/ → Worker pool
│   │
│   ├── server/
│   │   ├── router/    → Route groups & middleware stacks
│   │   ├── http/
│   │   │   ├── handler/ → Controller/handler
│   │   │   └── middleware/ → Security, logger, rate-limit
│   │
│   ├── repository/    → DB queries
│   ├── service/       → Business logic
│   └── migration/     → Migration registry system
│
├── migrations/         → SQL files (`*.sql`) auto-loaded
│
├── docs/
│   └── swagger.yaml    → API documentation
│
├── .env.example
├── go.mod
├── Makefile
└── README.md
Install Dependency
    - go mod tidy
Migration Management
    - go run ./cmd/migrate
Menjalankan Server
    - go run ./cmd/server
API Documentation (Swagger)
    http://localhost:8080/swagger