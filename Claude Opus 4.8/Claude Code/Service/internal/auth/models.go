package auth

import (
	"time"

	"github.com/google/uuid"
)

// User is a system account.
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RegisterRequest is the registration payload.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	Role     string `json:"role" validate:"omitempty,oneof=admin billing_manager support customer"`
}

// LoginRequest is the login payload.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshRequest carries a refresh token.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// TokenPair is returned on login/refresh.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}
