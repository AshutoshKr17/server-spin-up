package api

import (
	"net/http"
	"strconv"

	"gpu-cloud-manager/internal/services"
	"gpu-cloud-manager/pkg/types"

	"github.com/gin-gonic/gin"
)

// GPUHandler handles all GPU-related HTTP requests
type GPUHandler struct {
	gpuService *services.GPUService
}

// NewGPUHandler creates a new GPU handler
func NewGPUHandler(gpuService *services.GPUService) *GPUHandler {
	return &GPUHandler{
		gpuService: gpuService,
	}
}

// SearchOffers searches for available GPU offers
// @Summary Search for available GPU offers
// @Description Search for available GPU instances across different providers
// @Tags GPU
// @Accept json
// @Produce json
// @Param provider query string false "GPU provider (vast_ai)"
// @Param gpu_model query string false "GPU model filter"
// @Param min_gpu_count query int false "Minimum GPU count"
// @Param max_price query number false "Maximum price per hour"
// @Param region query string false "Region filter"
// @Param available query bool false "Show only available instances"
// @Success 200 {object} types.APIResponse{data=[]types.GPUInstance}
// @Failure 400 {object} types.APIResponse
// @Failure 500 {object} types.APIResponse
// @Router /api/v1/offers/search [get]
func (h *GPUHandler) SearchOffers(c *gin.Context) {
	var filter types.SearchFilter
	
	// Parse query parameters
	if provider := c.Query("provider"); provider != "" {
		filter.Provider = types.GPUProvider(provider)
	}
	if gpuModel := c.Query("gpu_model"); gpuModel != "" {
		filter.GPUModel = gpuModel
	}
	if minGPUCount := c.Query("min_gpu_count"); minGPUCount != "" {
		if count, err := strconv.Atoi(minGPUCount); err == nil {
			filter.MinGPUCount = count
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filter.MaxPrice = price
		}
	}
	if region := c.Query("region"); region != "" {
		filter.Region = region
	}
	if available := c.Query("available"); available != "" {
		if avail, err := strconv.ParseBool(available); err == nil {
			filter.Available = avail
		}
	}
	
	offers, err := h.gpuService.SearchOffers(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Message: "Offers retrieved successfully",
		Data:    offers,
	})
}

// GetInstances retrieves all instances for the user
// @Summary Get user's GPU instances
// @Description Retrieve all GPU instances owned by the user
// @Tags GPU
// @Produce json
// @Success 200 {object} types.APIResponse{data=[]types.GPUInstance}
// @Failure 500 {object} types.APIResponse
// @Router /api/v1/instances [get]
func (h *GPUHandler) GetInstances(c *gin.Context) {
	instances, err := h.gpuService.GetInstances()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Message: "Instances retrieved successfully",
		Data:    instances,
	})
}

// GetInstance retrieves a specific instance by ID
// @Summary Get GPU instance details
// @Description Retrieve details of a specific GPU instance
// @Tags GPU
// @Produce json
// @Param id path string true "Instance ID"
// @Success 200 {object} types.APIResponse{data=types.GPUInstance}
// @Failure 404 {object} types.APIResponse
// @Failure 500 {object} types.APIResponse
// @Router /api/v1/instances/{id} [get]
func (h *GPUHandler) GetInstance(c *gin.Context) {
	instanceID := c.Param("id")
	
	instance, err := h.gpuService.GetInstance(instanceID)
	if err != nil {
		c.JSON(http.StatusNotFound, types.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Message: "Instance retrieved successfully",
		Data:    instance,
	})
}

// CreateInstance creates a new GPU instance
// @Summary Create a new GPU instance
// @Description Create and launch a new GPU instance
// @Tags GPU
// @Accept json
// @Produce json
// @Param request body types.CreateInstanceRequest true "Instance creation request"
// @Success 201 {object} types.APIResponse{data=types.GPUInstance}
// @Failure 400 {object} types.APIResponse
// @Failure 500 {object} types.APIResponse
// @Router /api/v1/instances [post]
func (h *GPUHandler) CreateInstance(c *gin.Context) {
	var req types.CreateInstanceRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}
	
	instance, err := h.gpuService.CreateInstance(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated, types.APIResponse{
		Success: true,
		Message: "Instance created successfully",
		Data:    instance,
	})
}

// DestroyInstance terminates a GPU instance
// @Summary Destroy a GPU instance
// @Description Permanently terminate a GPU instance
// @Tags GPU
// @Produce json
// @Param id path string true "Instance ID"
// @Success 200 {object} types.APIResponse
// @Failure 404 {object} types.APIResponse
// @Failure 500 {object} types.APIResponse
// @Router /api/v1/instances/{id} [delete]
func (h *GPUHandler) DestroyInstance(c *gin.Context) {
	instanceID := c.Param("id")
	
	err := h.gpuService.DestroyInstance(instanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Message: "Instance destroyed successfully",
	})
}

// StartInstance starts a stopped instance
// @Summary Start a GPU instance
// @Description Start a stopped GPU instance
// @Tags GPU
// @Produce json
// @Param id path string true "Instance ID"
// @Success 200 {object} types.APIResponse
// @Failure 404 {object} types.APIResponse
// @Failure 500 {object} types.APIResponse
// @Router /api/v1/instances/{id}/start [post]
func (h *GPUHandler) StartInstance(c *gin.Context) {
	instanceID := c.Param("id")
	
	err := h.gpuService.StartInstance(instanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Message: "Instance started successfully",
	})
}

// StopInstance stops a running instance
// @Summary Stop a GPU instance
// @Description Stop a running GPU instance
// @Tags GPU
// @Produce json
// @Param id path string true "Instance ID"
// @Success 200 {object} types.APIResponse
// @Failure 404 {object} types.APIResponse
// @Failure 500 {object} types.APIResponse
// @Router /api/v1/instances/{id}/stop [post]
func (h *GPUHandler) StopInstance(c *gin.Context) {
	instanceID := c.Param("id")
	
	err := h.gpuService.StopInstance(instanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Message: "Instance stopped successfully",
	})
}

// GetProviders returns supported GPU providers
// @Summary Get supported providers
// @Description Get list of supported GPU cloud providers
// @Tags GPU
// @Produce json
// @Success 200 {object} types.APIResponse{data=[]types.GPUProvider}
// @Router /api/v1/providers [get]
func (h *GPUHandler) GetProviders(c *gin.Context) {
	providers := h.gpuService.GetSupportedProviders()
	
	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Message: "Providers retrieved successfully",
		Data:    providers,
	})
}
