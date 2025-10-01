package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Name      string         `gorm:"not null" json:"name"`
	APIKey    string         `gorm:"uniqueIndex;not null" json:"-"` // Hidden from JSON
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relationships
	Instances []Instance `gorm:"foreignKey:UserID" json:"-"`
	
	// Provider-specific configurations
	Providers []UserProvider `gorm:"foreignKey:UserID" json:"-"`
}

// TableName overrides the table name for the User model
func (User) TableName() string {
	return "users"
}

// UserProvider stores user's configuration for each GPU provider
type UserProvider struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	Provider   string    `gorm:"not null" json:"provider"`
	APIKey     string    `gorm:"not null" json:"-"` // Hidden from JSON
	IsEnabled  bool      `gorm:"default:true" json:"is_enabled"`
	Config     JSONMap   `gorm:"type:jsonb" json:"config,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	
	// Foreign key relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// JSONMap is a custom type for storing JSON data
type JSONMap map[string]interface{}

// TableName overrides the table name for the UserProvider model
func (UserProvider) TableName() string {
	return "user_providers"
}
