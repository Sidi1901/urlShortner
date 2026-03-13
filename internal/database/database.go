package database

import (
	"database/sql"
	"log"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB


func ConnectDB(cfg config) {

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.host, cfg.port, cfg.username, cfg.password, cfg.dbname, cfg.sslmode)

	db, err := sql.Open("postgres", connstr)

	if err != nil {
		log.Fatal("Database Connection error: ", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Datase unreachable: ", err)
	}


	DB = db

	log.Println(Connected to PostgreSQL)
}

