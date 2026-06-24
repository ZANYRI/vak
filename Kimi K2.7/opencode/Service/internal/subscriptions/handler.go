package subscriptions

import (
	"net/http"
	"strconv"

	"billing-service/internal/api"
	"billing-service/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles subscription HTTP requests.
type Handler struct {
	svc *Service
}

// NewHandler creates a subscription handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers subscription routes.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
	r.Patch("/{id}", h.Update)
	r.Post("/{id}/cancel", h.Cancel)
	r.Post("/{id}/pause", h.Pause)
	r.Post("/{id}/resume", h.Resume)
	r.Post("/{id}/change-plan", h.ChangePlan)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if !api.DecodeJSON(w, r, &req) {
		return
	}
	sub, err := h.svc.Create(r.Context(), req)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusCreated, sub)
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
	id, err := h.parseID(r)
	if err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid id", nil)
		return
	}
	sub, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, sub)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid id", nil)
		return
	}
	var req UpdateRequest
	if !api.DecodeJSON(w, r, &req) {
		return
	}
	sub, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, sub)
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid id", nil)
		return
	}
	sub, err := h.svc.Cancel(r.Context(), id)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, map[string]interface{}{"subscription": sub})
}

func (h *Handler) Pause(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid id", nil)
		return
	}
	sub, err := h.svc.Pause(r.Context(), id)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, map[string]interface{}{"subscription": sub})
}

func (h *Handler) Resume(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid id", nil)
		return
	}
	sub, err := h.svc.Resume(r.Context(), id)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, map[string]interface{}{"subscription": sub})
}

func (h *Handler) ChangePlan(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "invalid id", nil)
		return
	}
	var req ChangePlanRequest
	if !api.DecodeJSON(w, r, &req) {
		return
	}
	sub, proration, err := h.svc.ChangePlan(r.Context(), id, req)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, map[string]interface{}{"subscription": sub, "proration": proration})
}

func (h *Handler) parseID(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(chi.URLParam(r, "id"))
}
