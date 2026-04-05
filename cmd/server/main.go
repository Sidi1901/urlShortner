package main

import (
	"log"

	"github.com/Sidi1901/urlShortner/internal/config"
	"github.com/Sidi1901/urlShortner/internal/database"
	"github.com/Sidi1901/urlShortner/internal/handler"
	"github.com/Sidi1901/urlShortner/internal/logger"
	"github.com/Sidi1901/urlShortner/internal/middleware"
	"github.com/Sidi1901/urlShortner/internal/repository"
	"github.com/Sidi1901/urlShortner/internal/routes"
	"github.com/Sidi1901/urlShortner/internal/service"
	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	cfg := config.Config{}

	godotenv.Load()

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Unable to load config: %v", err)
	}

	db := database.ConnectDB(&cfg)
	if db == nil {
		log.Fatal("Failed to connect to database")
	}

	repo := repository.NewRepository(db, &cfg)
	service := service.NewService(repo, &cfg)
	handler := handler.NewHandler(service, &cfg)

	middleware := middleware.NewMiddleware(&cfg)

	logger.Init()

	r := gin.New()
	r.Use(middleware.LoggerMiddleware())
	routes.SetupRoutes(r, handler, middleware)

	r.Run(":" + cfg.AppPort)

}
