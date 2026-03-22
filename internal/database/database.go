package database

import (
	"database/sql"
	"log"
	"fmt"
	"github.com/Sidi1901/urlShortner/internal/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB


func ConnectDB(cfg *config.Config) {

	fmt.Println("Connecting to PostgreSQL database...")
	fmt.Printf("DB Config: Host=%s Port=%s User=%s DBName=%s SSLMode=%s\n", cfg.DBHost, cfg.DBPort, cfg.Username, cfg.DBName, cfg.SSLMode)

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.DBHost, cfg.DBPort, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Database Connection error: ", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Datase unreachable: ", err)
	}


	DB = db

	log.Println("Connected to PostgreSQL")
}

