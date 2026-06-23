package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/example/billing-service/internal/auth"
	"github.com/example/billing-service/internal/config"
	"github.com/example/billing-service/internal/queue"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	DB     *pgxpool.Pool
	Queue  *queue.Client
	Tokens auth.TokenService
	Config config.Config
	Log    *slog.Logger
	metric *metrics
	rate   *rateLimiter
}
type metrics struct{ requests, requestDurationNanos, jobsProcessed, jobsFailed, invoices, paymentsOK, paymentsFailed atomic.Int64 }
type ctxKey string

const claimsKey ctxKey = "claims"

type principal struct{ ID, Role string }

func New(db *pgxpool.Pool, q *queue.Client, tokens auth.TokenService, cfg config.Config, log *slog.Logger) *App {
	return &App{DB: db, Queue: q, Tokens: tokens, Config: cfg, Log: log, metric: &metrics{}, rate: newRateLimiter()}
}

func (a *App) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(a.requestContext, a.securityHeaders, a.cors, a.recover)
	r.Get("/healthz", a.health)
	r.Get("/readyz", a.ready)
	r.Get("/metrics", a.metrics)
	r.Get("/api/docs", a.docs)
	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Route("/auth", func(ar chi.Router) {
			ar.With(a.authRateLimit).Post("/register", a.register)
			ar.With(a.authRateLimit).Post("/login", a.login)
			ar.Post("/refresh", a.refresh)
			ar.Post("/logout", a.logout)
			ar.With(a.authenticate).Get("/me", a.me)
		})
		v1.Group(func(pr chi.Router) {
			pr.Use(a.authenticate)
			pr.Get("/plans", a.listPlans)
			pr.Get("/plans/{id}", a.getPlan)
			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager)).Post("/plans", a.createPlan)
			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager)).Patch("/plans/{id}", a.updatePlan)
			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager)).Delete("/plans/{id}", a.deletePlan)

			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager, auth.RoleSupport, auth.RoleCustomer)).Get("/customers", a.listCustomers)
			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager)).Post("/customers", a.createCustomer)
			pr.Get("/customers/{id}", a.getCustomer)
			pr.Patch("/customers/{id}", a.updateCustomer)
			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager)).Delete("/customers/{id}", a.deleteCustomer)

			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager, auth.RoleCustomer), a.idempotent("subscription.create")).Post("/subscriptions", a.createSubscription)
			pr.Get("/subscriptions", a.listSubscriptions)
			pr.Get("/subscriptions/{id}", a.getSubscription)
			pr.Patch("/subscriptions/{id}", a.updateSubscription)
			pr.Post("/subscriptions/{id}/cancel", a.cancelSubscription)
			pr.Post("/subscriptions/{id}/pause", a.pauseSubscription)
			pr.Post("/subscriptions/{id}/resume", a.resumeSubscription)
			pr.Post("/subscriptions/{id}/change-plan", a.changePlan)
			pr.Post("/subscriptions/{id}/apply-coupon", a.applyCoupon)

			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager, auth.RoleCustomer), a.idempotent("usage.record")).Post("/usage", a.recordUsage)
			pr.Get("/usage", a.listUsage)
			pr.Get("/usage/summary", a.usageSummary)

			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager)).Post("/coupons", a.createCoupon)
			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager)).Get("/coupons", a.listCoupons)
			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager)).Get("/coupons/{id}", a.getCoupon)
			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager)).Patch("/coupons/{id}", a.updateCoupon)
			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager)).Delete("/coupons/{id}", a.deleteCoupon)
			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager)).Post("/tax-rules", a.createTaxRule)
			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager)).Get("/tax-rules", a.listTaxRules)
			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager)).Patch("/tax-rules/{id}", a.updateTaxRule)
			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager)).Delete("/tax-rules/{id}", a.deleteTaxRule)

			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager, auth.RoleCustomer), a.idempotent("invoice.generate")).Post("/invoices/generate", a.generateInvoice)
			pr.Get("/invoices", a.listInvoices)
			pr.Get("/invoices/{id}", a.getInvoice)
			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager)).Post("/invoices/{id}/finalize", a.finalizeInvoice)
			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager)).Post("/invoices/{id}/void", a.voidInvoice)
			pr.With(a.idempotent("payment.simulate")).Post("/invoices/{id}/pay", a.payInvoice)
			pr.With(a.require(auth.RoleAdmin, auth.RoleBillingManager, auth.RoleCustomer), a.idempotent("payment.simulate")).Post("/payments/simulate", a.simulatePayment)
			pr.Get("/payments", a.listPayments)
			pr.Get("/payments/{id}", a.getPayment)
		})
	})
	return r
}

func (a *App) requestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		w.Header().Set("X-Request-ID", id)
		start := time.Now()
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey("request_id"), id)))
		a.metric.requests.Add(1)
		duration := time.Since(start)
		a.metric.requestDurationNanos.Add(duration.Nanoseconds())
		a.Log.Info("request", "method", r.Method, "path", r.URL.Path, "duration", duration.String(), "request_id", id)
	})
}
func (a *App) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}
func (a *App) cors(next http.Handler) http.Handler {
	allowed := strings.Split(a.Config.CORSOrigins, ",")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o := r.Header.Get("Origin")
		for _, v := range allowed {
			if strings.TrimSpace(v) == o {
				w.Header().Set("Access-Control-Allow-Origin", o)
				w.Header().Set("Vary", "Origin")
				break
			}
		}
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Idempotency-Key, X-Request-ID")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func (a *App) recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				a.Log.Error("panic recovered", "panic", x)
				fail(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func (a *App) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			fail(w, 401, "UNAUTHORIZED", "Missing bearer token", nil)
			return
		}
		c, e := a.Tokens.ParseAccess(strings.TrimPrefix(h, "Bearer "))
		if e != nil {
			fail(w, 401, "UNAUTHORIZED", "Invalid or expired access token", nil)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), claimsKey, principal{c.Subject, c.Role})))
	})
}
func (a *App) require(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := user(r)
			if !auth.Allowed(p.Role, roles...) {
				fail(w, 403, "INSUFFICIENT_PERMISSIONS", "Insufficient permissions", nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
func user(r *http.Request) principal { p, _ := r.Context().Value(claimsKey).(principal); return p }

type recorder struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (r *recorder) WriteHeader(s int) { r.status = s; r.ResponseWriter.WriteHeader(s) }
func (r *recorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return r.ResponseWriter.Write(b)
}
func (a *App) idempotent(scope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := strings.TrimSpace(r.Header.Get("Idempotency-Key"))
			if key == "" {
				fail(w, 400, "VALIDATION_ERROR", "Idempotency-Key header is required", nil)
				return
			}
			if len(key) > 255 {
				fail(w, 400, "VALIDATION_ERROR", "Idempotency-Key is too long", nil)
				return
			}
			body, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
			if e != nil {
				fail(w, 400, "VALIDATION_ERROR", "Request body is too large", nil)
				return
			}
			r.Body.Close()
			r.Body = io.NopCloser(strings.NewReader(string(body)))
			h := sha256.Sum256(body)
			hash := hex.EncodeToString(h[:])
			p := user(r)
			var oldHash string
			var status int
			var oldBody []byte
			e = a.DB.QueryRow(r.Context(), `SELECT request_hash,response_status,response_body::text FROM idempotency_records WHERE scope=$1 AND key=$2 AND actor_user_id=$3`, scope, key, p.ID).Scan(&oldHash, &status, &oldBody)
			if e == nil {
				if oldHash != hash {
					fail(w, 409, "IDEMPOTENCY_CONFLICT", "Key was used for a different request", nil)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				_, _ = w.Write(oldBody)
				return
			}
			if !errors.Is(e, pgx.ErrNoRows) {
				fail(w, 500, "INTERNAL_ERROR", "Unable to check idempotency", nil)
				return
			}
			rec := &recorder{ResponseWriter: w, status: 200}
			next.ServeHTTP(rec, r)
			if rec.status >= 200 && rec.status < 300 {
				var bodyJSON any
				if json.Unmarshal(rec.body, &bodyJSON) != nil {
					bodyJSON = map[string]string{"raw": string(rec.body)}
				}
				_, e = a.DB.Exec(r.Context(), `INSERT INTO idempotency_records(scope,key,actor_user_id,request_hash,response_status,response_body) VALUES($1,$2,$3,$4,$5,$6) ON CONFLICT DO NOTHING`, scope, key, p.ID, hash, rec.status, bodyJSON)
				if e != nil {
					a.Log.Error("store idempotency", "error", e)
				}
			}
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func fail(w http.ResponseWriter, status int, code, msg string, details any) {
	writeJSON(w, status, map[string]any{"error": map[string]any{"code": code, "message": msg, "details": details}})
}
func decode(r *http.Request, v any) error {
	d := json.NewDecoder(io.LimitReader(r.Body, 1<<20+1))
	d.DisallowUnknownFields()
	if e := d.Decode(v); e != nil {
		return e
	}
	if d.Decode(&struct{}{}) != io.EOF {
		return errors.New("request must contain one JSON object")
	}
	return nil
}
func id(r *http.Request) string { return chi.URLParam(r, "id") }
func clientIP(r *http.Request) string {
	host, _, e := net.SplitHostPort(r.RemoteAddr)
	if e == nil {
		return host
	}
	return r.RemoteAddr
}

type rateLimiter struct {
	mu   sync.Mutex
	hits map[string][]time.Time
}

func newRateLimiter() *rateLimiter { return &rateLimiter{hits: map[string][]time.Time{}} }
func (a *App) authRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := clientIP(r)
		now := time.Now()
		a.rate.mu.Lock()
		xs := a.rate.hits[key]
		cut := now.Add(-time.Minute)
		n := xs[:0]
		for _, x := range xs {
			if x.After(cut) {
				n = append(n, x)
			}
		}
		if len(n) >= 10 {
			a.rate.hits[key] = n
			a.rate.mu.Unlock()
			fail(w, 429, "RATE_LIMITED", "Too many authentication requests", nil)
			return
		}
		a.rate.hits[key] = append(n, now)
		a.rate.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func (a *App) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok"})
}
func (a *App) ready(w http.ResponseWriter, r *http.Request) {
	if e := a.DB.Ping(r.Context()); e != nil {
		fail(w, 503, "NOT_READY", "Database is unavailable", nil)
		return
	}
	writeJSON(w, 200, map[string]string{"status": "ready"})
}
func (a *App) metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	var depth, processed, failed int64
	if err := a.DB.QueryRow(r.Context(), `SELECT count(*) FILTER (WHERE status IN ('queued','retrying','running')),count(*) FILTER (WHERE status='completed'),count(*) FILTER (WHERE status IN ('failed','dead')) FROM jobs`).Scan(&depth, &processed, &failed); err != nil {
		depth, processed, failed = -1, -1, -1
	}
	durationSeconds := float64(a.metric.requestDurationNanos.Load()) / float64(time.Second)
	_, _ = w.Write([]byte("http_requests_total " + itoa(a.metric.requests.Load()) + "\nhttp_request_duration_seconds_sum " + fmt.Sprintf("%.9f", durationSeconds) + "\nhttp_request_duration_seconds_count " + itoa(a.metric.requests.Load()) + "\njobs_processed_total " + itoa(processed) + "\njobs_failed_total " + itoa(failed) + "\ninvoices_generated_total " + itoa(a.metric.invoices.Load()) + "\npayments_succeeded_total " + itoa(a.metric.paymentsOK.Load()) + "\npayments_failed_total " + itoa(a.metric.paymentsFailed.Load()) + "\nqueue_depth " + itoa(depth) + "\n"))
}
func itoa(v int64) string { return fmt.Sprintf("%d", v) }
