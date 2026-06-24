package payments

import (
	"testing"

	"billing-service/internal/models"
)

func TestDetermineOutcome(t *testing.T) {
	status, _ := determineOutcome("0000")
	if status != models.PaymentSucceeded {
		t.Errorf("0000 should succeed, got %s", status)
	}
	status, _ = determineOutcome("9999")
	if status != models.PaymentFailed {
		t.Errorf("9999 should fail, got %s", status)
	}
}
