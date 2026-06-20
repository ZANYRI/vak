package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/example/billing-service/internal/audit"
	"github.com/example/billing-service/internal/httpx"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service implements authentication use-cases.
type Service struct {
	repo       *Repository
	tokens     *TokenManager
	audit      audit.Recorder
	bcryptCost int
}

func NewService(repo *Repository, tokens *TokenManager, recorder audit.Recorder, bcryptCost int) *Service {
	return &Service{repo: repo, tokens: tokens, audit: recorder, bcryptCost: bcryptCost}
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// Register creates a new user. The first registration defaults to customer role
// unless an explicit role is supplied (admin-bootstrapping is done separately).
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	role := req.Role
	if role == "" {
		role = RoleCustomer
	}
	if !ValidRoles[role] {
		return nil, httpx.ErrValidation("invalid role")
	}

	if _, err := s.repo.GetUserByEmail(ctx, email); err == nil {
		return nil, httpx.ErrConflict("email already registered")
	} else if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.bcryptCost)
	if err != nil {
		return nil, err
	}
	return s.repo.CreateUser(ctx, email, string(hash), role)
}

// Login verifies credentials and issues a token pair.
func (s *Service) Login(ctx context.Context, req LoginRequest, ip, ua string) (*TokenPair, *User, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil, httpx.ErrUnauthorized("invalid credentials")
		}
		return nil, nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return nil, nil, httpx.ErrUnauthorized("invalid credentials")
	}

	pair, err := s.issueTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	actor := user.ID
	s.audit.Record(ctx, audit.Entry{
		ActorUserID: &actor, Action: audit.ActionUserLogin,
		ResourceType: "user", ResourceID: user.ID.String(),
		IPAddress: ip, UserAgent: ua,
	})
	return pair, user, nil
}

func (s *Service) issueTokens(ctx context.Context, user *User) (*TokenPair, error) {
	access, err := s.tokens.GenerateAccess(user.ID, user.Role)
	if err != nil {
		return nil, err
	}
	refresh, _, err := s.tokens.GenerateRefresh(user.ID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.StoreRefreshToken(ctx, user.ID, hashToken(refresh), time.Now().Add(s.tokens.RefreshTTL())); err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.tokens.AccessTTL().Seconds()),
	}, nil
}

// Refresh rotates a refresh token: the presented token is revoked and a new
// pair is issued (refresh token rotation).
func (s *Service) Refresh(ctx context.Context, req RefreshRequest) (*TokenPair, error) {
	userID, _, err := s.tokens.ParseRefresh(req.RefreshToken)
	if err != nil {
		return nil, httpx.ErrUnauthorized("invalid refresh token")
	}
	h := hashToken(req.RefreshToken)
	rec, err := s.repo.GetRefreshToken(ctx, h)
	if err != nil {
		return nil, httpx.ErrUnauthorized("refresh token not recognized")
	}
	if rec.RevokedAt != nil {
		// Token reuse after rotation: revoke everything for safety.
		_ = s.repo.RevokeAllForUser(ctx, userID)
		return nil, httpx.ErrUnauthorized("refresh token already used")
	}
	if time.Now().After(rec.ExpiresAt) {
		return nil, httpx.ErrUnauthorized("refresh token expired")
	}

	if err := s.repo.RevokeRefreshToken(ctx, h); err != nil {
		return nil, err
	}
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, httpx.ErrUnauthorized("user no longer exists")
	}
	return s.issueTokens(ctx, user)
}

// Logout revokes the presented refresh token.
func (s *Service) Logout(ctx context.Context, req RefreshRequest) error {
	return s.repo.RevokeRefreshToken(ctx, hashToken(req.RefreshToken))
}

// Me returns the current user by id.
func (s *Service) Me(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, httpx.ErrNotFound("user not found")
		}
		return nil, err
	}
	return user, nil
}

// EnsureAdmin bootstraps an admin account if it does not yet exist.
func (s *Service) EnsureAdmin(ctx context.Context, email, password string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || password == "" {
		return nil
	}
	if _, err := s.repo.GetUserByEmail(ctx, email); err == nil {
		return nil
	} else if !errors.Is(err, ErrNotFound) {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	if err != nil {
		return err
	}
	if _, err = s.repo.CreateUser(ctx, email, string(hash), RoleAdmin); err != nil {
		// Benign race: another process (worker/scheduler) created it first.
		if isUniqueViolation(err) {
			return nil
		}
		return err
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) {
		return pgErr.SQLState() == "23505"
	}
	return false
}
