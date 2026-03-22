package handler

import(
	"fmt"
	"net/http"
	"time"
	"github.com/Sidi1901/urlShortner/internal/service"
	"github.com/Sidi1901/urlShortner/internal/config"
	"github.com/gin-gonic/gin"
)

var cfg *config.Config

type request struct {
	URL	          string             `json:"url"`
	CustomShort   string             `json:"short"`
	Expiry        string             `json:"expiry"`
}

type response struct {
	URL                 string          `json:"url"`
	CustomShortURL      string          `json:"shorturl"`
	Expiry              string          `json:"expiry"`
	XRateRemaining      int             `json:"rate_limit"`
	XRateLimitReset     time.Duration   `json:"rate_limit_reset"`
}

func ShortenURL(c *gin.Context) {
	req := request{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}


	var expiry_time time.Duration

	if req.Expiry != "" {
		if _, err := time.ParseDuration(req.Expiry); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiry format"})
			return
		}

		expiry_time, _ = time.ParseDuration(req.Expiry)

	} else {
		expiry_time = 24 * time.Hour
	} 

	serviceResponse, err := service.CreateShortURL(req.URL, c.ClientIP(), expiry_time, req.CustomShort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not call short URL service"})
		return
	}


	shortURL := fmt.Sprintf("https://%s/%s",cfg.DOMAIN, serviceResponse.CustomShortCodeURL)

	resp := response{
		URL: req.URL,
		CustomShortURL: shortURL,
		Expiry: req.Expiry,
		XRateRemaining: serviceResponse.Quota,
		XRateLimitReset: serviceResponse.ResetTime,
	}

	c.JSON(http.StatusOK, resp)
}


func ResolveURL(c *gin.Context) {
	code := c.Param("code")

	url, err := service.GetOriginalURL(code)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error":"URL not found"})

		return
	}

	c.Redirect(http.StatusMovedPermanently, url)
}

