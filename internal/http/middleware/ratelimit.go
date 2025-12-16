package middleware

import (
	"net/http"
	"sync"
	"time"

	"go-api-starter/internal/http/response"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Cleanup goroutine
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()
		requests := rl.requests[ip]

		// Remove old requests outside the window
		var valid []time.Time
		for _, req := range requests {
			if now.Sub(req) < rl.window {
				valid = append(valid, req)
			}
		}

		if len(valid) >= rl.limit {
			response.TooManyRequests(w)
			return
		}

		rl.requests[ip] = append(valid, now)
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()

		for ip, requests := range rl.requests {
			var valid []time.Time
			for _, req := range requests {
				if now.Sub(req) < rl.window {
					valid = append(valid, req)
				}
			}

			if len(valid) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = valid
			}
		}

		rl.mu.Unlock()
	}
}
