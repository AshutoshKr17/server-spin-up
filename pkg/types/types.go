package types

import "time"

// GPUProvider represents a GPU cloud provider
type GPUProvider string

const (
	VastAI     GPUProvider = "vast_ai"
	RunPod     GPUProvider = "runpod"
	LambdaLabs GPUProvider = "lambda_labs"
	Paperspace GPUProvider = "paperspace"
	AWS        GPUProvider = "aws"
	GCP        GPUProvider = "gcp"
	Azure      GPUProvider = "azure"
)

// GPUModel represents different GPU models with their specifications
type GPUModel struct {
	Name        string  `json:"name"`
	Memory      int     `json:"memory_gb"`
	Compute     float64 `json:"compute_capability"`
	Architecture string `json:"architecture"`
	Category    string  `json:"category"`
	Performance int     `json:"performance_score"`
}

// Predefined GPU models
var GPUModels = map[string]GPUModel{
	// NVIDIA RTX Series
	"RTX 4090": {Name: "RTX 4090", Memory: 24, Compute: 8.9, Architecture: "Ada Lovelace", Category: "consumer", Performance: 100},
	"RTX 4080": {Name: "RTX 4080", Memory: 16, Compute: 8.9, Architecture: "Ada Lovelace", Category: "consumer", Performance: 85},
	"RTX 4070": {Name: "RTX 4070", Memory: 12, Compute: 8.9, Architecture: "Ada Lovelace", Category: "consumer", Performance: 70},
	"RTX 3090": {Name: "RTX 3090", Memory: 24, Compute: 8.6, Architecture: "Ampere", Category: "consumer", Performance: 90},
	"RTX 3080": {Name: "RTX 3080", Memory: 10, Compute: 8.6, Architecture: "Ampere", Category: "consumer", Performance: 80},
	"RTX 3070": {Name: "RTX 3070", Memory: 8, Compute: 8.6, Architecture: "Ampere", Category: "consumer", Performance: 65},
	
	// NVIDIA Professional Series
	"A100": {Name: "A100", Memory: 80, Compute: 8.0, Architecture: "Ampere", Category: "datacenter", Performance: 120},
	"H100": {Name: "H100", Memory: 80, Compute: 9.0, Architecture: "Hopper", Category: "datacenter", Performance: 150},
	"V100": {Name: "V100", Memory: 32, Compute: 7.0, Architecture: "Volta", Category: "datacenter", Performance: 95},
	"A40":  {Name: "A40", Memory: 48, Compute: 8.6, Architecture: "Ampere", Category: "professional", Performance: 85},
	"A6000": {Name: "A6000", Memory: 48, Compute: 8.6, Architecture: "Ampere", Category: "professional", Performance: 90},
	
	// NVIDIA Gaming/Entry Level
	"GTX 1080 Ti": {Name: "GTX 1080 Ti", Memory: 11, Compute: 6.1, Architecture: "Pascal", Category: "consumer", Performance: 45},
	"RTX 2080 Ti": {Name: "RTX 2080 Ti", Memory: 11, Compute: 7.5, Architecture: "Turing", Category: "consumer", Performance: 60},
	
	// AMD GPUs
	"RX 7900 XTX": {Name: "RX 7900 XTX", Memory: 24, Compute: 0.0, Architecture: "RDNA3", Category: "consumer", Performance: 85},
	"RX 6900 XT":  {Name: "RX 6900 XT", Memory: 16, Compute: 0.0, Architecture: "RDNA2", Category: "consumer", Performance: 75},
}

// GPUCategory represents different categories of GPUs
type GPUCategory string

const (
	CategoryConsumer   GPUCategory = "consumer"
	CategoryProfessional GPUCategory = "professional"
	CategoryDatacenter GPUCategory = "datacenter"
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
	
	// Enhanced GPU information
	GPUInfo        *GPUModel              `json:"gpu_info,omitempty"`
	Performance    int                    `json:"performance_score,omitempty"`
	Reliability    float64                `json:"reliability,omitempty"`
	NetworkSpeed   *NetworkInfo           `json:"network_info,omitempty"`
}

// NetworkInfo represents network capabilities
type NetworkInfo struct {
	DownloadMbps int `json:"download_mbps"`
	UploadMbps   int `json:"upload_mbps"`
}

// InstanceStatus represents the status of a GPU instance
type InstanceStatus string

const (
	StatusOffline     InstanceStatus = "offline"
	StatusRunning     InstanceStatus = "running"
	StatusLoading     InstanceStatus = "loading"
	StatusRented      InstanceStatus = "rented"
	StatusUnavailable InstanceStatus = "unavailable"
	StatusStarting    InstanceStatus = "starting"
	StatusStopping    InstanceStatus = "stopping"
	StatusError       InstanceStatus = "error"
)

// CreateInstanceRequest represents a request to create a new GPU instance
type CreateInstanceRequest struct {
	Provider      GPUProvider        `json:"provider" binding:"required"`
	OfferID       string             `json:"offer_id" binding:"required"`
	Image         string             `json:"image"`
	OnStartScript string             `json:"onstart_script,omitempty"`
	SSHKey        string             `json:"ssh_key,omitempty"`
	Label         string             `json:"label,omitempty"`
	Environment   map[string]string  `json:"environment,omitempty"`
	Ports         []PortMapping      `json:"ports,omitempty"`
	Resources     *ResourceRequests  `json:"resources,omitempty"`
}

// PortMapping represents port forwarding configuration
type PortMapping struct {
	ContainerPort int    `json:"container_port"`
	HostPort      int    `json:"host_port,omitempty"`
	Protocol      string `json:"protocol,omitempty"` // tcp, udp
}

// ResourceRequests represents resource requirements
type ResourceRequests struct {
	MinRAM     int `json:"min_ram_gb,omitempty"`
	MinStorage int `json:"min_storage_gb,omitempty"`
	MinGPUs    int `json:"min_gpus,omitempty"`
}

// AdvancedSearchFilter represents filters for searching GPU instances
type AdvancedSearchFilter struct {
	Provider       GPUProvider   `json:"provider,omitempty"`
	GPUModel       string        `json:"gpu_model,omitempty"`
	GPUModels      []string      `json:"gpu_models,omitempty"`
	GPUCategory    GPUCategory   `json:"gpu_category,omitempty"`
	MinGPUCount    int           `json:"min_gpu_count,omitempty"`
	MaxGPUCount    int           `json:"max_gpu_count,omitempty"`
	MinPrice       float64       `json:"min_price,omitempty"`
	MaxPrice       float64       `json:"max_price,omitempty"`
	MinRAM         int           `json:"min_ram_gb,omitempty"`
	MaxRAM         int           `json:"max_ram_gb,omitempty"`
	MinStorage     int           `json:"min_storage_gb,omitempty"`
	Region         string        `json:"region,omitempty"`
	Regions        []string      `json:"regions,omitempty"`
	Available      bool          `json:"available,omitempty"`
	MinReliability float64       `json:"min_reliability,omitempty"`
	MinPerformance int           `json:"min_performance,omitempty"`
	SortBy         string        `json:"sort_by,omitempty"` // price, performance, reliability
	SortOrder      string        `json:"sort_order,omitempty"` // asc, desc
}

// SearchFilter represents basic filters for backward compatibility
type SearchFilter struct {
	Provider    GPUProvider `json:"provider,omitempty"`
	GPUModel    string      `json:"gpu_model,omitempty"`
	MinGPUCount int         `json:"min_gpu_count,omitempty"`
	MaxPrice    float64     `json:"max_price,omitempty"`
	Region      string      `json:"region,omitempty"`
	Available   bool        `json:"available,omitempty"`
}

// GPUModelInfo represents detailed information about available GPU models
type GPUModelInfo struct {
	Model       GPUModel `json:"model"`
	Available   int      `json:"available_count"`
	MinPrice    float64  `json:"min_price"`
	MaxPrice    float64  `json:"max_price"`
	AvgPrice    float64  `json:"avg_price"`
	Providers   []GPUProvider `json:"providers"`
}

// ProviderInfo represents information about a GPU provider
type ProviderInfo struct {
	Name         GPUProvider `json:"name"`
	DisplayName  string      `json:"display_name"`
	Website      string      `json:"website"`
	Regions      []string    `json:"regions"`
	GPUModels    []string    `json:"gpu_models"`
	Features     []string    `json:"features"`
	IsConfigured bool        `json:"is_configured"`
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

// MarketplaceStats represents marketplace statistics
type MarketplaceStats struct {
	TotalInstances    int                        `json:"total_instances"`
	AvailableInstances int                       `json:"available_instances"`
	ModelStats        map[string]GPUModelInfo   `json:"model_stats"`
	ProviderStats     map[string]int            `json:"provider_stats"`
	AveragePrice      float64                   `json:"average_price"`
	PriceRange        PriceRange                `json:"price_range"`
}

// PriceRange represents price range information
type PriceRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
	Avg float64 `json:"avg"`
}
