CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY, email citext NOT NULL UNIQUE, password_hash text NOT NULL,
  role text NOT NULL CHECK (role IN ('admin','billing_manager','customer','support')),
  created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS users_email_idx ON users(email);
CREATE TABLE IF NOT EXISTS refresh_tokens (
  id uuid PRIMARY KEY, user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash text NOT NULL UNIQUE, expires_at timestamptz NOT NULL, revoked_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS customers (
  id uuid PRIMARY KEY, user_id uuid UNIQUE REFERENCES users(id) ON DELETE SET NULL,
  email citext NOT NULL, name text NOT NULL, company_name text NOT NULL DEFAULT '', billing_address jsonb NOT NULL DEFAULT '{}',
  country text NOT NULL DEFAULT '', tax_id text NOT NULL DEFAULT '', currency char(3) NOT NULL DEFAULT 'USD',
  created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS customers_user_id_idx ON customers(user_id);
CREATE TABLE IF NOT EXISTS plans (
  id uuid PRIMARY KEY, name text NOT NULL, description text NOT NULL DEFAULT '', currency char(3) NOT NULL,
  billing_interval text NOT NULL CHECK (billing_interval IN ('monthly','yearly')), pricing_model text NOT NULL CHECK (pricing_model IN ('flat','per_seat','usage_based','tiered','hybrid')),
  base_price_cents bigint NOT NULL CHECK(base_price_cents >= 0), trial_days integer NOT NULL DEFAULT 0 CHECK(trial_days >= 0),
  pricing jsonb NOT NULL DEFAULT '{}', is_active boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS plan_prices (id uuid PRIMARY KEY, plan_id uuid NOT NULL REFERENCES plans(id) ON DELETE CASCADE, currency char(3) NOT NULL, amount_cents bigint NOT NULL, created_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE IF NOT EXISTS pricing_tiers (id uuid PRIMARY KEY, plan_id uuid NOT NULL REFERENCES plans(id) ON DELETE CASCADE, from_units bigint NOT NULL, to_units bigint, unit_price_cents bigint NOT NULL, created_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE IF NOT EXISTS subscriptions (
  id uuid PRIMARY KEY, customer_id uuid NOT NULL REFERENCES customers(id), plan_id uuid NOT NULL REFERENCES plans(id),
  status text NOT NULL CHECK(status IN ('trialing','active','past_due','paused','cancelled','expired')), quantity bigint NOT NULL DEFAULT 1 CHECK(quantity > 0),
  current_period_start timestamptz NOT NULL, current_period_end timestamptz NOT NULL, trial_start timestamptz, trial_end timestamptz,
  cancel_at_period_end boolean NOT NULL DEFAULT false, created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS subscriptions_customer_id_idx ON subscriptions(customer_id);
CREATE INDEX IF NOT EXISTS subscriptions_status_idx ON subscriptions(status);
CREATE TABLE IF NOT EXISTS usage_events (
  id uuid PRIMARY KEY, customer_id uuid NOT NULL REFERENCES customers(id), subscription_id uuid NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
  metric text NOT NULL, quantity bigint NOT NULL CHECK(quantity >= 0), idempotency_key text NOT NULL, recorded_at timestamptz NOT NULL, created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE(subscription_id, idempotency_key)
);
CREATE INDEX IF NOT EXISTS usage_events_subscription_id_idx ON usage_events(subscription_id);
CREATE INDEX IF NOT EXISTS usage_events_idempotency_key_idx ON usage_events(idempotency_key);
CREATE TABLE IF NOT EXISTS usage_summaries (id uuid PRIMARY KEY, subscription_id uuid NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE, metric text NOT NULL, period_start timestamptz NOT NULL, period_end timestamptz NOT NULL, quantity bigint NOT NULL, updated_at timestamptz NOT NULL DEFAULT now(), UNIQUE(subscription_id,metric,period_start,period_end));
CREATE TABLE IF NOT EXISTS coupons (
  id uuid PRIMARY KEY, code citext NOT NULL UNIQUE, type text NOT NULL CHECK(type IN ('percentage','fixed_amount')),
  percent_off integer, amount_off_cents bigint, currency char(3), max_redemptions integer, times_redeemed integer NOT NULL DEFAULT 0,
  valid_from timestamptz, valid_until timestamptz, is_active boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now(),
  CHECK ((type='percentage' AND percent_off BETWEEN 1 AND 100 AND amount_off_cents IS NULL) OR (type='fixed_amount' AND amount_off_cents > 0 AND percent_off IS NULL))
);
CREATE INDEX IF NOT EXISTS coupons_code_idx ON coupons(code);
CREATE TABLE IF NOT EXISTS subscription_coupons (subscription_id uuid NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE, coupon_id uuid NOT NULL REFERENCES coupons(id), created_at timestamptz NOT NULL DEFAULT now(), PRIMARY KEY(subscription_id,coupon_id));
CREATE TABLE IF NOT EXISTS tax_rules (id uuid PRIMARY KEY, country text NOT NULL, region text NOT NULL DEFAULT '', tax_name text NOT NULL, tax_rate_basis_points integer NOT NULL CHECK(tax_rate_basis_points >= 0), is_active boolean NOT NULL DEFAULT true, created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE IF NOT EXISTS invoices (
  id uuid PRIMARY KEY, customer_id uuid NOT NULL REFERENCES customers(id), subscription_id uuid REFERENCES subscriptions(id), status text NOT NULL CHECK(status IN ('draft','open','paid','void','uncollectible')),
  currency char(3) NOT NULL, subtotal_cents bigint NOT NULL, discount_cents bigint NOT NULL, tax_cents bigint NOT NULL, total_cents bigint NOT NULL, amount_due_cents bigint NOT NULL, amount_paid_cents bigint NOT NULL DEFAULT 0,
  period_start timestamptz NOT NULL, period_end timestamptz NOT NULL, issued_at timestamptz, due_at timestamptz, paid_at timestamptz, created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS invoices_customer_id_idx ON invoices(customer_id);
CREATE INDEX IF NOT EXISTS invoices_status_idx ON invoices(status);
CREATE TABLE IF NOT EXISTS invoice_lines (id uuid PRIMARY KEY, invoice_id uuid NOT NULL REFERENCES invoices(id) ON DELETE CASCADE, type text NOT NULL, description text NOT NULL, quantity bigint NOT NULL DEFAULT 1, unit_amount_cents bigint NOT NULL, amount_cents bigint NOT NULL, metadata jsonb NOT NULL DEFAULT '{}', created_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE IF NOT EXISTS payments (id uuid PRIMARY KEY, invoice_id uuid NOT NULL REFERENCES invoices(id), status text NOT NULL CHECK(status IN ('pending','succeeded','failed','refunded')), amount_cents bigint NOT NULL, card_last4 char(4), failure_reason text, created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE IF NOT EXISTS jobs (id uuid PRIMARY KEY, type text NOT NULL, payload jsonb NOT NULL DEFAULT '{}', status text NOT NULL CHECK(status IN ('queued','running','completed','failed','retrying','dead')), attempts integer NOT NULL DEFAULT 0, max_attempts integer NOT NULL DEFAULT 5, available_at timestamptz NOT NULL DEFAULT now(), last_error text, idempotency_key text UNIQUE, created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now());
CREATE INDEX IF NOT EXISTS jobs_status_idx ON jobs(status);
CREATE TABLE IF NOT EXISTS audit_logs (id uuid PRIMARY KEY, actor_user_id uuid REFERENCES users(id) ON DELETE SET NULL, action text NOT NULL, resource_type text NOT NULL, resource_id uuid, metadata jsonb NOT NULL DEFAULT '{}', ip_address inet, user_agent text, created_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE IF NOT EXISTS idempotency_records (scope text NOT NULL, key text NOT NULL, actor_user_id uuid, request_hash text NOT NULL, response_status integer NOT NULL, response_body jsonb NOT NULL, created_at timestamptz NOT NULL DEFAULT now(), PRIMARY KEY(scope,key,actor_user_id));
