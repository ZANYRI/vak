package config

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

// Config holds all application configuration sourced from environment variables.
type Config struct {
	AppEnv          string
	HTTPPort        string
	MetricsPort     string
	DatabaseURL     string
	MigrationsPath  string
	QueueURL        string
	QueueStream     string
	QueueConsumerAPI string
	QueueConsumerWorker string
	SchedulerInterval string
	JWT             JWTConfig
	LogLevel        string
}

// JWTConfig holds JWT settings.
type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

// Load reads configuration from environment variables and validates required fields.
func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:              getEnvDefault("APP_ENV", "local"),
		HTTPPort:            getEnvDefault("HTTP_PORT", "8080"),
		MetricsPort:         getEnvDefault("METRICS_PORT", "9090"),
		DatabaseURL:         getEnvDefault("DATABASE_URL", "postgres://billing:billing@localhost:5432/billing?sslmode=disable"),
		MigrationsPath:      getEnvDefault("MIGRATIONS_PATH", "migrations"),
		QueueURL:            getEnvDefault("QUEUE_URL", "nats://localhost:4222"),
		QueueStream:         getEnvDefault("QUEUE_STREAM", "BILLING"),
		QueueConsumerAPI:    getEnvDefault("QUEUE_CONSUMER_API", "api"),
		QueueConsumerWorker: getEnvDefault("QUEUE_CONSUMER_WORKER", "worker"),
		SchedulerInterval:   getEnvDefault("SCHEDULER_INTERVAL", "*/1 * * * *"),
		JWT: JWTConfig{
			AccessSecret:  getEnvDefault("JWT_ACCESS_SECRET", "change-me-access-secret-minimum-32-characters"),
			RefreshSecret: getEnvDefault("JWT_REFRESH_SECRET", "change-me-refresh-secret-minimum-32-characters"),
			AccessTTL:     parseDuration(getEnvDefault("JWT_ACCESS_TTL", "15m")),
			RefreshTTL:    parseDuration(getEnvDefault("JWT_REFRESH_TTL", "7d")),
		},
		LogLevel: getEnvDefault("LOG_LEVEL", "info"),
	}

	if cfg.JWT.AccessSecret == "" || cfg.JWT.RefreshSecret == "" {
		return nil, fmt.Errorf("JWT secrets are required")
	}
	if _, err := uuid.Parse("00000000-0000-0000-0000-000000000000"); err != nil {
		return nil, err
	}
	_ = uuid.New()
	return cfg, nil
}

func getEnvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 15 * time.Minute
	}
	return d
}
