package model

import (
	"time"
)

type User struct {
	ID           int       `db:"id"`
	Email        *string   `db:"email"`
	Name         *string   `db:"name"`
	Password     *string   `db:"password"`
	UserType     string    `db:"user_type"`
	UserRole     string    `db:"user_role"`
	Token        string    `db:"token"`
	RefreshToken string    `db:"refresh_token"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	UserID       string    `db:"user_id"`
}
