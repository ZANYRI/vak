package auth

import (
	"testing"
	"time"
)

func TestAccessToken(t *testing.T) {
	s := New("12345678901234567890123456789012", "12345678901234567890123456789012", time.Minute, time.Hour)
	token, e := s.AccessToken("u", RoleAdmin)
	if e != nil {
		t.Fatal(e)
	}
	claims, e := s.ParseAccess(token)
	if e != nil || claims.Subject != "u" || claims.Role != RoleAdmin {
		t.Fatalf("%+v %v", claims, e)
	}
}
func TestPasswordPolicy(t *testing.T) {
	if _, e := HashPassword("short"); e == nil {
		t.Fatal("expected error")
	}
	h, e := HashPassword("very-long-pass")
	if e != nil || ComparePassword(h, "very-long-pass") != nil {
		t.Fatal("hashing failed")
	}
}
