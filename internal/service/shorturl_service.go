package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Sidi1901/urlShortner/internal/config"
	"github.com/Sidi1901/urlShortner/internal/domain"
	errs "github.com/Sidi1901/urlShortner/internal/errors"
	"github.com/Sidi1901/urlShortner/internal/logger"
	"github.com/Sidi1901/urlShortner/internal/model"
	"github.com/Sidi1901/urlShortner/internal/repository"
	"github.com/Sidi1901/urlShortner/pkg/utils"
	"github.com/asaskevich/govalidator"
)

type ShortURLService interface {
	CreateShortURL(ctx context.Context, url string, ip string, expirySec *int, shortCode string, email string) (string, int, error)
	ResolveShortURL(ctx context.Context, shortcode string) (string, error)
	GetShortURLInfo(ctx context.Context, shortcode string) (domain.ShortURLInfo, error)
	DeleteShortCode(ctx context.Context, shortcode string) error
	UpdateShortURLInfo(ctx context.Context, shortcode string, url *string, expiryDuration *int, isActive *bool) error
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
		"userID": userData.ID,
	}).Info("User data retrieved successfully")

	// Ensure expirySec is non-nil before dereferencing to avoid nil pointer panic.
	if expirySec == nil {
		defaultExpiry := 24 * 3600
		expirySec = &defaultExpiry
	}

	if userData.UserType == "free" && *expirySec > 24*3600 {
		*expirySec = 24 * 3600
	}

	if shortCode == "" {

		retryCount := 0

		for retryCount < 5 {
			shortCode = utils.GenerateShortCode(6)

			_, err := s.urlRepo.GetShortCode(ctx, shortCode)
			if errors.Is(err, sql.ErrNoRows) {
				break // available
			}
			if err != nil {
				return "", 0, err // real DB error
			}

			retryCount++
		}

		if retryCount == 5 {
			return "", 0, fmt.Errorf("failed to generate unique shortcode")
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

		// TODO make the as a backgroud job to avoid latency in redirection, mark as inactive and then update in DB, for now doing in same request to avoid complexity of background job
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

func (s *shortURLService) GetShortURLInfo(ctx context.Context, shortcode string) (domain.ShortURLInfo, error) {

	logger.Log.WithFields(map[string]interface{}{
		"shortcode": shortcode,
	}).Info("Fetching short URL info")

	shortURLData, err := s.urlRepo.GetShortCode(ctx, shortcode)

	var shortURLInfo domain.ShortURLInfo

	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"shortcode": shortcode,
			"error":     err.Error(),
		}).Error("Failed to retrieve short URL info")
		return shortURLInfo, err
	}

	shortURL := fmt.Sprintf("https://%s/%s", s.cfg.Domain, shortURLData.ShortCode)

	shortURLInfo = domain.ShortURLInfo{
		URL:            shortURLData.OriginalURL,
		ShortCode:      shortURLData.ShortCode,
		ShortURL:       shortURL,
		ExpiryDuration: shortURLData.ExpiryDuration,
		CreatedAt:      shortURLData.CreatedAt,
		LastUpdatedAt:  shortURLData.UpdatedAt,
		IsActive:       shortURLData.IsActive,
		UserID:         shortURLData.UserID,
	}

	return shortURLInfo, nil
}

func (s *shortURLService) DeleteShortCode(ctx context.Context, shortcode string) error {

	// Fetch existing short URL record
	shortURLData, _ := s.urlRepo.GetShortCode(ctx, shortcode)

	// Check resource ownership
	userID := ctx.Value("user_id").(int)
	if shortURLData.UserID != userID {
		return errs.ErrUnauthorized
	}

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

func (s *shortURLService) UpdateShortURLInfo(ctx context.Context, shortcode string, url *string, expiryDuration *int, isActive *bool) error {

	// Fetch existing short URL record
	shortURLData, err := s.urlRepo.GetShortCode(ctx, shortcode)

	// Check resource ownership
	userID := ctx.Value("user_id").(int)
	if shortURLData.UserID != userID {
		return errs.ErrUnauthorized
	}

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
	if isActive != nil {
		shortURLData.IsActive = *isActive
	}
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
