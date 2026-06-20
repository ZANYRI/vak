package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims is the JWT payload for access tokens.
type Claims struct {
	UserID uuid.UUID `json:"uid"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// TokenManager issues and verifies access/refresh JWTs.
type TokenManager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewTokenManager(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *TokenManager {
	return &TokenManager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (m *TokenManager) AccessTTL() time.Duration  { return m.accessTTL }
func (m *TokenManager) RefreshTTL() time.Duration { return m.refreshTTL }

// GenerateAccess returns a signed access token.
func (m *TokenManager) GenerateAccess(userID uuid.UUID, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.accessSecret)
}

// GenerateRefresh returns a signed opaque refresh token (a random jti embedded).
func (m *TokenManager) GenerateRefresh(userID uuid.UUID) (string, string, error) {
	jti := uuid.NewString()
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ID:        jti,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTTL)),
	}
	tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.refreshSecret)
	return tok, jti, err
}

// ParseAccess validates an access token and returns its claims.
func (m *TokenManager) ParseAccess(token string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.accessSecret, nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// ParseRefresh validates a refresh token and returns the user id and jti.
func (m *TokenManager) ParseRefresh(token string) (uuid.UUID, string, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return uuid.Nil, errors.New("unexpected signing method")
		}
		return m.refreshSecret, nil
	})
	if err != nil {
		return uuid.Nil, "", err
	}
	uid, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, "", err
	}
	return uid, claims.ID, nil
}
