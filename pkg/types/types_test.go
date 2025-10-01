package types

import (
	"testing"
	"time"
)

func TestGPUInstanceCreation(t *testing.T) {
	instance := GPUInstance{
		ID:           "test_123",
		Provider:     VastAI,
		ProviderID:   "123",
		Name:         "Test Instance",
		Status:       StatusRunning,
		GPUModel:     "RTX 4090",
		GPUCount:     1,
		CPUCount:     8,
		RAM:          32,
		Storage:      100,
		PricePerHour: 1.50,
		Region:       "US-East",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if instance.ID != "test_123" {
		t.Errorf("Expected ID to be 'test_123', got %s", instance.ID)
	}

	if instance.Provider != VastAI {
		t.Errorf("Expected provider to be VastAI, got %s", instance.Provider)
	}

	if instance.Status != StatusRunning {
		t.Errorf("Expected status to be StatusRunning, got %s", instance.Status)
	}
}

func TestInstanceStatus(t *testing.T) {
	testCases := []struct {
		status   InstanceStatus
		expected string
	}{
		{StatusOffline, "offline"},
		{StatusRunning, "running"},
		{StatusLoading, "loading"},
		{StatusRented, "rented"},
		{StatusUnavailable, "unavailable"},
	}

	for _, tc := range testCases {
		if string(tc.status) != tc.expected {
			t.Errorf("Expected status %s, got %s", tc.expected, string(tc.status))
		}
	}
}

func TestSearchFilter(t *testing.T) {
	filter := SearchFilter{
		Provider:    VastAI,
		GPUModel:    "RTX 4090",
		MinGPUCount: 1,
		MaxPrice:    2.0,
		Region:      "US-East",
		Available:   true,
	}

	if filter.Provider != VastAI {
		t.Errorf("Expected provider to be VastAI, got %s", filter.Provider)
	}

	if filter.MinGPUCount != 1 {
		t.Errorf("Expected MinGPUCount to be 1, got %d", filter.MinGPUCount)
	}
}

func TestCreateInstanceRequest(t *testing.T) {
	req := CreateInstanceRequest{
		Provider:      VastAI,
		OfferID:       "12345",
		Image:         "pytorch/pytorch:latest",
		OnStartScript: "#!/bin/bash\necho 'Starting'",
		SSHKey:        "ssh-rsa AAAAB3...",
		Label:         "Test Instance",
	}

	if req.Provider != VastAI {
		t.Errorf("Expected provider to be VastAI, got %s", req.Provider)
	}

	if req.OfferID != "12345" {
		t.Errorf("Expected OfferID to be '12345', got %s", req.OfferID)
	}

	if req.Image != "pytorch/pytorch:latest" {
		t.Errorf("Expected image to be 'pytorch/pytorch:latest', got %s", req.Image)
	}
}

func TestAPIResponse(t *testing.T) {
	response := APIResponse{
		Success: true,
		Message: "Operation successful",
		Data:    map[string]string{"key": "value"},
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	if response.Message != "Operation successful" {
		t.Errorf("Expected message to be 'Operation successful', got %s", response.Message)
	}
}

func TestPagination(t *testing.T) {
	pagination := Pagination{
		Page:       1,
		Limit:      10,
		Total:      100,
		TotalPages: 10,
	}

	if pagination.Page != 1 {
		t.Errorf("Expected page to be 1, got %d", pagination.Page)
	}

	if pagination.TotalPages != 10 {
		t.Errorf("Expected total pages to be 10, got %d", pagination.TotalPages)
	}
}
