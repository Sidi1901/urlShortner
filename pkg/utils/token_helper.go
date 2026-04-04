package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID, email, userRole, userType, jwtSecret string) (string, error) {
	fmt.Println("Inside GenerateAccessToken with userID:", userID, "email:", email, "userRole:", userRole, "userType:", userType)
	claims := jwt.MapClaims{
		"user_id":   userID,
		"email":     email,
		"user_role": userRole,
		"user_type": userType,
		"type":      "access",
		"exp":       time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func GenerateRefreshToken(userID, jwtSecret string) (string, error) {
	fmt.Println("Inside GenerateRefreshToken with userID:", userID)
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     time.Now().Add(12 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
