package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter provides a simple per-IP fixed-window rate limiter.
// For production, prefer a distributed store such as Redis.
type RateLimiter struct {
	maxRequests int
	window      time.Duration
	mu          sync.Mutex
	clients     map[string]*clientWindow
}

type clientWindow struct {
	count  int
	resetAt time.Time
}

// NewRateLimiter creates a per-IP rate limiter.
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests: maxRequests,
		window:      window,
		clients:     make(map[string]*clientWindow),
	}
}

// Limit returns HTTP middleware enforcing the rate limit.
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if !rl.allow(ip) {
			RespondError(w, r, http.StatusTooManyRequests, "RATE_LIMITED", "too many requests", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	cw, ok := rl.clients[ip]
	if !ok || now.After(cw.resetAt) {
		rl.clients[ip] = &clientWindow{count: 1, resetAt: now.Add(rl.window)}
		return true
	}
	if cw.count >= rl.maxRequests {
		return false
	}
	cw.count++
	return true
}
