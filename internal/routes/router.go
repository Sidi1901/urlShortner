package routes

import (
	"github.com/Sidi1901/urlShortner/internal/handler"
	"github.com/gin-gonic/gin"
	
)

func setupRoutes(app *gin.Context) *gin.Engine{

	r := gin.Default()

	app.Get("/:url",handler.ResolveURL)
	app.Post("/api/v1",handler.ShortenURL)

	return r;
}