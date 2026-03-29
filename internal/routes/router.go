package routes

import (
	"github.com/Sidi1901/urlShortner/internal/handler"
	"github.com/Sidi1901/urlShortner/internal/middleware"
	"github.com/Sidi1901/urlShortner/internal/ratelimiter"
	"github.com/gin-gonic/gin"
	
)

func SetupRoutes(r *gin.Engine){

	r.GET("/:url", handler.ResolveURL)

	api := r.Group("/api")
	api.Use(middleware.RateLimit(limiter))
	api.POST("/api/v1",handler.ShortenURL)
}