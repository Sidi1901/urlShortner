package domain

import (
	"github.com/golang-jwt/jwt/v5"
)

// Create Claims for JWT
type Claims struct {
	UserID   int
	Email    string
	Type     string
	UserRole string
	UserType string
	jwt.RegisteredClaims
}
