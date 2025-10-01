package services

import (
	"fmt"
	"strconv"

	"gpu-cloud-manager/internal/config"
	"gpu-cloud-manager/pkg/types"
	"gpu-cloud-manager/pkg/vastai"
	"gorm.io/gorm"
)

// GPUService handles all GPU-related business logic
type GPUService struct {
	db         *gorm.DB
	config     *config.Config
	vastClient *vastai.Client
}

// NewGPUService creates a new GPU service
func NewGPUService(db *gorm.DB, cfg *config.Config) *GPUService {
	var vastClient *vastai.Client
	if cfg.VastAIAPIKey != "" {
		vastClient = vastai.NewClient(cfg.VastAIAPIKey)
	}

	return &GPUService{
		db:         db,
		config:     cfg,
		vastClient: vastClient,
	}
}

// SearchOffers searches for available GPU offers across providers
func (s *GPUService) SearchOffers(filter *types.SearchFilter) ([]types.GPUInstance, error) {
	var allOffers []types.GPUInstance
	
	// Search Vast.ai offers
	if s.vastClient != nil && (filter.Provider == "" || filter.Provider == types.VastAI) {
		vastFilter := &vastai.SearchOffersRequest{
			AvailableOnly: filter.Available,
		}
		
		if filter.GPUModel != "" {
			vastFilter.GPUName = filter.GPUModel
		}
		if filter.MinGPUCount > 0 {
			vastFilter.MinGPUCount = filter.MinGPUCount
		}
		if filter.MaxPrice > 0 {
			vastFilter.MaxPrice = filter.MaxPrice
		}
		
		offers, err := s.vastClient.SearchOffers(vastFilter)
		if err != nil {
			return nil, fmt.Errorf("error searching Vast.ai offers: %v", err)
		}
		
		// Convert to our internal type
		for _, offer := range offers {
			allOffers = append(allOffers, vastai.ConvertOfferToGPUInstance(offer))
		}
	}
	
	// TODO: Add other providers (RunPod, Lambda Labs, etc.)
	
	return allOffers, nil
}

// GetInstances retrieves all instances for the user
func (s *GPUService) GetInstances() ([]types.GPUInstance, error) {
	var allInstances []types.GPUInstance
	
	// Get Vast.ai instances
	if s.vastClient != nil {
		instances, err := s.vastClient.GetInstances()
		if err != nil {
			return nil, fmt.Errorf("error getting Vast.ai instances: %v", err)
		}
		
		// Convert to our internal type
		for _, instance := range instances {
			allInstances = append(allInstances, vastai.ConvertInstanceToGPUInstance(instance))
		}
	}
	
	// TODO: Add other providers
	
	return allInstances, nil
}

// CreateInstance creates a new GPU instance
func (s *GPUService) CreateInstance(req *types.CreateInstanceRequest) (*types.GPUInstance, error) {
	switch req.Provider {
	case types.VastAI:
		if s.vastClient == nil {
			return nil, fmt.Errorf("Vast.ai client not configured")
		}
		
		offerID, err := strconv.Atoi(req.OfferID)
		if err != nil {
			return nil, fmt.Errorf("invalid offer ID: %v", err)
		}
		
		vastReq := &vastai.CreateInstanceRequest{
			OfferID:       offerID,
			Price:         0, // Use default price from offer
			DiskSizeGB:    10, // Default disk size
			Image:         req.Image,
			Label:         req.Label,
			OnStartScript: req.OnStartScript,
		}
		
		if vastReq.Image == "" {
			vastReq.Image = "pytorch/pytorch:latest" // Default image
		}
		
		instance, err := s.vastClient.CreateInstance(vastReq)
		if err != nil {
			return nil, fmt.Errorf("error creating Vast.ai instance: %v", err)
		}
		
		result := vastai.ConvertInstanceToGPUInstance(*instance)
		return &result, nil
		
	default:
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}
}

// DestroyInstance terminates an instance
func (s *GPUService) DestroyInstance(instanceID string) error {
	// Parse provider from instance ID
	provider, providerID, err := s.parseInstanceID(instanceID)
	if err != nil {
		return err
	}
	
	switch provider {
	case types.VastAI:
		if s.vastClient == nil {
			return fmt.Errorf("Vast.ai client not configured")
		}
		
		id, err := strconv.Atoi(providerID)
		if err != nil {
			return fmt.Errorf("invalid provider ID: %v", err)
		}
		
		return s.vastClient.DestroyInstance(id)
		
	default:
		return fmt.Errorf("unsupported provider: %s", provider)
	}
}

// StartInstance starts a stopped instance
func (s *GPUService) StartInstance(instanceID string) error {
	provider, providerID, err := s.parseInstanceID(instanceID)
	if err != nil {
		return err
	}
	
	switch provider {
	case types.VastAI:
		if s.vastClient == nil {
			return fmt.Errorf("Vast.ai client not configured")
		}
		
		id, err := strconv.Atoi(providerID)
		if err != nil {
			return fmt.Errorf("invalid provider ID: %v", err)
		}
		
		return s.vastClient.StartInstance(id)
		
	default:
		return fmt.Errorf("unsupported provider: %s", provider)
	}
}

// StopInstance stops a running instance
func (s *GPUService) StopInstance(instanceID string) error {
	provider, providerID, err := s.parseInstanceID(instanceID)
	if err != nil {
		return err
	}
	
	switch provider {
	case types.VastAI:
		if s.vastClient == nil {
			return fmt.Errorf("Vast.ai client not configured")
		}
		
		id, err := strconv.Atoi(providerID)
		if err != nil {
			return fmt.Errorf("invalid provider ID: %v", err)
		}
		
		return s.vastClient.StopInstance(id)
		
	default:
		return fmt.Errorf("unsupported provider: %s", provider)
	}
}

// GetInstance retrieves details of a specific instance
func (s *GPUService) GetInstance(instanceID string) (*types.GPUInstance, error) {
	provider, providerID, err := s.parseInstanceID(instanceID)
	if err != nil {
		return nil, err
	}
	
	switch provider {
	case types.VastAI:
		if s.vastClient == nil {
			return nil, fmt.Errorf("Vast.ai client not configured")
		}
		
		id, err := strconv.Atoi(providerID)
		if err != nil {
			return nil, fmt.Errorf("invalid provider ID: %v", err)
		}
		
		instance, err := s.vastClient.GetInstance(id)
		if err != nil {
			return nil, fmt.Errorf("error getting Vast.ai instance: %v", err)
		}
		
		result := vastai.ConvertInstanceToGPUInstance(*instance)
		return &result, nil
		
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// parseInstanceID parses our internal instance ID format (provider_id)
func (s *GPUService) parseInstanceID(instanceID string) (types.GPUProvider, string, error) {
	if len(instanceID) < 5 {
		return "", "", fmt.Errorf("invalid instance ID format")
	}
	
	if instanceID[:5] == "vast_" {
		return types.VastAI, instanceID[5:], nil
	}
	
	// TODO: Add other providers
	
	return "", "", fmt.Errorf("unknown provider in instance ID")
}

// GetSupportedProviders returns a list of configured providers
func (s *GPUService) GetSupportedProviders() []types.GPUProvider {
	var providers []types.GPUProvider
	
	if s.vastClient != nil {
		providers = append(providers, types.VastAI)
	}
	
	// TODO: Add other providers when configured
	
	return providers
}
