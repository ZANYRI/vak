package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/example/billing-service/internal/httpx"
)

// Recover converts panics into 500 responses and logs the stack trace.
func Recover(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic recovered",
						"error", rec,
						"path", r.URL.Path,
						"stack", string(debug.Stack()),
					)
					httpx.Error(w, r, httpx.ErrInternal(""))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
