package config

import (
	"os"
)

// Config holds environment-based settings. Load from env to avoid hardcoding secrets.
type Config struct {
	Port      string
	JWTSecret string
	DBPath    string
}

// Load reads config from environment with safe defaults for local dev.
// In production, set PORT, JWT_SECRET, and DB_PATH.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "change-me-in-production"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "app.db"
	}
	return &Config{
		Port:      port,
		JWTSecret: jwtSecret,
		DBPath:    dbPath,
	}
}
