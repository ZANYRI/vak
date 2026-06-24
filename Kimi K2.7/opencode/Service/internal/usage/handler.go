package usage

import (
	"net/http"
	"strconv"
	"time"

	"billing-service/internal/api"
	"billing-service/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles usage HTTP requests.
type Handler struct {
	svc *Service
}

// NewHandler creates a usage handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers usage routes.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.Report)
	r.Get("/", h.List)
	r.Get("/summary", h.Summary)
}

func (h *Handler) Report(w http.ResponseWriter, r *http.Request) {
	var req ReportRequest
	if !api.DecodeJSON(w, r, &req) {
		return
	}
	event, err := h.svc.Report(r.Context(), req)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusCreated, event)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	subID, err := uuid.Parse(r.URL.Query().Get("subscription_id"))
	if err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "subscription_id required", nil)
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	list, total, err := h.svc.List(r.Context(), subID, limit, offset)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, map[string]interface{}{"data": list, "total": total})
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	subID, err := uuid.Parse(r.URL.Query().Get("subscription_id"))
	if err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "subscription_id required", nil)
		return
	}
	metric := r.URL.Query().Get("metric")
	if metric == "" {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "metric required", nil)
		return
	}
	from := parseTimeOrDefault(r.URL.Query().Get("from"), time.Now().AddDate(0, -1, 0))
	to := parseTimeOrDefault(r.URL.Query().Get("to"), time.Now())
	summary, err := h.svc.Summary(r.Context(), subID, metric, from, to)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, summary)
}

func parseTimeOrDefault(s string, def time.Time) time.Time {
	if s == "" {
		return def
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return def
	}
	return t
}
