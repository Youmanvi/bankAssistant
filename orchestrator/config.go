package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Server
	Host string
	Port string

	// Retell AI
	RetellAPIKey string

	// Python Backend
	PythonBackendURL string

	// Logging
	LogLevel string

	// Database (call state storage)
	CallStateDB string // "memory" or "redis"
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	_ = godotenv.Load(".env")

	return &Config{
		Host:             getEnv("ORCHESTRATOR_HOST", "0.0.0.0"),
		Port:             getEnv("ORCHESTRATOR_PORT", "8001"),
		RetellAPIKey:     getEnv("RETELL_API_KEY", ""),
		PythonBackendURL: getEnv("PYTHON_BACKEND_URL", "http://localhost:8000"),
		LogLevel:         getEnv("LOG_LEVEL", "INFO"),
		CallStateDB:      getEnv("CALL_STATE_DB", "memory"),
	}
}

// getEnv gets an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if defaultValue == "" {
		log.Printf("Warning: Required environment variable %s not set\n", key)
	}
	return defaultValue
}
