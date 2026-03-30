package dto

import "time"

// Request
type CreateShortURLRequest struct {
	URL        string  `json:"url" binding:"required,url"`
	ShortCode  *string `json:"short_code,omitempty"`
	ExpirySec  *int    `json:"expiry_seconds,omitempty"` 
}

// Response
type CreateShortURLResponse struct {
	URL       string     `json:"url"`
	ShortCode string     `json:"short_code"`
	ShortURL  string     `json:"short_url"`
	Expiry    *time.Time `json:"expiry,omitempty"`
}
