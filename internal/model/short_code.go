package model

import (
	"time"
)

// ShortURL represents a shortened URL entry
type ShortURL struct {
	ShortCode      string    `db:"short_code"` // UNIQUE, NON NULL, IMMUTABLE INDEX
	OriginalURL    string    `db:"original_url"`
	CreatedAt      time.Time `db:"created_at"`
	ExpiryDuration int       `db:"expiry_duration"`
	IPAddress      string    `db:"ip_address"`
	IsActive       bool      `db:"is_active"`
}
