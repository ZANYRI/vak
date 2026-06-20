package observability

import (
	"log/slog"
	"os"
)

// NewLogger builds a structured JSON logger. In local env it uses text for readability.
func NewLogger(env string) *slog.Logger {
	level := slog.LevelInfo
	if env == "local" || env == "test" {
		level = slog.LevelDebug
	}
	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if env == "local" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	return slog.New(handler)
}
