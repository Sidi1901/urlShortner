package main

import(
	"fmt"
	"os"
	"log"
	"github.com/Sidi1901/urlShortner/internal/routes"
	"github.com/Sidi1901/urlShortner/internal/database"
	"github.com/Sidi1901/urlShortner/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/caarlos0/env/v6"
	
)

func main(){
	
	cfg := Config{}

	fmt.Println("Loading config from environment variables...")

	if err := env.Parse(); err != nil {
		log.Fatalf("Unable to load config: %v", err)
	}

	database.ConnectDB(cfg)

	r := router.SetupRouter()

	r.Run(":"+os.Getenv("APP_PORT"))

	fmt.Println("Server is running on http://localhost:" + os.Getenv("APP_PORT"))

}
