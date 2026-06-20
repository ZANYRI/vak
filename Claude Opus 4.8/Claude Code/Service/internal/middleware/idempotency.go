package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/example/billing-service/internal/httpx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Idempotency provides replay protection for unsafe POST endpoints via the
// Idempotency-Key header. Same key + same body => cached response is replayed.
// Same key + different body => IDEMPOTENCY_CONFLICT.
type Idempotency struct {
	db *pgxpool.Pool
}

func NewIdempotency(db *pgxpool.Pool) *Idempotency { return &Idempotency{db: db} }

type captureWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (c *captureWriter) WriteHeader(code int) {
	c.status = code
	c.ResponseWriter.WriteHeader(code)
}

func (c *captureWriter) Write(b []byte) (int, error) {
	c.body.Write(b)
	return c.ResponseWriter.Write(b)
}

// Middleware enforces idempotency when the header is present; otherwise passes through.
func (i *Idempotency) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			next.ServeHTTP(w, r)
			return
		}

		bodyBytes, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		hash := hashBytes(bodyBytes)
		endpoint := r.Method + " " + r.URL.Path

		// Look for a previous response.
		var storedHash string
		var statusCode int
		var respBody string
		err := i.db.QueryRow(r.Context(),
			`SELECT request_hash, status_code, response_body FROM idempotency_keys
			 WHERE key = $1 AND endpoint = $2`, key, endpoint,
		).Scan(&storedHash, &statusCode, &respBody)
		if err == nil {
			if storedHash != hash {
				httpx.Error(w, r, httpx.ErrIdempotencyConflict("idempotency key reused with a different request body"))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Idempotent-Replayed", "true")
			w.WriteHeader(statusCode)
			_, _ = w.Write([]byte(respBody))
			return
		}
		if err != pgx.ErrNoRows {
			httpx.Error(w, r, httpx.ErrInternal(""))
			return
		}

		// First time: capture the response and persist it.
		cw := &captureWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(cw, r)

		// Only cache successful (2xx) responses.
		if cw.status >= 200 && cw.status < 300 {
			i.store(r.Context(), key, endpoint, hash, cw.status, cw.body.Bytes())
		}
	})
}

func (i *Idempotency) store(ctx context.Context, key, endpoint, hash string, status int, body []byte) {
	// Stored verbatim in a TEXT column so replays are byte-identical.
	_, _ = i.db.Exec(ctx, `
		INSERT INTO idempotency_keys (key, endpoint, request_hash, status_code, response_body)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (key, endpoint) DO NOTHING`,
		key, endpoint, hash, status, string(body))
}

func hashBytes(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
