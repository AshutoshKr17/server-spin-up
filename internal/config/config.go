package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Port        string
	Environment string
	
	// Database configuration
	DatabaseURL string
	
	// GPU Provider APIs
	VastAIAPIKey   string
	RunPodAPIKey   string
	LambdaAPIKey   string
	PaperspaceKey  string
	
	// Feature flags
	EnableMetrics bool
	EnableCORS    bool
	
	// Rate limiting
	RateLimitEnabled bool
	RateLimitRPM     int // Requests per minute
	
	// Logging
	LogLevel string
}

// Load loads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost/gpu_cloud_manager?sslmode=disable"),
		
		VastAIAPIKey:   getEnv("VAST_AI_API_KEY", ""),
		RunPodAPIKey:   getEnv("RUNPOD_API_KEY", ""),
		LambdaAPIKey:   getEnv("LAMBDA_API_KEY", ""),
		PaperspaceKey:  getEnv("PAPERSPACE_API_KEY", ""),
		
		EnableMetrics: getBoolEnv("ENABLE_METRICS", true),
		EnableCORS:    getBoolEnv("ENABLE_CORS", true),
		
		RateLimitEnabled: getBoolEnv("RATE_LIMIT_ENABLED", true),
		RateLimitRPM:     getIntEnv("RATE_LIMIT_RPM", 100),
		
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
	
	return cfg
}

// getEnv gets an environment variable with a fallback default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getBoolEnv gets a boolean environment variable with a fallback default
func getBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	
	return boolValue
}

// getIntEnv gets an integer environment variable with a fallback default
func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	
	return intValue
}
