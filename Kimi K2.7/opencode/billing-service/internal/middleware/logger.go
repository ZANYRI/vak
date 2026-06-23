package middleware

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type loggerKey struct{}

// WithLogger adds a logger to the request context.
func WithLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loggerKey{}, logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// LoggerFromContext retrieves the logger from context.
func LoggerFromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(loggerKey{}).(*zap.Logger); ok && l != nil {
		return l
	}
	return zap.NewNop()
}
