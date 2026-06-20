package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/example/billing-service/internal/observability"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type ctxKey int

const requestIDKey ctxKey = iota

// RequestID assigns/propagates a request id and a correlation id.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get("X-Request-ID")
		if rid == "" {
			rid = uuid.NewString()
		}
		correlation := r.Header.Get("X-Correlation-ID")
		if correlation == "" {
			correlation = rid
		}
		w.Header().Set("X-Request-ID", rid)
		w.Header().Set("X-Correlation-ID", correlation)
		ctx := context.WithValue(r.Context(), requestIDKey, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFrom returns the request id stored in context.
func RequestIDFrom(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}

// SecureHeaders sets a baseline of security headers.
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		next.ServeHTTP(w, r)
	})
}

// Metrics records request counts and latency.
func Metrics(m *observability.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			routePattern := chiRoutePattern(r)
			m.HTTPRequestsTotal.WithLabelValues(r.Method, routePattern, statusClass(ww.Status())).Inc()
			m.HTTPRequestDuration.WithLabelValues(r.Method, routePattern).Observe(time.Since(start).Seconds())
		})
	}
}

func chiRoutePattern(r *http.Request) string {
	if rctx := chi.RouteContext(r.Context()); rctx != nil {
		if p := rctx.RoutePattern(); p != "" {
			return p
		}
	}
	return r.URL.Path
}

func statusClass(code int) string {
	switch {
	case code >= 500:
		return "5xx"
	case code >= 400:
		return "4xx"
	case code >= 300:
		return "3xx"
	case code >= 200:
		return "2xx"
	default:
		return "1xx"
	}
}
