package payments

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
	manage := h.mw.RequireRole(auth.RoleAdmin, auth.RoleBillingManager)

	r.With(manage).Post("/simulate", h.simulate)
	r.Get("/", h.list)
	r.Get("/{id}", h.get)
	return r
}

func (h *Handler) simulate(w http.ResponseWriter, r *http.Request) {
	var req SimulateRequest
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	p, err := h.svc.Simulate(r.Context(), req)
	if err != nil {
		// A declined payment returns the payment record alongside the error.
		if p != nil {
			httpx.JSON(w, http.StatusPaymentRequired, map[string]any{
				"payment": p,
				"error": map[string]string{"code": httpx.CodePaymentFailed, "message": err.Error()},
			})
			return
		}
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, p)
}

// PayInvoiceHandler implements POST /invoices/{id}/pay (mounted by the invoices module).
func (h *Handler) PayInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	p, err := h.svc.PayInvoice(r.Context(), id)
	if err != nil {
		if p != nil {
			httpx.JSON(w, http.StatusPaymentRequired, map[string]any{"payment": p})
			return
		}
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, p)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	limit, offset := httpx.Pagination(r)
	var invoiceID *uuid.UUID
	if v := r.URL.Query().Get("invoice_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			invoiceID = &id
		}
	}
	items, err := h.svc.List(r.Context(), invoiceID, limit, offset)
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
	p, err := h.svc.Get(r.Context(), id)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, p)
}
