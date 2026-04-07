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

	shortURLRepo := repository.NewShortURLRepository(db)
	userRepo := repository.NewUserRepository(db)

	shortURLService := service.NewShortURLService(shortURLRepo, userRepo, &cfg)
	userService := service.NewUserService(userRepo, &cfg)

	shortURLHandler := handler.NewURLHandler(shortURLService)
	userHandler := handler.NewUserHandler(userService)
	healthHandler := handler.NewHealthHandler()

	middleware := middleware.NewMiddleware(&cfg)

	logger.Init()

	r := gin.New()
	r.Use(middleware.LoggerMiddleware())

	routeReg := []routes.RouteRegistrar{
		userHandler,
		shortURLHandler,
		healthHandler,
	}

	routes.SetupRoutes(r, middleware, routeReg)

	r.Run(":" + cfg.AppPort)

}
