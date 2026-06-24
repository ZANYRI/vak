package auth

import (
	"net/http"

	"billing-service/internal/api"
	"billing-service/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler provides HTTP handlers for authentication.
type Handler struct {
	svc *Service
}

// NewHandler creates a new auth handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes adds auth routes to the router.
func (h *Handler) RegisterRoutes(r chi.Router, auth *middleware.Authenticator, limiter *middleware.RateLimiter) {
	r.Use(limiter.Limit)
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
	r.With(auth.Middleware).Post("/logout", h.Logout)
	r.With(auth.Middleware).Get("/me", h.Me)
}

// Register creates a new user.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if !api.DecodeJSON(w, r, &req) {
		return
	}
	user, err := h.svc.Register(r.Context(), req)
	if err != nil {
		if err == ErrEmailTaken {
			middleware.RespondError(w, r, http.StatusConflict, "CONFLICT", "email already registered", nil)
			return
		}
		middleware.RespondError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusCreated, map[string]interface{}{"user": user})
}

// Login authenticates a user.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if !api.DecodeJSON(w, r, &req) {
		return
	}
	_, tokens, err := h.svc.Login(r.Context(), req)
	if err != nil {
		if err == ErrInvalidCredentials {
			middleware.RespondUnauthorized(w, r, "invalid credentials")
			return
		}
		middleware.RespondError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, tokens)
}

// Refresh swaps a refresh token for a new pair.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if !api.DecodeJSON(w, r, &req) {
		return
	}
	tokens, _, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		middleware.RespondUnauthorized(w, r, "invalid refresh token")
		return
	}
	api.RespondJSON(w, r, http.StatusOK, tokens)
}

// Logout invalidates refresh tokens.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		middleware.RespondUnauthorized(w, r, "")
		return
	}
	id, err := uuid.Parse(user.ID)
	if err != nil {
		middleware.RespondUnauthorized(w, r, "")
		return
	}
	if err := h.svc.Logout(r.Context(), id); err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, map[string]interface{}{"message": "logged out"})
}

// Me returns the current authenticated user.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		middleware.RespondUnauthorized(w, r, "")
		return
	}
	id, err := uuid.Parse(user.ID)
	if err != nil {
		middleware.RespondUnauthorized(w, r, "")
		return
	}
	u, err := h.svc.GetUserByID(r.Context(), id)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, map[string]interface{}{"user": u})
}
