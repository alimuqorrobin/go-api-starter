# 🚀 Golang API Starter - Enterprise Edition

Enterprise-grade REST API starter template dengan Go (Golang) yang dilengkapi dengan fitur production-ready!

## ✨ Features

### 🔐 Authentication & Authorization
- ✅ JWT-based authentication
- ✅ Access token & refresh token
- ✅ Password hashing dengan bcrypt
- ✅ Protected routes dengan middleware

### 🚦 Rate Limiting
- ✅ Token bucket algorithm
- ✅ Per-IP rate limiting
- ✅ Configurable limits
- ✅ Automatic cleanup

### 📝 Logging
- ✅ Daily log rotation (lumberjack)
- ✅ Structured logging (zap)
- ✅ Request/response logging
- ✅ Error tracking dengan stack trace
- ✅ Multiple output (console & file)

### 🛡️ Error Handling
- ✅ Recovery middleware untuk panic
- ✅ Generic JSON response
- ✅ Consistent error format
- ✅ Graceful shutdown

### 🔄 Concurrency
- ✅ Worker pool pattern
- ✅ Bulk operations
- ✅ Goroutine management
- ✅ Context-aware cancellation

### 📚 Documentation
- ✅ Auto-generated Swagger docs
- ✅ Interactive API testing
- ✅ Request/response examples

### 🗄️ Database
- ✅ Database abstraction layer
- ✅ Support PostgreSQL & MySQL
- ✅ Migration system
- ✅ Connection pooling
- ✅ **Tidak tergantung ORM!**

## 🔧 Prerequisites

- **Go** 1.21 or higher
- **PostgreSQL** 12+ atau **MySQL** 8+
- **Make** (optional, tapi recommended)
- **Docker** (optional)

## 📦 Installation

### Quick Start (5 menit)
```bash
# 1. Clone repository
git clone https://golang-api-starter.git
cd golang-api-starter

# 2. Setup environment
make init

# 3. Edit .env file
nano .env

# 4. Run setup (install tools, create db, migrate)
make setup

# 5. Run application
make run
```

### Manual Setup
```bash
# 1. Clone
git clone https://golang-api-starter.git
cd golang-api-starter

# 2. Copy .env
cp .env.example .env

# 3. Edit .env dengan konfigurasi database Anda
nano .env

# 4. Install dependencies
go mod download

# 5. Install tools
go install github.com/swaggo/swag/cmd/swag@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 6. Create database
createdb golang_api_db

# 7. Run migrations
make migrate-up

# 8. Generate Swagger docs
make swagger

# 9. Run application
make run
```

## 🚀 Usage

### Server akan berjalan di:
- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health

### Available Endpoints

#### Authentication
```bash
POST   /api/v1/auth/register   # Register new user
POST   /api/v1/auth/login      # Login user
POST   /api/v1/auth/refresh    # Refresh access token
POST   /api/v1/auth/logout     # Logout user (protected)
```

#### Users
```bash
GET    /api/v1/users           # List all users (protected, paginated)
POST   /api/v1/users           # Create new user (protected)
GET    /api/v1/users/:id       # Get user by ID (protected)
PUT    /api/v1/users/:id       # Update user (protected)
DELETE /api/v1/users/:id       # Delete user (protected)
GET    /api/v1/users/profile   # Get current user profile (protected)
POST   /api/v1/users/bulk      # Bulk create users (protected)
```

## 💡 Usage Examples

### 1. Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123",
    "full_name": "John Doe"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com"
    }
  }
}
```

### 3. Get Profile (Protected)
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 4. List Users dengan Pagination
```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 🗄️ Database Migration

### Create New Migration
```bash
make migrate-create name=create_products_table
```

Akan membuat 2 file:
- `migrations/000003_create_products_table.up.sql`
- `migrations/000003_create_products_table.down.sql`

### Run Migrations
```bash
make migrate-up
```

### Rollback
```bash
make migrate-down
```

### Check Version
```bash
make migrate-version
```

## 📝 Configuration

Edit file `.env`:
```env
# Application
APP_NAME=golang-api-starter
APP_PORT=8080

# Database
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=golang_api_db

# JWT
JWT_SECRET=your-super-secret-jwt-key-min-32-chars
JWT_EXPIRATION_HOURS=24

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=1m

# Logging
LOG_LEVEL=info
LOG_FILE_PATH=logs/app.log
```

## 🧪 Testing
```bash
# Run all tests
make test

# Run with coverage
make test-coverage
```

## 📁 Project Structure
```
golang-api-starter/
├── cmd/server/main.go           # Entry point
├── config/config.go             # Configuration
├── internal/
│   ├── database/                # Database connection & migration
│   ├── domain/                  # Business entities & DTOs
│   ├── repository/              # Data access layer
│   ├── service/                 # Business logic
│   ├── handler/                 # HTTP handlers
│   ├── middleware/              # HTTP middleware
│   └── router/                  # Route definitions
├── pkg/
│   ├── response/                # Generic JSON response
│   ├── logger/                  # Logger utility
│   ├── jwt/                     # JWT token service
│   └── validator/               # Input validation
├── migrations/                   # SQL migration files
├── docs/                        # Swagger documentation
├── Makefile                     # Build automation
└── README.md
```

## 🛠️ Development

### Hot Reload
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
make dev
```

### Generate Swagger
```bash
make swagger
```

### Format Code
```bash
make fmt
```

### Lint
```bash
make lint
```

## 🐳 Docker

### Build Image
```bash
make docker-build
```

### Run Container
```bash
make docker-run
```

### Docker Compose
```bash
make docker-compose-up
```

## 📖 Makefile Commands

| Command | Description |
|---------|-------------|
| `make help` | Show all commands |
| `make init` | Initialize project |
| `make setup` | Complete setup |
| `make run` | Run application |
| `make build` | Build binary |
| `make test` | Run tests |
| `make swagger` | Generate Swagger docs |
| `make migrate-up` | Run migrations |
| `make migrate-down` | Rollback migration |
| `make migrate-create` | Create new migration |
| `make db-create` | Create database |
| `make db-reset` | Reset database |

## 🔒 Security

- ✅ Password hashing dengan bcrypt
- ✅ JWT token validation
- ✅ Rate limiting per IP
- ✅ CORS middleware
- ✅ Input validation
- ✅ SQL injection prevention (prepared statements)

## 📊 Monitoring & Logging

### View Logs
```bash
# Real-time
tail -f logs/app.log

# Last 100 lines
tail -n 100 logs/app.log
```

### Log Format
```json
{
  "level": "info",
  "timestamp": "2024-01-15T10:00:00Z",
  "method": "GET",
  "path": "/api/v1/users",
  "status": 200,
  "latency": "15ms",
  "client_ip": "192.168.1.1"
}
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📝 License

MIT License - see [LICENSE](LICENSE) file

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [golang-jwt](https://github.com/golang-jwt/jwt)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [Swagger](https://swagger.io/)
- [Zap Logger](https://github.com/uber-go/zap)

## 📞 Support

- Email: support@example.com
- GitHub Issues: [Create Issue](https://golang-api-starter/issues)

---

**⭐ If you find this helpful, please give it a star!**