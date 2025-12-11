package middleware

import (
    "sync"
    "time"

    "go-api-starter/pkg/response"
    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

type RateLimiterMiddleware struct {
    limiters map[string]*rate.Limiter
    mu       sync.RWMutex
    rate     rate.Limit
    burst    int
}

func NewRateLimiterMiddleware(r rate.Limit, b int) *RateLimiterMiddleware {
    return &RateLimiterMiddleware{
        limiters: make(map[string]*rate.Limiter),
        rate:     r,
        burst:    b,
    }
}

func (rl *RateLimiterMiddleware) getLimiter(key string) *rate.Limiter {
    rl.mu.RLock()
    limiter, exists := rl.limiters[key]
    rl.mu.RUnlock()

    if !exists {
        rl.mu.Lock()
        limiter = rate.NewLimiter(rl.rate, rl.burst)
        rl.limiters[key] = limiter
        rl.mu.Unlock()
    }

    return limiter
}

// Limit middleware untuk rate limiting
func (rl *RateLimiterMiddleware) Limit() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Use IP address as key
        key := c.ClientIP()

        limiter := rl.getLimiter(key)

        if !limiter.Allow() {
            response.Error(c, 429, "Too many requests", "Rate limit exceeded")
            c.Abort()
            return
        }

        c.Next()
    }
}

// CleanupRoutine cleans up old limiters periodically
func (rl *RateLimiterMiddleware) CleanupRoutine(interval time.Duration) {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            rl.mu.Lock()
            // Reset the map to clean up old entries
            rl.limiters = make(map[string]*rate.Limiter)
            rl.mu.Unlock()
        }
    }()
}
package middleware

import (
    "sync"
    "time"

    "go-api-starter/pkg/response"
    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

type RateLimiterMiddleware struct {
    limiters map[string]*rate.Limiter
    mu       sync.RWMutex
    rate     rate.Limit
    burst    int
}

func NewRateLimiterMiddleware(r rate.Limit, b int) *RateLimiterMiddleware {
    return &RateLimiterMiddleware{
        limiters: make(map[string]*rate.Limiter),
        rate:     r,
        burst:    b,
    }
}

func (rl *RateLimiterMiddleware) getLimiter(key string) *rate.Limiter {
    rl.mu.RLock()
    limiter, exists := rl.limiters[key]
    rl.mu.RUnlock()

    if !exists {
        rl.mu.Lock()
        limiter = rate.NewLimiter(rl.rate, rl.burst)
        rl.limiters[key] = limiter
        rl.mu.Unlock()
    }

    return limiter
}

// Limit middleware untuk rate limiting
func (rl *RateLimiterMiddleware) Limit() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Use IP address as key
        key := c.ClientIP()

        limiter := rl.getLimiter(key)

        if !limiter.Allow() {
            response.Error(c, 429, "Too many requests", "Rate limit exceeded")
            c.Abort()
            return
        }

        c.Next()
    }
}

// CleanupRoutine cleans up old limiters periodically
func (rl *RateLimiterMiddleware) CleanupRoutine(interval time.Duration) {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            rl.mu.Lock()
            // Reset the map to clean up old entries
            rl.limiters = make(map[string]*rate.Limiter)
            rl.mu.Unlock()
        }
    }()
}
```

**Penjelasan:**
- Rate limiting per IP address
- Token bucket algorithm (golang.org/x/time/rate)
- Thread-safe dengan RWMutex
- Auto cleanup untuk memory efficiency

**Rate Limiting:**
```
Config: 100 requests per minute

Client IP: 192.168.1.1
├─ Request 1-100: ✅ Allowed
├─ Request 101: ❌ Blocked (429 Too Many Requests)
│
└─ After 1 minute: Bucket refilled
   └─ Request 102: ✅ Allowed
```

**Per-IP Tracking:**
```
limiters map:
├─ "192.168.1.1" → rate.Limiter{tokens: 95}
├─ "192.168.1.2" → rate.Limiter{tokens: 100}
└─ "192.168.1.3" → rate.Limiter{tokens: 50}