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

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

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

