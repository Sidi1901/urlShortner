package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
		return fmt.Errorf("failed to save short code: %w", err)
	}

	fmt.Printf("Short code saved successfully => %s\n", sURL)

	return nil
}

func (r *Repository) GetByShortCode(ctx context.Context, code string) (*model.ShortURL, error) {
	query := `SELECT short_code, original_url, created_at, expiry_duration, ip_address, is_active FROM url_shortner.short_urls WHERE short_code = $1`

	var shortURL model.ShortURL
	err := r.db.GetContext(ctx, &shortURL, query, code)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("short code not found")
		} else {
			return nil, fmt.Errorf("error retrieving short code: %w", err)
		}
	}

	fmt.Printf("Short code retrieved successfully => %s\n", shortURL)

	return &shortURL, nil

}

func (r *Repository) UpdateShortCode(ctx context.Context, sURL *model.ShortURL) error {
	query := `UPDATE url_shortner.short_urls SET expiry_duration = $1 WHERE short_code = $2`

	_, err := r.db.ExecContext(ctx, query, sURL.ExpiryDuration, sURL.ShortCode)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteShortCode(ctx context.Context, code string) error {
	query := `DELETE FROM url_shortner.short_urls WHERE short_code = $1`

	_, err := r.db.ExecContext(ctx, query, code)

	return err
}
