package utils

import (
	"errors"
	"time"

	"github.com/Sidi1901/urlShortner/internal/dto"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID, email, userRole, userType, jwtSecret string) (string, error) {
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
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     time.Now().Add(12 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ValidateJWT(tokenStr, jwtSecret string) (*dto.Claims, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &dto.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*dto.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
