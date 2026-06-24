# Billing Service

Backend service for subscription management, billing calculations, invoices, usage metering and simulated payments. It is a Go modular monolith with PostgreSQL persistence and NATS JetStream background jobs. There is no frontend.

## Run locally

```bash
cp .env.example .env
docker compose up --build
```

The API is available at `http://localhost:8080`; the OpenAPI document is served at `http://localhost:8080/api/docs`. Docker creates the development administrator from `BOOTSTRAP_ADMIN_EMAIL` and `BOOTSTRAP_ADMIN_PASSWORD` (default: `admin@example.com` / `admin-password-change-me`). Change both secrets outside local development.

```bash
go mod download
go test ./...
go run ./cmd/api       # with PostgreSQL and NATS already running
```

## Load test

The included read-only load command logs in once and concurrently invokes a protected endpoint. It does not create billing records.

```bash
go run ./cmd/loadtest -requests 2000 -concurrency 40
```

Use `-endpoint`, `-base-url`, `-timeout`, or the matching `LOAD_*` environment variables to change the target. The command exits non-zero on any non-2xx response or transport error.

## Architecture

```text
HTTP API ── PostgreSQL (migrations, domain state, idempotency, audit)
   │
   └── NATS JetStream ── worker (retries, backoff, DLQ state) 
             ▲
        scheduler (publishes periodic jobs only)
```

`cmd/api` applies migrations and exposes the REST API. `cmd/worker` consumes persisted jobs and records `queued`, `running`, `completed`, `retrying`, or `dead` state. `cmd/scheduler` emits renewal, trial-expiry and usage jobs every minute; it does not perform heavy work itself.

## Main features

- JWT access tokens, rotating hashed refresh tokens, bcrypt passwords, auth rate limiting, CORS/security headers and role-based access control.
- Roles: `admin`, `billing_manager`, `support`, `customer`. Customers are limited to records linked to their account.
- Flat, per-seat, usage, tiered and hybrid price models. Values are always integer minor units (`*_cents`); floating point is never used.
- Coupons (percentage/fixed), country tax rules in basis points, trials, upgrades/downgrades with invoice proration, invoices and mock payments.
- Idempotency keys on subscription creation, usage recording, invoice generation and payment simulation. Reusing a key with a distinct body returns `IDEMPOTENCY_CONFLICT`.
- Structured JSON logs, request IDs, `/healthz`, `/readyz`, `/metrics`, audit records and graceful shutdown.

## Billing rules

Base price is combined with the selected price-model charges. Per-seat charges use `max(quantity - included_seats, 0)`. Usage charges use `max(usage - included_units, 0)`. Tier charges calculate each tier's covered unit range. A valid coupon is applied before tax. Tax is calculated exactly as:

```text
tax_cents = (subtotal_cents - discount_cents) * tax_rate_basis_points / 10000
```

All integer division rounds toward zero. Proration uses the portion of the current period remaining: the old-plan credit and new-plan charge are recorded as separate invoice lines.

## Example workflow

Register or log in, copy `access_token`, then use it as a Bearer token. The bootstrap administrator can create plans.

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@example.com","password":"admin-password-change-me"}'

curl -X POST http://localhost:8080/api/v1/plans \
  -H 'Authorization: Bearer <access_token>' -H 'Content-Type: application/json' \
  -d '{"name":"Starter","currency":"USD","billing_interval":"monthly","pricing_model":"flat","base_price_cents":1900}'

curl -X POST http://localhost:8080/api/v1/usage \
  -H 'Authorization: Bearer <access_token>' -H 'Idempotency-Key: usage-001' \
  -H 'Content-Type: application/json' \
  -d '{"customer_id":"<uuid>","subscription_id":"<uuid>","metric":"api_calls","quantity":120}'
```

For payment simulation, a card number ending in `0000` succeeds, one ending in `9999` fails, and every other last-four value has a random outcome. Only the last four digits are stored.

## Configuration

| Variable | Purpose |
| --- | --- |
| `DATABASE_URL` | PostgreSQL connection string |
| `QUEUE_URL` | NATS server URL |
| `JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET` | Separate secrets, minimum 32 bytes |
| `ACCESS_TOKEN_TTL`, `REFRESH_TOKEN_TTL` | Token lifetimes |
| `CORS_ORIGINS` | Comma-separated allowed origins |
| `BOOTSTRAP_ADMIN_EMAIL`, `BOOTSTRAP_ADMIN_PASSWORD` | One-time development administrator |

## API and tests

The complete endpoint contract and request schemas are in [`docs/openapi.yaml`](docs/openapi.yaml). Tests are run with `go test ./...`; calculation tests cover tiers, discounts/tax and proration, and authentication tests cover tokens and password handling.

## Future improvements

Production deployments would typically add an external secrets manager, distributed rate limiting, a real payment-provider adapter, email delivery adapter, OpenTelemetry export, a generated interactive Swagger UI, and a dedicated migration job. The current implementation keeps those integrations deliberately replaceable without changing the domain model.
