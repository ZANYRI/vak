package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus collectors used across the service.
type Metrics struct {
	HTTPRequestsTotal   *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec
	JobsProcessedTotal  *prometheus.CounterVec
	JobsFailedTotal     *prometheus.CounterVec
	InvoicesGenerated   prometheus.Counter
	PaymentsSucceeded   prometheus.Counter
	PaymentsFailed      prometheus.Counter
	QueueDepth          *prometheus.GaugeVec
	registry            *prometheus.Registry
}

// NewMetrics creates a dedicated registry and registers every collector.
func NewMetrics() *Metrics {
	reg := prometheus.NewRegistry()
	reg.MustRegister(prometheus.NewGoCollector())
	factory := promauto.With(reg)

	return &Metrics{
		registry: reg,
		HTTPRequestsTotal: factory.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		}, []string{"method", "path", "status"}),
		HTTPRequestDuration: factory.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "path"}),
		JobsProcessedTotal: factory.NewCounterVec(prometheus.CounterOpts{
			Name: "jobs_processed_total",
			Help: "Total number of jobs processed.",
		}, []string{"type"}),
		JobsFailedTotal: factory.NewCounterVec(prometheus.CounterOpts{
			Name: "jobs_failed_total",
			Help: "Total number of jobs that failed.",
		}, []string{"type"}),
		InvoicesGenerated: factory.NewCounter(prometheus.CounterOpts{
			Name: "invoices_generated_total",
			Help: "Total number of invoices generated.",
		}),
		PaymentsSucceeded: factory.NewCounter(prometheus.CounterOpts{
			Name: "payments_succeeded_total",
			Help: "Total number of successful payments.",
		}),
		PaymentsFailed: factory.NewCounter(prometheus.CounterOpts{
			Name: "payments_failed_total",
			Help: "Total number of failed payments.",
		}),
		QueueDepth: factory.NewGaugeVec(prometheus.GaugeOpts{
			Name: "queue_depth",
			Help: "Approximate number of pending jobs per stream.",
		}, []string{"stream"}),
	}
}

// Registry exposes the underlying Prometheus registry for the /metrics handler.
func (m *Metrics) Registry() *prometheus.Registry { return m.registry }
