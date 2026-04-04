package routes

import (
	"github.com/Sidi1901/urlShortner/internal/handler"
	// "github.com/Sidi1901/urlShortner/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *handler.Handler) {

	// Public redirect
	r.GET("/:shortcode", h.ResolveURL)

	user := r.Group("/user")
	{
		user.POST("/signup", h.Signup)
		user.POST("/login", h.Login)
		user.POST("/refresh", h.RefreshToken)
	}

	// API Group
	api := r.Group("/api/v1")
	// api.Use(middleware.Authenticate())
	{
		urls := api.Group("/urls")
		{
			urls.POST("", h.CreateShortURL)
			urls.GET("/:shortcode", h.GetShortURL)
			urls.DELETE("/:shortcode", h.DeleteShortURL)
			urls.PUT("", h.UpdateShortURLInfo)
			// urls.GET("/:shortCode/stats", h.GetStats)
		}

	}

	// Health Check
	r.GET("/health", h.HealthCheck)
}
