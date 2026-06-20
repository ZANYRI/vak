package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func newTM() *TokenManager {
	return NewTokenManager("access-secret", "refresh-secret", 15*time.Minute, time.Hour)
}

func TestAccessTokenRoundTrip(t *testing.T) {
	tm := newTM()
	uid := uuid.New()
	tok, err := tm.GenerateAccess(uid, RoleAdmin)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	claims, err := tm.ParseAccess(tok)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if claims.UserID != uid {
		t.Fatalf("uid mismatch: got %s want %s", claims.UserID, uid)
	}
	if claims.Role != RoleAdmin {
		t.Fatalf("role mismatch: got %s", claims.Role)
	}
}

func TestExpiredAccessTokenRejected(t *testing.T) {
	tm := NewTokenManager("a", "r", -time.Minute, time.Hour) // already expired
	tok, _ := tm.GenerateAccess(uuid.New(), RoleCustomer)
	if _, err := tm.ParseAccess(tok); err == nil {
		t.Fatal("expected expired token to be rejected")
	}
}

func TestRefreshTokenRoundTrip(t *testing.T) {
	tm := newTM()
	uid := uuid.New()
	tok, jti, err := tm.GenerateRefresh(uid)
	if err != nil {
		t.Fatalf("generate refresh: %v", err)
	}
	if jti == "" {
		t.Fatal("expected a non-empty jti")
	}
	parsedUID, parsedJTI, err := tm.ParseRefresh(tok)
	if err != nil {
		t.Fatalf("parse refresh: %v", err)
	}
	if parsedUID != uid || parsedJTI != jti {
		t.Fatalf("refresh mismatch: uid %s/%s jti %s/%s", parsedUID, uid, parsedJTI, jti)
	}
}

func TestAccessTokenWrongSecretRejected(t *testing.T) {
	tm := newTM()
	tok, _ := tm.GenerateAccess(uuid.New(), RoleAdmin)
	other := NewTokenManager("different", "refresh-secret", time.Minute, time.Hour)
	if _, err := other.ParseAccess(tok); err == nil {
		t.Fatal("expected token signed with a different secret to be rejected")
	}
}

func TestRBACHelpers(t *testing.T) {
	if !CanManageBilling(RoleAdmin) || !CanManageBilling(RoleBillingManager) {
		t.Fatal("admin/billing_manager should manage billing")
	}
	if CanManageBilling(RoleSupport) || CanManageBilling(RoleCustomer) {
		t.Fatal("support/customer must not manage billing")
	}
	if !CanViewBilling(RoleSupport) {
		t.Fatal("support should view billing")
	}
}
