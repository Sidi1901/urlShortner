package dto

import (
	"time"
)

// Create short URL Request
type CreateShortURLRequest struct {
	Email          string  `json:"email" binding:"required,email"`
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
	ExpiryDuration int       `json:"expiry_duration"`
	CreatedAt      time.Time `json:"created_at"`
	LastUpdatedAt  time.Time `json:"last_updated_at"`
}

// Request Update short URL data
type UpdateShortURLRequest struct {
	Shortcode      string  `json:"short_code" binding:"required"`
	ExpiryDuration *int    `json:"expiry_seconds,omitempty"`
	URL            *string `json:"url,omitempty"`
	IsActive       *bool   `json:"is_active,omitempty"`
}

// USER DTOs

// Create User Request
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=2"`
	Password string `json:"password" binding:"required,min=8"`
	UserType string `json:"user_type" binding:"required,oneof=free premium"`
}

// Create Login Request
type CreateLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
