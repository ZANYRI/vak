# Billing Service

A production-ready backend for **subscription management and billing**, written in Go.
It supports flat, per-seat, usage-based, tiered and hybrid pricing, coupons, taxes,
proration, invoice generation, simulated payments, a Redis Streams job queue with
background workers, and a scheduler — all behind a JWT-authenticated REST API with
role-based access control.

> Backend only. No frontend. All money is stored and computed in **integer minor
> units (cents)** — no floating point is ever used for money.

---

## Table of contents

- [Features](#features)
- [Architecture](#architecture)
- [Tech stack](#tech-stack)
- [Quick start (Docker)](#quick-start-docker)
- [Local development](#local-development)
- [Configuration](#configuration)
- [Authentication & authorization](#authentication--authorization)
- [How billing calculations work](#how-billing-calculations-work)
- [How proration works](#how-proration-works)
- [How usage-based billing works](#how-usage-based-billing-works)
- [Queues & workers](#queues--workers)
- [Scheduler](#scheduler)
- [Idempotency](#idempotency)
- [Observability](#observability)
- [API documentation](#api-documentation)
- [Example API requests](#example-api-requests)
- [Testing](#testing)
- [Project structure](#project-structure)
- [Future improvements](#future-improvements)

---

## Features

- **Authentication** — register, login, refresh-token rotation, logout, current user; bcrypt password hashing; JWT access tokens.
- **RBAC** — `admin`, `billing_manager`, `support`, `customer` roles enforced per endpoint.
- **Customers, Plans, Subscriptions** — full lifecycle (trial, active, past_due, paused, cancelled, expired).
- **Pricing models** — flat, per-seat, usage-based, tiered, hybrid.
- **Usage metering** — idempotent usage events, aggregation, summaries.
- **Coupons & discounts** — percentage and fixed-amount with validation rules.
- **Taxes** — basis-point tax rules with country/region resolution.
- **Invoices** — automatic & manual generation, finalize, void, pay; itemized lines.
- **Proration** — time-based credit/charge for mid-period plan changes.
- **Payments** — mock provider with deterministic test cards.
- **Queue & workers** — Redis Streams broker, retries with exponential backoff, dead-letter queue, idempotency, graceful shutdown, panic recovery.
- **Scheduler** — periodic renewals, trial expiry, overdue retries, usage aggregation.
- **Observability** — structured logs, request/correlation IDs, Prometheus metrics, health & readiness probes.
- **Security** — JWT expiry, refresh rotation, RBAC, auth rate limiting, request validation, parameterized SQL, CORS, secure headers, non-root Docker user, env-based secrets.
- **OpenAPI docs** served at `/api/docs`.

---

## Architecture

```
                      ┌─────────────┐         publish jobs        ┌──────────────┐
   HTTP clients ────▶ │   cmd/api   │ ──────────────────────────▶ │ Redis Streams│
                      │  (REST API) │                             │   (queue)    │
                      └──────┬──────┘                             └──────┬───────┘
                             │ SQL                                       │ consume
                             ▼                                           ▼
                      ┌─────────────┐   publish jobs            ┌──────────────┐
                      │ PostgreSQL  │ ◀──────────────────────── │ cmd/worker   │
                      │             │                           │ (job consumer)│
                      └──────▲──────┘                           └──────────────┘
                             │ SQL                                       ▲
                             │                                           │ publish jobs
                      ┌──────┴───────┐                                   │
                      │ cmd/scheduler│ ──────────────────────────────────┘
                      │ (cron-like)  │
                      └──────────────┘
```

All three binaries share a single dependency-injection container
(`internal/app`) that builds the database pool, queue, and every service.

- **`cmd/api`** serves the REST API and publishes jobs.
- **`cmd/worker`** consumes jobs from Redis Streams and updates the `jobs` table.
- **`cmd/scheduler`** periodically finds due work and publishes jobs (it never does heavy work itself).

The codebase follows a modular-monolith layout: each domain (`auth`, `customers`,
`plans`, `subscriptions`, `invoices`, …) owns its models, repository, service, and
HTTP handler. The pure billing engine (`internal/billing`) has no I/O and is fully
unit-tested.

---

## Tech stack

| Concern        | Choice                                  |
|----------------|-----------------------------------------|
| Language       | Go 1.26                                 |
| HTTP router    | `go-chi/chi`                            |
| Database       | PostgreSQL 16 via `jackc/pgx/v5`        |
| Migrations     | Embedded SQL runner (`go:embed`)        |
| Queue broker   | Redis Streams via `redis/go-redis/v9`   |
| Auth           | `golang-jwt/jwt/v5`, `x/crypto/bcrypt`  |
| Validation     | `go-playground/validator/v10`           |
| Logging        | `log/slog` (structured)                 |
| Metrics        | `prometheus/client_golang`              |
| Docs           | OpenAPI 3.0 + Swagger UI                 |

---

## Quick start (Docker)

```bash
# from the project root
docker compose up --build
```

This starts `postgres`, `queue` (Redis), `api`, `worker`, and `scheduler`.
Migrations run automatically on API startup. The API is available at:

```
http://localhost:8080
```

A default admin account is bootstrapped on first start:

```
email:    admin@billing.local
password: admin12345
```

Open the interactive docs at <http://localhost:8080/api/docs>.

---

## Local development

Requirements: Go 1.26+, a PostgreSQL instance, and a Redis instance.

```bash
cp .env.example .env          # then edit as needed
go mod download

# Start Postgres + Redis only:
docker compose up -d postgres queue

# Run each process in its own terminal:
make run-api
make run-worker
make run-scheduler
```

Useful targets: `make build`, `make test`, `make vet`, `make up`, `make down`.

---

## Configuration

All configuration comes from environment variables (see `.env.example`):

| Variable                 | Default                                   | Description                          |
|--------------------------|-------------------------------------------|--------------------------------------|
| `APP_ENV`                | `local`                                   | `local`, `docker`, `production`, …   |
| `HTTP_PORT`              | `8080`                                    | API listen port                      |
| `DATABASE_URL`           | `postgres://billing:billing@…/billing`    | PostgreSQL DSN                       |
| `QUEUE_URL`              | `redis://localhost:6379/0`                | Redis Streams URL                    |
| `JWT_ACCESS_SECRET`      | `change-me-access`                        | Access-token signing secret          |
| `JWT_REFRESH_SECRET`     | `change-me-refresh`                       | Refresh-token signing secret         |
| `JWT_ACCESS_TTL`         | `15m`                                     | Access-token lifetime                |
| `JWT_REFRESH_TTL`        | `720h`                                    | Refresh-token lifetime               |
| `BCRYPT_COST`            | `12`                                      | bcrypt cost                          |
| `AUTH_RATE_LIMIT`        | `10`                                      | Auth requests per window per IP      |
| `AUTH_RATE_LIMIT_WINDOW` | `1m`                                      | Rate-limit window                    |
| `WORKER_CONCURRENCY`     | `4`                                       | Concurrent consumers per worker      |
| `JOB_MAX_ATTEMPTS`       | `5`                                       | Retries before dead-letter           |
| `SCHEDULER_INTERVAL`     | `60s`                                     | Scheduler tick interval              |
| `CORS_ALLOWED_ORIGINS`   | `*`                                       | Comma-separated origins              |
| `ADMIN_EMAIL` / `ADMIN_PASSWORD` | `admin@billing.local` / `admin12345` | Bootstrap admin account          |

Secrets are never committed; outside `local`/`test` the service refuses to start with default JWT secrets.

---

## Authentication & authorization

- Passwords are hashed with **bcrypt**; plaintext is never stored or logged.
- **Login** returns an `access_token` (short-lived JWT) and a `refresh_token`.
- Send the access token as `Authorization: Bearer <token>` on protected endpoints.
- **Refresh-token rotation**: each refresh revokes the presented token and issues a new pair. Re-using a revoked refresh token revokes the whole chain.
- **Roles & permissions**:
  - `admin` — everything.
  - `billing_manager` — manage plans, invoices, subscriptions, coupons, tax rules.
  - `support` — view customers and invoices; cannot modify billing rules.
  - `customer` — create/report usage and view resources.
- Auth credential endpoints are **rate-limited** per client IP.

---

## How billing calculations work

The engine lives in `internal/billing` and is pure (no DB, fully tested). Given a
plan, subscription quantity, period usage, coupons and a tax rate it produces
itemized lines and totals — all in integer cents.

Order of operations:

1. **Charges** by pricing model:
   - `flat` → base price.
   - `per_seat` → base + `(quantity − included_seats) × seat_price`.
   - `usage_based` → base + `(usage − included_units) × unit_price`.
   - `tiered` → graduated per-tier charges over the cumulative usage ranges.
   - `hybrid` → base + seats + (tiered or usage).
2. **Subtotal** = sum of charge lines (plus any proration lines).
3. **Discounts** from coupons, applied to the subtotal and capped so the result is never negative.
4. **Tax** = `taxable × tax_rate_basis_points / 10000` on the discounted amount.
5. **Total** = discounted subtotal + tax.

Worked example (matches the spec): subtotal `7900`, a `1000`-cent fixed coupon, and
a `2000` bp (20%) tax rule →

```
subtotal  = 7900
discount  = 1000
taxable   = 6900
tax       = 6900 * 2000 / 10000 = 1380
total     = 8280
```

---

## How proration works

When a subscription changes plan mid-period (`internal/billing/proration.go`), the
engine computes a **time-weighted** credit and charge:

```
credit = old_price × remaining_seconds / total_seconds   (unused old plan)
charge = new_price × remaining_seconds / total_seconds    (remaining new plan)
difference = charge − credit
```

Example — a `3000`-cent monthly plan upgraded **halfway** to a `6000`-cent plan:

```
credit     = 3000 × 15/30 = 1500
charge     = 6000 × 15/30 = 3000
difference = 1500
```

Two proration lines (a negative credit and a positive charge) are stored on a new
invoice. Downgrades produce a net credit (negative difference).

---

## How usage-based billing works

- Report usage via `POST /api/v1/usage` with a unique `idempotency_key`.
- Re-submitting the same key returns the **original** event and never double-counts (enforced by a unique DB index).
- At invoice generation, usage is summed for the metric over the billing period and charged according to the plan's included units / unit price (or tiers).
- `GET /api/v1/usage/summary?subscription_id=…` returns per-metric totals.

---

## Queues & workers

- The broker is **Redis Streams** (`billing:jobs`), with a dead-letter stream (`billing:jobs:dead`).
- Every job is also recorded in the `jobs` table with a lifecycle status (`queued → running → completed | retrying | failed | dead`) for observability.
- Workers (`cmd/worker`) provide **retries with exponential backoff**, a **dead-letter queue**, **idempotency** (completed jobs are skipped), **structured logs**, **panic recovery**, **context timeouts**, and **graceful shutdown**.

Job types: `invoice.generate`, `invoice.finalize`, `payment.process`,
`subscription.renew`, `subscription.expire_trial`, `email.invoice_created`,
`email.payment_failed`, `usage.aggregate`.

---

## Scheduler

`cmd/scheduler` ticks on `SCHEDULER_INTERVAL` and **only publishes jobs** (no heavy
work). Each tick it:

- finds subscriptions due for renewal → `subscription.renew`
- finds expired trials → `subscription.expire_trial`
- finds overdue invoices → `payment.process` (retry)
- emits a `usage.aggregate` tick

Renewal and trial jobs use idempotency keys so a job is not enqueued twice for the same period.

---

## Idempotency

Two layers:

1. **HTTP** — send an `Idempotency-Key` header on `POST /usage`, `POST /invoices/generate`, `POST /payments/simulate`, `POST /subscriptions`. The same key + same body replays the stored response; the same key + a different body returns `IDEMPOTENCY_CONFLICT`.
2. **Domain** — usage events are idempotent by `idempotency_key`; queue jobs may carry an idempotency key enforced by a unique index.

---

## Observability

- **Structured logs** (`slog`) with request & correlation IDs.
- **Health**: `GET /healthz` (liveness), `GET /readyz` (checks Postgres + Redis).
- **Metrics**: `GET /metrics` (Prometheus) including `http_requests_total`,
  `http_request_duration_seconds`, `jobs_processed_total`, `jobs_failed_total`,
  `invoices_generated_total`, `payments_succeeded_total`, `payments_failed_total`,
  `queue_depth`.
- **Graceful shutdown** on SIGINT/SIGTERM across all processes.

Errors use a consistent envelope:

```json
{ "error": { "code": "VALIDATION_ERROR", "message": "…", "details": { } } }
```

---

## API documentation

- Swagger UI: <http://localhost:8080/api/docs>
- Raw spec: <http://localhost:8080/api/docs/openapi.yaml> (also at `docs/openapi.yaml`)

---

## Example API requests

```bash
BASE=http://localhost:8080/api/v1

# 1. Log in as the bootstrap admin
TOKENS=$(curl -s -X POST $BASE/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@billing.local","password":"admin12345"}')
ACCESS=$(echo "$TOKENS" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
AUTH="Authorization: Bearer $ACCESS"

# 2. Create a customer
curl -s -X POST $BASE/customers -H "$AUTH" -H 'Content-Type: application/json' \
  -d '{"email":"acme@example.com","name":"Acme","country":"US","currency":"USD"}'

# 3. Create a flat plan
curl -s -X POST $BASE/plans -H "$AUTH" -H 'Content-Type: application/json' \
  -d '{"name":"Starter","currency":"USD","billing_interval":"monthly","pricing_model":"flat","base_price_cents":1900}'

# 4. Create a subscription (use the ids returned above)
curl -s -X POST $BASE/subscriptions -H "$AUTH" -H 'Content-Type: application/json' \
  -H 'Idempotency-Key: sub-001' \
  -d '{"customer_id":"<CUSTOMER_ID>","plan_id":"<PLAN_ID>","quantity":1}'

# 5. Generate and pay an invoice
curl -s -X POST $BASE/invoices/generate -H "$AUTH" -H 'Content-Type: application/json' \
  -d '{"subscription_id":"<SUBSCRIPTION_ID>"}'
curl -s -X POST $BASE/payments/simulate -H "$AUTH" -H 'Content-Type: application/json' \
  -d '{"invoice_id":"<INVOICE_ID>","card_number":"4242424242420000"}'   # ends 0000 → succeeds
```

Payment test cards: ends in `0000` → **succeeds**, ends in `9999` → **fails**, anything else → random.

---

## Testing

```bash
go test ./...                       # unit tests (no infra required)
go test -tags=integration ./tests/...   # integration tests (needs TEST_DATABASE_URL)
```

Unit tests cover money & tax calculations, proration, coupon validation, JWT,
RBAC helpers and payment simulation. The integration suite exercises usage
idempotency against a real PostgreSQL.

---

## Project structure

```
billing-service/
├── cmd/
│   ├── api/          # REST API process
│   ├── worker/       # queue consumer process
│   └── scheduler/    # periodic job publisher
├── internal/
│   ├── app/          # dependency-injection container
│   ├── auth/         # users, JWT, RBAC middleware
│   ├── customers/ plans/ subscriptions/
│   ├── usage/ coupons/ taxes/ invoices/ payments/
│   ├── billing/      # pure pricing & proration engine (tested)
│   ├── queue/        # Redis Streams broker + consumer
│   ├── workers/      # job handlers
│   ├── scheduler/    # periodic enqueuer
│   ├── audit/        # audit logging
│   ├── database/     # pool + embedded migrations
│   ├── middleware/   # request id, recover, rate limit, idempotency, …
│   ├── observability/# logging, metrics, health
│   ├── httpx/        # request/response/error helpers
│   ├── config/       # env configuration
│   └── server/       # router + OpenAPI docs
├── migrations/       # SQL migrations (also embedded in the binary)
├── docs/             # openapi.yaml
├── tests/            # integration tests (build tag: integration)
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---

## Future improvements

- Per-customer row-level scoping so `customer` users see only their own data.
- Replace the in-memory auth rate limiter with a Redis-backed limiter for multi-instance deployments.
- Webhooks/event bus for invoice & payment lifecycle events.
- Real payment-provider adapters behind the existing payment interface.
- A dead-letter inspection/replay API and per-job retry policies.
- `sqlc`-generated type-safe queries and a richer test suite (testcontainers).
