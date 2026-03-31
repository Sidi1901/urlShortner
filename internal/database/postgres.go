package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Sidi1901/urlShortner/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDB(cfg *config.Config) *sqlx.DB {

	fmt.Println("Connecting to PostgreSQL database...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Printf("DB Config: Host=%s Port=%s User=%s DBName=%s SSLMode=%s\n", cfg.DBHost, cfg.DBPort, cfg.Username, cfg.DBName, cfg.SSLMode)

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.DBHost, cfg.DBPort, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Open("postgres", connStr)

	if err != nil {
		fmt.Println("Database Connection error: ", err)
		return db
	}

	if err = db.Ping(); err != nil {
		fmt.Println("Database unreachable: ", err)
		return db
	}

	// 	PingContext is used to check if the database connection is alive within the specified context timeout.

	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Connected to PostgreSQL")

	return db
}
