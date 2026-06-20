package invoices

import (
	"net/http"

	"github.com/example/billing-service/internal/auth"
	"github.com/example/billing-service/internal/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
	mw  *auth.Middleware
	pay http.HandlerFunc // mounted from the payments module
}

func NewHandler(svc *Service, mw *auth.Middleware, pay http.HandlerFunc) *Handler {
	return &Handler{svc: svc, mw: mw, pay: pay}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(h.mw.Authenticate)
	manage := h.mw.RequireRole(auth.RoleAdmin, auth.RoleBillingManager)

	r.With(manage).Post("/generate", h.generate)
	r.Get("/", h.list)
	r.Get("/{id}", h.get)
	r.With(manage).Post("/{id}/finalize", h.finalize)
	r.With(manage).Post("/{id}/void", h.void)
	if h.pay != nil {
		r.With(manage).Post("/{id}/pay", h.pay)
	}
	return r
}

func (h *Handler) generate(w http.ResponseWriter, r *http.Request) {
	var req GenerateRequest
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	actor, _ := auth.FromContext(r.Context())
	inv, err := h.svc.Generate(r.Context(), actor.UserID, req.SubscriptionID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, inv)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	limit, offset := httpx.Pagination(r)
	var customerID *uuid.UUID
	if v := r.URL.Query().Get("customer_id"); v != "" {
		if cid, err := uuid.Parse(v); err == nil {
			customerID = &cid
		}
	}
	items, err := h.svc.List(r.Context(), customerID, limit, offset)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"data": items})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	inv, err := h.svc.Get(r.Context(), id)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, inv)
}

func (h *Handler) finalize(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	inv, err := h.svc.Finalize(r.Context(), id)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, inv)
}

func (h *Handler) void(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	inv, err := h.svc.Void(r.Context(), id)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, inv)
}
