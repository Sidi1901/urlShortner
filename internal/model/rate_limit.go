package model

import (
    "time"
)

// RateLimit represents API quota per IP
type RateLimit struct {
    IPAddress       string     `db:"ip_address"` // INDEX
    RemainingQuota  int        `db:"remaining_quota"`
    ResetAt         time.Time  `db:"reset_at"`
    UpdatedAt       time.Time  `db:"updated_at"`
}