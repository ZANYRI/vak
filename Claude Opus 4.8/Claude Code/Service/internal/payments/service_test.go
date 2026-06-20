package payments

import "testing"

func newSvc() *Service {
	return NewService(nil, nil, nil, nil, nil, nil)
}

func TestDecideOutcomeSuccessCard(t *testing.T) {
	s := newSvc()
	ok, reason := s.decideOutcome("4242424242420000")
	if !ok || reason != "" {
		t.Fatalf("card ending 0000 should succeed, got ok=%v reason=%q", ok, reason)
	}
}

func TestDecideOutcomeFailCard(t *testing.T) {
	s := newSvc()
	ok, reason := s.decideOutcome("4000000000009999")
	if ok || reason == "" {
		t.Fatalf("card ending 9999 should fail, got ok=%v reason=%q", ok, reason)
	}
}

func TestLast4(t *testing.T) {
	if got := last4("4242424242424242"); got != "4242" {
		t.Fatalf("last4 got %q want 4242", got)
	}
	if got := last4("12"); got != "12" {
		t.Fatalf("last4 short got %q want 12", got)
	}
}
