# Golang API Starter - Enterprise Template

Sebuah template API production-ready dengan arsitektur enterprise, built dengan Go 1.22+.

**Status:** ✅ Ready to Use | 🚀 Production-Ready | 📦 Zero External Router Dependencies

---

## 📋 Table of Contents

1. [Quick Start](#quick-start)
2. [Sistem Architecture](#sistem-architecture)
3. [Project Structure](#project-structure)
4. [Features](#features)
5. [API Endpoints](#api-endpoints)
6. [Database Setup](#database-setup)
7. [Migration Guide](#migration-guide)
8. [Logger System](#logger-system)
9. [Scheduler Guide](#scheduler-guide)
10. [Menambah Library Baru](#menambah-library-baru)
11. [Troubleshooting](#troubleshooting)
12. [Development Tips](#development-tips)

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

# 6. Test health check
curl http://localhost:8080/health
```

**Expected Output:**
```
Starting application in development environment
Database connected successfully
Migrations completed successfully
Server starting on port 8080
```

✅ Aplikasi sudah running! 🎉

---

## 🏗️ Sistem Architecture

### **Layers Overview**

```
┌─────────────────────────────────────────┐
│         HTTP Layer (Handlers)           │
│  • Request parsing                      │
│  • Response formatting                  │
│  • Middleware (Auth, CORS, Rate Limit)  │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│       Service Layer (Business Logic)    │
│  • Validation                           │
│  • Business rules                       │
│  • Orchestration                        │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│    Repository Layer (Data Access)       │
│  • Database queries (via interfaces)    │
│  • Data mapping                         │
│  • Error handling                       │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│      Database Layer (Abstraction)       │
│  • GORM ORM                             │
│  • Driver: MySQL/PostgreSQL             │
│  • Connection pooling                   │
└─────────────────────────────────────────┘
```

### **Why This Architecture?**

| Layer | Benefit |
|-------|---------|
| **Handler** | Easy to test, clean separation |
| **Service** | Reusable logic, business rules |
| **Repository** | Database agnostic (swap MySQL ↔ PostgreSQL) |
| **Database** | Abstraction, prevent vendor lock-in |

---

## 📁 Project Structure

```
go-api-starter/
│
├── cmd/
│   └── main.go                    # Application entry point
│
├── config/
│   └── config.go                  # Load .env & config
│
├── internal/
│   ├── domain/
│   │   ├── models.go              # Data structures (User, Product, etc)
│   │   └── errors.go              # Custom errors
│   │
│   ├── repository/                # Data access layer
│   │   ├── interfaces.go          # Interfaces (contracts)
│   │   ├── user_repository.go     # User CRUD
│   │   ├── product_repository.go  # Product CRUD
│   │   └── jwt_token_repository.go # JWT token management
│   │
│   ├── service/                   # Business logic layer
│   │   ├── interfaces.go
│   │   ├── user_service.go
│   │   ├── product_service.go
│   │   └── auth_service.go
│   │
│   ├── http/
│   │   ├── handlers/              # HTTP request handlers
│   │   │   ├── user_handler.go
│   │   │   ├── product_handler.go
│   │   │   └── auth_handler.go
│   │   │
│   │   ├── middleware/            # HTTP middleware
│   │   │   ├── auth.go            # JWT validation
│   │   │   ├── cors.go            # CORS handling
│   │   │   ├── logging.go         # Request logging
│   │   │   └── ratelimit.go       # Rate limiting
│   │   │
│   │   ├── response/
│   │   │   └── response.go        # Standardized responses
│   │   │
│   │   └── router.go              # Route definition
│   │
│   ├── database/
│   │   └── database.go            # Database abstraction
│   │
│   ├── logger/
│   │   └── logger.go              # Zap logger (daily rotation)
│   │
│   ├── scheduler/
│   │   ├── main.go                # Scheduler manager
│   │   └── task.go                # Task implementations
│   │
│   └── migration/
│       ├── migration.go           # Migration manager
│       └── migrations/
│           ├── 001_create_users.go
│           ├── 002_create_products.go
│           └── 003_create_jwt_tokens.go
│
├── pkg/
│   ├── jwt/
│   │   └── jwt.go                 # JWT utilities
│   └── utils/
│       └── helpers.go             # Helper functions
│
├── logs/                          # Generated logs (daily)
│   ├── app.log
│   ├── app-2024-01-15.log.gz
│   └── ...
│
├── .env                           # Environment variables (GITIGNORE)
├── .env.example                   # Template .env
├── docker-compose.yml             # Docker setup
├── Makefile                       # Common commands
├── go.mod                         # Go modules
├── go.sum                         # Dependencies lock
└── README.md                      # This file
```

---

## ✨ Features

### ✅ Core Features

- [x] **Clean Architecture** - Layered with interfaces for decoupling
- [x] **Database Agnostic** - Switch MySQL ↔ PostgreSQL tanpa code change
- [x] **JWT Authentication** - Access token (24h) + Refresh token (7d)
- [x] **CRUD Operations** - User, Product, JWT Token
- [x] **Rate Limiting** - 100 req/minute per IP
- [x] **CORS Support** - Configurable origins
- [x] **Graceful Shutdown** - Clean shutdown with timeout
- [x] **Logger** - Daily rotation with 7-day retention
- [x] **Migrations** - Version control for schema
- [x] **Scheduler** - Non-overlapping background tasks
- [x] **Error Handling** - Centralized error management
- [x] **Request Logging** - All requests logged with timing

### 📊 Configuration

| Feature | Default | Customizable |
|---------|---------|--------------|
| Port | 8080 | ✅ PORT in .env |
| JWT Expiry | 24h | ✅ JWT_EXPIRATION_HOURS |
| Refresh Expiry | 7d | ✅ REFRESH_EXPIRATION_DAYS |
| Rate Limit | 100/min | ✅ In code |
| Log Retention | 7d | ✅ In logger.go |
| CORS Origins | localhost:3000 | ✅ CORS_ALLOWED_ORIGINS |

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
# Create user
POST /api/users
{
  "email": "newuser@example.com",
  "password": "Pass123",
  "name": "John Doe"
}

# Get user
GET /api/users/{id}

# Get all users
GET /api/users

# Update user (requires auth)
PUT /api/users/{id}
Header: Authorization: Bearer <token>
{
  "name": "Jane Doe"
}

# Delete user (requires auth)
DELETE /api/users/{id}
Header: Authorization: Bearer <token>
```

### Products (CRUD)

```bash
# Create product (requires auth)
POST /api/products
Header: Authorization: Bearer <token>
{
  "name": "MacBook Pro",
  "price": 1999.99,
  "stock": 50
}

# Get product
GET /api/products/{id}

# Get all products (paginated)
GET /api/products?page=1&limit=10

# Update product (requires auth)
PUT /api/products/{id}
Header: Authorization: Bearer <token>

# Delete product (requires auth)
DELETE /api/products/{id}
Header: Authorization: Bearer <token>
```

### Health Check

```bash
GET /health
Response: {"status":"healthy"}
```

---

## 💾 Database Setup

### Option 1: Docker (Recommended)

```bash
docker-compose up -d mysql
```

Ini akan:
- Start MySQL container
- Create database `golang_api`
- Port: 3306

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
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=golang_api
```

Atau manual:
```bash
psql postgres
CREATE DATABASE golang_api;
\q
```

### Verify Connection

```bash
# MySQL
mysql -u root -p golang_api -e "SELECT 1"

# PostgreSQL
psql -U postgres -d golang_api -c "SELECT 1"
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
For each migration file:
  - Already in DB? → SKIP (safe!)
  - Not in DB? → RUN & record
    ↓
Tables created ✅
```

**Result:** Data TIDAK pernah dihapus! 🔒

### View Migration Status

```bash
mysql -u root golang_api

# Lihat semua migrations yang sudah jalan
SELECT * FROM migration_records;

# Lihat structure
DESCRIBE users;
DESCRIBE products;
DESCRIBE jwt_tokens;
```

### Add New Migration

**Step 1: Create migration file**

File: `internal/migration/migrations/004_create_categories.go`

```go
package migrations

import (
	"gorm.io/gorm"
)

type CreateCategoriesTable struct{}

func NewCreateCategoriesTable() *CreateCategoriesTable {
	return &CreateCategoriesTable{}
}

func (m *CreateCategoriesTable) Name() string {
	return "004_create_categories_table"
}

func (m *CreateCategoriesTable) Up(db *gorm.DB) error {
	type Category struct {
		ID   uint   `gorm:"primaryKey"`
		Name string `gorm:"type:varchar(255);uniqueIndex"`
	}
	return db.AutoMigrate(&Category{})
}
```

**Step 2: Add model**

File: `internal/domain/models.go`

```go
type Category struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(255);uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

**Step 3: Register migration**

File: `internal/migration/migration.go`

```go
func NewMigrator(db *gorm.DB) *Migrator {
	// ... existing code ...
	
	// Add this line:
	m.migrations = append(m.migrations, migrations.NewCreateCategoriesTable())
	
	return m
}
```

**Step 4: Run**

```bash
go run cmd/main.go
```

Output:
```
Running migration: 004_create_categories_table
✅ Migration completed: 004_create_categories_table
```

**Done!** ✅

---

## 📝 Logger System

### Log Location

```
logs/
├── app.log                  # Today's log (readable)
├── app-2024-01-15.log.gz   # Yesterday (compressed)
├── app-2024-01-14.log.gz
└── ... (max 7 days)
```

### View Logs

```bash
# Real-time log
tail -f logs/app.log

# Old logs (compressed)
zcat logs/app-2024-01-15.log.gz

# Search
grep "error" logs/app.log
grep "user_id=1" logs/app.log
```

### Log Format

```json
{
  "timestamp": "2024-01-16T10:45:17.279+0700",
  "level": "info",
  "caller": "cmd/main.go:31",
  "message": "Starting application in production environment"
}
```

### Customize Retention

Edit `internal/logger/logger.go`:

```go
w := &lumberjack.Logger{
    Filename:   logPath + "/app.log",
    MaxSize:    100,   // Rotate jika >100MB
    MaxBackups: 7,     // Keep 7 files
    MaxAge:     7,     // Keep 7 hari (EDIT INI)
    Compress:   true,  // Compress old files
}
```

Ganti `7` dengan jumlah hari yang Anda inginkan! 📅

---

## ⏰ Scheduler Guide

### How It Works

Scheduler otomatis jalankan tasks tanpa overlapping:

```
Start scheduler
    ↓
Register task (interval: 24 hours)
    ↓
Wait 24 hours
    ↓
Execute task in goroutine
    ↓
Wait untuk selesai
    ↓
Repeat
```

### Built-in Tasks

**Cleanup Expired Tokens** (every 24 hours)
- Delete JWT tokens yang sudah expired
- Prevent database bloat
- Otomatis jalan

### Add New Scheduler Task

**Step 1: Create task file**

File: `internal/scheduler/tasks/email_notification_task.go`

```go
package tasks

import (
	"go-api-starter/internal/scheduler"
	"go.uber.org/zap"
)

type EmailNotificationTask struct {
	logger *zap.SugaredLogger
}

func NewEmailNotificationTask(logger *zap.SugaredLogger) scheduler.Task {
	return &EmailNotificationTask{
		logger: logger,
	}
}

func (t *EmailNotificationTask) Name() string {
	return "send-email-notifications"
}

func (t *EmailNotificationTask) Execute() error {
	t.logger.Info("Executing email notification task")
	
	// Your logic here
	// Example: Send pending notifications
	
	t.logger.Info("Email notification task completed")
	return nil
}
```

**Step 2: Register in main.go**

File: `cmd/main.go`

```go
func setupSchedulerTasks(sch *scheduler.Scheduler, db database.Database, logger *zap.SugaredLogger) {
	// Existing task
	tokenRepo := repository.NewJWTTokenRepository(db)
	cleanupTask := scheduler.NewCleanupExpiredTokensTask(tokenRepo, logger)
	sch.AddTask(24*time.Hour, cleanupTask)

	// Add new task (every 1 hour)
	emailTask := tasks.NewEmailNotificationTask(logger)
	sch.AddTask(1*time.Hour, emailTask)  // ← ADD THIS

	logger.Info("Scheduler tasks registered successfully")
}
```

**Step 3: Run**

```bash
go run cmd/main.go
```

Output:
```
Scheduler started
Scheduler tasks registered successfully
Executing send-email-notifications
Email notification task completed
```

**Done!** ✅

### Schedule Intervals

```go
// Common intervals
time.Second       // 1 second
time.Minute       // 1 minute
5 * time.Minute   // 5 minutes
1 * time.Hour     // 1 hour
24 * time.Hour    // 1 day
7 * 24 * time.Hour // 1 week
```

---

## 📦 Menambah Library Baru

### Scenario: Tambah Email Service (Sendgrid)

### Step 1: Add Dependency

```bash
go get github.com/sendgrid/sendgrid-go
go mod tidy
```

### Step 2: Create Service

File: `internal/service/email_service.go`

```go
package service

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService interface {
	SendEmail(to, subject, body string) error
}

type emailService struct {
	apiKey string
}

func NewEmailService(apiKey string) EmailService {
	return &emailService{apiKey: apiKey}
}

func (s *emailService) SendEmail(to, subject, body string) error {
	from := mail.NewEmail("noreply", "noreply@example.com")
	toEmail := mail.NewEmail("Recipient", to)
	message := mail.NewSingleEmail(from, subject, toEmail, body, "")
	
	client := sendgrid.NewSendClient(s.apiKey)
	_, err := client.Send(message)
	return err
}
```

### Step 3: Add to Config

File: `config/config.go`

```go
type Config struct {
	// ... existing fields ...
	SendgridAPIKey string
}

func LoadConfig() *Config {
	return &Config{
		// ... existing ...
		SendgridAPIKey: getEnv("SENDGRID_API_KEY", ""),
	}
}
```

### Step 4: Add to .env

File: `.env`

```env
SENDGRID_API_KEY=SG.xxxxx...
```

### Step 5: Inject & Use

File: `cmd/main.go`

```go
// After database init
emailService := service.NewEmailService(cfg.SendgridAPIKey)

// Pass to handler atau scheduler
```

### Step 6: Update go.mod

```bash
go mod tidy
go mod verify
```

**Verify dependency:**

```bash
go list -m all | grep sendgrid
# Output: github.com/sendgrid/sendgrid-go v3.x.x
```

---

## 🔍 Troubleshooting

### Database Connection Error

```
Error: Access denied for user 'root'@'localhost'
```

**Solutions:**
1. Check .env file exists di root folder
2. Verify MySQL running: `mysql -u root -p`
3. Update DB_USER & DB_PASSWORD di .env
4. Create database: `CREATE DATABASE golang_api;`

### Migration Error

```
Error 1170: BLOB/TEXT column used in key specification without a key length
```

**Solution:**
- Check `internal/domain/models.go`
- Use `VARCHAR(255)` instead of TEXT for indexed columns
- Text fields hanya untuk non-indexed columns

### Port Already in Use

```
Address already in use: :[8080]
```

**Solution:**
1. Change port di .env: `PORT=8081`
2. Or kill existing process: `lsof -i :8080` (Mac/Linux)

### Logger Not Writing

```
logs/ folder empty
```

**Solution:**
1. Check LOG_PATH di .env
2. Create logs folder: `mkdir logs`
3. Check file permissions: `chmod 755 logs/`

---

## 💡 Development Tips

### Environment Variables

```bash
# Development
ENVIRONMENT=development

# Production
ENVIRONMENT=production
```

### Database Switching

```bash
# From MySQL
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306

# To PostgreSQL (no code change needed!)
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
```

### Debug Mode

Add logging:
```go
log.Debugf("User ID: %d", userID)
```

View logs:
```bash
tail -f logs/app.log | jq 'select(.level=="debug")'
```

### Test Endpoints

Use Postman collection atau curl:

```bash
# Create user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Pass123","name":"Test"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Pass123"}'

# Get token dari response, gunakan di protected endpoints
```

---

## 📚 Best Practices

### Code Organization

✅ **DO:**
- Keep logic di service layer
- Use interfaces untuk decoupling
- Handle errors explicitly
- Log important events

❌ **DON'T:**
- Put logic di handlers
- Import directly, use interfaces
- Ignore errors
- Log sensitive data

### Database

✅ **DO:**
- Use transactions untuk batch operations
- Index frequently queried fields
- Backup regularly
- Use migrations untuk schema changes

❌ **DON'T:**
- Manual SQL jika bisa pakai GORM
- Delete data langsung di MySQL
- Skip migrations

### Security

✅ **DO:**
- Hash passwords (bcrypt)
- Validate inputs
- Use HTTPS in production
- Rotate JWT secrets

❌ **DON'T:**
- Store passwords in plain text
- Trust user input
- Use HTTP in production
- Hardcode secrets

---

## 🚀 Production Checklist

Before deploying:

```
□ Update .env.production dengan actual values
□ Set ENVIRONMENT=production
□ Change JWT_SECRET ke strong random string
□ Set CORS_ALLOWED_ORIGINS ke production domains
□ Setup database backups
□ Configure log rotation
□ Test all endpoints
□ Check error logs
□ Verify rate limiting
□ Test graceful shutdown
```

---

## 📞 Support & Resources

### Documentation
- [Go Language](https://golang.org/doc/)
- [GORM](https://gorm.io/)
- [Zap Logger](https://github.com/uber-go/zap)

### Useful Commands

```bash
# Build
go build ./...

# Format code
go fmt ./...

# Run tests
go test ./...

# View dependencies
go list -m all

# Clean cache
go clean -modcache
```

---

## 📄 License

MIT License - Feel free to use for personal or commercial projects.

---

**Happy coding! 🎉**

Untuk pertanyaan atau issue, check documentation di atas atau review code comments!