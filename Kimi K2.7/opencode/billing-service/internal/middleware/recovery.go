package middleware

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// Recovery catches panics and returns a structured 500 response.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger := LoggerFromContext(r.Context())
				logger.Error("panic recovered", zap.Any("error", rec), zap.Stack("stack"))
				RespondError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", fmt.Sprintf("internal server error: %v", rec), nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
