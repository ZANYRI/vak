package invoices

import (
	"context"
	"net/http"
	"strconv"

	"billing-service/internal/api"
	"billing-service/internal/middleware"
	"billing-service/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles invoice HTTP requests.
type Handler struct {
	svc *Service
}

// NewHandler creates an invoice handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers invoice routes.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/generate", h.Generate)
	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
	r.Post("/{id}/finalize", h.Finalize)
	r.Post("/{id}/void", h.Void)
	r.Post("/{id}/pay", h.Pay)
}

func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {
	var req GenerateRequest
	if !api.DecodeJSON(w, r, &req) {
		return
	}
	key := api.IdempotencyKey(r)
	inv, err := h.svc.Generate(r.Context(), req, key)
	if err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusCreated, inv)
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
	inv, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, inv)
}

func (h *Handler) Finalize(w http.ResponseWriter, r *http.Request) {
	h.transition(w, r, h.svc.Finalize)
}

func (h *Handler) Void(w http.ResponseWriter, r *http.Request) {
	h.transition(w, r, h.svc.Void)
}

func (h *Handler) Pay(w http.ResponseWriter, r *http.Request) {
	h.transition(w, r, h.svc.Pay)
}

func (h *Handler) transition(w http.ResponseWriter, r *http.Request, fn func(context.Context, uuid.UUID) (*models.Invoice, error)) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid id", nil)
		return
	}
	inv, err := fn(r.Context(), id)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, inv)
}
