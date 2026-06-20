package usage

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
}

func NewHandler(svc *Service, mw *auth.Middleware) *Handler { return &Handler{svc: svc, mw: mw} }

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(h.mw.Authenticate)
	report := h.mw.RequireRole(auth.RoleAdmin, auth.RoleBillingManager, auth.RoleCustomer)

	r.With(report).Post("/", h.record)
	r.Get("/", h.list)
	r.Get("/summary", h.summary)
	return r
}

func (h *Handler) record(w http.ResponseWriter, r *http.Request) {
	var req RecordRequest
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	event, duplicate, err := h.svc.Record(r.Context(), req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	status := http.StatusCreated
	if duplicate {
		status = http.StatusOK // idempotent replay
	}
	httpx.JSON(w, status, event)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	limit, offset := httpx.Pagination(r)
	var subID *uuid.UUID
	if v := r.URL.Query().Get("subscription_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			subID = &id
		}
	}
	items, err := h.svc.List(r.Context(), subID, limit, offset)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"data": items})
}

func (h *Handler) summary(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query().Get("subscription_id")
	id, err := uuid.Parse(v)
	if err != nil {
		httpx.Error(w, r, httpx.ErrValidation("subscription_id query parameter is required"))
		return
	}
	items, err := h.svc.Summary(r.Context(), id)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"data": items})
}
