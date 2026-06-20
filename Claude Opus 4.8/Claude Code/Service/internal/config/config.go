package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all runtime configuration, loaded from environment variables.
type Config struct {
	AppEnv string
	HTTPPort string

	DatabaseURL string
	QueueURL    string

	JWTAccessSecret  string
	JWTRefreshSecret string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration

	BcryptCost int

	AuthRateLimit       int           // requests
	AuthRateLimitWindow time.Duration // per window

	WorkerConcurrency int
	JobMaxAttempts    int

	SchedulerInterval time.Duration

	CORSAllowedOrigins string

	adminEmail    string
	adminPassword string
}

// AdminEmail / AdminPassword bootstrap an initial admin account if both are set.
func (c *Config) AdminEmail() string    { return c.adminEmail }
func (c *Config) AdminPassword() string { return c.adminPassword }

// Load reads configuration from the environment, applying sane defaults.
func Load() (*Config, error) {
	c := &Config{
		AppEnv:              getenv("APP_ENV", "local"),
		HTTPPort:            getenv("HTTP_PORT", "8080"),
		DatabaseURL:         getenv("DATABASE_URL", "postgres://billing:billing@localhost:5432/billing?sslmode=disable"),
		QueueURL:            getenv("QUEUE_URL", "redis://localhost:6379/0"),
		JWTAccessSecret:     getenv("JWT_ACCESS_SECRET", "change-me-access"),
		JWTRefreshSecret:    getenv("JWT_REFRESH_SECRET", "change-me-refresh"),
		AccessTokenTTL:      getdur("JWT_ACCESS_TTL", 15*time.Minute),
		RefreshTokenTTL:     getdur("JWT_REFRESH_TTL", 720*time.Hour),
		BcryptCost:          getint("BCRYPT_COST", 12),
		AuthRateLimit:       getint("AUTH_RATE_LIMIT", 10),
		AuthRateLimitWindow: getdur("AUTH_RATE_LIMIT_WINDOW", time.Minute),
		WorkerConcurrency:   getint("WORKER_CONCURRENCY", 4),
		JobMaxAttempts:      getint("JOB_MAX_ATTEMPTS", 5),
		SchedulerInterval:   getdur("SCHEDULER_INTERVAL", time.Minute),
		CORSAllowedOrigins:  getenv("CORS_ALLOWED_ORIGINS", "*"),
		adminEmail:          getenv("ADMIN_EMAIL", "admin@billing.local"),
		adminPassword:       getenv("ADMIN_PASSWORD", "admin12345"),
	}

	if c.AppEnv != "local" && c.AppEnv != "test" {
		if c.JWTAccessSecret == "change-me-access" || c.JWTRefreshSecret == "change-me-refresh" {
			return nil, fmt.Errorf("JWT secrets must be set outside local/test environments")
		}
	}
	return c, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getint(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getdur(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
