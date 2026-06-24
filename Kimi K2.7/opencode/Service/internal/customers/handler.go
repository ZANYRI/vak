package customers

import (
	"net/http"
	"strconv"

	"billing-service/internal/api"
	"billing-service/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles customer HTTP requests.
type Handler struct {
	svc *Service
}

// NewHandler creates a handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers customer routes.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
	r.Patch("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if !api.DecodeJSON(w, r, &req) {
		return
	}
	c, err := h.svc.Create(r.Context(), req)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusCreated, c)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	list, total, err := h.svc.List(r.Context(), limit, offset)
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
	c, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, c)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid id", nil)
		return
	}
	var req UpdateRequest
	if !api.DecodeJSON(w, r, &req) {
		return
	}
	c, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, c)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid id", nil)
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusNoContent, nil)
}
