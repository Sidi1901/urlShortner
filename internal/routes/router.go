package routes

import (
	"github.com/Sidi1901/urlShortner/internal/handler"
	"github.com/Sidi1901/urlShortner/internal/middleware"
	"github.com/Sidi1901/urlShortner/internal/ratelimiter"
	"github.com/gin-gonic/gin"
	
)

func SetupRoutes(r *gin.Engine){

	// Public redirect
	r.GET("/:shortcode", handler.ResolveURL)

	// API Group
	api := r.Group("/api/v1")
	{
		urls := api.Group("/urls")
		{
			urls.POST("", handler.CreateShortURL)
			urls.GET("/:shortCode", handler.GetURLDetails)
			urls.DELETE("/:shortCode", handler.DeleteURL)
			urls.PUT("/:shortCode", handler.UpdateURL)
			urls.GET("/:shortCode/stats", handler.GetStats)
		}

	}

	// Health

	router.GET("/health",handler.HealthCheck)
}