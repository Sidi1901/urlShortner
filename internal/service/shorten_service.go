package service

import (
	"context"
	"fmt"
	"time"

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

	fmt.Printf("Generated short code => %s\n", shortCode)

	// 4. Check if the custom short code is already in use. If it is, return an error message to the user.

	resp, err := s.repo.GetByShortCode(ctx, shortCode)
	if err == nil {
		return "", fmt.Errorf("Custom short url is already in use. Please submit request with different custom short code")
	}

	fmt.Printf("Retrieved short code data => %v\n", resp)

	fmt.Printf("Custom short code is available => %s\n", shortCode)

	// 5. Save data in table ShortURL.

	// Save data in table ShortURL.
	shortURL := &model.ShortURL{
		ShortCode:      shortCode,
		OriginalURL:    url,
		CreatedAt:      time.Now(),
		ExpiryDuration: expirySec,
		IPAddress:      ip,
		IsActive:       true}

	if err := s.repo.SaveShortCode(ctx, shortURL); err != nil {
		return "", fmt.Errorf("Error occured - %s", err)
	}

	fmt.Printf("https://%s:%s/:%s", s.cfg.Domain, s.cfg.AppPort, shortCode)
	shortFDQN := fmt.Sprintf("https://%s:%s/%s", s.cfg.Domain, s.cfg.AppPort, shortCode)

	return shortFDQN, nil
}

// func GetShortURL(shortcode string) {

// 	GetByShortCode(ctx context.Context, code string)

// 	shortURLData := shortURLData{
// 		URL           :
// 		ShortURL      :
// 		Expiry        :
// 		CreatedAt     :
// 		LastUpdatedAt :
// 	}

// }
