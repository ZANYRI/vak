package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/example/billing-service/internal/app"
	"github.com/example/billing-service/internal/auth"
	"github.com/example/billing-service/internal/coupons"
	"github.com/example/billing-service/internal/customers"
	"github.com/example/billing-service/internal/invoices"
	"github.com/example/billing-service/internal/middleware"
	"github.com/example/billing-service/internal/observability"
	"github.com/example/billing-service/internal/payments"
	"github.com/example/billing-service/internal/plans"
	"github.com/example/billing-service/internal/subscriptions"
	"github.com/example/billing-service/internal/taxes"
	"github.com/example/billing-service/internal/usage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// NewRouter builds the HTTP handler for the API process.
func NewRouter(c *app.Container) http.Handler {
	r := chi.NewRouter()

	// Global middleware.
	r.Use(middleware.RequestID)
	r.Use(middleware.Recover(c.Log))
	r.Use(middleware.SecureHeaders)
	r.Use(corsMiddleware(c.Cfg.CORSAllowedOrigins))
	r.Use(middleware.Metrics(c.Metrics))

	idem := middleware.NewIdempotency(c.Pool)
	r.Use(idem.Middleware)

	// Observability endpoints.
	health := observability.NewHealth(c.Metrics, map[string]observability.ReadyCheck{
		"postgres": func(ctx context.Context) error { return c.Pool.Ping(ctx) },
		"redis":    func(ctx context.Context) error { return c.Queue.Ping(ctx) },
	})
	r.Get("/healthz", health.Healthz)
	r.Get("/readyz", health.Readyz)
	r.Handle("/metrics", health.MetricsHandler())

	// API documentation.
	mountDocs(r)

	// Auth rate limiter.
	rl := middleware.NewRateLimiter(c.Cfg.AuthRateLimit, c.Cfg.AuthRateLimitWindow)

	// Handlers.
	authHandler := auth.NewHandler(c.Auth, c.AuthMW)
	customersHandler := customers.NewHandler(c.Customers, c.AuthMW)
	plansHandler := plans.NewHandler(c.Plans, c.AuthMW)
	couponsHandler := coupons.NewHandler(c.Coupons, c.AuthMW)
	taxesHandler := taxes.NewHandler(c.Taxes, c.AuthMW)
	usageHandler := usage.NewHandler(c.Usage, c.AuthMW)
	paymentsHandler := payments.NewHandler(c.Payments, c.AuthMW)

	subsHandler := subscriptions.NewHandler(c.Subs, c.AuthMW, couponsHandler.ApplyCouponHandler)
	invoicesHandler := invoices.NewHandler(c.Invoices, c.AuthMW, paymentsHandler.PayInvoiceHandler)

	r.Route("/api/v1", func(api chi.Router) {
		api.Mount("/auth", authHandler.Routes(rl.Middleware))
		api.Mount("/customers", customersHandler.Routes())
		api.Mount("/plans", plansHandler.Routes())
		api.Mount("/subscriptions", subsHandler.Routes())
		api.Mount("/usage", usageHandler.Routes())
		api.Mount("/coupons", couponsHandler.Routes())
		api.Mount("/tax-rules", taxesHandler.Routes())
		api.Mount("/invoices", invoicesHandler.Routes())
		api.Mount("/payments", paymentsHandler.Routes())
	})

	return r
}

func corsMiddleware(origins string) func(http.Handler) http.Handler {
	allowed := strings.Split(origins, ",")
	for i := range allowed {
		allowed[i] = strings.TrimSpace(allowed[i])
	}
	return cors.Handler(cors.Options{
		AllowedOrigins:   allowed,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Idempotency-Key", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID", "X-Correlation-ID"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}
