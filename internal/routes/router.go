package routes

import (
	"github.com/Sidi1901/urlShortner/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *handler.Handler) {

	// Public redirect
	r.GET("/:shortcode", h.ResolveURL)

	// API Group
	api := r.Group("/api/v1")
	{
		urls := api.Group("/urls")
		{
			urls.POST("", h.CreateShortURL)
			// urls.GET("/:shortCode", h.GetURLDetails)
			// urls.DELETE("/:shortCode", h.DeleteURL)
			// urls.PUT("/:shortCode", h.UpdateURL)
			// urls.GET("/:shortCode/stats", h.GetStats)
		}

	}

	// Health Check
	r.GET("/health", h.HealthCheck)
}
