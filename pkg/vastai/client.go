package vastai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gpu-cloud-manager/pkg/types"
)

const (
	BaseURL = "https://console.vast.ai/api/v0"
)

// Client represents a Vast.ai API client
type Client struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewClient creates a new Vast.ai API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: BaseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// VastOffer represents a GPU offer from Vast.ai
type VastOffer struct {
	ID                int     `json:"id"`
	MachineID         int     `json:"machine_id"`
	ComputeCap        int     `json:"compute_cap"`
	CPUCores          float64 `json:"cpu_cores"`
	CPUName           string  `json:"cpu_name"`
	GPUName           string  `json:"gpu_name"`
	GPUMemoryGB       float64 `json:"gpu_ram"`
	NumGPUs           int     `json:"num_gpus"`
	RAMMemoryGB       float64 `json:"ram"`
	DiskSpaceGB       float64 `json:"disk_space"`
	DiskName          string  `json:"disk_name"`
	InternetDown      float64 `json:"inet_down"`
	InternetUp        float64 `json:"inet_up"`
	DirectPortCount   int     `json:"direct_port_count"`
	PricePerHour      float64 `json:"dph_total"`
	Reliability       float64 `json:"reliability2"`
	IsAvailable       bool    `json:"rentable"`
	Rented            bool    `json:"rented"`
	PublicIPv4        string  `json:"public_ipaddr"`
	GeolocationString string  `json:"geolocation"`
	Datacenter        string  `json:"datacenter_name"`
	HostRunTime       float64 `json:"host_run_time"`
	Score             float64 `json:"score"`
}

// VastInstance represents a rented instance from Vast.ai
type VastInstance struct {
	ID               int    `json:"id"`
	MachineID        int    `json:"machine_id"`
	StatusMsg        string `json:"status_msg"`
	ActualStatus     string `json:"actual_status"`
	IntendedStatus   string `json:"intended_status"`
	PublicIPAddress  string `json:"public_ipaddr"`
	SSHHost          string `json:"ssh_host"`
	SSHPort          int    `json:"ssh_port"`
	Label            string `json:"label"`
	Image            string `json:"image"`
	OnStartScript    string `json:"onstart"`
	PricePerHour     float64 `json:"dph_total"`
	StartDate        string `json:"start_date"`
	Duration         float64 `json:"duration"`
}

// SearchOffers searches for available GPU offers
func (c *Client) SearchOffers(filter *SearchOffersRequest) ([]VastOffer, error) {
	endpoint := "/bundles"
	
	// Build query parameters
	params := url.Values{}
	if filter != nil {
		if filter.GPUName != "" {
			params.Set("q", fmt.Sprintf(`gpu_name:%s`, filter.GPUName))
		}
		if filter.MinGPUCount > 0 {
			params.Set("min_num_gpus", strconv.Itoa(filter.MinGPUCount))
		}
		if filter.MaxPrice > 0 {
			params.Set("max_dph", fmt.Sprintf("%.2f", filter.MaxPrice))
		}
		if filter.MinRAM > 0 {
			params.Set("min_ram", strconv.Itoa(filter.MinRAM))
		}
		if filter.AvailableOnly {
			params.Set("rentable", "true")
		}
		if filter.Datacenter != "" {
			params.Set("datacenter", filter.Datacenter)
		}
		// Default sorting by price per hour
		params.Set("order", "dph_total-")
	}

	var offers []VastOffer
	err := c.makeRequest("GET", endpoint+"?"+params.Encode(), nil, &offers)
	return offers, err
}

// GetInstances retrieves user's rented instances
func (c *Client) GetInstances() ([]VastInstance, error) {
	endpoint := "/instances"
	
	var response struct {
		Instances []VastInstance `json:"instances"`
	}
	
	err := c.makeRequest("GET", endpoint, nil, &response)
	return response.Instances, err
}

// CreateInstance creates a new GPU instance
func (c *Client) CreateInstance(request *CreateInstanceRequest) (*VastInstance, error) {
	endpoint := "/asks/" + strconv.Itoa(request.OfferID) + "/"
	
	payload := map[string]interface{}{
		"price":  request.Price,
		"disk":   request.DiskSizeGB,
		"image":  request.Image,
		"label":  request.Label,
	}
	
	if request.OnStartScript != "" {
		payload["onstart"] = request.OnStartScript
	}
	
	var instance VastInstance
	err := c.makeRequest("PUT", endpoint, payload, &instance)
	return &instance, err
}

// DestroyInstance terminates a GPU instance
func (c *Client) DestroyInstance(instanceID int) error {
	endpoint := fmt.Sprintf("/instances/%d/", instanceID)
	return c.makeRequest("DELETE", endpoint, nil, nil)
}

// StartInstance starts a stopped instance
func (c *Client) StartInstance(instanceID int) error {
	endpoint := fmt.Sprintf("/instances/%d/", instanceID)
	payload := map[string]string{"state": "running"}
	return c.makeRequest("PUT", endpoint, payload, nil)
}

// StopInstance stops a running instance
func (c *Client) StopInstance(instanceID int) error {
	endpoint := fmt.Sprintf("/instances/%d/", instanceID)
	payload := map[string]string{"state": "stopped"}
	return c.makeRequest("PUT", endpoint, payload, nil)
}

// GetInstance retrieves details of a specific instance
func (c *Client) GetInstance(instanceID int) (*VastInstance, error) {
	endpoint := fmt.Sprintf("/instances/%d/", instanceID)
	
	var instance VastInstance
	err := c.makeRequest("GET", endpoint, nil, &instance)
	return &instance, err
}

// makeRequest performs HTTP requests to Vast.ai API
func (c *Client) makeRequest(method, endpoint string, payload interface{}, result interface{}) error {
	var body io.Reader
	
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("error marshaling request body: %v", err)
		}
		body = bytes.NewBuffer(jsonData)
	}
	
	req, err := http.NewRequest(method, c.baseURL+endpoint, body)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	
	// Make the request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()
	
	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}
	
	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}
	
	// Parse response if result is provided
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("error parsing response: %v", err)
		}
	}
	
	return nil
}

// Request types for API calls

// SearchOffersRequest represents parameters for searching offers
type SearchOffersRequest struct {
	GPUName       string  `json:"gpu_name,omitempty"`
	MinGPUCount   int     `json:"min_gpu_count,omitempty"`
	MaxPrice      float64 `json:"max_price,omitempty"`
	MinRAM        int     `json:"min_ram,omitempty"`
	Datacenter    string  `json:"datacenter,omitempty"`
	AvailableOnly bool    `json:"available_only,omitempty"`
}

// CreateInstanceRequest represents parameters for creating an instance
type CreateInstanceRequest struct {
	OfferID       int     `json:"offer_id"`
	Price         float64 `json:"price"`
	DiskSizeGB    int     `json:"disk_size_gb"`
	Image         string  `json:"image"`
	Label         string  `json:"label,omitempty"`
	OnStartScript string  `json:"onstart_script,omitempty"`
}

// Helper functions to convert between Vast.ai types and our internal types

// ConvertOfferToGPUInstance converts a Vast.ai offer to our internal GPU instance type
func ConvertOfferToGPUInstance(offer VastOffer) types.GPUInstance {
	status := types.StatusOffline
	if offer.IsAvailable && !offer.Rented {
		status = types.StatusOffline
	} else if offer.Rented {
		status = types.StatusRented
	} else {
		status = types.StatusUnavailable
	}

	return types.GPUInstance{
		ID:           fmt.Sprintf("vast_%d", offer.ID),
		Provider:     types.VastAI,
		ProviderID:   strconv.Itoa(offer.ID),
		Name:         fmt.Sprintf("Vast.ai Machine %d", offer.MachineID),
		Status:       status,
		GPUModel:     offer.GPUName,
		GPUCount:     offer.NumGPUs,
		CPUCount:     int(offer.CPUCores),
		RAM:          int(offer.RAMMemoryGB),
		Storage:      int(offer.DiskSpaceGB),
		PricePerHour: offer.PricePerHour,
		Region:       offer.Datacenter,
		ProviderData: map[string]interface{}{
			"machine_id":      offer.MachineID,
			"compute_cap":     offer.ComputeCap,
			"cpu_name":        offer.CPUName,
			"gpu_memory_gb":   offer.GPUMemoryGB,
			"disk_name":       offer.DiskName,
			"internet_down":   offer.InternetDown,
			"internet_up":     offer.InternetUp,
			"public_ipv4":     offer.PublicIPv4,
			"reliability":     offer.Reliability,
			"score":          offer.Score,
		},
	}
}

// ConvertInstanceToGPUInstance converts a Vast.ai instance to our internal type
func ConvertInstanceToGPUInstance(instance VastInstance) types.GPUInstance {
	status := types.StatusOffline
	switch instance.ActualStatus {
	case "running":
		status = types.StatusRunning
	case "loading":
		status = types.StatusLoading
	case "offline":
		status = types.StatusOffline
	}

	return types.GPUInstance{
		ID:           fmt.Sprintf("vast_%d", instance.ID),
		Provider:     types.VastAI,
		ProviderID:   strconv.Itoa(instance.ID),
		Name:         instance.Label,
		Status:       status,
		PricePerHour: instance.PricePerHour,
		ProviderData: map[string]interface{}{
			"machine_id":       instance.MachineID,
			"ssh_host":         instance.SSHHost,
			"ssh_port":         instance.SSHPort,
			"public_ipaddr":    instance.PublicIPAddress,
			"image":           instance.Image,
			"onstart_script":  instance.OnStartScript,
			"start_date":      instance.StartDate,
			"duration":        instance.Duration,
			"status_msg":      instance.StatusMsg,
			"intended_status": instance.IntendedStatus,
		},
	}
}