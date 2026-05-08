package routes

import (
	"4g-proxy/internal/controllers"
	"4g-proxy/internal/models"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the HTTP API routes
func SetupRouter(state *models.ProxyState) *gin.Engine {
	router := gin.Default()

	// Create controller
	controller := controllers.NewMMEController(state)

	// Health check
	router.GET("/health", controller.Health)

	// API group
	api := router.Group("/api/v1")
	{
		// Status and stats
		api.GET("/status", controller.GetStatus)
		api.GET("/stats", controller.GetStats)

		// Drop flags management
		drop := api.Group("/drop")
		{
			// Get all drop flags
			drop.GET("", controller.GetDropFlags)

			// Set multiple drop flags
			drop.PUT("", controller.SetDropFlags)

			// Reset all drop flags
			drop.DELETE("", controller.ResetDropFlags)

			// Set individual drop flag
			drop.PUT("/:signalType", controller.SetDropFlag)

			// Convenience endpoints
			drop.POST("/attach", controller.DropAttach)
			drop.DELETE("/attach", controller.AllowAttach)

			drop.POST("/detach", controller.DropDetach)
			drop.DELETE("/detach", controller.AllowDetach)

			drop.POST("/tau", controller.DropTAU)
			drop.DELETE("/tau", controller.AllowTAU)

			drop.POST("/service-request", controller.DropServiceRequest)
			drop.DELETE("/service-request", controller.AllowServiceRequest)

			drop.POST("/ue-context-release", controller.DropUEContextRelease)
			drop.DELETE("/ue-context-release", controller.AllowUEContextRelease)
		}

		// Delay configuration management
		delay := api.Group("/delay")
		{
			// Get all delay settings
			delay.GET("", controller.GetDelayConfig)

			// Set multiple delays
			delay.PUT("", controller.SetDelays)

			// Reset all delays
			delay.DELETE("", controller.ResetDelays)

			// Set individual delay
			delay.PUT("/:signalType", controller.SetDelay)
		}
	}

	return router
}
