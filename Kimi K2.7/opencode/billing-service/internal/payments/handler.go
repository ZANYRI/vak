package payments

import (
	"net/http"
	"strconv"

	"billing-service/internal/api"
	"billing-service/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles payment HTTP requests.
type Handler struct {
	svc *Service
}

// NewHandler creates a payment handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers payment routes.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/simulate", h.Simulate)
	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
}

func (h *Handler) Simulate(w http.ResponseWriter, r *http.Request) {
	var req SimulateRequest
	if !api.DecodeJSON(w, r, &req) {
		return
	}
	key := api.IdempotencyKey(r)
	pay, err := h.svc.Simulate(r.Context(), req, key)
	if err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusCreated, pay)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	list, total, err := h.svc.List(r.Context(), nil, limit, offset)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, map[string]interface{}{"data": list, "total": total})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid id", nil)
		return
	}
	p, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, p)
}
