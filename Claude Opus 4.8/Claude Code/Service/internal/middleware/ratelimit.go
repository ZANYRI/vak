package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/example/billing-service/internal/httpx"
)

// RateLimiter is a simple in-memory fixed-window limiter keyed by client IP.
// Suitable for protecting auth endpoints in a single-instance deployment.
type RateLimiter struct {
	mu      sync.Mutex
	hits    map[string]*window
	limit   int
	window  time.Duration
}

type window struct {
	count int
	reset time.Time
}

func NewRateLimiter(limit int, win time.Duration) *RateLimiter {
	rl := &RateLimiter{hits: make(map[string]*window), limit: limit, window: win}
	return rl
}

func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	w, ok := rl.hits[key]
	if !ok || now.After(w.reset) {
		rl.hits[key] = &window{count: 1, reset: now.Add(rl.window)}
		return true
	}
	if w.count >= rl.limit {
		return false
	}
	w.count++
	return true
}

// Middleware enforces the rate limit.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := clientIP(r)
		if !rl.allow(key) {
			w.Header().Set("Retry-After", "60")
			httpx.Error(w, r, &httpx.APIError{
				Status: http.StatusTooManyRequests, Code: "RATE_LIMITED",
				Message: "too many requests, slow down",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	host := r.RemoteAddr
	if i := indexColon(host); i >= 0 {
		return host[:i]
	}
	return host
}

// indexColon returns the index of the last colon (to strip the port).
func indexColon(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ':' {
			return i
		}
	}
	return -1
}
