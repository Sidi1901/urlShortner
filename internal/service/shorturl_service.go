package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Sidi1901/urlShortner/internal/dto"
	"github.com/Sidi1901/urlShortner/internal/model"
	"github.com/Sidi1901/urlShortner/pkg/utils"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

func (s *Service) CreateShortURL(ctx context.Context, url string, ip string, expirySec int, shortCode string) (string, error) {

	// 1. Check if it is an acual URL
	if !govalidator.IsURL(url) {
		return "", fmt.Errorf("Invalid URL")
	}

	// 2. Check for domain error
	if !utils.IsValidDomain(url) {
		return "", fmt.Errorf("Invalid Domain")
	}

	// 3. enforce ssl for https

	url = utils.EnforceHTTP(url)

	// 	After all checks have been passed, Create (or input from user) unique Custom short code for url = domain + customShortCode.

	if shortCode == "" {
		shortCode = uuid.New().String()[:6]
	}

	// 4. Check if the custom short code is already in use. If it is, return an error message to the user.

	_, err := s.repo.GetShortCode(ctx, shortCode)

	if err == nil {
		return "", fmt.Errorf("Custom short url is already in use. Please submit request with different custom short code")
	}

	// 5. Save data in table ShortURL.

	// Save data in table ShortURL.
	shortURL := &model.ShortURL{
		ShortCode:      shortCode,
		OriginalURL:    url,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiryDuration: expirySec,
		IPAddress:      ip,
		IsActive:       true,
	}

	if err := s.repo.SaveShortCode(ctx, shortURL); err != nil {
		return "", fmt.Errorf("Error occured - %s", err)
	}

	fmt.Printf("https://%s:%s/:%s", s.cfg.Domain, s.cfg.AppPort, shortCode)
	shortFDQN := fmt.Sprintf("https://%s:%s/%s", s.cfg.Domain, s.cfg.AppPort, shortCode)

	return shortFDQN, nil
}

func (s *Service) ResolveShortURL(ctx context.Context, shortcode string) (string, error) {

	shortURLData, err := s.repo.GetShortCode(ctx, shortcode)

	if err != nil {
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
		if err := s.repo.UpdateShortCode(ctx, shortURLData); err != nil {
			return "", fmt.Errorf("Error updating short URL status: %s", err)
		}
		return "", fmt.Errorf("Short URL has expired")
	}

	return shortURLData.OriginalURL, nil

}

func (s *Service) GetShortURLInfo(ctx context.Context, shortcode string) (dto.ShortURLInfo, error) {

	fmt.Printf("In service %s", shortcode)
	shortURLData, err := s.repo.GetShortCode(ctx, shortcode)

	var shortURLInfo dto.ShortURLInfo

	if err != nil {
		fmt.Println(err.Error())
		fmt.Errorf("Error occurred %s", err.Error())
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

func (s *Service) DeleteShortCode(ctx context.Context, shortcode string) error {
	err := s.repo.DeleteShortCode(ctx, shortcode)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateShortURLInfo(ctx context.Context, shortcode string, url *string, expiryDuration *int, isActive bool) error {
	// Fetch existing short URL record
	shortURLData, err := s.repo.GetShortCode(ctx, shortcode)
	if err != nil {
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

	if err := s.repo.UpdateShortCode(ctx, shortURLData); err != nil {
		return fmt.Errorf("Error updating short URL: %s", err)
	}

	return nil
}
