# Billing Service

A production-ready backend service for subscription management and billing calculations written in Go.

## Features

- JWT-based authentication with access & refresh tokens
- Role-based access control (admin, billing_manager, support, customer)
- Customer management
- Subscription plans with flat, per-seat, usage-based, tiered and hybrid pricing
- Subscription lifecycle (trialing, active, paused, cancelled, past_due, expired)
- Usage event reporting with idempotency
- Coupon & discount support (percentage and fixed amount)
- Tax rules with basis-point rates
- Invoice generation with line items, discounts, taxes, proration
- Mock payment simulation
- Background job queue (NATS JetStream) and workers (cmd/worker)
- Scheduler (cmd/scheduler) for renewals, trial expiry, usage aggregation
- PostgreSQL persistence with migrations
- Structured logging, request/correlation IDs, Prometheus metrics
- OpenAPI documentation at `/api/docs`
- Docker Compose deployment

## Architecture

```
┌────────┐  ┌────────────┐  ┌───────────┐
│  API   │  │   Worker   │  │ Scheduler │
└───┬────┘  └──────┬─────┘  └─────┬─────┘
    │              │              │
    ├──────────────┼──────────────┤
    │         PostgreSQL          │
    │         NATS JetStream      │
    └─────────────────────────────┘
```

## Tech Stack

- Go 1.24+
- PostgreSQL 16
- NATS JetStream
- chi router
- pgx
- golang-migrate
- golang-jwt/jwt/v5
- bcrypt
- Prometheus client
- Zap logging
- Docker & Docker Compose

## How Billing Calculations Work

All monetary values are stored and calculated in minor units (cents/pence).

For each invoice the engine calculates:

1. Base plan price
2. Per-seat charges (seats above included amount)
3. Usage charges (units above included amount)
4. Tiered usage charges
5. Proration adjustments
6. Coupon discounts
7. Taxes using basis points

Formula:

```
tax_cents = taxable_amount_cents * tax_rate_basis_points / 10000
```

## How Proration Works

When a customer changes plans mid-period:

```
unused_credit = old_plan_price * unused_days / total_days
new_charge    = new_plan_price * unused_days / total_days
proration     = new_charge - unused_credit
```

A positive value is charged; a negative value is a credit. The result is stored as a `proration` invoice line.

## How Usage-Based Billing Works

Usage events are submitted to `POST /api/v1/usage` with an idempotency key. During invoice generation, events in the billing period are summed for the metric `api_requests` and charged according to the plan price configuration.

## How Authentication Works

- Register at `POST /api/v1/auth/register`
- Login at `POST /api/v1/auth/login` to receive access and refresh tokens
- Include `Authorization: Bearer <access_token>` on protected endpoints
- Refresh tokens at `POST /api/v1/auth/refresh`
- Logout at `POST /api/v1/auth/logout`

Passwords are hashed with bcrypt. Refresh tokens are rotated on each use.

## How Authorization Works

Roles and permissions:

| Role             | Permissions                                  |
|------------------|----------------------------------------------|
| admin            | manage_all                                   |
| billing_manager  | plans, invoices, subscriptions, coupons, tax |
| support          | view customers and invoices                  |
| customer         | view own subscriptions/invoices              |

Each protected route checks the required permission.

## How Queues and Workers Work

The API and scheduler publish jobs to NATS JetStream. Jobs are persisted in PostgreSQL with status, attempts and retry scheduling. The worker process consumes each subject, executes handlers, records metrics and supports:
- Exponential backoff retries
- Dead-letter / dead job visibility
- Panic recovery
- Graceful shutdown

## Local Development

Prerequisites: Go 1.24+, Docker, Docker Compose.

1. Copy the environment file:

```bash
cp .env.example .env
```

2. Start dependencies and services:

```bash
docker compose up --build
```

The API is available at http://localhost:8080.

3. View API docs at http://localhost:8080/api/docs.

## Docker Setup

The compose file starts:

- `postgres` - database
- `queue` - NATS with JetStream
- `migrate` - golang-migrate up
- `api` - HTTP API
- `worker` - background worker
- `scheduler` - cron scheduler

Run migrations manually:

```bash
make migrate-up
```

## API Documentation

OpenAPI spec is at `docs/openapi.yaml` and served at `/api/docs`.

## Testing

```bash
go test ./...
```

Integration tests can be run with `docker compose up` and the example requests below.

## Environment Variables

| Variable               | Description                     |
|------------------------|---------------------------------|
| APP_ENV                | local / production              |
| HTTP_PORT              | API port                        |
| METRICS_PORT           | Prometheus metrics port         |
| DATABASE_URL           | PostgreSQL connection string    |
| JWT_ACCESS_SECRET      | JWT access token secret         |
| JWT_REFRESH_SECRET     | JWT refresh token secret        |
| JWT_ACCESS_TTL         | Access token TTL                |
| JWT_REFRESH_TTL        | Refresh token TTL               |
| QUEUE_URL              | NATS URL                        |
| QUEUE_STREAM           | JetStream stream name           |
| SCHEDULER_INTERVAL     | Cron expression for scheduler   |

## Example API Requests

### Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password123","name":"Admin","role":"admin"}'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password123"}'
```

### Create a plan

```bash
curl -X POST http://localhost:8080/api/v1/plans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name":"Starter","billing_interval":"monthly","base_price_cents":1900,"currency":"USD"}'
```

### Create a customer and subscription

```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"email":"customer@example.com","name":"Customer","country":"US","currency":"USD"}'

curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"customer_id":"CUSTOMER_ID","plan_id":"PLAN_ID","quantity":1,"seats":1}'
```

### Report usage

```bash
curl -X POST http://localhost:8080/api/v1/usage \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -H "Idempotency-Key: usage-1" \
  -d '{"customer_id":"CUSTOMER_ID","subscription_id":"SUB_ID","metric":"api_requests","quantity":5000,"idempotency_key":"usage-1"}'
```

### Generate and pay an invoice

```bash
curl -X POST http://localhost:8080/api/v1/invoices/generate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -H "Idempotency-Key: inv-1" \
  -d '{"subscription_id":"SUB_ID","customer_id":"CUSTOMER_ID"}'

curl -X POST http://localhost:8080/api/v1/invoices/INVOICE_ID/pay \
  -H "Authorization: Bearer <token>"
```

### Simulate payment

```bash
curl -X POST http://localhost:8080/api/v1/payments/simulate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -H "Idempotency-Key: pay-1" \
  -d '{"invoice_id":"INVOICE_ID","card_number":"4242424242420000"}'
```

## Future Improvements

- Webhook notifications for invoice/payment events
- Multi-currency support with exchange rates
- Stripe/Adyen payment provider integration
- Recurring billing with automated dunning
- Admin dashboard / analytic metrics
- gRPC API
