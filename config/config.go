package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all the application configuration parsed from the environment.
type Config struct {
	DBDSN     string
	JWTSecret string
	Port      string
}

// LoadConfig reads the .env file, extracts the necessary variables,
// and applies validation and fallbacks.
func Load() (*Config, error) {
	// Attempt to load the .env file.
	// We ignore the error here because in production environments (like Docker),
	// environment variables are often injected directly by the host rather than a file.
	_ = godotenv.Load()

	// Extract and validate DB DSN
	dbDsn := os.Getenv("DB_DSN")
	if dbDsn == "" {
		return nil, errors.New("DB_DSN environment variable is required but is empty")
	}

	// Extract and validate JWT Secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable is required but is empty")
	}

	// Extract Port and set fallback if empty
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Return the populated Config struct
	return &Config{
		DBDSN:     dbDsn,
		JWTSecret: jwtSecret,
		Port:      port,
	}, nil
}
