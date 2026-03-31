package repository

import (
	"github.com/Sidi1901/urlShortner/internal/config"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db  *sqlx.DB
	cfg *config.Config
}

func NewRepository(db *sqlx.DB, cfg *config.Config) *Repository {
	return &Repository{
		db:  db,
		cfg: cfg,
	}
}
