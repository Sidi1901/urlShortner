package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sidi1901/urlShortner/internal/logger"
	"github.com/Sidi1901/urlShortner/internal/model"

	"github.com/jmoiron/sqlx"
)

type ShortURLRepository interface {
	SaveShortCode(ctx context.Context, shortURL *model.ShortURL) error
	GetShortCode(ctx context.Context, code string) (*model.ShortURL, error)
	UpdateShortCode(ctx context.Context, shortURL *model.ShortURL) error
	DeleteShortCode(ctx context.Context, shortCode string) error
}

type shortURLRepository struct {
	db *sqlx.DB
}

func NewShortURLRepository(db *sqlx.DB) ShortURLRepository {
	return &shortURLRepository{db: db}
}

// ---------------- CREATE ----------------
func (r *shortURLRepository) SaveShortCode(ctx context.Context, sURL *model.ShortURL) error {
	query := `INSERT INTO url_shortner.short_urls
	(short_code, original_url, created_at, expiry_duration, ip_address, is_active, user_id) VALUES (:short_code, :original_url, :created_at, :expiry_duration, :ip_address, :is_active, :user_id)`

	_, err := r.db.NamedExecContext(ctx, query, sURL)
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"shortcode": sURL.ShortCode,
			"error":     err.Error(),
		}).Error("Failed to save short URL")
		return fmt.Errorf("failed to save short code: %w", err)
	}

	logger.Log.WithFields(map[string]interface{}{
		"shortcode": sURL.ShortCode,
	}).Info("Short code saved successfully")

	return nil
}

// ---------------- READ ----------------
func (r *shortURLRepository) GetShortCode(ctx context.Context, code string) (*model.ShortURL, error) {
	query := `SELECT short_code, original_url, created_at, expiry_duration, ip_address, is_active, user_id
	FROM url_shortner.short_urls
	WHERE short_code = :short_code`

	params := map[string]interface{}{
		"short_code": code,
	}

	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("short code not found")
	}

	var shortURL model.ShortURL
	if err := rows.StructScan(&shortURL); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	logger.Log.WithFields(map[string]interface{}{
		"shortcode": shortURL.ShortCode,
	}).Info("Short code retrieved successfully")

	return &shortURL, nil
}

// ---------------- UPDATE ----------------
func (r *shortURLRepository) UpdateShortCode(ctx context.Context, sURL *model.ShortURL) error {
	query := `UPDATE url_shortner.short_urls SET
	original_url = :original_url,
	expiry_duration = :expiry_duration,
	is_active = :is_active,
	updated_at = :updated_at
	WHERE short_code = :short_code`

	_, err := r.db.NamedExecContext(ctx, query, sURL)
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"shortcode": sURL.ShortCode,
			"error":     err.Error(),
		}).Error("Failed to update short URL")
		return err
	}

	logger.Log.WithFields(map[string]interface{}{
		"shortcode": sURL.ShortCode,
	}).Info("Short URL updated successfully")

	return nil
}

// ---------------- DELETE ----------------
func (r *shortURLRepository) DeleteShortCode(ctx context.Context, code string) error {
	query := `DELETE FROM url_shortner.short_urls WHERE short_code = :short_code`

	params := map[string]interface{}{
		"short_code": code,
	}

	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"shortcode": code,
			"error":     err.Error(),
		}).Error("Failed to delete short URL")
		return err
	}

	logger.Log.WithFields(map[string]interface{}{
		"shortcode": code,
	}).Info("Short URL deleted successfully")

	return nil
}
