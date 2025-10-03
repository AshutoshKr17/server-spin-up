package vastai

import (
	"testing"

	"gpu-cloud-manager/pkg/types"
)

func TestNewClient(t *testing.T) {
	apiKey := "test_api_key"
	client := NewClient(apiKey)

	if client.apiKey != apiKey {
		t.Errorf("Expected API key to be %s, got %s", apiKey, client.apiKey)
	}

	if client.baseURL != BaseURL {
		t.Errorf("Expected base URL to be %s, got %s", BaseURL, client.baseURL)
	}

	if client.client == nil {
		t.Error("Expected HTTP client to be initialized")
	}
}

func TestConvertOfferToGPUInstance(t *testing.T) {
	offer := VastOffer{
		ID:               12345,
		MachineID:        67890,
		ComputeCap:       89,
		CPUCores:         8.0,
		CPUName:          "Intel Xeon",
		GPUName:          "RTX 4090",
		GPUMemoryGB:      24.0,
		NumGPUs:          1,
		RAMMemoryGB:      32.0,
		DiskSpaceGB:      100.0,
		DiskName:         "NVMe SSD",
		InternetDown:     1000.0,
		InternetUp:       100.0,
		DirectPortCount:  1,
		PricePerHour:     1.50,
		Reliability:      0.95,
		IsAvailable:      true,
		Rented:           false,
		PublicIPv4:       "1.2.3.4",
		GeolocationString: "US",
		Datacenter:       "US-East",
		HostRunTime:      1000.0,
		Score:            8.5,
	}

	instance := ConvertOfferToGPUInstance(offer)

	if instance.ID != "vast_12345" {
		t.Errorf("Expected ID to be 'vast_12345', got %s", instance.ID)
	}

	if instance.Provider != types.VastAI {
		t.Errorf("Expected provider to be VastAI, got %s", instance.Provider)
	}

	if instance.ProviderID != "12345" {
		t.Errorf("Expected provider ID to be '12345', got %s", instance.ProviderID)
	}

	if instance.GPUModel != "RTX 4090" {
		t.Errorf("Expected GPU model to be 'RTX 4090', got %s", instance.GPUModel)
	}

	if instance.GPUCount != 1 {
		t.Errorf("Expected GPU count to be 1, got %d", instance.GPUCount)
	}

	if instance.CPUCount != 8 {
		t.Errorf("Expected CPU count to be 8, got %d", instance.CPUCount)
	}

	if instance.RAM != 32 {
		t.Errorf("Expected RAM to be 32, got %d", instance.RAM)
	}

	if instance.Storage != 100 {
		t.Errorf("Expected storage to be 100, got %d", instance.Storage)
	}

	if instance.PricePerHour != 1.50 {
		t.Errorf("Expected price per hour to be 1.50, got %.2f", instance.PricePerHour)
	}

	if instance.Region != "US-East" {
		t.Errorf("Expected region to be 'US-East', got %s", instance.Region)
	}

	// Test status conversion for available offer
	if instance.Status != types.StatusOffline {
		t.Errorf("Expected status to be StatusOffline for available offer, got %s", instance.Status)
	}
}

func TestConvertInstanceToGPUInstance(t *testing.T) {
	vastInstance := VastInstance{
		ID:              12345,
		MachineID:       67890,
		StatusMsg:       "Running",
		ActualStatus:    "running",
		IntendedStatus:  "running",
		PublicIPAddress: "1.2.3.4",
		SSHHost:         "ssh5.vast.ai",
		SSHPort:         12345,
		Label:           "My Training Instance",
		Image:           "pytorch/pytorch:latest",
		OnStartScript:   "#!/bin/bash\necho 'Starting'",
		PricePerHour:    1.50,
		StartDate:       "2024-01-15T10:30:00Z",
		Duration:        1.5,
	}

	instance := ConvertInstanceToGPUInstance(vastInstance)

	if instance.ID != "vast_12345" {
		t.Errorf("Expected ID to be 'vast_12345', got %s", instance.ID)
	}

	if instance.Provider != types.VastAI {
		t.Errorf("Expected provider to be VastAI, got %s", instance.Provider)
	}

	if instance.Name != "My Training Instance" {
		t.Errorf("Expected name to be 'My Training Instance', got %s", instance.Name)
	}

	if instance.Status != types.StatusRunning {
		t.Errorf("Expected status to be StatusRunning, got %s", instance.Status)
	}

	if instance.PricePerHour != 1.50 {
		t.Errorf("Expected price per hour to be 1.50, got %.2f", instance.PricePerHour)
	}

	// Check provider data
	if instance.ProviderData["ssh_host"] != "ssh5.vast.ai" {
		t.Errorf("Expected SSH host in provider data to be 'ssh5.vast.ai', got %v", instance.ProviderData["ssh_host"])
	}

	if instance.ProviderData["ssh_port"] != 12345 {
		t.Errorf("Expected SSH port in provider data to be 12345, got %v", instance.ProviderData["ssh_port"])
	}
}

func TestSearchOffersRequest(t *testing.T) {
	req := &SearchOffersRequest{
		GPUName:       "RTX 4090",
		MinGPUCount:   1,
		MaxPrice:      2.0,
		MinRAM:        16,
		Datacenter:    "US-East",
		AvailableOnly: true,
	}

	if req.GPUName != "RTX 4090" {
		t.Errorf("Expected GPU name to be 'RTX 4090', got %s", req.GPUName)
	}

	if req.MinGPUCount != 1 {
		t.Errorf("Expected min GPU count to be 1, got %d", req.MinGPUCount)
	}

	if !req.AvailableOnly {
		t.Error("Expected AvailableOnly to be true")
	}
}

func TestCreateInstanceRequest(t *testing.T) {
	req := &CreateInstanceRequest{
		OfferID:       12345,
		Price:         1.50,
		DiskSizeGB:    100,
		Image:         "pytorch/pytorch:latest",
		Label:         "Test Instance",
		OnStartScript: "#!/bin/bash\necho 'Starting'",
	}

	if req.OfferID != 12345 {
		t.Errorf("Expected offer ID to be 12345, got %d", req.OfferID)
	}

	if req.Price != 1.50 {
		t.Errorf("Expected price to be 1.50, got %.2f", req.Price)
	}

	if req.Image != "pytorch/pytorch:latest" {
		t.Errorf("Expected image to be 'pytorch/pytorch:latest', got %s", req.Image)
	}
}
