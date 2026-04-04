package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Sidi1901/urlShortner/internal/dto"
	errs "github.com/Sidi1901/urlShortner/internal/errors"
	"github.com/Sidi1901/urlShortner/internal/model"
	"github.com/Sidi1901/urlShortner/internal/repository"
	"github.com/Sidi1901/urlShortner/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) CreateUser(ctx context.Context, email, name, password, userType string) error {

	if email == "" {
		return fmt.Errorf("email is required %w", errs.ErrInvalidInput)
	}

	if name == "" {
		return fmt.Errorf("name is required %w", errs.ErrInvalidInput)
	}

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters %w", errs.ErrInvalidInput)
	}

	if !utils.IsValidPassword(password) {
		return fmt.Errorf("password must contain letter, number and special character %w", errs.ErrInvalidInput)
	}

	_, err := s.repo.GetUser(ctx, email)
	if err == nil {
		return fmt.Errorf("user already exists %w", errs.ErrUserAlreadyExists)
	}

	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return err
	}

	// -------- 3. HASH PASSWORD --------
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("failed to hash password %w", errs.ErrInternal)
	}

	now := time.Now()

	user := &model.User{
		Email:     email,
		Name:      name,
		Password:  string(hashedPassword),
		UserType:  userType,
		UserRole:  "user",
		UserID:    uuid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, string, error) {

	user, err := s.repo.GetUser(ctx, email)

	if err != nil {
		return "", "", errs.ErrUserNotFound
	}

	// password check
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errs.ErrInvalidInput
	}

	// generate tokens
	accessToken, err := utils.GenerateAccessToken(user.UserID, user.Email, user.UserRole, user.UserType, s.cfg.JwtSecret)

	if err != nil {
		return "", "", errs.ErrInternal
	}

	// generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.UserID, s.cfg.JwtSecret)

	if err != nil {
		return "", "", errs.ErrInternal
	}

	return accessToken, refreshToken, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (string, error) {

	claims := &dto.RefreshClaims{}

	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		// validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrInvalidInput
		}
		return []byte(s.cfg.JwtSecret), nil
	})

	if err != nil {
		return "", errs.ErrInvalidInput
	}

	// -------- EXPLICIT CHECKS --------

	if !token.Valid {
		return "", errs.ErrInvalidInput
	}

	if claims.Type != "refresh" {
		return "", errs.ErrInvalidInput
	}

	// expiry check (explicit)
	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		return "", errs.ErrInvalidInput
	}

	// -------- GENERATE ACCESS TOKEN --------
	accessToken, err := utils.GenerateAccessToken(claims.UserID, claims.Email, claims.UserRole, claims.UserType, s.cfg.JwtSecret)
	if err != nil {
		return "", errs.ErrInternal
	}

	return accessToken, nil
}
