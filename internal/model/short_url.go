package models

import (
    "time"
)

// ShortURL represents a shortened URL entry
type ShortURL struct {
    ID          string    `db:"id" json:"id"`
    ShortCode   string    `db:"short_code" json:"short_code"`
    OriginalURL string    `db:"original_url" json:"original_url"`
    CreatedAt   time.Time `db:"created_at" json:"created_at"`
    ExpiresAt   time.Time `db:"expires_at" json:"expires_at"`
    IPAddress   string    `db:"created_by" json:"created_by"` // IP or user
}

// RateLimit represents API quota per IP
type RateLimit struct {
    IPAddress       string    `db:"ip_address" json:"ip_address"`
    RemainingQuota  int       `db:"remaining_quota" json:"remaining_quota"`
    ResetAt         time.Time `db:"reset_at" json:"reset_at"`
    UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}
