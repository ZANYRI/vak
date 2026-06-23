package api

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/example/billing-service/internal/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type registerInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
type loginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type refreshInput struct {
	RefreshToken string `json:"refresh_token"`
}

func (a *App) register(w http.ResponseWriter, r *http.Request) {
	var in registerInput
	if e := decode(r, &in); e != nil || !validEmail(in.Email) || strings.TrimSpace(in.Name) == "" {
		fail(w, 400, "VALIDATION_ERROR", "Invalid registration payload", map[string]string{"field": "email, password or name"})
		return
	}
	hash, e := auth.HashPassword(in.Password)
	if e != nil {
		fail(w, 400, "VALIDATION_ERROR", e.Error(), map[string]string{"field": "password"})
		return
	}
	tx, e := a.DB.Begin(r.Context())
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to register user", nil)
		return
	}
	defer tx.Rollback(r.Context())
	uid := uuid.New()
	_, e = tx.Exec(r.Context(), `INSERT INTO users(id,email,password_hash,role) VALUES($1,$2,$3,'customer')`, uid, strings.ToLower(strings.TrimSpace(in.Email)), hash)
	if e != nil {
		if strings.Contains(e.Error(), "unique") {
			fail(w, 409, "CONFLICT", "Email is already registered", nil)
		} else {
			fail(w, 500, "INTERNAL_ERROR", "Unable to register user", nil)
		}
		return
	}
	// Every self-registered account owns a customer record, which lets it manage only its own subscriptions and invoices.
	_, e = tx.Exec(r.Context(), `INSERT INTO customers(id,user_id,email,name,currency) VALUES($1,$2,$3,$4,'USD')`, uuid.New(), uid, strings.ToLower(strings.TrimSpace(in.Email)), strings.TrimSpace(in.Name))
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to create customer", nil)
		return
	}
	if e = tx.Commit(r.Context()); e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to register user", nil)
		return
	}
	a.audit(r, "user.registered", "user", uid.String(), map[string]string{"email": strings.ToLower(strings.TrimSpace(in.Email))})
	a.issueTokens(w, r, uid.String(), auth.RoleCustomer, 201)
}
func (a *App) login(w http.ResponseWriter, r *http.Request) {
	var in loginInput
	if e := decode(r, &in); e != nil {
		fail(w, 400, "VALIDATION_ERROR", "Invalid request payload", nil)
		return
	}
	var uid, hash, role string
	e := a.DB.QueryRow(r.Context(), `SELECT id::text,password_hash,role FROM users WHERE email=$1`, strings.TrimSpace(in.Email)).Scan(&uid, &hash, &role)
	if e != nil || auth.ComparePassword(hash, in.Password) != nil {
		fail(w, 401, "UNAUTHORIZED", "Invalid email or password", nil)
		return
	}
	a.auditWithUser(r, uid, "user.login", "user", uid, nil)
	a.issueTokens(w, r, uid, role, 200)
}
func (a *App) issueTokens(w http.ResponseWriter, r *http.Request, uid, role string, status int) {
	access, e := a.Tokens.AccessToken(uid, role)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to create access token", nil)
		return
	}
	plain, hash, expiry, e := a.Tokens.NewRefreshToken()
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to create refresh token", nil)
		return
	}
	_, e = a.DB.Exec(r.Context(), `INSERT INTO refresh_tokens(id,user_id,token_hash,expires_at) VALUES($1,$2,$3,$4)`, uuid.New(), uid, hash, expiry)
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to store refresh token", nil)
		return
	}
	writeJSON(w, status, map[string]any{"access_token": access, "refresh_token": plain, "token_type": "Bearer", "expires_in": int(a.Config.AccessTokenTTL.Seconds()), "role": role})
}
func (a *App) refresh(w http.ResponseWriter, r *http.Request) {
	var in refreshInput
	if e := decode(r, &in); e != nil || in.RefreshToken == "" {
		fail(w, 400, "VALIDATION_ERROR", "refresh_token is required", nil)
		return
	}
	tx, e := a.DB.Begin(r.Context())
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to refresh token", nil)
		return
	}
	defer tx.Rollback(r.Context())
	var uid, role string
	e = tx.QueryRow(r.Context(), `SELECT u.id::text,u.role FROM refresh_tokens t JOIN users u ON u.id=t.user_id WHERE t.token_hash=$1 AND t.revoked_at IS NULL AND t.expires_at > now() FOR UPDATE`, auth.HashRefresh(in.RefreshToken)).Scan(&uid, &role)
	if e != nil {
		fail(w, 401, "UNAUTHORIZED", "Invalid or expired refresh token", nil)
		return
	}
	_, e = tx.Exec(r.Context(), `UPDATE refresh_tokens SET revoked_at=now() WHERE token_hash=$1`, auth.HashRefresh(in.RefreshToken))
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to rotate refresh token", nil)
		return
	}
	if e = tx.Commit(r.Context()); e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to refresh token", nil)
		return
	}
	a.issueTokens(w, r, uid, role, 200)
}
func (a *App) logout(w http.ResponseWriter, r *http.Request) {
	var in refreshInput
	if e := decode(r, &in); e != nil || in.RefreshToken == "" {
		fail(w, 400, "VALIDATION_ERROR", "refresh_token is required", nil)
		return
	}
	_, e := a.DB.Exec(r.Context(), `UPDATE refresh_tokens SET revoked_at=now() WHERE token_hash=$1 AND revoked_at IS NULL`, auth.HashRefresh(in.RefreshToken))
	if e != nil {
		fail(w, 500, "INTERNAL_ERROR", "Unable to logout", nil)
		return
	}
	writeJSON(w, 200, map[string]bool{"logged_out": true})
}
func (a *App) me(w http.ResponseWriter, r *http.Request) {
	p := user(r)
	var email, role string
	var created time.Time
	if e := a.DB.QueryRow(r.Context(), `SELECT email::text,role,created_at FROM users WHERE id=$1`, p.ID).Scan(&email, &role, &created); e != nil {
		fail(w, 401, "UNAUTHORIZED", "Account no longer exists", nil)
		return
	}
	writeJSON(w, 200, map[string]any{"id": p.ID, "email": email, "role": role, "created_at": created})
}
func (a *App) BootstrapAdmin(ctx context.Context) error {
	if a.Config.BootstrapAdminEmail == "" || a.Config.BootstrapAdminPassword == "" {
		return nil
	}
	hash, e := auth.HashPassword(a.Config.BootstrapAdminPassword)
	if e != nil {
		return e
	}
	_, e = a.DB.Exec(ctx, `INSERT INTO users(id,email,password_hash,role) VALUES($1,$2,$3,'admin') ON CONFLICT(email) DO NOTHING`, uuid.New(), strings.TrimSpace(a.Config.BootstrapAdminEmail), hash)
	return e
}
func validEmail(v string) bool {
	v = strings.TrimSpace(v)
	return len(v) <= 254 && strings.Count(v, "@") == 1 && strings.Index(v, "@") > 0 && strings.LastIndex(v, ".") > strings.Index(v, "@")+1
}
func (a *App) audit(r *http.Request, action, resource, id string, metadata any) {
	p := user(r)
	a.auditWithUser(r, p.ID, action, resource, id, metadata)
}
func (a *App) auditWithUser(r *http.Request, actor, action, resource, id string, metadata any) {
	if actor == "" {
		return
	}
	_, e := a.DB.Exec(r.Context(), `INSERT INTO audit_logs(id,actor_user_id,action,resource_type,resource_id,metadata,ip_address,user_agent) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`, uuid.New(), actor, action, resource, id, metadata, clientIP(r), r.UserAgent())
	if e != nil && !errors.Is(e, pgx.ErrNoRows) {
		a.Log.Warn("audit log failed", "error", e)
	}
}
