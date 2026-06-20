package subscriptions

import (
	"net/http"

	"github.com/example/billing-service/internal/auth"
	"github.com/example/billing-service/internal/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	svc         *Service
	mw          *auth.Middleware
	applyCoupon http.HandlerFunc // mounted from the coupons module
}

func NewHandler(svc *Service, mw *auth.Middleware, applyCoupon http.HandlerFunc) *Handler {
	return &Handler{svc: svc, mw: mw, applyCoupon: applyCoupon}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(h.mw.Authenticate)

	create := h.mw.RequireRole(auth.RoleAdmin, auth.RoleBillingManager, auth.RoleCustomer)
	manage := h.mw.RequireRole(auth.RoleAdmin, auth.RoleBillingManager)

	r.With(create).Post("/", h.create)
	r.Get("/", h.list)
	r.Get("/{id}", h.get)
	r.With(manage).Patch("/{id}", h.update)
	r.With(manage).Post("/{id}/cancel", h.cancel)
	r.With(manage).Post("/{id}/pause", h.pause)
	r.With(manage).Post("/{id}/resume", h.resume)
	r.With(manage).Post("/{id}/change-plan", h.changePlan)
	if h.applyCoupon != nil {
		r.With(manage).Post("/{id}/apply-coupon", h.applyCoupon)
	}
	return r
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	id, _ := auth.FromContext(r.Context())
	sub, err := h.svc.Create(r.Context(), id.UserID, req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, sub)
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
	sub, err := h.svc.Get(r.Context(), id)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, sub)
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
	sub, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, sub)
}

type cancelRequest struct {
	AtPeriodEnd bool `json:"at_period_end"`
}

func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	var req cancelRequest
	_ = httpx.DecodeAndValidate(r, &req) // body optional
	actor, _ := auth.FromContext(r.Context())
	sub, err := h.svc.Cancel(r.Context(), actor.UserID, id, req.AtPeriodEnd)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, sub)
}

func (h *Handler) pause(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	sub, err := h.svc.Pause(r.Context(), id)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, sub)
}

func (h *Handler) resume(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	sub, err := h.svc.Resume(r.Context(), id)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, sub)
}

func (h *Handler) changePlan(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathUUID(r, "id")
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	var req ChangePlanRequest
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	res, err := h.svc.ChangePlan(r.Context(), id, req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, res)
}
