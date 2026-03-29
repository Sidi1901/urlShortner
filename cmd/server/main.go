package main

import(
	"fmt"
	"log"
	"github.com/Sidi1901/urlShortner/internal/routes"
	"github.com/Sidi1901/urlShortner/internal/database"
	"github.com/Sidi1901/urlShortner/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	
)

func main(){
	
	cfg := config.Config{}

	fmt.Println("Loading config from environment variables...")

	godotenv.Load()

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Unable to load config: %v", err)
	}

	database.ConnectDB(&cfg)

	r := gin.Default()

	routes.SetupRoutes(r)

	fmt.Println("Server is running on http://localhost:" + cfg.APPPORT)

	r.Run(":"+cfg.APPPORT)

	

}