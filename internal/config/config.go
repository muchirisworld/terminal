package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the application configuration.
type Config struct {
	AppEnv             string
	HTTPPort           int
	DatabaseURL        string
	LogLevel           string
	ShutdownTimeout    time.Duration
	ClerkWebhookSecret string
	ClerkSecretKey     string
}

// New creates a new Config struct.
func New() *Config {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: .env file not found")
		}
	}

	return &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		HTTPPort:           getEnvAsInt("HTTP_PORT", 8080),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		LogLevel:           getEnv("LOG_LEVEL", "debug"),
		ShutdownTimeout:    getEnvAsDuration("SHUTDOWN_TIMEOUT", 5*time.Second),
		ClerkWebhookSecret: getEnv("CLERK_WEBHOOK_SECRET", ""),
		ClerkSecretKey:     getEnv("CLERK_SECRET_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		d, err := time.ParseDuration(value)
		if err != nil {
			return fallback
		}
		return d
	}
	return fallback
}
