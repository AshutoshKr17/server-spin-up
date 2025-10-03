package api

import (
	"gpu-cloud-manager/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, gpuService *services.GPUService) {
	// Create handlers
	gpuHandler := NewGPUHandler(gpuService)
	
	// API version 1
	v1 := router.Group("/api/v1")
	{
		// GPU Offers routes
		offers := v1.Group("/offers")
		{
			offers.GET("/search", gpuHandler.SearchOffers)
			offers.POST("/search/advanced", gpuHandler.SearchOffersAdvanced)
		}
		
		// GPU Instances routes
		instances := v1.Group("/instances")
		{
			instances.GET("", gpuHandler.GetInstances)
			instances.POST("", gpuHandler.CreateInstance)
			instances.GET("/:id", gpuHandler.GetInstance)
			instances.DELETE("/:id", gpuHandler.DestroyInstance)
			instances.POST("/:id/start", gpuHandler.StartInstance)
			instances.POST("/:id/stop", gpuHandler.StopInstance)
		}
		
		// Providers and Models routes
		v1.GET("/providers", gpuHandler.GetProviders)
		v1.GET("/gpu-models", gpuHandler.GetGPUModels)
		
		// Marketplace routes
		marketplace := v1.Group("/marketplace")
		{
			marketplace.GET("/stats", gpuHandler.GetMarketplaceStats)
		}
	}
	
	// CORS middleware (if enabled)
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
}
