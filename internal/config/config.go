package config

import (
	"os"
)

type Config struct {
	Env      string
	DBUrl    string
	LogLevel string
}

func Load() (*Config, error) {
	cfg := &Config{
		Env:      getEnv("APP_ENV", "development"),
		DBUrl:    getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/gofinancialsystem?sslmode=disable"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
