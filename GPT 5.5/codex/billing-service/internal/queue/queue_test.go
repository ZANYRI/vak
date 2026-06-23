package queue

import (
	"context"
	"errors"
	"testing"
)

func TestBackoffAndPanicRecovery(t *testing.T) {
	if backoff(0).Seconds() != 1 || backoff(3).Seconds() != 8 {
		t.Fatal("unexpected exponential backoff")
	}
	if err := runSafely(context.Background(), func(context.Context, Message) error { panic("boom") }, Message{}); err == nil {
		t.Fatal("panic was not recovered")
	}
	if err := runSafely(context.Background(), func(context.Context, Message) error { return errors.New("expected") }, Message{}); err == nil {
		t.Fatal("handler error was lost")
	}
}
