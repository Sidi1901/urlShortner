package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Sidi1901/urlShortner/internal/config"
	"github.com/Sidi1901/urlShortner/internal/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDB(cfg *config.Config) *sqlx.DB {

	logger.Log.Info("Connecting to PostgreSQL database...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Log.WithFields(map[string]interface{}{
		"host":    cfg.DBHost,
		"port":    cfg.DBPort,
		"user":    cfg.Username,
		"dbname":  cfg.DBName,
		"sslmode": cfg.SSLMode,
	}).Info("DB Config")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.DBHost, cfg.DBPort, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Open("postgres", connStr)

	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Database Connection error")
		return db
	}

	if err = db.Ping(); err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Database unreachable")
	}

	// 	PingContext is used to check if the database connection is alive within the specified context timeout.

	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}

	logger.Log.Info("Connected to PostgreSQL")

	return db
}
