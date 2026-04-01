package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Sidi1901/urlShortner/internal/logger"
	"github.com/Sidi1901/urlShortner/internal/model"
)

type ShortURLRepository interface {
	// Save Short code for an original URL
	SaveShortCode(ctx context.Context, shortURL *model.ShortURL) error

	// Read operation
	GetByShortCode(ctx context.Context, code string) (*model.ShortURL, error)

	// Update expiry time for a short code
	UpdateShortCode(ctx context.Context, shortURL *model.ShortURL) error

	// Delete Short code foran original URL
	DeleteShortCode(ctx context.Context, shortCode string) error
}

// Getters and setters for ShortURLRepository

func (r *Repository) SaveShortCode(ctx context.Context, sURL *model.ShortURL) error {
	query := `INSERT into url_shortner.short_urls (short_code, original_url, created_at, expiry_duration, ip_address, is_active) VALUES($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query, sURL.ShortCode, sURL.OriginalURL, sURL.CreatedAt, sURL.ExpiryDuration, sURL.IPAddress, sURL.IsActive)

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

func (r *Repository) GetShortCode(ctx context.Context, code string) (*model.ShortURL, error) {
	query := `SELECT short_code, original_url, created_at, expiry_duration, ip_address, is_active FROM url_shortner.short_urls WHERE short_code = $1`

	var shortURL model.ShortURL
	err := r.db.GetContext(ctx, &shortURL, query, code)

	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"shortcode": code,
			"error":     err.Error(),
		}).Error("Failed to retrieve short URL data")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("short code not found")
		} else {
			return nil, fmt.Errorf("error retrieving short code: %w", err)
		}
	}

	logger.Log.WithFields(map[string]interface{}{
		"shortcode": shortURL.ShortCode,
	}).Info("Short code retrieved successfully")

	return &shortURL, nil

}

func (r *Repository) UpdateShortCode(ctx context.Context, shortURLData *model.ShortURL) error {
	query := `UPDATE url_shortner.short_urls SET original_url = $1, expiry_duration = $2, is_active = $3, updated_at = $4 WHERE short_code = $5`

	logger.Log.WithFields(map[string]interface{}{
		"shortcode": shortURLData.ShortCode,
	}).Info("Updating short URL")

	_, err := r.db.ExecContext(ctx, query, shortURLData.OriginalURL, shortURLData.ExpiryDuration, shortURLData.IsActive, shortURLData.UpdatedAt, shortURLData.ShortCode)

	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"shortcode": shortURLData.ShortCode,
			"error":     err.Error(),
		}).Error("Failed to update short URL")
		return err
	}

	logger.Log.WithFields(map[string]interface{}{
		"shortcode": shortURLData.ShortCode,
	}).Info("Short URL updated successfully")

	return nil
}

func (r *Repository) DeleteShortCode(ctx context.Context, code string) error {
	query := `DELETE FROM url_shortner.short_urls WHERE short_code = $1`

	_, err := r.db.ExecContext(ctx, query, code)

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
