# Prompt: Subscription & Billing Backend Service in Go

## Task

Create a complete production-ready backend service for **subscription management and billing calculations**.

The service must be written in **Go** and must include:

* User authentication
* Role-based authorization
* Subscription plans
* Customers
* Subscriptions
* Usage-based billing
* Invoice generation
* Billing calculations
* Coupons and discounts
* Taxes
* Proration
* Payment simulation
* Queue broker
* Background workers
* PostgreSQL persistence
* Dockerized deployment
* OpenAPI documentation
* Tests

This is a backend-only application. Do not create a frontend.

---

## Main Idea

Build a backend system that allows companies to manage subscription billing.

The system should support different pricing models:

* Fixed monthly subscription
* Fixed yearly subscription
* Usage-based billing
* Tiered pricing
* Per-seat pricing
* Discounts
* Coupons
* Tax calculation
* Trial periods
* Subscription upgrades and downgrades
* Invoice generation
* Payment status simulation

The backend must expose a REST API.

---

## Tech Stack

Use modern stable tools and best practices.

Required stack:

* Go
* PostgreSQL
* Queue broker: NATS JetStream, RabbitMQ, or Redis Streams
* Docker
* Docker Compose
* JWT authentication
* Refresh tokens
* Role-based access control
* SQL migrations
* Structured logging
* OpenAPI / Swagger documentation
* Unit and integration tests

Preferred Go libraries:

* `chi` or `gin` for HTTP routing
* `pgx` for PostgreSQL
* `sqlc` or repository pattern for DB access
* `golang-migrate` for migrations
* `go-playground/validator` for request validation
* `zap`, `zerolog`, or `slog` for logging
* `testcontainers-go` for integration tests
* `swaggo` or OpenAPI YAML for API documentation

Do not use floating-point numbers for money calculations.

All money values must be stored and calculated in minor units:

```txt
USD cents
EUR cents
GBP pence
```

Example:

```json
{
  "amount_cents": 1999,
  "currency": "USD"
}
```

---

## Core Features

### 1. Authentication

Implement authentication with:

* User registration
* User login
* Password hashing
* Access token
* Refresh token
* Logout
* Token refresh
* Current user endpoint

Endpoints:

```txt
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
GET  /api/v1/auth/me
```

Use secure password hashing.

Do not store plain-text passwords.

---

### 2. Authorization

Implement role-based access control.

Roles:

```txt
admin
billing_manager
customer
support
```

Permissions:

* `admin` can manage everything
* `billing_manager` can manage plans, invoices, subscriptions, coupons
* `support` can view customers and invoices but cannot modify billing rules
* `customer` can view only their own subscriptions and invoices

Every protected endpoint must check permissions.

---

### 3. Customers

Implement customer management.

Customer fields:

```txt
id
user_id
email
name
company_name
billing_address
country
tax_id
currency
created_at
updated_at
```

Endpoints:

```txt
POST   /api/v1/customers
GET    /api/v1/customers
GET    /api/v1/customers/{id}
PATCH  /api/v1/customers/{id}
DELETE /api/v1/customers/{id}
```

---

### 4. Subscription Plans

Create billing plans.

Plan fields:

```txt
id
name
description
currency
billing_interval
base_price_cents
trial_days
is_active
created_at
updated_at
```

Supported billing intervals:

```txt
monthly
yearly
```

Plan pricing models:

```txt
flat
per_seat
usage_based
tiered
hybrid
```

Endpoints:

```txt
POST   /api/v1/plans
GET    /api/v1/plans
GET    /api/v1/plans/{id}
PATCH  /api/v1/plans/{id}
DELETE /api/v1/plans/{id}
```

---

### 5. Pricing Rules

Support complex pricing rules.

Examples:

#### Flat Monthly Plan

```json
{
  "name": "Starter",
  "billing_interval": "monthly",
  "base_price_cents": 1900,
  "currency": "USD"
}
```

#### Per-seat Plan

```json
{
  "name": "Team",
  "base_price_cents": 1200,
  "seat_price_cents": 700,
  "included_seats": 3
}
```

#### Usage-based Plan

```json
{
  "name": "API Pro",
  "base_price_cents": 2900,
  "included_units": 10000,
  "unit_price_cents": 2
}
```

#### Tiered Plan

```json
{
  "tiers": [
    {
      "from": 0,
      "to": 1000,
      "unit_price_cents": 5
    },
    {
      "from": 1001,
      "to": 10000,
      "unit_price_cents": 3
    },
    {
      "from": 10001,
      "to": null,
      "unit_price_cents": 1
    }
  ]
}
```

---

### 6. Subscriptions

Implement subscription lifecycle.

Subscription statuses:

```txt
trialing
active
past_due
paused
cancelled
expired
```

Subscription fields:

```txt
id
customer_id
plan_id
status
quantity
current_period_start
current_period_end
trial_start
trial_end
cancel_at_period_end
created_at
updated_at
```

Endpoints:

```txt
POST   /api/v1/subscriptions
GET    /api/v1/subscriptions
GET    /api/v1/subscriptions/{id}
PATCH  /api/v1/subscriptions/{id}
POST   /api/v1/subscriptions/{id}/cancel
POST   /api/v1/subscriptions/{id}/pause
POST   /api/v1/subscriptions/{id}/resume
POST   /api/v1/subscriptions/{id}/change-plan
```

Changing a subscription plan must support proration.

---

## Billing Calculations

### 1. Invoice Calculation

The service must calculate invoice totals.

Invoice calculation must include:

```txt
base price
seat quantity
usage charges
tiered usage
discounts
coupons
taxes
proration
final total
```

Invoice totals:

```txt
subtotal_cents
discount_cents
tax_cents
total_cents
amount_due_cents
amount_paid_cents
```

Example response:

```json
{
  "invoice_id": "inv_123",
  "currency": "USD",
  "subtotal_cents": 7900,
  "discount_cents": 1000,
  "tax_cents": 1380,
  "total_cents": 8280,
  "amount_due_cents": 8280
}
```

---

### 2. Proration

Implement proration for plan upgrades and downgrades.

Example:

A customer has a monthly plan for `3000` cents.

They upgrade halfway through the billing period to a plan that costs `6000` cents.

The system should calculate:

```txt
unused old plan credit
remaining new plan charge
proration difference
```

Store proration lines on the invoice.

Example invoice line:

```json
{
  "type": "proration",
  "description": "Plan upgrade proration",
  "amount_cents": 1500
}
```

---

### 3. Usage-based Billing

Support reporting usage events.

Usage event fields:

```txt
id
customer_id
subscription_id
metric
quantity
idempotency_key
recorded_at
created_at
```

Endpoints:

```txt
POST /api/v1/usage
GET  /api/v1/usage
GET  /api/v1/usage/summary
```

Usage reporting must be idempotent.

If the same `idempotency_key` is submitted twice, do not double-charge.

---

### 4. Coupons and Discounts

Implement coupons.

Coupon types:

```txt
percentage
fixed_amount
```

Coupon fields:

```txt
id
code
type
percent_off
amount_off_cents
currency
max_redemptions
times_redeemed
valid_from
valid_until
is_active
created_at
updated_at
```

Endpoints:

```txt
POST   /api/v1/coupons
GET    /api/v1/coupons
GET    /api/v1/coupons/{id}
PATCH  /api/v1/coupons/{id}
DELETE /api/v1/coupons/{id}
POST   /api/v1/subscriptions/{id}/apply-coupon
```

Validation rules:

* Expired coupons cannot be applied
* Inactive coupons cannot be applied
* Coupon currency must match invoice currency for fixed amount coupons
* Percentage coupons must be between 1 and 100
* Fixed amount discount cannot make invoice total negative

---

### 5. Tax Calculation

Implement simple tax rules.

Tax rule fields:

```txt
id
country
region
tax_name
tax_rate_basis_points
is_active
created_at
updated_at
```

Use basis points for tax rates.

Example:

```txt
2100 basis points = 21.00%
```

Tax calculation:

```txt
tax_cents = taxable_amount_cents * tax_rate_basis_points / 10000
```

Endpoints:

```txt
POST   /api/v1/tax-rules
GET    /api/v1/tax-rules
PATCH  /api/v1/tax-rules/{id}
DELETE /api/v1/tax-rules/{id}
```

---

## Invoice System

Create invoices automatically and manually.

Invoice fields:

```txt
id
customer_id
subscription_id
status
currency
subtotal_cents
discount_cents
tax_cents
total_cents
amount_due_cents
amount_paid_cents
period_start
period_end
issued_at
due_at
paid_at
created_at
updated_at
```

Invoice statuses:

```txt
draft
open
paid
void
uncollectible
```

Invoice line fields:

```txt
id
invoice_id
type
description
quantity
unit_amount_cents
amount_cents
metadata
created_at
```

Endpoints:

```txt
POST /api/v1/invoices/generate
GET  /api/v1/invoices
GET  /api/v1/invoices/{id}
POST /api/v1/invoices/{id}/finalize
POST /api/v1/invoices/{id}/void
POST /api/v1/invoices/{id}/pay
```

---

## Payment Simulation

Do not integrate a real payment provider.

Implement a mock payment provider.

Payment statuses:

```txt
pending
succeeded
failed
refunded
```

Endpoints:

```txt
POST /api/v1/payments/simulate
GET  /api/v1/payments
GET  /api/v1/payments/{id}
```

Payment simulation rules:

* Card number ending in `0000` succeeds
* Card number ending in `9999` fails
* Any other card randomly succeeds or fails
* On successful payment, mark invoice as `paid`
* On failed payment, mark subscription as `past_due`

---

## Queue Broker and Workers

The service must include a queue broker.

Use it for background jobs:

```txt
invoice.generate
invoice.finalize
payment.process
subscription.renew
subscription.expire_trial
email.invoice_created
email.payment_failed
usage.aggregate
```

Create separate worker process:

```txt
cmd/api
cmd/worker
```

The API process must publish jobs.

The worker process must consume jobs.

Workers must support:

* retries
* exponential backoff
* dead-letter queue
* idempotency
* structured logs
* graceful shutdown
* panic recovery
* context timeouts

Job statuses:

```txt
queued
running
completed
failed
retrying
dead
```

---

## Scheduler

Implement a scheduler process or scheduler worker.

It must periodically:

* find subscriptions that need renewal
* generate invoices
* expire trials
* mark overdue invoices
* retry failed payments
* aggregate usage
* send billing notifications

Suggested command:

```txt
cmd/scheduler
```

The scheduler should not directly execute heavy work.

It should publish jobs to the queue.

---

## Database

Use PostgreSQL.

Create migrations for all tables.

Required tables:

```txt
users
refresh_tokens
customers
plans
plan_prices
pricing_tiers
subscriptions
usage_events
usage_summaries
coupons
subscription_coupons
tax_rules
invoices
invoice_lines
payments
jobs
audit_logs
```

Use UUID primary keys.

Add indexes for:

```txt
users.email
customers.user_id
subscriptions.customer_id
subscriptions.status
invoices.customer_id
invoices.status
usage_events.subscription_id
usage_events.idempotency_key
coupons.code
jobs.status
```

---

## Idempotency

Implement idempotency for critical operations.

Required idempotent operations:

```txt
POST /api/v1/usage
POST /api/v1/invoices/generate
POST /api/v1/payments/simulate
POST /api/v1/subscriptions
```

Use `Idempotency-Key` header.

If the same key is reused with the same request body, return the previous response.

If the same key is reused with a different request body, return an error.

---

## Audit Logs

Create audit logs for important actions.

Audit log fields:

```txt
id
actor_user_id
action
resource_type
resource_id
metadata
ip_address
user_agent
created_at
```

Log actions such as:

```txt
user.login
plan.created
subscription.created
subscription.cancelled
invoice.generated
invoice.paid
payment.failed
coupon.applied
tax_rule.updated
```

---

## API Documentation

Generate OpenAPI documentation.

Include:

* Auth endpoints
* Customer endpoints
* Plan endpoints
* Subscription endpoints
* Usage endpoints
* Coupon endpoints
* Tax endpoints
* Invoice endpoints
* Payment endpoints
* Admin endpoints
* Error responses
* Auth requirements
* Example requests
* Example responses

Expose documentation at:

```txt
/api/docs
```

---

## Error Handling

Use consistent error format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request payload",
    "details": {
      "field": "email",
      "reason": "Invalid email format"
    }
  }
}
```

Common error codes:

```txt
VALIDATION_ERROR
UNAUTHORIZED
FORBIDDEN
NOT_FOUND
CONFLICT
IDEMPOTENCY_CONFLICT
PAYMENT_FAILED
INSUFFICIENT_PERMISSIONS
INTERNAL_ERROR
```

---

## Observability

Add:

* structured logs
* request ID
* correlation ID
* metrics endpoint
* health endpoint
* readiness endpoint
* graceful shutdown

Endpoints:

```txt
GET /healthz
GET /readyz
GET /metrics
```

Metrics should include:

```txt
http_requests_total
http_request_duration_seconds
jobs_processed_total
jobs_failed_total
invoices_generated_total
payments_succeeded_total
payments_failed_total
queue_depth
```

---

## Security Requirements

Implement:

* Secure password hashing
* JWT expiration
* Refresh token rotation
* RBAC
* Rate limiting for auth endpoints
* Request validation
* SQL injection protection
* CORS configuration
* Secure headers
* Environment-based secrets
* No secrets committed to repository
* Non-root Docker user

Do not log passwords, tokens, card numbers, or sensitive payment data.

---

## Docker Requirements

Create:

```txt
Dockerfile
docker-compose.yml
.dockerignore
.env.example
README.md
```

The whole system must start with:

```bash
docker compose up --build
```

Services:

```txt
api
worker
scheduler
postgres
queue
```

The API must be available at:

```txt
http://localhost:8080
```

Use environment variables for configuration.

Example:

```env
APP_ENV=local
HTTP_PORT=8080
DATABASE_URL=postgres://billing:billing@postgres:5432/billing?sslmode=disable
JWT_ACCESS_SECRET=change-me
JWT_REFRESH_SECRET=change-me
QUEUE_URL=nats://queue:4222
```

---

## Project Structure

Use clean architecture or modular monolith structure:

```txt
billing-service/
├── cmd/
│   ├── api/
│   │   └── main.go
│   ├── worker/
│   │   └── main.go
│   └── scheduler/
│       └── main.go
├── internal/
│   ├── auth/
│   ├── customers/
│   ├── plans/
│   ├── subscriptions/
│   ├── billing/
│   ├── invoices/
│   ├── payments/
│   ├── coupons/
│   ├── taxes/
│   ├── usage/
│   ├── queue/
│   ├── workers/
│   ├── scheduler/
│   ├── audit/
│   ├── database/
│   ├── middleware/
│   ├── config/
│   └── observability/
├── migrations/
├── docs/
├── tests/
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── README.md
└── .env.example
```

---

## Testing

Add tests for:

* Auth
* RBAC
* Money calculations
* Tax calculations
* Coupon validation
* Proration
* Usage aggregation
* Invoice generation
* Payment simulation
* Idempotency
* Queue publishing
* Worker handling
* API endpoints

Include:

```txt
unit tests
integration tests
repository tests
worker tests
```

Tests must be runnable with:

```bash
go test ./...
```

---

## README Requirements

Create a complete README.

Include:

* Project description
* Features
* Architecture diagram
* Tech stack
* How billing calculations work
* How proration works
* How usage-based billing works
* How authentication works
* How authorization works
* How queues and workers work
* Local development setup
* Docker setup
* API documentation link
* Testing commands
* Environment variables
* Example API requests
* Future improvements

---

## Acceptance Criteria

The project is complete only if:

1. `docker compose up --build` starts all services.
2. API is available at `http://localhost:8080`.
3. PostgreSQL migrations run successfully.
4. Users can register and log in.
5. JWT auth works.
6. RBAC works.
7. Admin can create plans.
8. Customer can create a subscription.
9. Usage events can be submitted.
10. Invoices can be generated.
11. Billing totals are calculated correctly.
12. Taxes are calculated correctly.
13. Coupons are applied correctly.
14. Proration works for plan changes.
15. Payments can be simulated.
16. Workers consume queue jobs.
17. Failed jobs are retried.
18. Dead jobs are stored or visible.
19. Idempotency works for critical endpoints.
20. Tests pass with `go test ./...`.
21. OpenAPI documentation is available.
22. README explains how to run and use the service.

---

## Final Output Format

Return the complete project code.

For every file, include the path first:

```txt
/path/to/file
```

Then include the file content in a fenced code block.

Do not skip important files.

Do not provide pseudocode.

Do not leave TODO comments for core functionality.

The result must be ready to copy into a real project and run immediately.
