package database

import (
	"fmt"

	"gpu-cloud-manager/internal/config"
	"gpu-cloud-manager/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Initialize sets up the database connection and runs migrations
func Initialize(cfg *config.Config) (*gorm.DB, error) {
	// Configure GORM logger based on environment
	var logLevel logger.LogLevel
	switch cfg.Environment {
	case "production":
		logLevel = logger.Error
	case "development":
		logLevel = logger.Info
	default:
		logLevel = logger.Warn
	}
	
	// Open database connection
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogLevel(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	
	// Get underlying SQL DB for connection configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}
	
	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	
	// Run auto migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}
	
	return db, nil
}

// runMigrations runs all database migrations
func runMigrations(db *gorm.DB) error {
	// List of models to migrate
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.UserProvider{},
		&models.Instance{},
	}
	
	for _, model := range modelsToMigrate {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %v", model, err)
		}
	}
	
	// Create indexes
	if err := createIndexes(db); err != nil {
		return fmt.Errorf("failed to create indexes: %v", err)
	}
	
	return nil
}

// createIndexes creates additional database indexes for performance
func createIndexes(db *gorm.DB) error {
	// Create compound indexes for better query performance
	indexes := []struct {
		table   string
		index   string
		columns []string
	}{
		{
			table:   "instances",
			index:   "idx_user_provider_status",
			columns: []string{"user_id", "provider", "status"},
		},
		{
			table:   "instances",
			index:   "idx_provider_status_created",
			columns: []string{"provider", "status", "created_at"},
		},
		{
			table:   "user_providers",
			index:   "idx_user_provider_enabled",
			columns: []string{"user_id", "provider", "is_enabled"},
		},
	}
	
	for _, idx := range indexes {
		query := fmt.Sprintf(
			"CREATE INDEX IF NOT EXISTS %s ON %s (%s)",
			idx.index,
			idx.table,
			joinStrings(idx.columns, ", "),
		)
		
		if err := db.Exec(query).Error; err != nil {
			return fmt.Errorf("failed to create index %s: %v", idx.index, err)
		}
	}
	
	return nil
}

// Helper function to join strings
func joinStrings(strings []string, separator string) string {
	if len(strings) == 0 {
		return ""
	}
	
	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += separator + strings[i]
	}
	
	return result
}
