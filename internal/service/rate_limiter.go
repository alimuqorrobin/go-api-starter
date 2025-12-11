package service

import (
    "context"
    "sync"
    "time"
)

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
    rate       int
    capacity   int
    tokens     int
    lastRefill time.Time
    mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, capacity int) *RateLimiter {
    return &RateLimiter{
        rate:       rate,
        capacity:   capacity,
        tokens:     capacity,
        lastRefill: time.Now(),
    }
}

// Allow checks if a request is allowed
func (rl *RateLimiter) Allow() bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    // Refill tokens based on elapsed time
    now := time.Now()
    elapsed := now.Sub(rl.lastRefill)
    tokensToAdd := int(elapsed.Seconds()) * rl.rate

    if tokensToAdd > 0 {
        rl.tokens = min(rl.capacity, rl.tokens+tokensToAdd)
        rl.lastRefill = now
    }

    // Check if we have tokens available
    if rl.tokens > 0 {
        rl.tokens--
        return true
    }

    return false
}

// Wait waits until a token is available
func (rl *RateLimiter) Wait(ctx context.Context) error {
    for {
        if rl.Allow() {
            return nil
        }

        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(time.Millisecond * 100):
            // Continue checking
        }
    }
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

**Penjelasan:**
- Token Bucket Algorithm untuk rate limiting
- `rate` = Tokens per second
- `capacity` = Max tokens in bucket
- Thread-safe dengan mutex

**Token Bucket Algorithm:**
```
// Bucket (capacity: 100)
// ├─ Start: 100 tokens
// ├─ Request 1: 99 tokens (allowed)
// ├─ Request 2: 98 tokens (allowed)
// ├─ ... (many requests)
// ├─ Request 100: 0 tokens (allowed)
// ├─ Request 101: 0 tokens (BLOCKED! 429)
// │
// └─ After 1 second: Refill +100 tokens
//    └─ Request 102: 99 tokens (allowed again)


limiter := NewRateLimiter(100, 100) // 100 req/sec

if limiter.Allow() {
    // Process request
} else {
    // Return 429 Too Many Requests
}
```

---

### ✅ Part 3 Selesai!

// Sudah selesai semua file di Part 3:
// - ✅ internal/domain/user.go (Entity & DTOs)
// - ✅ internal/domain/auth.go (Auth DTOs)
// - ✅ internal/domain/errors.go (Custom errors)
// - ✅ internal/repository/interface.go (Repository contract)
// - ✅ internal/repository/user_repository.go (Data access layer)
// - ✅ internal/service/user_service.go (User business logic)
// - ✅ internal/service/auth_service.go (Auth business logic)
// - ✅ internal/service/worker_pool.go (Concurrency)
// - ✅ internal/service/rate_limiter.go (Rate limiting)

// **Struktur sekarang:**
// ```
// golang-api-starter/
// ├── cmd/server/main.go
// ├── config/config.go
// ├── internal/
// │   ├── database/
// │   │   ├── connection.go
// │   │   ├── postgres.go
// │   │   ├── mysql.go
// │   │   └── migration.go
// │   ├── domain/
// │   │   ├── user.go
// │   │   ├── auth.go
// │   │   └── errors.go
// │   ├── repository/
// │   │   ├── interface.go
// │   │   └── user_repository.go
// │   └── service/
// │       ├── user_service.go
// │       ├── auth_service.go
// │       ├── worker_pool.go
// │       └── rate_limiter.go
// ├── go.mod
// └── .env.example