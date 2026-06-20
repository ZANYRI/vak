-- Billing service schema. All money is stored as int64 minor units (cents).
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ── Auth ─────────────────────────────────────────────────────────────────────
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role          TEXT NOT NULL DEFAULT 'customer'
                  CHECK (role IN ('admin','billing_manager','support','customer')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_users_email ON users (email);

CREATE TABLE refresh_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);

-- ── Customers ────────────────────────────────────────────────────────────────
CREATE TABLE customers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    email           TEXT NOT NULL,
    name            TEXT NOT NULL,
    company_name    TEXT,
    billing_address TEXT,
    country         TEXT NOT NULL DEFAULT '',
    region          TEXT NOT NULL DEFAULT '',
    tax_id          TEXT,
    currency        TEXT NOT NULL DEFAULT 'USD',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_customers_user_id ON customers (user_id);

-- ── Plans ────────────────────────────────────────────────────────────────────
CREATE TABLE plans (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name             TEXT NOT NULL,
    description      TEXT NOT NULL DEFAULT '',
    currency         TEXT NOT NULL DEFAULT 'USD',
    billing_interval TEXT NOT NULL CHECK (billing_interval IN ('monthly','yearly')),
    pricing_model    TEXT NOT NULL DEFAULT 'flat'
                     CHECK (pricing_model IN ('flat','per_seat','usage_based','tiered','hybrid')),
    base_price_cents BIGINT NOT NULL DEFAULT 0,
    seat_price_cents BIGINT NOT NULL DEFAULT 0,
    included_seats   BIGINT NOT NULL DEFAULT 0,
    included_units   BIGINT NOT NULL DEFAULT 0,
    unit_price_cents BIGINT NOT NULL DEFAULT 0,
    usage_metric     TEXT NOT NULL DEFAULT '',
    trial_days       INT NOT NULL DEFAULT 0,
    is_active        BOOLEAN NOT NULL DEFAULT true,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE pricing_tiers (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id          UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    from_unit        BIGINT NOT NULL,
    to_unit          BIGINT,            -- NULL = unbounded
    unit_price_cents BIGINT NOT NULL,
    sort_order       INT NOT NULL DEFAULT 0
);
CREATE INDEX idx_pricing_tiers_plan_id ON pricing_tiers (plan_id);

-- ── Subscriptions ────────────────────────────────────────────────────────────
CREATE TABLE subscriptions (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id          UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    plan_id              UUID NOT NULL REFERENCES plans(id),
    status               TEXT NOT NULL DEFAULT 'active'
                         CHECK (status IN ('trialing','active','past_due','paused','cancelled','expired')),
    quantity             BIGINT NOT NULL DEFAULT 1,
    current_period_start TIMESTAMPTZ NOT NULL DEFAULT now(),
    current_period_end   TIMESTAMPTZ NOT NULL,
    trial_start          TIMESTAMPTZ,
    trial_end            TIMESTAMPTZ,
    cancel_at_period_end BOOLEAN NOT NULL DEFAULT false,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_subscriptions_customer_id ON subscriptions (customer_id);
CREATE INDEX idx_subscriptions_status ON subscriptions (status);

-- ── Usage ────────────────────────────────────────────────────────────────────
CREATE TABLE usage_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id     UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    metric          TEXT NOT NULL,
    quantity        BIGINT NOT NULL,
    idempotency_key TEXT NOT NULL UNIQUE,
    recorded_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_usage_events_subscription_id ON usage_events (subscription_id);
CREATE INDEX idx_usage_events_idempotency_key ON usage_events (idempotency_key);

CREATE TABLE usage_summaries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    metric          TEXT NOT NULL,
    period_start    TIMESTAMPTZ NOT NULL,
    period_end      TIMESTAMPTZ NOT NULL,
    total_quantity  BIGINT NOT NULL DEFAULT 0,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (subscription_id, metric, period_start, period_end)
);

-- ── Coupons ──────────────────────────────────────────────────────────────────
CREATE TABLE coupons (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code             TEXT NOT NULL UNIQUE,
    type             TEXT NOT NULL CHECK (type IN ('percentage','fixed_amount')),
    percent_off      INT,
    amount_off_cents BIGINT,
    currency         TEXT,
    max_redemptions  INT,
    times_redeemed   INT NOT NULL DEFAULT 0,
    valid_from       TIMESTAMPTZ,
    valid_until      TIMESTAMPTZ,
    is_active        BOOLEAN NOT NULL DEFAULT true,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_coupons_code ON coupons (code);

CREATE TABLE subscription_coupons (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    coupon_id       UUID NOT NULL REFERENCES coupons(id) ON DELETE CASCADE,
    applied_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (subscription_id, coupon_id)
);

-- ── Taxes ────────────────────────────────────────────────────────────────────
CREATE TABLE tax_rules (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country               TEXT NOT NULL,
    region                TEXT NOT NULL DEFAULT '',
    tax_name              TEXT NOT NULL,
    tax_rate_basis_points INT NOT NULL,
    is_active             BOOLEAN NOT NULL DEFAULT true,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_tax_rules_country_region ON tax_rules (country, region);

-- ── Invoices ─────────────────────────────────────────────────────────────────
CREATE TABLE invoices (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id       UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    subscription_id   UUID REFERENCES subscriptions(id) ON DELETE SET NULL,
    status            TEXT NOT NULL DEFAULT 'draft'
                      CHECK (status IN ('draft','open','paid','void','uncollectible')),
    currency          TEXT NOT NULL DEFAULT 'USD',
    subtotal_cents    BIGINT NOT NULL DEFAULT 0,
    discount_cents    BIGINT NOT NULL DEFAULT 0,
    tax_cents         BIGINT NOT NULL DEFAULT 0,
    total_cents       BIGINT NOT NULL DEFAULT 0,
    amount_due_cents  BIGINT NOT NULL DEFAULT 0,
    amount_paid_cents BIGINT NOT NULL DEFAULT 0,
    period_start      TIMESTAMPTZ,
    period_end        TIMESTAMPTZ,
    issued_at         TIMESTAMPTZ,
    due_at            TIMESTAMPTZ,
    paid_at           TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_invoices_customer_id ON invoices (customer_id);
CREATE INDEX idx_invoices_status ON invoices (status);

CREATE TABLE invoice_lines (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id        UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    type              TEXT NOT NULL,
    description       TEXT NOT NULL,
    quantity          BIGINT NOT NULL DEFAULT 1,
    unit_amount_cents BIGINT NOT NULL DEFAULT 0,
    amount_cents      BIGINT NOT NULL DEFAULT 0,
    metadata          JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_invoice_lines_invoice_id ON invoice_lines (invoice_id);

-- ── Payments ─────────────────────────────────────────────────────────────────
CREATE TABLE payments (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id    UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    amount_cents  BIGINT NOT NULL,
    currency      TEXT NOT NULL,
    status        TEXT NOT NULL DEFAULT 'pending'
                  CHECK (status IN ('pending','succeeded','failed','refunded')),
    card_last4    TEXT NOT NULL DEFAULT '',
    failure_reason TEXT NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_payments_invoice_id ON payments (invoice_id);

-- ── Jobs (queue observability / lifecycle) ───────────────────────────────────
CREATE TABLE jobs (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type         TEXT NOT NULL,
    payload      JSONB NOT NULL DEFAULT '{}'::jsonb,
    status       TEXT NOT NULL DEFAULT 'queued'
                 CHECK (status IN ('queued','running','completed','failed','retrying','dead')),
    attempts     INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 5,
    last_error   TEXT NOT NULL DEFAULT '',
    idempotency_key TEXT,
    available_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_jobs_status ON jobs (status);
CREATE UNIQUE INDEX idx_jobs_idempotency_key ON jobs (idempotency_key) WHERE idempotency_key IS NOT NULL;

-- ── Idempotency (HTTP-level) ─────────────────────────────────────────────────
CREATE TABLE idempotency_keys (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key           TEXT NOT NULL,
    endpoint      TEXT NOT NULL,
    request_hash  TEXT NOT NULL,
    status_code   INT NOT NULL,
    response_body TEXT NOT NULL,  -- raw bytes, so replays are byte-identical
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (key, endpoint)
);

-- ── Audit logs ───────────────────────────────────────────────────────────────
CREATE TABLE audit_logs (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_user_id UUID,
    action        TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id   TEXT NOT NULL DEFAULT '',
    metadata      JSONB NOT NULL DEFAULT '{}'::jsonb,
    ip_address    TEXT NOT NULL DEFAULT '',
    user_agent    TEXT NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_audit_logs_actor ON audit_logs (actor_user_id);
CREATE INDEX idx_audit_logs_resource ON audit_logs (resource_type, resource_id);
