package config

import (
	"os"
	"testing"
)

func TestLoadDefaultConfig(t *testing.T) {
	// Clear environment variables
	os.Clearenv()

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("Expected default port to be '8080', got %s", cfg.Port)
	}

	if cfg.Environment != "development" {
		t.Errorf("Expected default environment to be 'development', got %s", cfg.Environment)
	}

	if cfg.LogLevel != "info" {
		t.Errorf("Expected default log level to be 'info', got %s", cfg.LogLevel)
	}

	if !cfg.EnableMetrics {
		t.Error("Expected default EnableMetrics to be true")
	}

	if !cfg.EnableCORS {
		t.Error("Expected default EnableCORS to be true")
	}

	if cfg.RateLimitRPM != 100 {
		t.Errorf("Expected default rate limit RPM to be 100, got %d", cfg.RateLimitRPM)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("PORT", "9000")
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost/test")
	os.Setenv("VAST_AI_API_KEY", "test_key_123")
	os.Setenv("ENABLE_METRICS", "false")
	os.Setenv("ENABLE_CORS", "false")
	os.Setenv("RATE_LIMIT_ENABLED", "false")
	os.Setenv("RATE_LIMIT_RPM", "200")
	os.Setenv("LOG_LEVEL", "debug")

	cfg := Load()

	if cfg.Port != "9000" {
		t.Errorf("Expected port to be '9000', got %s", cfg.Port)
	}

	if cfg.Environment != "production" {
		t.Errorf("Expected environment to be 'production', got %s", cfg.Environment)
	}

	if cfg.DatabaseURL != "postgres://test:test@localhost/test" {
		t.Errorf("Expected database URL to match, got %s", cfg.DatabaseURL)
	}

	if cfg.VastAIAPIKey != "test_key_123" {
		t.Errorf("Expected Vast AI API key to be 'test_key_123', got %s", cfg.VastAIAPIKey)
	}

	if cfg.EnableMetrics {
		t.Error("Expected EnableMetrics to be false")
	}

	if cfg.EnableCORS {
		t.Error("Expected EnableCORS to be false")
	}

	if cfg.RateLimitEnabled {
		t.Error("Expected RateLimitEnabled to be false")
	}

	if cfg.RateLimitRPM != 200 {
		t.Errorf("Expected rate limit RPM to be 200, got %d", cfg.RateLimitRPM)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("Expected log level to be 'debug', got %s", cfg.LogLevel)
	}

	// Clean up
	os.Clearenv()
}

func TestGetEnv(t *testing.T) {
	// Test with existing environment variable
	os.Setenv("TEST_VAR", "test_value")
	result := getEnv("TEST_VAR", "default")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got %s", result)
	}

	// Test with non-existing environment variable
	result = getEnv("NON_EXISTENT_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got %s", result)
	}

	os.Unsetenv("TEST_VAR")
}

func TestGetBoolEnv(t *testing.T) {
	// Test with valid boolean values
	testCases := []struct {
		envValue string
		expected bool
	}{
		{"true", true},
		{"True", true},
		{"TRUE", true},
		{"1", true},
		{"false", false},
		{"False", false},
		{"FALSE", false},
		{"0", false},
	}

	for _, tc := range testCases {
		os.Setenv("TEST_BOOL_VAR", tc.envValue)
		result := getBoolEnv("TEST_BOOL_VAR", false)
		if result != tc.expected {
			t.Errorf("For env value '%s', expected %t, got %t", tc.envValue, tc.expected, result)
		}
	}

	// Test with invalid boolean value (should return default)
	os.Setenv("TEST_BOOL_VAR", "invalid")
	result := getBoolEnv("TEST_BOOL_VAR", true)
	if result != true {
		t.Errorf("Expected default value true for invalid boolean, got %t", result)
	}

	// Test with non-existent environment variable
	result = getBoolEnv("NON_EXISTENT_BOOL_VAR", false)
	if result != false {
		t.Errorf("Expected default value false for non-existent var, got %t", result)
	}

	os.Unsetenv("TEST_BOOL_VAR")
}

func TestGetIntEnv(t *testing.T) {
	// Test with valid integer
	os.Setenv("TEST_INT_VAR", "42")
	result := getIntEnv("TEST_INT_VAR", 0)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	// Test with invalid integer (should return default)
	os.Setenv("TEST_INT_VAR", "invalid")
	result = getIntEnv("TEST_INT_VAR", 100)
	if result != 100 {
		t.Errorf("Expected default value 100 for invalid integer, got %d", result)
	}

	// Test with non-existent environment variable
	result = getIntEnv("NON_EXISTENT_INT_VAR", 50)
	if result != 50 {
		t.Errorf("Expected default value 50 for non-existent var, got %d", result)
	}

	os.Unsetenv("TEST_INT_VAR")
}
