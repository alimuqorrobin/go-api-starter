# Golang API Starter - Enterprise Template

Sebuah template API production-ready dengan arsitektur enterprise, built dengan Go 1.22+.

**Status:** ✅ Ready to Use | 🚀 Production-Ready | 🔒 Race Condition Protected | 🔄 Transaction Safe

**Latest Updates:**
- ✅ Transaction Management dengan automatic rollback
- ✅ Race Condition Protection dengan pessimistic locking
- ✅ Safe stock operations (DeductStock, AddStock, TransferStock)
- ✅ Comprehensive logging dengan daily rotation
- ✅ Scheduler untuk background tasks

---

## 📋 Table of Contents

1. [Quick Start](#quick-start)
2. [Features](#features)
3. [Architecture](#architecture)
4. [Project Structure](#project-structure)
5. [API Endpoints](#api-endpoints)
6. [Database Setup](#database-setup)
7. [Transaction Management](#transaction-management)
8. [Locking & Race Condition](#locking--race-condition)
9. [Migration Guide](#migration-guide)
10. [Logger System](#logger-system)
11. [Scheduler Guide](#scheduler-guide)
12. [Menambah Feature Baru](#menambah-feature-baru)
13. [Troubleshooting](#troubleshooting)

---

## 🚀 Quick Start

### Prerequisites
- Go 1.22+
- MySQL 8.0+ atau PostgreSQL 13+
- Docker & Docker Compose (optional)

### Installation

```bash
# 1. Clone atau buat folder project
mkdir go-api-starter
cd go-api-starter

# 2. Initialize module
go mod init go-api-starter

# 3. Setup database
docker-compose up -d mysql
# atau manual: mysql -u root -p
# CREATE DATABASE golang_api;

# 4. Create .env
cp .env.example .env
# Edit .env sesuai config Anda

# 5. Run application
go run cmd/main.go

# 6. Test
curl http://localhost:8080/health
```

**Expected Output:**
```
2025-12-16T10:45:17.279+0700    info    Starting application in development environment
2025-12-16T10:45:17.279+0700    info    Database connected successfully
2025-12-16T10:45:17.279+0700    info    Migrations completed successfully
2025-12-16T10:45:17.279+0700    info    Server starting on port 8080
```

✅ Aplikasi siap digunakan! 🎉

---

## ✨ Features

### ✅ Core Features

| Feature | Status | Detail |
|---------|--------|--------|
| **Clean Architecture** | ✅ | Layered with interfaces |
| **Database Agnostic** | ✅ | Switch MySQL ↔ PostgreSQL |
| **JWT Authentication** | ✅ | Access + Refresh token |
| **CRUD Operations** | ✅ | User, Product, JWT Token |
| **Transaction Support** | ✅ | Atomic operations dengan WithTx() |
| **Pessimistic Locking** | ✅ | Race condition protected |
| **Rate Limiting** | ✅ | 100 req/minute per IP |
| **CORS Support** | ✅ | Configurable origins |
| **Graceful Shutdown** | ✅ | Clean shutdown |
| **Logger** | ✅ | Daily rotation, 7-day retention |
| **Migrations** | ✅ | Version control schema |
| **Scheduler** | ✅ | Non-overlapping tasks |
| **Error Handling** | ✅ | Centralized management |
| **Request Logging** | ✅ | All requests logged |

---

## 🏗️ Architecture

### Layered Architecture

```
┌─────────────────────────────────────┐
│      HTTP Layer (Handlers)          │
│  • Request/Response handling        │
│  • Middleware (Auth, CORS, Rate)    │
│  • Validation                       │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│    Service Layer (Business Logic)   │
│  • Validation & transformation      │
│  • Business rules                   │
│  • Transaction management           │
│  • Locking logic                    │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│  Repository Layer (Data Access)     │
│  • Database queries                 │
│  • Data mapping                     │
│  • Error handling                   │
│  • (Database agnostic)              │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│    Database Layer (Abstraction)     │
│  • GORM ORM                         │
│  • Driver: MySQL/PostgreSQL         │
│  • Connection pooling               │
└─────────────────────────────────────┘
```

### Benefits

✅ **Separation of Concerns** - Each layer has single responsibility
✅ **Easy Testing** - Mock layers independently
✅ **Maintainability** - Changes isolated to specific layer
✅ **Scalability** - Add features without affecting others
✅ **Database Independence** - Switch database driver easily

---

## 📁 Project Structure

```
go-api-starter/
├── cmd/
│   └── main.go                    # Entry point
│
├── config/
│   └── config.go                  # .env loading & config
│
├── internal/
│   ├── domain/
│   │   ├── models.go              # User, Product, JWTToken
│   │   └── errors.go              # Custom errors
│   │
│   ├── repository/
│   │   ├── interfaces.go          # Contracts
│   │   ├── user_repository.go
│   │   ├── product_repository.go
│   │   └── jwt_token_repository.go
│   │
│   ├── service/
│   │   ├── interfaces.go
│   │   ├── user_service.go
│   │   ├── product_service.go     # ✨ NEW: dengan locking
│   │   └── auth_service.go
│   │
│   ├── http/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   ├── response/
│   │   └── router.go
│   │
│   ├── database/
│   │   └── transaction.go         # ✨ NEW: dengan locking helpers
│   │
│   ├── logger/
│   │   └── logger.go              # Daily rotation
│   │
│   ├── scheduler/
│   │   ├── main.go
│   │   └── task.go
│   │
│   └── migration/
│       ├── migration.go
│       └── migrations/
│
├── pkg/
│   ├── jwt/
│   │   └── jwt.go
│   └── utils/
│       └── helpers.go
│
├── logs/                          # Auto-generated
├── .env                           # Git ignored
├── go.mod
├── go.sum
└── README.md
```

---

## 🔌 API Endpoints

### Authentication

```bash
# Login
POST /api/auth/login
{
  "email": "user@example.com",
  "password": "password123"
}

# Refresh token
POST /api/auth/refresh
{
  "refresh_token": "eyJhbGc..."
}

# Logout (requires auth)
POST /api/auth/logout
Header: Authorization: Bearer <token>
```

### Users (CRUD)

```bash
# Create
POST /api/users
{
  "email": "user@example.com",
  "password": "Pass123",
  "name": "John Doe"
}

# Get by ID
GET /api/users/{id}

# Get all
GET /api/users

# Update (requires auth)
PUT /api/users/{id}

# Delete (requires auth)
DELETE /api/users/{id}
```

### Products (CRUD + Safe Stock Operations)

```bash
# Create (requires auth)
POST /api/products

# Get by ID
GET /api/products/{id}

# Get all (paginated)
GET /api/products?page=1&limit=10

# Update (requires auth)
PUT /api/products/{id}

# Delete (requires auth)
DELETE /api/products/{id}
```

---

## 💾 Database Setup

### Option 1: Docker (Recommended)

```bash
docker-compose up -d mysql
```

### Option 2: Manual MySQL

```bash
mysql -u root -p
CREATE DATABASE golang_api;
exit
```

### Option 3: PostgreSQL

Edit `.env`:
```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=golang_api
```

---

## 🔄 Transaction Management

### Apa Itu Transaction?

Transaction adalah operasi atomic (all-or-nothing) yang memastikan data consistency.

### Usage

```go
// Simple transaction
err := database.WithTx(db, ctx, func(tx *gorm.DB) error {
    // Operation 1
    if err := tx.Create(&user).Error; err != nil {
        return err  // Auto rollback
    }

    // Operation 2
    if err := tx.Update("field", "value").Error; err != nil {
        return err  // Auto rollback
    }

    return nil  // Auto commit all
})
```

### With Panic Recovery

```go
// More robust untuk production
err := database.WithTxRecovery(db, ctx, func(tx *gorm.DB) error {
    // operations
    return nil
})
```

### Example: Safe Order Processing

```go
// Create order + update stock atomically
err := database.WithTx(db, ctx, func(tx *gorm.DB) error {
    // 1. Create order
    if err := tx.Create(&order).Error; err != nil {
        return err
    }

    // 2. Update stock (must succeed or rollback all)
    if err := tx.Model(&product).
        Update("stock", gorm.Expr("stock - ?", qty)).Error; err != nil {
        return err
    }

    return nil
})

// If any step fails: rollback all ✅
// If all succeed: commit all ✅
```

---

## 🔒 Locking & Race Condition Protection

### Apa Itu Race Condition?

```
Tanpa locking:
  Request 1: Read stock (100) → Deduct 60 → Save (40)
  Request 2: Read stock (100) → Deduct 50 → Save (50) ❌ WRONG!

Dengan locking:
  Request 1: LOCK stock → Read (100) → Deduct 60 → Save (40) → UNLOCK
  Request 2: WAIT (locked) → LOCK stock → Read (40) → Error (insufficient) ✅
```

### Locking Functions

```go
// Lock untuk update (pessimistic locking)
database.LockForUpdate(tx, &product, productID)

// Lock untuk read (shared)
database.LockForShare(tx, &product, productID)

// Custom where clause
database.LockForUpdateWhere(tx, &product, "id = ?", id)

// No wait (fail immediately)
database.NoWait(tx, &product, productID)

// Skip locked rows
database.SkipLocked(tx, &product)
```

### Safe Stock Operations

Product service sudah provide safe operations dengan built-in locking:

```go
// Safely deduct stock (race condition protected)
err := productService.DeductStock(ctx, productID, quantity)

// Safely add stock
err := productService.AddStock(ctx, productID, quantity)

// Safely transfer stock between products
err := productService.TransferStock(ctx, fromID, toID, quantity)
```

### Example: Safe Stock Deduction

```go
// Handler: Create order
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    // 1. Safely deduct stock (locked, atomic, race condition protected)
    if err := h.productService.DeductStock(ctx, productID, quantity); err != nil {
        response.BadRequest(w, "Insufficient stock")
        return
    }

    // 2. Create order
    order := &domain.Order{...}
    
    // 3. Success
    response.Created(w, order)
}
```

### When to Use Locking

**✅ Use locking untuk:**
- Stock update (inventory)
- Balance transfer (money)
- Counter increment
- Critical business operations

**❌ Tidak perlu locking untuk:**
- Simple read operations
- Single insert (no concurrent risk)
- Independent operations

### Test Race Condition

```bash
# Go race detector
go run -race cmd/main.go

# Or test
go test -race ./...
```

---

## 🗃️ Migration Guide

### How Migrations Work

Migrations otomatis jalan saat aplikasi start:

```
Application start
  ↓
Check migration_records table
  ↓
For each migration:
  - Already ran? → SKIP (safe!)
  - Not ran? → RUN & record
  ↓
Tables created ✅
```

### View Migration History

```bash
mysql -u root golang_api

# See all migrations
SELECT * FROM migration_records;

# See table structure
DESCRIBE users;
DESCRIBE products;
DESCRIBE jwt_tokens;
```

### Add New Migration

**Step 1: Create migration file**

File: `internal/migration/migrations/004_create_categories.go`

```go
package migrations

import "gorm.io/gorm"

type CreateCategoriesTable struct{}

func NewCreateCategoriesTable() *CreateCategoriesTable {
	return &CreateCategoriesTable{}
}

func (m *CreateCategoriesTable) Name() string {
	return "004_create_categories_table"
}

func (m *CreateCategoriesTable) Up(db *gorm.DB) error {
	type Category struct {
		ID   uint
		Name string `gorm:"uniqueIndex"`
	}
	return db.AutoMigrate(&Category{})
}
```

**Step 2: Add model** (if needed)

**Step 3: Register migration**

File: `internal/migration/migration.go`

```go
m.migrations = append(m.migrations, migrations.NewCreateCategoriesTable())
```

**Step 4: Run**

```bash
go run cmd/main.go
```

---

## 📝 Logger System

### Log Location

```
logs/
├── app.log                  # Today (readable)
├── app-2024-01-15.log.gz   # Yesterday (compressed)
└── ... (max 7 files)
```

### View Logs

```bash
# Real-time
tail -f logs/app.log

# Old logs
zcat logs/app-2024-01-15.log.gz

# Search
grep "error" logs/app.log
```

### Log Format

```json
{
  "timestamp": "2024-01-16T10:45:17.279+0700",
  "level": "info",
  "caller": "cmd/main.go:31",
  "message": "Starting application"
}
```

### Customize Retention

Edit `internal/logger/logger.go`:

```go
MaxAge: 7,  // Change ini untuk hari retention
```

---

## ⏰ Scheduler Guide

### How It Works

Scheduler otomatis jalankan tasks tanpa overlapping:

```
Start scheduler
  ↓
Register task (interval)
  ↓
Wait interval
  ↓
Execute task
  ↓
Repeat
```

### Built-in Tasks

- **Cleanup Expired Tokens** (every 24 hours)

### Add New Task

**Step 1: Create task**

File: `internal/scheduler/tasks/email_task.go`

```go
package tasks

import (
	"golang-api-starter/internal/scheduler"
	"go.uber.org/zap"
)

type EmailTask struct {
	logger *zap.SugaredLogger
}

func NewEmailTask(logger *zap.SugaredLogger) scheduler.Task {
	return &EmailTask{logger: logger}
}

func (t *EmailTask) Name() string {
	return "send-emails"
}

func (t *EmailTask) Execute() error {
	t.logger.Info("Sending emails...")
	// Your logic
	return nil
}
```

**Step 2: Register in main.go**

```go
func setupSchedulerTasks(sch *scheduler.Scheduler, ...) {
	// Existing tasks...
	
	// Add new task
	emailTask := tasks.NewEmailTask(logger)
	sch.AddTask(1*time.Hour, emailTask)
}
```

**Step 3: Run**

```bash
go run cmd/main.go
```

---

## 🆕 Menambah Feature Baru

### Contoh: Tambah Model `Category`

**Step 1: Create model**

File: `internal/domain/models.go`

```go
type Category struct {
	ID        uint
	Name      string `gorm:"uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
```

**Step 2: Create repository interface**

File: `internal/repository/interfaces.go`

```go
type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	GetByID(ctx context.Context, id uint) (*domain.Category, error)
	GetAll(ctx context.Context) ([]domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id uint) error
}
```

**Step 3: Implement repository**

File: `internal/repository/category_repository.go` (copy pattern dari product_repository.go)

**Step 4: Create service interface**

File: `internal/service/interfaces.go`

```go
type CategoryService interface {
	CreateCategory(ctx context.Context, req *domain.CreateCategoryRequest) (*domain.Category, error)
	GetCategory(ctx context.Context, id uint) (*domain.Category, error)
	// ... etc
}
```

**Step 5: Implement service**

File: `internal/service/category_service.go` (copy pattern dari product_service.go)

**Step 6: Create handler**

File: `internal/http/handlers/category_handler.go`

**Step 7: Register routes**

File: `internal/http/router.go`

```go
categoryService := service.NewCategoryService(gormDB, categoryRepo)
categoryHandler := handlers.NewCategoryHandler(categoryService)

mux.HandleFunc("POST /api/categories", categoryHandler.Create)
// ... etc
```

**Step 8: Create migration**

File: `internal/migration/migrations/004_create_categories.go`

**Step 9: Register migration**

File: `internal/migration/migration.go`

```go
m.migrations = append(m.migrations, migrations.NewCreateCategoriesTable())
```

**Step 10: Run**

```bash
go run cmd/main.go
```

---

## 🔧 Troubleshooting

### Database Connection Error

```
Error: Access denied for user 'root'@'localhost'
```

**Solutions:**
1. Check `.env` file exists di root folder
2. Verify MySQL running: `mysql -u root -p`
3. Create database: `CREATE DATABASE golang_api;`

### Migration Error

```
Error: BLOB/TEXT column used in key specification
```

**Solution:** Use `VARCHAR(255)` untuk indexed columns, bukan TEXT

### Port Already in Use

```
Address already in use: :[8080]
```

**Solution:** Change PORT di `.env` atau kill existing process

### Redeclared Interface

```
ProductService redeclared in this block
```

**Solution:** Declare interface hanya di `interfaces.go`, bukan di implementation file

### Missing Module Import

```
package not found
```

**Solution:**
```bash
go mod tidy
go mod download
```

---

## 📊 Production Checklist

Before deploying:

```
□ Update .env.production dengan actual values
□ Set ENVIRONMENT=production
□ Change JWT_SECRET ke strong random string
□ Update CORS_ALLOWED_ORIGINS
□ Setup database backups
□ Test all endpoints
□ Check error logs
□ Test rate limiting
□ Test graceful shutdown
□ Test with -race flag
□ Monitor transaction logs
```

---

## 📞 Common Commands

```bash
# Build
go build ./...

# Format code
go fmt ./...

# Run tests
go test ./...

# Test with race detection
go test -race ./...

# Run with race detection
go run -race cmd/main.go

# View dependencies
go list -m all

# Clean cache
go clean -modcache

# Tidy modules
go mod tidy
```

---

## 🎯 Best Practices

### Code Organization
✅ DO:
- Keep logic di service layer
- Use interfaces untuk decoupling
- Handle errors explicitly
- Use transactions untuk multi-step operations
- Use locking untuk concurrent access

❌ DON'T:
- Put logic di handlers
- Import directly, use interfaces
- Ignore errors
- Skip transactions
- Forget to lock critical sections

### Database
✅ DO:
- Use transactions untuk atomicity
- Lock rows saat concurrent access
- Index frequently queried fields
- Backup regularly
- Use migrations untuk schema changes

❌ DON'T:
- Manual SQL jika bisa pakai GORM
- Skip locking untuk stock operations
- Ignore race conditions
- Delete data langsung di MySQL

---

## 📚 Resources

- [Go Documentation](https://golang.org/doc/)
- [GORM](https://gorm.io/)
- [Zap Logger](https://github.com/uber-go/zap)
- [Go Race Detector](https://golang.org/doc/articles/race_detector)

---

## 📝 Version History

### v1.1.0 (Latest)
- ✨ Added Transaction Management dengan WithTx()
- ✨ Added Pessimistic Locking dengan LockForUpdate()
- ✨ Added Safe Stock Operations (DeductStock, AddStock, TransferStock)
- 🐛 Fixed race condition pada concurrent operations
- 📚 Comprehensive documentation

### v1.0.0
- Initial release
- Basic CRUD operations
- JWT authentication
- Rate limiting
- Logger dengan daily rotation
- Migrations

---

**Happy coding! 🎉**

Untuk pertanyaan atau issue, check documentation di atas atau review code comments!