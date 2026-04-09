package domain

import (
	"time"
)

// Get Short URL data
type ShortURLInfo struct {
	URL            string
	ShortCode      string
	ShortURL       string
	ExpiryDuration int
	CreatedAt      time.Time
	LastUpdatedAt  time.Time
	IsActive       bool
	UserID         int
}
