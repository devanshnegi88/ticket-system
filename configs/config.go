package configs

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration loaded from environment
// variables (optionally via a .env file).
type Config struct {
	Port         string
	JWTSecret    string
	JWTExpiry    time.Duration
	DatabasePath string
}

// Load reads a .env file if present (it is fine if it does not exist,
// e.g. in a container where env vars are injected directly) and returns
// a populated Config, applying sane defaults for anything left unset.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on environment variables")
	}

	expiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil || expiryHours <= 0 {
		expiryHours = 24
	}

	return &Config{
		Port:         getEnv("PORT", "8080"),
		JWTSecret:    getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiry:    time.Duration(expiryHours) * time.Hour,
		DatabasePath: getEnv("DATABASE_PATH", "ticket_system.db"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
