package dto

import "time"

// Create short URL Request
type CreateShortURLRequest struct {
	URL            string  `json:"url" binding:"required,url"`
	ShortCode      *string `json:"short_code,omitempty"`
	ExpiryDuration *int    `json:"expiry_seconds,omitempty"`
}

// Create short URL Response
type CreateShortURLResponse struct {
	URL            string    `json:"url"`
	ShortCode      string    `json:"short_code"`
	ShortURL       string    `json:"short_url"`
	ExpiryDuration int       `json:"expiry_duration"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Get Short URL data Response
type ShortURLInfoResponse struct {
	URL            string    `json:"url"`
	ShortCode      string    `json:"short_code"`
	ShortURL       string    `json:"short_url"`
	ExpiryDuration int       `json:"expiry"`
	CreatedAt      time.Time `json:"created_at"`
	LastUpdatedAt  time.Time `json:"last_updated_at"`
}

// Get Short URL data
type ShortURLInfo struct {
	URL            string
	ShortCode      string
	ShortURL       string
	ExpiryDuration int
	CreatedAt      time.Time
	LastUpdatedAt  time.Time
	IsActive       bool
}

// Request Update short URL data
type UpdateShortURLRequest struct {
	Shortcode      string  `json:"short_code" binding:"required"`
	ExpiryDuration *int    `json:"expiry_seconds,omitempty"`
	URL            *string `json:"url,omitempty"`
	IsActive       bool    `json:"is_active,omitempty"`
}
