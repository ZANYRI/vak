package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"billing-service/internal/api"
	"billing-service/internal/audit"
	"billing-service/internal/auth"
	"billing-service/internal/config"
	"billing-service/internal/coupons"
	"billing-service/internal/customers"
	"billing-service/internal/database"
	"billing-service/internal/invoices"
	"billing-service/internal/middleware"
	"billing-service/internal/models"
	"billing-service/internal/observability"
	"billing-service/internal/payments"
	"billing-service/internal/plans"
	"billing-service/internal/queue"
	"billing-service/internal/subscriptions"
	"billing-service/internal/taxes"
	"billing-service/internal/usage"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		observability.Must(err)
	}
	logger := observability.NewLogger(cfg.LogLevel)
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("database connection failed", zap.Error(err))
	}
	defer pool.Close()

	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		logger.Fatal("migrations failed", zap.Error(err))
	}
	logger.Info("migrations applied")

	// Services
	authSvc := auth.NewService(pool, cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	customerSvc := customers.NewService(pool)
	planSvc := plans.NewService(pool)
	subSvc := subscriptions.NewService(pool, planSvc)
	usageSvc := usage.NewService(pool)
	couponSvc := coupons.NewService(pool)
	couponHandler := coupons.NewHandler(couponSvc)
	taxSvc := taxes.NewService(pool)
	invoiceSvc := invoices.NewService(pool, customerSvc, subSvc, planSvc, usageSvc, couponSvc, taxSvc)
	paymentSvc := payments.NewService(pool, invoiceSvc, subSvc)
	auditSvc := audit.NewService(pool)

	// Queue client for publishing jobs
	jobStore := queue.NewJobStore(pool)
	queueClient, err := queue.NewClient(cfg.QueueURL, cfg.QueueStream, jobStore)
	if err != nil {
		logger.Fatal("queue connection failed", zap.Error(err))
	}
	defer queueClient.Close()
	_ = auditSvc

	authenticator := middleware.NewAuthenticator(cfg.JWT.AccessSecret)
	rateLimiter := middleware.NewRateLimiter(10, time.Minute)

	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(observability.PrometheusMiddleware)
	r.Use(middleware.WithLogger(logger))
	r.Use(middleware.Recovery)
	r.Use(middleware.RequestLogger)
	r.Use(middleware.CommonMiddleware)
	r.Use(middleware.CORS([]string{"*"}))

	// Health and observability
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		api.RespondJSON(w, r, http.StatusOK, map[string]string{"status": "healthy"})
	})
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			middleware.RespondError(w, r, http.StatusServiceUnavailable, "INTERNAL_ERROR", "database unreachable", nil)
			return
		}
		api.RespondJSON(w, r, http.StatusOK, map[string]string{"status": "ready"})
	})
	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r.Get("/api/docs", docsHandler)

	// API v1 routes
	apiRouter := chi.NewRouter()
	apiRouter.Route("/auth", func(r chi.Router) {
		auth.NewHandler(authSvc).RegisterRoutes(r, authenticator, rateLimiter)
	})
	apiRouter.Route("/customers", func(r chi.Router) {
		r.Use(authenticator.Middleware)
		r.Use(middleware.RequirePermission(models.PermViewCustomers))
		customers.NewHandler(customerSvc).RegisterRoutes(r)
	})
	apiRouter.Route("/plans", func(r chi.Router) {
		r.Use(authenticator.Middleware)
		r.Use(middleware.RequirePermission(models.PermManagePlans))
		plans.NewHandler(planSvc).RegisterRoutes(r)
	})
	apiRouter.Route("/subscriptions", func(r chi.Router) {
		r.Use(authenticator.Middleware)
		r.Use(middleware.RequirePermission(models.PermManageSubscriptions))
		subscriptions.NewHandler(subSvc).RegisterRoutes(r)
		couponHandler.RegisterApplyCoupon(r)
	})
	apiRouter.Route("/usage", func(r chi.Router) {
		r.Use(authenticator.Middleware)
		usage.NewHandler(usageSvc).RegisterRoutes(r)
	})
	apiRouter.Route("/coupons", func(r chi.Router) {
		r.Use(authenticator.Middleware)
		r.Use(middleware.RequirePermission(models.PermManageCoupons))
		couponHandler.RegisterRoutes(r)
	})
	apiRouter.Route("/tax-rules", func(r chi.Router) {
		r.Use(authenticator.Middleware)
		r.Use(middleware.RequirePermission(models.PermManageTaxRules))
		taxes.NewHandler(taxSvc).RegisterRoutes(r)
	})
	apiRouter.Route("/invoices", func(r chi.Router) {
		r.Use(authenticator.Middleware)
		r.Use(middleware.RequirePermission(models.PermManageInvoices))
		invoices.NewHandler(invoiceSvc).RegisterRoutes(r)
	})
	apiRouter.Route("/payments", func(r chi.Router) {
		r.Use(authenticator.Middleware)
		payments.NewHandler(paymentSvc).RegisterRoutes(r)
	})

	r.Mount("/api/v1", apiRouter)

	// Metrics server
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{Addr: ":" + cfg.MetricsPort, Handler: metricsMux}
	go func() {
		logger.Info("metrics server listening", zap.String("port", cfg.MetricsPort))
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("metrics server error", zap.Error(err))
		}
	}()

	srv := &http.Server{Addr: ":" + cfg.HTTPPort, Handler: r}
	go func() {
		logger.Info("api server listening", zap.String("port", cfg.HTTPPort))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("api server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", zap.Error(err))
	}
	_ = metricsServer.Shutdown(shutdownCtx)
}

func docsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "docs/openapi.yaml")
}
