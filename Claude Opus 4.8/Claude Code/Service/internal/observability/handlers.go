package observability

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ReadyCheck reports whether a dependency is reachable.
type ReadyCheck func(ctx context.Context) error

// Health bundles liveness/readiness/metrics HTTP handlers.
type Health struct {
	metrics *Metrics
	checks  map[string]ReadyCheck
}

func NewHealth(metrics *Metrics, checks map[string]ReadyCheck) *Health {
	return &Health{metrics: metrics, checks: checks}
}

// Healthz is a liveness probe.
func (h *Health) Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// Readyz verifies all dependencies are reachable.
func (h *Health) Readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	for name, check := range h.checks {
		if err := check(ctx); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"unavailable","failed":"` + name + `"}`))
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ready"}`))
}

// MetricsHandler serves Prometheus metrics.
func (h *Health) MetricsHandler() http.Handler {
	return promhttp.HandlerFor(h.metrics.Registry(), promhttp.HandlerOpts{})
}
