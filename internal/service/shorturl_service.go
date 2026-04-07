package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Sidi1901/urlShortner/internal/config"
	"github.com/Sidi1901/urlShortner/internal/dto"
	"github.com/Sidi1901/urlShortner/internal/logger"
	"github.com/Sidi1901/urlShortner/internal/model"
	"github.com/Sidi1901/urlShortner/internal/repository"
	"github.com/Sidi1901/urlShortner/pkg/utils"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type ShortURLService interface {
	CreateShortURL(ctx context.Context, url string, ip string, expirySec *int, shortCode string, email string) (string, int, error)
	ResolveShortURL(ctx context.Context, shortcode string) (string, error)
	GetShortURLInfo(ctx context.Context, shortcode string) (dto.ShortURLInfo, error)
	DeleteShortCode(ctx context.Context, shortcode string) error
	UpdateShortURLInfo(ctx context.Context, shortcode string, url *string, expiryDuration *int, isActive bool) error
}

type shortURLService struct {
	urlRepo  repository.ShortURLRepository
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewShortURLService(urlRepo repository.ShortURLRepository, userRepo repository.UserRepository, cfg *config.Config) ShortURLService {
	return &shortURLService{urlRepo: urlRepo, userRepo: userRepo, cfg: cfg}
}

func (s *shortURLService) CreateShortURL(ctx context.Context, url string, ip string, expirySec *int, shortCode string, email string) (string, int, error) {

	logger.Log.WithFields(map[string]interface{}{
		"url":       url,
		"ip":        ip,
		"expirySec": expirySec,
		"shortCode": shortCode,
		"email":     email,
	}).Info("Creating short URL")

	// 1. Check if it is an acual URL
	if !govalidator.IsURL(url) {
		logger.Log.WithFields(map[string]interface{}{
			"url":   url,
			"error": "Invalid URL",
		}).Error("Failed to validate URL")
		return "", 0, fmt.Errorf("Invalid URL")
	}

	// 2. Check for domain error
	if !utils.IsValidDomain(url) {
		logger.Log.WithFields(map[string]interface{}{
			"url":   url,
			"error": "Invalid Domain",
		}).Error("Failed to validate domain")
		return "", 0, fmt.Errorf("Invalid Domain")
	}

	// 3. enforce ssl for https

	url = utils.EnforceHTTP(url)

	userData, err := s.userRepo.GetUser(ctx, email)
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"email": email,
			"error": err.Error(),
		}).Error("Failed to retrieve user data")
		return "", 0, fmt.Errorf("Error occurred while fetching user data: %s", err)
	}

	logger.Log.WithFields(map[string]interface{}{
		"email":     email,
		"userData":  userData,
		"expirySec": expirySec,
	}).Info("User data retrieved successfully")

	// Ensure expirySec is non-nil before dereferencing to avoid nil pointer panic.
	if userData.UserType == "free" || expirySec == nil {
		defaultExpiry := 24 * 3600
		if expirySec == nil {
			expirySec = &defaultExpiry
		} else {
			*expirySec = defaultExpiry
		}
	}
	fmt.Println("Expiry is ", *expirySec)

	if shortCode == "" {

		var err error
		retryCount := 0

		for err == nil && retryCount < 5 {
			shortCode = uuid.New().String()[:6]
			// 4. Check if the custom short code is already in use. If it is, return an error message to the user.

			_, err := s.urlRepo.GetShortCode(ctx, shortCode)

			if err == nil {
				return "", 0, fmt.Errorf("Custom short url is already in use. Please submit request with different custom short code")
			}
			retryCount++
		}
	} else {
		_, err := s.urlRepo.GetShortCode(ctx, shortCode)

		if err == nil {
			return "", 0, fmt.Errorf("Custom short url is already in use. Please submit request with different custom short code")
		}
	}

	// 5. Save data in table ShortURL.

	// Save data in table ShortURL.
	shortURL := &model.ShortURL{
		ShortCode:      shortCode,
		OriginalURL:    url,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiryDuration: *expirySec,
		IPAddress:      ip,
		IsActive:       true,
		UserID:         userData.ID, // Set to 0 for now, can be updated later when user authentication is implemented
	}

	if err := s.urlRepo.SaveShortCode(ctx, shortURL); err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"shortcode": shortCode,
			"error":     err.Error(),
		}).Error("Failed to save short URL")
		return "", 0, fmt.Errorf("Error occured - %s", err)
	}

	logger.Log.WithFields(map[string]interface{}{
		"shortcode": shortCode,
		"shortURL":  shortURL,
	}).Info("Short URL created successfully")

	shortFDQN := fmt.Sprintf("https://%s:%s/%s", s.cfg.Domain, s.cfg.AppPort, shortCode)

	return shortFDQN, *expirySec, nil
}

func (s *shortURLService) ResolveShortURL(ctx context.Context, shortcode string) (string, error) {

	shortURLData, err := s.urlRepo.GetShortCode(ctx, shortcode)

	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"shortcode": shortcode,
			"error":     err.Error(),
		}).Error("Failed to retrieve short URL data")
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("Short URL not found")
		}
		return "", fmt.Errorf("Error occured - %s", err)
	}

	if !shortURLData.IsActive {
		return "", fmt.Errorf("Short URL has expired")
	}

	// Check if the short URL has expired based on the expiry duration and created at time
	if time.Since(shortURLData.CreatedAt) > time.Duration(shortURLData.ExpiryDuration)*time.Second {
		// Mark the short URL as inactive in the database
		shortURLData.UpdatedAt = time.Now()
		shortURLData.IsActive = false
		if err := s.urlRepo.UpdateShortCode(ctx, shortURLData); err != nil {
			logger.Log.WithFields(map[string]interface{}{
				"shortcode": shortcode,
				"error":     err.Error(),
			}).Error("Failed to update short URL status")
			return "", fmt.Errorf("Error updating short URL status: %s", err)
		}
		logger.Log.WithFields(map[string]interface{}{
			"shortcode": shortcode,
		}).Info("Short URL has expired and marked as inactive")
		return "", fmt.Errorf("Short URL has expired")
	}

	return shortURLData.OriginalURL, nil

}

func (s *shortURLService) GetShortURLInfo(ctx context.Context, shortcode string) (dto.ShortURLInfo, error) {

	logger.Log.WithFields(map[string]interface{}{
		"shortcode": shortcode,
	}).Info("Fetching short URL info")

	shortURLData, err := s.urlRepo.GetShortCode(ctx, shortcode)

	var shortURLInfo dto.ShortURLInfo

	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"shortcode": shortcode,
			"error":     err.Error(),
		}).Error("Failed to retrieve short URL info")
		return shortURLInfo, err
	}

	shortURL := fmt.Sprintf("https://%s:%s/%s", s.cfg.Domain, s.cfg.AppPort, shortURLData.ShortCode)

	shortURLInfo = dto.ShortURLInfo{
		URL:            shortURLData.OriginalURL,
		ShortCode:      shortURLData.ShortCode,
		ShortURL:       shortURL,
		ExpiryDuration: shortURLData.ExpiryDuration,
		CreatedAt:      shortURLData.CreatedAt,
		LastUpdatedAt:  shortURLData.UpdatedAt,
		IsActive:       shortURLData.IsActive,
	}

	return shortURLInfo, nil
}

func (s *shortURLService) DeleteShortCode(ctx context.Context, shortcode string) error {
	err := s.urlRepo.DeleteShortCode(ctx, shortcode)

	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"shortcode": shortcode,
			"error":     err.Error(),
		}).Error("Failed to delete short URL")
		return err
	}

	logger.Log.WithFields(map[string]interface{}{
		"shortcode": shortcode,
	}).Info("Short URL deleted successfully")

	return nil
}

func (s *shortURLService) UpdateShortURLInfo(ctx context.Context, shortcode string, url *string, expiryDuration *int, isActive bool) error {
	// Fetch existing short URL record
	shortURLData, err := s.urlRepo.GetShortCode(ctx, shortcode)
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"shortcode": shortcode,
			"error":     err.Error(),
		}).Error("Failed to retrieve short URL data")
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("Short URL not found")
		}
		return fmt.Errorf("Error occurred - %s", err)
	}

	// Update only the fields provided
	if url != nil {
		shortURLData.OriginalURL = *url
	}
	if expiryDuration != nil {
		shortURLData.ExpiryDuration = *expiryDuration
	}
	shortURLData.IsActive = isActive
	shortURLData.UpdatedAt = time.Now()

	if err := s.urlRepo.UpdateShortCode(ctx, shortURLData); err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"shortcode": shortcode,
			"error":     err.Error(),
		}).Error("Failed to update short URL")
		return fmt.Errorf("Error updating short URL: %s", err)
	}

	logger.Log.WithFields(map[string]interface{}{
		"shortcode": shortcode,
	}).Info("Short URL updated successfully")

	return nil
}
