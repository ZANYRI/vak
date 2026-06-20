package auth

import (
	"net/http"

	"github.com/example/billing-service/internal/httpx"
	"github.com/go-chi/chi/v5"
)

// Handler exposes auth HTTP endpoints.
type Handler struct {
	svc *Service
	mw  *Middleware
}

func NewHandler(svc *Service, mw *Middleware) *Handler {
	return &Handler{svc: svc, mw: mw}
}

// Routes mounts auth routes. rateLimit is applied to credential endpoints.
func (h *Handler) Routes(rateLimit func(http.Handler) http.Handler) http.Handler {
	r := chi.NewRouter()
	r.With(rateLimit).Post("/register", h.register)
	r.With(rateLimit).Post("/login", h.login)
	r.With(rateLimit).Post("/refresh", h.refresh)
	r.Post("/logout", h.logout)
	r.With(h.mw.Authenticate).Get("/me", h.me)
	return r
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	user, err := h.svc.Register(r.Context(), req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, user)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	pair, _, err := h.svc.Login(r.Context(), req, clientIP(r), r.UserAgent())
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, pair)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	pair, err := h.svc.Refresh(r.Context(), req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, pair)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	if err := h.svc.Logout(r.Context(), req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]string{"status": "logged_out"})
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	id, _ := FromContext(r.Context())
	user, err := h.svc.Me(r.Context(), id.UserID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, user)
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	return r.RemoteAddr
}
