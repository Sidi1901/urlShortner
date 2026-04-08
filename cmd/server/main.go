package main

import (
	"log"

	"github.com/Sidi1901/urlShortner/internal/config"
	"github.com/Sidi1901/urlShortner/internal/database"
	"github.com/Sidi1901/urlShortner/internal/handler"
	"github.com/Sidi1901/urlShortner/internal/infra"
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

	PostgresDB := database.ConnectDB(&cfg)
	if PostgresDB == nil {
		log.Fatal("Failed to connect to Postges Database")
	}

	RedisClient := database.NewRedisClient(&cfg)

	if RedisClient == nil {
		log.Fatal("Failed to connect to Redis Databse")
	}

	RateLimiter := infra.NewRateLimiter(RedisClient.Client, infra.TokenBucketLua)

	rateLimitMiddleware := middleware.NewRateLimitMiddlware(RateLimiter)
	loggerMiddleware := middleware.NewLoggerMiddleware("info")
	authMiddleware := middleware.NewAuthMiddleware(&cfg)

	shortURLRepo := repository.NewShortURLRepository(PostgresDB.DB)
	userRepo := repository.NewUserRepository(PostgresDB.DB)

	shortURLService := service.NewShortURLService(shortURLRepo, userRepo, &cfg)
	userService := service.NewUserService(userRepo, &cfg)

	shortURLHandler := handler.NewURLHandler(shortURLService)
	userHandler := handler.NewUserHandler(userService)
	healthHandler := handler.NewHealthHandler()

	r := gin.New()
	r.Use(loggerMiddleware.Logger())

	// Serve static files
	r.Static("/static", "./web/static")

	// Load HTML templates
	r.LoadHTMLGlob("web/templates/*")

	routeReg := []routes.RouteRegistrar{
		userHandler,
		shortURLHandler,
		healthHandler,
	}

	routes.SetupRoutes(r, rateLimitMiddleware, authMiddleware, routeReg)

	r.Run(":" + cfg.AppPort)

}
