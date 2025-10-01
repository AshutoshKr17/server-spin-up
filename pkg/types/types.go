package types

import "time"

// GPUProvider represents a GPU cloud provider
type GPUProvider string

const (
	VastAI GPUProvider = "vast_ai"
	// Add more providers later: RunPod, Lambda Labs, etc.
)

// GPUInstance represents a GPU instance across different providers
type GPUInstance struct {
	ID             string                 `json:"id"`
	Provider       GPUProvider            `json:"provider"`
	ProviderID     string                 `json:"provider_id"`
	Name           string                 `json:"name"`
	Status         InstanceStatus         `json:"status"`
	GPUModel       string                 `json:"gpu_model"`
	GPUCount       int                    `json:"gpu_count"`
	CPUCount       int                    `json:"cpu_count"`
	RAM            int                    `json:"ram_gb"`
	Storage        int                    `json:"storage_gb"`
	PricePerHour   float64                `json:"price_per_hour"`
	Region         string                 `json:"region"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	ProviderData   map[string]interface{} `json:"provider_data,omitempty"`
}

// InstanceStatus represents the status of a GPU instance
type InstanceStatus string

const (
	StatusOffline    InstanceStatus = "offline"
	StatusRunning    InstanceStatus = "running"
	StatusLoading    InstanceStatus = "loading"
	StatusRented     InstanceStatus = "rented"
	StatusUnavailable InstanceStatus = "unavailable"
)

// CreateInstanceRequest represents a request to create a new GPU instance
type CreateInstanceRequest struct {
	Provider     GPUProvider `json:"provider" binding:"required"`
	OfferID      string      `json:"offer_id" binding:"required"`
	Image        string      `json:"image"`
	OnStartScript string      `json:"onstart_script,omitempty"`
	SSHKey       string      `json:"ssh_key,omitempty"`
	Label        string      `json:"label,omitempty"`
}

// SearchFilter represents filters for searching GPU instances
type SearchFilter struct {
	Provider    GPUProvider `json:"provider,omitempty"`
	GPUModel    string      `json:"gpu_model,omitempty"`
	MinGPUCount int         `json:"min_gpu_count,omitempty"`
	MaxPrice    float64     `json:"max_price,omitempty"`
	Region      string      `json:"region,omitempty"`
	Available   bool        `json:"available,omitempty"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	APIResponse
	Pagination Pagination `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
