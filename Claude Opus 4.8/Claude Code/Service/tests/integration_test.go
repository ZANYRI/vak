//go:build integration

// Package tests holds integration tests that require a live PostgreSQL.
// Run with: go test -tags=integration ./tests/...
// Requires TEST_DATABASE_URL (or DATABASE_URL) pointing at an empty database.
package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/example/billing-service/internal/database"
	"github.com/example/billing-service/internal/observability"
	"github.com/example/billing-service/internal/usage"
	"github.com/google/uuid"
)

func dsn(t *testing.T) string {
	if v := os.Getenv("TEST_DATABASE_URL"); v != "" {
		return v
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v
	}
	t.Skip("no TEST_DATABASE_URL set; skipping integration test")
	return ""
}

func TestUsageIdempotencyIntegration(t *testing.T) {
	ctx := context.Background()
	pool, err := database.Connect(ctx, dsn(t), 5)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	log := observability.NewLogger("test")
	if err := database.Migrate(ctx, pool, log); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	// Seed a customer, plan, subscription so foreign keys are satisfied.
	var customerID, planID, subID uuid.UUID
	if err := pool.QueryRow(ctx, `INSERT INTO customers (email, name, currency)
		VALUES ('it@example.com','IT','USD') RETURNING id`).Scan(&customerID); err != nil {
		t.Fatalf("seed customer: %v", err)
	}
	if err := pool.QueryRow(ctx, `INSERT INTO plans (name, currency, billing_interval, pricing_model, base_price_cents)
		VALUES ('IT Plan','USD','monthly','usage_based',0) RETURNING id`).Scan(&planID); err != nil {
		t.Fatalf("seed plan: %v", err)
	}
	if err := pool.QueryRow(ctx, `INSERT INTO subscriptions (customer_id, plan_id, current_period_end)
		VALUES ($1,$2, now() + interval '30 days') RETURNING id`, customerID, planID).Scan(&subID); err != nil {
		t.Fatalf("seed subscription: %v", err)
	}

	svc := usage.NewService(usage.NewRepository(pool))
	key := "idem-" + uuid.NewString()
	req := usage.RecordRequest{
		CustomerID: customerID, SubscriptionID: subID, Metric: "api_calls",
		Quantity: 100, IdempotencyKey: key, RecordedAt: ptr(time.Now()),
	}

	first, dup1, err := svc.Record(ctx, req)
	if err != nil || dup1 {
		t.Fatalf("first record: err=%v dup=%v", err, dup1)
	}
	second, dup2, err := svc.Record(ctx, req)
	if err != nil {
		t.Fatalf("second record: %v", err)
	}
	if !dup2 {
		t.Fatal("re-submitting same idempotency key should be detected as duplicate")
	}
	if first.ID != second.ID {
		t.Fatalf("duplicate should return original event id: %s != %s", first.ID, second.ID)
	}

	total, err := svc.UsageForPeriod(ctx, subID, "api_calls", time.Now().Add(-time.Hour), time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("sum: %v", err)
	}
	if total != 100 {
		t.Fatalf("usage should be counted once: got %d want 100", total)
	}
}

func ptr[T any](v T) *T { return &v }
