package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env                    string
	HTTPPort               string
	DatabaseURL            string
	QueueURL               string
	JWTAccessSecret        string
	JWTRefreshSecret       string
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	CORSOrigins            string
	BootstrapAdminEmail    string
	BootstrapAdminPassword string
}

func Load() (Config, error) {
	c := Config{
		Env: env("APP_ENV", "local"), HTTPPort: env("HTTP_PORT", "8080"),
		DatabaseURL:     env("DATABASE_URL", "postgres://billing:billing@localhost:5432/billing?sslmode=disable"),
		QueueURL:        env("QUEUE_URL", "nats://localhost:4222"),
		JWTAccessSecret: os.Getenv("JWT_ACCESS_SECRET"), JWTRefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
		AccessTokenTTL: duration("ACCESS_TOKEN_TTL", 15*time.Minute), RefreshTokenTTL: duration("REFRESH_TOKEN_TTL", 30*24*time.Hour),
		CORSOrigins:         env("CORS_ORIGINS", "http://localhost:3000"),
		BootstrapAdminEmail: os.Getenv("BOOTSTRAP_ADMIN_EMAIL"), BootstrapAdminPassword: os.Getenv("BOOTSTRAP_ADMIN_PASSWORD"),
	}
	if c.JWTAccessSecret == "" || c.JWTRefreshSecret == "" {
		return Config{}, fmt.Errorf("JWT_ACCESS_SECRET and JWT_REFRESH_SECRET must be set")
	}
	if len(c.JWTAccessSecret) < 32 || len(c.JWTRefreshSecret) < 32 {
		return Config{}, fmt.Errorf("JWT secrets must be at least 32 bytes")
	}
	return c, nil
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func duration(k string, d time.Duration) time.Duration {
	if v, e := time.ParseDuration(os.Getenv(k)); e == nil && v > 0 {
		return v
	}
	return d
}
func Int(k string, d int) int {
	if v, e := strconv.Atoi(os.Getenv(k)); e == nil {
		return v
	}
	return d
}
