package models

import (
	"time"

	"gpu-cloud-manager/pkg/types"
	"gorm.io/gorm"
)

// Instance represents a GPU instance in the database
type Instance struct {
	ID           uint                     `gorm:"primaryKey" json:"id"`
	UserID       uint                     `gorm:"not null;index" json:"user_id"`
	Provider     types.GPUProvider        `gorm:"not null;index" json:"provider"`
	ProviderID   string                   `gorm:"not null;uniqueIndex:idx_provider_instance" json:"provider_id"`
	Name         string                   `gorm:"not null" json:"name"`
	Status       types.InstanceStatus     `gorm:"not null;index" json:"status"`
	GPUModel     string                   `gorm:"not null" json:"gpu_model"`
	GPUCount     int                      `gorm:"not null" json:"gpu_count"`
	CPUCount     int                      `gorm:"not null" json:"cpu_count"`
	RAM          int                      `gorm:"not null" json:"ram_gb"`
	Storage      int                      `gorm:"not null" json:"storage_gb"`
	PricePerHour float64                  `gorm:"not null" json:"price_per_hour"`
	Region       string                   `gorm:"not null" json:"region"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
	DeletedAt    gorm.DeletedAt           `gorm:"index" json:"-"`
	
	// Foreign key relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
	
	// JSON fields for provider-specific data
	ProviderData ProviderData `gorm:"type:jsonb" json:"provider_data,omitempty"`
}

// ProviderData is a custom type for storing provider-specific data as JSON
type ProviderData map[string]interface{}

// TableName overrides the table name for the Instance model
func (Instance) TableName() string {
	return "instances"
}

// ToGPUInstance converts database model to API type
func (i *Instance) ToGPUInstance() types.GPUInstance {
	return types.GPUInstance{
		ID:           i.generateAPIID(),
		Provider:     i.Provider,
		ProviderID:   i.ProviderID,
		Name:         i.Name,
		Status:       i.Status,
		GPUModel:     i.GPUModel,
		GPUCount:     i.GPUCount,
		CPUCount:     i.CPUCount,
		RAM:          i.RAM,
		Storage:      i.Storage,
		PricePerHour: i.PricePerHour,
		Region:       i.Region,
		CreatedAt:    i.CreatedAt,
		UpdatedAt:    i.UpdatedAt,
		ProviderData: map[string]interface{}(i.ProviderData),
	}
}

// generateAPIID creates the external API ID format
func (i *Instance) generateAPIID() string {
	switch i.Provider {
	case types.VastAI:
		return "vast_" + i.ProviderID
	default:
		return string(i.Provider) + "_" + i.ProviderID
	}
}

// FromGPUInstance populates database model from API type
func (i *Instance) FromGPUInstance(gpu types.GPUInstance, userID uint) {
	i.UserID = userID
	i.Provider = gpu.Provider
	i.ProviderID = gpu.ProviderID
	i.Name = gpu.Name
	i.Status = gpu.Status
	i.GPUModel = gpu.GPUModel
	i.GPUCount = gpu.GPUCount
	i.CPUCount = gpu.CPUCount
	i.RAM = gpu.RAM
	i.Storage = gpu.Storage
	i.PricePerHour = gpu.PricePerHour
	i.Region = gpu.Region
	i.ProviderData = ProviderData(gpu.ProviderData)
}
