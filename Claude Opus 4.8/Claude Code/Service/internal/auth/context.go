package auth

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey int

const identityKey ctxKey = iota

// Identity is the authenticated principal attached to a request context.
type Identity struct {
	UserID uuid.UUID
	Role   string
}

// WithIdentity returns a context carrying the identity.
func WithIdentity(ctx context.Context, id Identity) context.Context {
	return context.WithValue(ctx, identityKey, id)
}

// FromContext returns the identity, if present.
func FromContext(ctx context.Context) (Identity, bool) {
	id, ok := ctx.Value(identityKey).(Identity)
	return id, ok
}
