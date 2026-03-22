package routes

import (
	"github.com/Sidi1901/urlShortner/internal/handler"
	"github.com/gin-gonic/gin"
	
)

func SetupRoutes(r *gin.Engine){

	r.GET("/:url", handler.ResolveURL)
	r.POST("/api/v1",handler.ShortenURL)
}