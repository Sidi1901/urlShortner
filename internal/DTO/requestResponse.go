package dto

import "time"

// Request Create short URL
type CreateShortURLRequest struct {
	URL        string  `json:"url" binding:"required,url"`
	ShortCode  *string `json:"short_code,omitempty"`
	ExpirySec  *int    `json:"expiry_seconds,omitempty"` 
}

// Response Create short URL
type CreateShortURLResponse struct {
	URL       string     `json:"url"`
	ShortCode string     `json:"short_code"`
	ShortURL  string     `json:"short_url"`
	Expiry    *time.Time `json:"expiry,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Get Short URL data response
type ShortURLDataResponse struct {
	URL           string     `json:"url"`
	ShortCode     string     `json:"short_code"`
	ShortURL      string     `json:"short_url"`
	Expiry        *time.Time `json:"expiry"`
	CreatedAt     time.Time  `json:"created_at"`
	LastUpdatedAt time.Time  `json:"last_updated_at"`
}

// Get Short URL data
type ShortURLData struct {
	URL           string     `json:"url"`
	ShortURL      string     `json:"short_url"`
	Expiry        *time.Time `json:"expiry"`
	CreatedAt     time.Time  `json:"created_at"`
	LastUpdatedAt time.Time  `json:"last_updated_at"`
}


// Request Update short URL data
type UpdateShortURLRequest struct {
	ShortURL        string  `json:"short_url" binding:"required,url"`
	ExpirySec      *int    `json:"expiry_seconds,omitempty"` 
	URL			   *string `json:"url,omitempty"`
}

// Response Update short URL data
type UpdateShortURLResponse struct {
	ShortURL       string     `json:"short_url" binding:"required,url"`
	ExpirySec      *int       `json:"expiry_seconds,omitempty"` 
	URL			   *string    `json:"url,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}