package models

import (
	"reflect"
	"testing"
	"time"

	"gpu-cloud-manager/pkg/types"
)

func TestTableName(t *testing.T) {
	var i Instance
	expected := "instances"
	if i.TableName() != expected {
		t.Errorf("expected %s, got %s", expected, i.TableName())
	}
}

func TestGenerateAPIID(t *testing.T) {
	tests := []struct {
		provider   types.GPUProvider
		providerID string
		expected   string
	}{
		{types.VastAI, "12345", "vast_12345"},
		{types.GCP, "abc", "GCP_abc"},
		{types.AWS, "xyz", "AWS_xyz"},
	}

	for _, tt := range tests {
		i := Instance{
			Provider:   tt.provider,
			ProviderID: tt.providerID,
		}
		got := i.generateAPIID()
		if got != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, got)
		}
	}
}

func TestToGPUInstance(t *testing.T) {
	now := time.Now()
	instance := Instance{
		ID:           1,
		UserID:       100,
		Provider:     types.AWS,
		ProviderID:   "i-123456",
		Name:         "test-instance",
		Status:       types.Running,
		GPUModel:     "NVIDIA A100",
		GPUCount:     2,
		CPUCount:     16,
		RAM:          64,
		Storage:      1000,
		PricePerHour: 2.5,
		Region:       "us-east-1",
		CreatedAt:    now,
		UpdatedAt:    now,
		ProviderData: ProviderData{"key": "value"},
	}

	got := instance.ToGPUInstance()

	if got.ID != "AWS_i-123456" {
		t.Errorf("expected ID %s, got %s", "AWS_i-123456", got.ID)
	}
	if got.Provider != types.AWS {
		t.Errorf("expected provider AWS, got %v", got.Provider)
	}
	if got.ProviderData["key"] != "value" {
		t.Errorf("expected provider_data.key=value, got %v", got.ProviderData["key"])
	}
}

func TestFromGPUInstance(t *testing.T) {
	now := time.Now()
	gpu := types.GPUInstance{
		ID:           "AWS_i-654321",
		Provider:     types.AWS,
		ProviderID:   "i-654321",
		Name:         "new-instance",
		Status:       types.Stopped,
		GPUModel:     "NVIDIA V100",
		GPUCount:     4,
		CPUCount:     32,
		RAM:          128,
		Storage:      2000,
		PricePerHour: 5.0,
		Region:       "us-west-2",
		CreatedAt:    now,
		UpdatedAt:    now,
		ProviderData: map[string]interface{}{"custom": "data"},
	}

	var i Instance
	i.FromGPUInstance(gpu, 200)

	if i.UserID != 200 {
		t.Errorf("expected UserID 200, got %d", i.UserID)
	}
	if i.ProviderID != "i-654321" {
		t.Errorf("expected ProviderID i-654321, got %s", i.ProviderID)
	}
	if !reflect.DeepEqual(i.ProviderData, ProviderData{"custom": "data"}) {
		t.Errorf("expected ProviderData map[custom:data], got %v", i.ProviderData)
	}
}
