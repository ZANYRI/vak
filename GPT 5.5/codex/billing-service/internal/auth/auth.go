package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	RoleAdmin          = "admin"
	RoleBillingManager = "billing_manager"
	RoleCustomer       = "customer"
	RoleSupport        = "support"
)

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}
type TokenService struct {
	accessSecret, refreshSecret []byte
	accessTTL, refreshTTL       time.Duration
}

func New(access, refresh string, accessTTL, refreshTTL time.Duration) TokenService {
	return TokenService{[]byte(access), []byte(refresh), accessTTL, refreshTTL}
}
func HashPassword(password string) (string, error) {
	if len(password) < 12 {
		return "", fmt.Errorf("password must contain at least 12 characters")
	}
	b, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), e
}
func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
func (s TokenService) AccessToken(userID, role string) (string, error) {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{role, jwt.RegisteredClaims{Subject: userID, IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)), ID: randomID()}}).SignedString(s.accessSecret)
}
func (s TokenService) ParseAccess(token string) (*Claims, error) {
	parsed, e := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.accessSecret, nil
	})
	if e != nil {
		return nil, e
	}
	c, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, fmt.Errorf("invalid access token")
	}
	return c, nil
}
func (s TokenService) NewRefreshToken() (plain, hash string, expiry time.Time, err error) {
	b := make([]byte, 48)
	if _, err = rand.Read(b); err != nil {
		return
	}
	plain = base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(plain))
	hash = base64.RawURLEncoding.EncodeToString(sum[:])
	expiry = time.Now().Add(s.refreshTTL)
	return
}
func HashRefresh(v string) string {
	sum := sha256.Sum256([]byte(v))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
func randomID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
func Allowed(role string, permitted ...string) bool {
	for _, p := range permitted {
		if role == RoleAdmin || role == p {
			return true
		}
	}
	return false
}
