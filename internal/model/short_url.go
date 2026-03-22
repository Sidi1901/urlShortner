package model

import (
    "time"
    "github.com/google/uuid"
)

// ShortURL represents a shortened URL entry
type ShortURL struct {
    ID          uuid.UUID    `db:"id"`
    ShortCode   string    `db:"short_code"`
    OriginalURL string    `db:"original_url"`
    CreatedAt   time.Time `db:"created_at"` // no need for pointer type here since it will always have a value ie a mandatory field
    ExpiresAt   *time.Time `db:"expires_at"`
    IPAddress   string    `db:"ip_address"` // IP or user
}


// RateLimit represents API quota per IP
type RateLimit struct {
    IPAddress       string     `db:"ip_address" json:"ip_address"`
    RemainingQuota  int        `db:"remaining_quota" json:"remaining_quota"`
    ResetAt         time.Time  `db:"reset_at" json:"reset_at"`
    UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}
