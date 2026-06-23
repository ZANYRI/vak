package middleware

import (
	"context"
	"net/http"
	"strings"

	"billing-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type authKey struct{}

// AuthUser is the authenticated user stored in context.
type AuthUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Authenticator validates JWT access tokens.
type Authenticator struct {
	Secret []byte
}

// NewAuthenticator creates an Authenticator.
func NewAuthenticator(secret string) *Authenticator {
	return &Authenticator{Secret: []byte(secret)}
}

// Middleware validates the Authorization header and adds user info to context.
func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			RespondUnauthorized(w, r, "missing authorization header")
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			RespondUnauthorized(w, r, "invalid authorization header format")
			return
		}
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return a.Secret, nil
		})
		if err != nil || !token.Valid {
			RespondUnauthorized(w, r, "invalid or expired token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			RespondUnauthorized(w, r, "invalid token claims")
			return
		}
		user := AuthUser{
			ID:    stringValue(claims["sub"]),
			Email: stringValue(claims["email"]),
			Role:  stringValue(claims["role"]),
		}
		ctx := context.WithValue(r.Context(), authKey{}, user)
		LoggerFromContext(ctx).Debug("authenticated request", zap.String("user_id", user.ID), zap.String("role", user.Role))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserFromContext returns the authenticated user from the request context.
func UserFromContext(ctx context.Context) (AuthUser, bool) {
	u, ok := ctx.Value(authKey{}).(AuthUser)
	return u, ok
}

// RequireRole returns a middleware that checks the user has any of the given roles.
func RequireRole(roles ...models.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok {
				RespondUnauthorized(w, r, "unauthorized")
				return
			}
			for _, role := range roles {
				if user.Role == string(role) {
					next.ServeHTTP(w, r)
					return
				}
			}
			RespondError(w, r, http.StatusForbidden, "FORBIDDEN", "insufficient role", nil)
		})
	}
}

// RequirePermission returns a middleware that checks the user has the required permission.
func RequirePermission(perm models.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok {
				RespondUnauthorized(w, r, "unauthorized")
				return
			}
			role := models.UserRole(user.Role)
			if !models.HasPermission(role, perm) {
				RespondError(w, r, http.StatusForbidden, "FORBIDDEN", "insufficient permissions", nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func stringValue(v interface{}) string {
	s, _ := v.(string)
	return s
}
