package plans

import (
	"net/http"

	"github.com/example/billing-service/internal/auth"
	"github.com/example/billing-service/internal/httpx"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
	mw  *auth.Middleware
}

func NewHandler(svc *Service, mw *auth.Middleware) *Handler { return &Handler{svc: svc, mw: mw} }

// Routes mounts plan routes. All require authentication; writes require billing roles.
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(h.mw.Authenticate)

	manage := h.mw.RequireRole(auth.RoleAdmin, auth.RoleBillingManager)

	r.Get("/", h.list)
	r.Get("/{id}", h.get)
	r.With(manage).Post("/", h.create)
	r.With(manage).Patch("/{id}", h.update)
	r.With(manage).Delete("/{id}", h.delete)
	return r
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	id, _ := auth.FromContext(r.Context())
	p, err := h.svc.Create(r.Context(), id.UserID, req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, p)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	p, err := h.svc.Get(r.Context(), id)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, p)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	limit, offset := httpx.Pagination(r)
	items, err := h.svc.List(r.Context(), limit, offset)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"data": items})
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	var req UpdateRequest
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	p, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, p)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusNoContent, nil)
}
