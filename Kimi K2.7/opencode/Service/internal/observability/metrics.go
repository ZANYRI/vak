package observability

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests by method, path and status.",
		},
		[]string{"method", "path", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
	jobsProcessedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "jobs_processed_total", Help: "Total processed jobs."},
		[]string{"queue", "status"},
	)
	jobsFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "jobs_failed_total", Help: "Total failed jobs."},
		[]string{"queue"},
	)
	invoicesGeneratedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "invoices_generated_total", Help: "Total invoices generated."},
		[]string{"status"},
	)
	paymentsSucceededTotal = prometheus.NewCounter(prometheus.CounterOpts{Name: "payments_succeeded_total", Help: "Total successful payments."})
	paymentsFailedTotal = prometheus.NewCounter(prometheus.CounterOpts{Name: "payments_failed_total", Help: "Total failed payments."})
	queueDepth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "queue_depth", Help: "Approximate queued jobs."},
		[]string{"queue"},
	)
)

func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		jobsProcessedTotal,
		jobsFailedTotal,
		invoicesGeneratedTotal,
		paymentsSucceededTotal,
		paymentsFailedTotal,
		queueDepth,
	)
}

// MetricsHandler exposes the /metrics endpoint.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// PrometheusMiddleware records HTTP request metrics.
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		duration := time.Since(start).Seconds()
		path := r.URL.Path
		status := strconv.Itoa(ww.Status())
		httpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
	})
}

// IncJobProcessed increments the jobs processed counter.
func IncJobProcessed(queue, status string) {
	jobsProcessedTotal.WithLabelValues(queue, status).Inc()
}

// IncJobFailed increments the jobs failed counter.
func IncJobFailed(queue string) {
	jobsFailedTotal.WithLabelValues(queue).Inc()
}

// IncInvoiceGenerated increments the invoices generated counter.
func IncInvoiceGenerated(status string) {
	invoicesGeneratedTotal.WithLabelValues(status).Inc()
}

// IncPaymentSucceeded increments the payments succeeded counter.
func IncPaymentSucceeded() {
	paymentsSucceededTotal.Inc()
}

// IncPaymentFailed increments the payments failed counter.
func IncPaymentFailed() {
	paymentsFailedTotal.Inc()
}

// SetQueueDepth sets the queue depth gauge.
func SetQueueDepth(queue string, depth float64) {
	queueDepth.WithLabelValues(queue).Set(depth)
}

// RequestID retrieves the chi request id from context.
func RequestID(ctx context.Context) string {
	if reqID := middleware.GetReqID(ctx); reqID != "" {
		return reqID
	}
	return "-"
}
