package auth

import (
	"net/http"
	"strings"

	"github.com/example/billing-service/internal/httpx"
)

// Middleware provides authentication and authorization HTTP middleware.
type Middleware struct {
	tokens *TokenManager
}

func NewMiddleware(tokens *TokenManager) *Middleware { return &Middleware{tokens: tokens} }

// Authenticate requires a valid Bearer access token and attaches the Identity.
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			httpx.Error(w, r, httpx.ErrUnauthorized("missing bearer token"))
			return
		}
		raw := strings.TrimPrefix(header, "Bearer ")
		claims, err := m.tokens.ParseAccess(raw)
		if err != nil {
			httpx.Error(w, r, httpx.ErrUnauthorized("invalid or expired token"))
			return
		}
		ctx := WithIdentity(r.Context(), Identity{UserID: claims.UserID, Role: claims.Role})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole authorizes the request if the caller has one of the allowed roles.
func (m *Middleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, ok := FromContext(r.Context())
			if !ok {
				httpx.Error(w, r, httpx.ErrUnauthorized("authentication required"))
				return
			}
			if !allowed[id.Role] {
				httpx.Error(w, r, httpx.ErrInsufficientPermission("role not permitted for this action"))
				return
			}
			next.ServeHTTP(w, r.WithContext(r.Context()))
		})
	}
}
