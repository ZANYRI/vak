// Package app wires together every component of the billing service so the
// api, worker, and scheduler binaries share a single construction path.
package app

import (
	"context"
	"log/slog"

	"github.com/example/billing-service/internal/audit"
	"github.com/example/billing-service/internal/auth"
	"github.com/example/billing-service/internal/config"
	"github.com/example/billing-service/internal/coupons"
	"github.com/example/billing-service/internal/customers"
	"github.com/example/billing-service/internal/database"
	"github.com/example/billing-service/internal/invoices"
	"github.com/example/billing-service/internal/observability"
	"github.com/example/billing-service/internal/payments"
	"github.com/example/billing-service/internal/plans"
	"github.com/example/billing-service/internal/queue"
	"github.com/example/billing-service/internal/subscriptions"
	"github.com/example/billing-service/internal/taxes"
	"github.com/example/billing-service/internal/usage"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Container holds all shared dependencies and services.
type Container struct {
	Cfg     *config.Config
	Log     *slog.Logger
	Metrics *observability.Metrics
	Pool    *pgxpool.Pool
	Queue   *queue.Queue
	Audit   *audit.Logger

	Tokens *auth.TokenManager
	AuthMW *auth.Middleware

	Auth      *auth.Service
	Customers *customers.Service
	Plans     *plans.Service
	Subs      *subscriptions.Service
	Coupons   *coupons.Service
	Taxes     *taxes.Service
	Usage     *usage.Service
	Invoices  *invoices.Service
	Payments  *payments.Service
}

// Build constructs the full dependency graph. The caller owns Close().
func Build(ctx context.Context, cfg *config.Config, log *slog.Logger) (*Container, error) {
	metrics := observability.NewMetrics()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	if err := database.Migrate(ctx, pool, log); err != nil {
		pool.Close()
		return nil, err
	}

	q, err := queue.New(ctx, cfg.QueueURL, pool, cfg.JobMaxAttempts)
	if err != nil {
		pool.Close()
		return nil, err
	}

	auditLog := audit.NewLogger(pool, log)

	tokens := auth.NewTokenManager(cfg.JWTAccessSecret, cfg.JWTRefreshSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	authMW := auth.NewMiddleware(tokens)

	// Independent services.
	authSvc := auth.NewService(auth.NewRepository(pool), tokens, auditLog, cfg.BcryptCost)
	customersSvc := customers.NewService(customers.NewRepository(pool), auditLog)
	plansSvc := plans.NewService(plans.NewRepository(pool), auditLog)
	taxesSvc := taxes.NewService(taxes.NewRepository(pool), auditLog)
	usageSvc := usage.NewService(usage.NewRepository(pool))
	couponsSvc := coupons.NewService(coupons.NewRepository(pool), auditLog)

	// Subscriptions need plans + queue; invoicer is wired after invoices exist.
	subsSvc := subscriptions.NewService(subscriptions.NewRepository(pool), plansSvc, nil, q, auditLog)

	// Invoices integrate everything.
	invoicesSvc := invoices.NewService(invoices.NewRepository(pool), subsSvc, plansSvc, customersSvc,
		usageSvc, couponsSvc, taxesSvc, q, auditLog, metrics)

	subsSvc.SetInvoicer(invoicesSvc)

	paymentsSvc := payments.NewService(payments.NewRepository(pool), invoicesSvc, subsSvc, q, auditLog, metrics)

	// Bootstrap an admin account if configured.
	if err := authSvc.EnsureAdmin(ctx, cfg.AdminEmail(), cfg.AdminPassword()); err != nil {
		log.Warn("admin bootstrap failed", "error", err.Error())
	}

	return &Container{
		Cfg: cfg, Log: log, Metrics: metrics, Pool: pool, Queue: q, Audit: auditLog,
		Tokens: tokens, AuthMW: authMW,
		Auth: authSvc, Customers: customersSvc, Plans: plansSvc, Subs: subsSvc,
		Coupons: couponsSvc, Taxes: taxesSvc, Usage: usageSvc, Invoices: invoicesSvc, Payments: paymentsSvc,
	}, nil
}

// Close releases all resources.
func (c *Container) Close() {
	if c.Queue != nil {
		_ = c.Queue.Close()
	}
	if c.Pool != nil {
		c.Pool.Close()
	}
}
