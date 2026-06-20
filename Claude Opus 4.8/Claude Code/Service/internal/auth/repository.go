package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a record does not exist.
var ErrNotFound = errors.New("not found")

// Repository persists users and refresh tokens.
type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

func (r *Repository) CreateUser(ctx context.Context, email, passwordHash, role string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, role, created_at, updated_at`,
		email, passwordHash, role,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return u, err
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	u := &User{}
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return u, err
}

// StoreRefreshToken persists a hashed refresh token tied to a jti.
func (r *Repository) StoreRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)`, userID, tokenHash, expiresAt)
	return err
}

// RefreshTokenRecord describes a stored refresh token.
type RefreshTokenRecord struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt *time.Time
}

func (r *Repository) GetRefreshToken(ctx context.Context, tokenHash string) (*RefreshTokenRecord, error) {
	rec := &RefreshTokenRecord{}
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, expires_at, revoked_at
		FROM refresh_tokens WHERE token_hash = $1`, tokenHash,
	).Scan(&rec.ID, &rec.UserID, &rec.ExpiresAt, &rec.RevokedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return rec, err
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE refresh_tokens SET revoked_at = now()
		WHERE token_hash = $1 AND revoked_at IS NULL`, tokenHash)
	return err
}

func (r *Repository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE refresh_tokens SET revoked_at = now()
		WHERE user_id = $1 AND revoked_at IS NULL`, userID)
	return err
}
