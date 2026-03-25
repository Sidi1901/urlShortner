package model

import (
    "time"
    "github.com/google/uuid"
)

// ShortURL represents a shortened URL entry
type ShortURL struct {
    ID          uuid.UUID    `db:"id"`
    ShortCode   string       `db:"short_code"`  // UNIUQUE INDEX
    OriginalURL string       `db:"original_url"`
    CreatedAt   time.Time    `db:"created_at"` 
    ExpiresAt   *time.Time   `db:"expires_at"`
    IPAddress   string       `db:"ip_address"` 
}