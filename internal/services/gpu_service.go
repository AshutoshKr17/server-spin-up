package services

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"gpu-cloud-manager/internal/config"
	"gpu-cloud-manager/pkg/runpod"
	"gpu-cloud-manager/pkg/types"
	"gpu-cloud-manager/pkg/vastai"
	"gorm.io/gorm"
)

// GPUService handles all GPU-related business logic
type GPUService struct {
	db           *gorm.DB
	config       *config.Config
	vastClient   *vastai.Client
	runpodClient *runpod.Client
}

// NewGPUService creates a new GPU service
func NewGPUService(db *gorm.DB, cfg *config.Config) *GPUService {
	var vastClient *vastai.Client
	var runpodClient *runpod.Client

	if cfg.VastAIAPIKey != "" {
		vastClient = vastai.NewClient(cfg.VastAIAPIKey)
	}

	if cfg.RunPodAPIKey != "" {
		runpodClient = runpod.NewClient(cfg.RunPodAPIKey)
	}

	return &GPUService{
		db:           db,
		config:       cfg,
		vastClient:   vastClient,
		runpodClient: runpodClient,
	}
}

// SearchOffers searches for available GPU offers across providers with advanced filtering
func (s *GPUService) SearchOffers(filter *types.SearchFilter) ([]types.GPUInstance, error) {
	// Convert basic filter to advanced filter for backward compatibility
	advancedFilter := &types.AdvancedSearchFilter{
		Provider:    filter.Provider,
		GPUModel:    filter.GPUModel,
		MinGPUCount: filter.MinGPUCount,
		MaxPrice:    filter.MaxPrice,
		Region:      filter.Region,
		Available:   filter.Available,
		SortBy:      "price",
		SortOrder:   "asc",
	}

	return s.SearchOffersAdvanced(advancedFilter)
}

// SearchOffersAdvanced searches for available GPU offers with advanced filtering
func (s *GPUService) SearchOffersAdvanced(filter *types.AdvancedSearchFilter) ([]types.GPUInstance, error) {
	var allOffers []types.GPUInstance

	// Search Vast.ai offers
	if s.vastClient != nil && (filter.Provider == "" || filter.Provider == types.VastAI) {
		vastOffers, err := s.searchVastAI(filter)
		if err != nil {
			return nil, fmt.Errorf("error searching Vast.ai offers: %v", err)
		}
		allOffers = append(allOffers, vastOffers...)
	}

	// Search RunPod offers
	if s.runpodClient != nil && (filter.Provider == "" || filter.Provider == types.RunPod) {
		runpodOffers, err := s.searchRunPod(filter)
		if err != nil {
			return nil, fmt.Errorf("error searching RunPod offers: %v", err)
		}
		allOffers = append(allOffers, runpodOffers...)
	}

	// Apply advanced filters
	allOffers = s.applyAdvancedFilters(allOffers, filter)

	// Sort results
	s.sortOffers(allOffers, filter.SortBy, filter.SortOrder)

	return allOffers, nil
}

// searchVastAI searches Vast.ai for offers
func (s *GPUService) searchVastAI(filter *types.AdvancedSearchFilter) ([]types.GPUInstance, error) {
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
	if filter.MinRAM > 0 {
		vastFilter.MinRAM = filter.MinRAM
	}
	if filter.Region != "" {
		vastFilter.Datacenter = filter.Region
	}

	offers, err := s.vastClient.SearchOffers(vastFilter)
	if err != nil {
		return nil, err
	}

	var instances []types.GPUInstance
	for _, offer := range offers {
		instance := vastai.ConvertOfferToGPUInstance(offer)
		// Enhance with GPU model information
		if gpuInfo, exists := types.GPUModels[instance.GPUModel]; exists {
			instance.GPUInfo = &gpuInfo
			instance.Performance = gpuInfo.Performance
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// searchRunPod searches RunPod for offers
func (s *GPUService) searchRunPod(filter *types.AdvancedSearchFilter) ([]types.GPUInstance, error) {
	gpuTypes, err := s.runpodClient.GetGPUTypes()
	if err != nil {
		return nil, err
	}

	var instances []types.GPUInstance
	for _, gpuType := range gpuTypes {
		instance := runpod.ConvertGPUTypeToGPUInstance(gpuType)
		
		// Enhance with GPU model information
		if gpuInfo, exists := types.GPUModels[instance.GPUModel]; exists {
			instance.GPUInfo = &gpuInfo
			instance.Performance = gpuInfo.Performance
		}
		
		instances = append(instances, instance)
	}

	return instances, nil
}

// applyAdvancedFilters applies advanced filtering to the results
func (s *GPUService) applyAdvancedFilters(offers []types.GPUInstance, filter *types.AdvancedSearchFilter) []types.GPUInstance {
	var filtered []types.GPUInstance

	for _, offer := range offers {
		// GPU Models filter (multiple models)
		if len(filter.GPUModels) > 0 {
			found := false
			for _, model := range filter.GPUModels {
				if strings.Contains(strings.ToLower(offer.GPUModel), strings.ToLower(model)) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// GPU Category filter
		if filter.GPUCategory != "" && offer.GPUInfo != nil {
			if types.GPUCategory(offer.GPUInfo.Category) != filter.GPUCategory {
				continue
			}
		}

		// GPU Count range
		if filter.MaxGPUCount > 0 && offer.GPUCount > filter.MaxGPUCount {
			continue
		}

		// Price range
		if filter.MinPrice > 0 && offer.PricePerHour < filter.MinPrice {
			continue
		}

		// RAM range
		if filter.MaxRAM > 0 && offer.RAM > filter.MaxRAM {
			continue
		}

		// Storage filter
		if filter.MinStorage > 0 && offer.Storage < filter.MinStorage {
			continue
		}

		// Regions filter (multiple regions)
		if len(filter.Regions) > 0 {
			found := false
			for _, region := range filter.Regions {
				if strings.Contains(strings.ToLower(offer.Region), strings.ToLower(region)) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Reliability filter
		if filter.MinReliability > 0 && offer.Reliability < filter.MinReliability {
			continue
		}

		// Performance filter
		if filter.MinPerformance > 0 && offer.Performance < filter.MinPerformance {
			continue
		}

		filtered = append(filtered, offer)
	}

	return filtered
}

// sortOffers sorts the offers based on the specified criteria
func (s *GPUService) sortOffers(offers []types.GPUInstance, sortBy, sortOrder string) {
	if sortBy == "" {
		sortBy = "price"
	}
	if sortOrder == "" {
		sortOrder = "asc"
	}

	sort.Slice(offers, func(i, j int) bool {
		var less bool
		
		switch sortBy {
		case "price":
			less = offers[i].PricePerHour < offers[j].PricePerHour
		case "performance":
			less = offers[i].Performance < offers[j].Performance
		case "reliability":
			less = offers[i].Reliability < offers[j].Reliability
		case "memory":
			less = offers[i].RAM < offers[j].RAM
		case "gpu_count":
			less = offers[i].GPUCount < offers[j].GPUCount
		default:
			less = offers[i].PricePerHour < offers[j].PricePerHour
		}

		if sortOrder == "desc" {
			return !less
		}
		return less
	})
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

		for _, instance := range instances {
			gpuInstance := vastai.ConvertInstanceToGPUInstance(instance)
			// Enhance with GPU model information
			if gpuInfo, exists := types.GPUModels[gpuInstance.GPUModel]; exists {
				gpuInstance.GPUInfo = &gpuInfo
				gpuInstance.Performance = gpuInfo.Performance
			}
			allInstances = append(allInstances, gpuInstance)
		}
	}

	// Get RunPod instances
	if s.runpodClient != nil {
		pods, err := s.runpodClient.SearchPods()
		if err != nil {
			return nil, fmt.Errorf("error getting RunPod instances: %v", err)
		}

		for _, pod := range pods {
			gpuInstance := runpod.ConvertPodToGPUInstance(pod)
			// Enhance with GPU model information
			if gpuInfo, exists := types.GPUModels[gpuInstance.GPUModel]; exists {
				gpuInstance.GPUInfo = &gpuInfo
				gpuInstance.Performance = gpuInfo.Performance
			}
			allInstances = append(allInstances, gpuInstance)
		}
	}

	return allInstances, nil
}

// CreateInstance creates a new GPU instance
func (s *GPUService) CreateInstance(req *types.CreateInstanceRequest) (*types.GPUInstance, error) {
	switch req.Provider {
	case types.VastAI:
		return s.createVastAIInstance(req)
	case types.RunPod:
		return s.createRunPodInstance(req)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}
}

// createVastAIInstance creates a Vast.ai instance
func (s *GPUService) createVastAIInstance(req *types.CreateInstanceRequest) (*types.GPUInstance, error) {
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

	// Override with resource requirements
	if req.Resources != nil {
		if req.Resources.MinStorage > 0 {
			vastReq.DiskSizeGB = req.Resources.MinStorage
		}
	}

	instance, err := s.vastClient.CreateInstance(vastReq)
	if err != nil {
		return nil, fmt.Errorf("error creating Vast.ai instance: %v", err)
	}

	result := vastai.ConvertInstanceToGPUInstance(*instance)
	// Enhance with GPU model information
	if gpuInfo, exists := types.GPUModels[result.GPUModel]; exists {
		result.GPUInfo = &gpuInfo
		result.Performance = gpuInfo.Performance
	}

	return &result, nil
}

// createRunPodInstance creates a RunPod instance
func (s *GPUService) createRunPodInstance(req *types.CreateInstanceRequest) (*types.GPUInstance, error) {
	if s.runpodClient == nil {
		return nil, fmt.Errorf("RunPod client not configured")
	}

	runpodReq := &runpod.CreatePodRequest{
		Name:            req.Label,
		ImageName:       req.Image,
		GPUTypeID:       req.OfferID,
		CloudType:       "COMMUNITY", // Default to community cloud
		SupportPublicIp: true,
		StartJupyter:    false,
		StartSsh:        true,
		ContainerDisk:   10, // Default container disk
		VolumeInGb:      0,  // No additional volume by default
		VolumeMountPath: "/workspace",
		Ports:           "22/tcp,8888/tcp", // SSH and Jupyter
	}

	if runpodReq.ImageName == "" {
		runpodReq.ImageName = "pytorch/pytorch:latest"
	}

	// Convert environment variables
	if req.Environment != nil {
		for key, value := range req.Environment {
			runpodReq.Env = append(runpodReq.Env, runpod.EnvVar{
				Key:   key,
				Value: value,
			})
		}
	}

	// Override with resource requirements
	if req.Resources != nil {
		if req.Resources.MinStorage > 0 {
			runpodReq.VolumeInGb = req.Resources.MinStorage
		}
	}

	// Handle ports
	if len(req.Ports) > 0 {
		var ports []string
		for _, port := range req.Ports {
			protocol := "tcp"
			if port.Protocol != "" {
				protocol = port.Protocol
			}
			ports = append(ports, fmt.Sprintf("%d/%s", port.ContainerPort, protocol))
		}
		runpodReq.Ports = strings.Join(ports, ",")
	}

	pod, err := s.runpodClient.CreatePod(runpodReq)
	if err != nil {
		return nil, fmt.Errorf("error creating RunPod instance: %v", err)
	}

	result := runpod.ConvertPodToGPUInstance(*pod)
	// Enhance with GPU model information
	if gpuInfo, exists := types.GPUModels[result.GPUModel]; exists {
		result.GPUInfo = &gpuInfo
		result.Performance = gpuInfo.Performance
	}

	return &result, nil
}

// DestroyInstance terminates an instance
func (s *GPUService) DestroyInstance(instanceID string) error {
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

	case types.RunPod:
		if s.runpodClient == nil {
			return fmt.Errorf("RunPod client not configured")
		}

		return s.runpodClient.TerminatePod(providerID)

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

	case types.RunPod:
		if s.runpodClient == nil {
			return fmt.Errorf("RunPod client not configured")
		}

		return s.runpodClient.ResumePod(providerID)

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

	case types.RunPod:
		if s.runpodClient == nil {
			return fmt.Errorf("RunPod client not configured")
		}

		return s.runpodClient.StopPod(providerID)

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
		// Enhance with GPU model information
		if gpuInfo, exists := types.GPUModels[result.GPUModel]; exists {
			result.GPUInfo = &gpuInfo
			result.Performance = gpuInfo.Performance
		}

		return &result, nil

	case types.RunPod:
		// RunPod doesn't have a direct get instance by ID, so we search through all pods
		pods, err := s.runpodClient.SearchPods()
		if err != nil {
			return nil, fmt.Errorf("error getting RunPod instances: %v", err)
		}

		for _, pod := range pods {
			if pod.ID == providerID {
				result := runpod.ConvertPodToGPUInstance(pod)
				// Enhance with GPU model information
				if gpuInfo, exists := types.GPUModels[result.GPUModel]; exists {
					result.GPUInfo = &gpuInfo
					result.Performance = gpuInfo.Performance
				}
				return &result, nil
			}
		}

		return nil, fmt.Errorf("RunPod instance not found: %s", providerID)

	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// parseInstanceID parses our internal instance ID format (provider_id)
func (s *GPUService) parseInstanceID(instanceID string) (types.GPUProvider, string, error) {
	if len(instanceID) < 5 {
		return "", "", fmt.Errorf("invalid instance ID format")
	}

	if strings.HasPrefix(instanceID, "vast_") {
		return types.VastAI, instanceID[5:], nil
	}
	if strings.HasPrefix(instanceID, "runpod_") {
		return types.RunPod, instanceID[7:], nil
	}

	return "", "", fmt.Errorf("unknown provider in instance ID: %s", instanceID)
}

// GetSupportedProviders returns a list of configured providers with details
func (s *GPUService) GetSupportedProviders() []types.ProviderInfo {
	var providers []types.ProviderInfo

	// Vast.ai
	vastInfo := types.ProviderInfo{
		Name:        types.VastAI,
		DisplayName: "Vast.ai",
		Website:     "https://vast.ai",
		Regions:     []string{"US-East", "US-West", "Europe", "Asia"},
		Features:    []string{"SSH Access", "Docker Support", "Jupyter Notebooks", "Custom Images"},
		IsConfigured: s.vastClient != nil,
	}
	
	if s.vastClient != nil {
		// Get available GPU models from static list for now
		for model := range types.GPUModels {
			vastInfo.GPUModels = append(vastInfo.GPUModels, model)
		}
	}
	
	providers = append(providers, vastInfo)

	// RunPod
	runpodInfo := types.ProviderInfo{
		Name:        types.RunPod,
		DisplayName: "RunPod",
		Website:     "https://runpod.io",
		Regions:     []string{"Global", "US", "Europe", "Asia"},
		Features:    []string{"GraphQL API", "Jupyter Support", "SSH Access", "Community & Secure Cloud"},
		IsConfigured: s.runpodClient != nil,
	}
	
	if s.runpodClient != nil {
		// Get available GPU models from static list for now
		for model := range types.GPUModels {
			runpodInfo.GPUModels = append(runpodInfo.GPUModels, model)
		}
	}
	
	providers = append(providers, runpodInfo)

	return providers
}

// GetGPUModels returns information about all available GPU models
func (s *GPUService) GetGPUModels() map[string]types.GPUModel {
	return types.GPUModels
}

// GetMarketplaceStats returns marketplace statistics
func (s *GPUService) GetMarketplaceStats() (*types.MarketplaceStats, error) {
	// This would typically aggregate data from multiple providers
	// For now, we'll return basic statistics
	
	stats := &types.MarketplaceStats{
		ModelStats:    make(map[string]types.GPUModelInfo),
		ProviderStats: make(map[string]int),
	}

	// Get offers from all providers
	filter := &types.AdvancedSearchFilter{
		Available: true,
	}
	
	offers, err := s.SearchOffersAdvanced(filter)
	if err != nil {
		return nil, err
	}

	stats.TotalInstances = len(offers)
	stats.AvailableInstances = len(offers)

	var totalPrice float64
	minPrice := float64(999999)
	maxPrice := float64(0)

	// Aggregate statistics
	modelCounts := make(map[string]int)
	modelPrices := make(map[string][]float64)
	
	for _, offer := range offers {
		// Provider stats
		stats.ProviderStats[string(offer.Provider)]++
		
		// Price stats
		totalPrice += offer.PricePerHour
		if offer.PricePerHour < minPrice {
			minPrice = offer.PricePerHour
		}
		if offer.PricePerHour > maxPrice {
			maxPrice = offer.PricePerHour
		}
		
		// Model stats
		modelCounts[offer.GPUModel]++
		modelPrices[offer.GPUModel] = append(modelPrices[offer.GPUModel], offer.PricePerHour)
	}

	// Calculate model statistics
	for model, count := range modelCounts {
		prices := modelPrices[model]
		var minModelPrice, maxModelPrice, totalModelPrice float64
		
		minModelPrice = prices[0]
		maxModelPrice = prices[0]
		
		for _, price := range prices {
			totalModelPrice += price
			if price < minModelPrice {
				minModelPrice = price
			}
			if price > maxModelPrice {
				maxModelPrice = price
			}
		}
		
		avgModelPrice := totalModelPrice / float64(len(prices))
		
		modelInfo := types.GPUModelInfo{
			Available: count,
			MinPrice:  minModelPrice,
			MaxPrice:  maxModelPrice,
			AvgPrice:  avgModelPrice,
		}
		
		if gpuModel, exists := types.GPUModels[model]; exists {
			modelInfo.Model = gpuModel
		}
		
		// Determine which providers have this model
		for _, offer := range offers {
			if offer.GPUModel == model {
				found := false
				for _, p := range modelInfo.Providers {
					if p == offer.Provider {
						found = true
						break
					}
				}
				if !found {
					modelInfo.Providers = append(modelInfo.Providers, offer.Provider)
				}
			}
		}
		
		stats.ModelStats[model] = modelInfo
	}

	if len(offers) > 0 {
		stats.AveragePrice = totalPrice / float64(len(offers))
		stats.PriceRange = types.PriceRange{
			Min: minPrice,
			Max: maxPrice,
			Avg: stats.AveragePrice,
		}
	}

	return stats, nil
}
