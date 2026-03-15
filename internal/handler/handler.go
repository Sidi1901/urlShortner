package handler

import(
	"fmt"
	"time"
	"github.com/Sidi1901/urlShortner/database"
	"github.com/Sidi1901/urlShortner/internal/service/shorten_service"
	"github.com/Sidi1901/urlShortner/internal/service/resolve_service"
)


type request struct {
	URL	          string             `json:"url"`
	CustomShort   string             `json:"short"`
	Expiry        time.Duration      `json:"expiry"`
}

type response struct {
	URL              string          `json:"url"`
	CustomShort      string          `json:"shorturl"`
	Expiry           time.Duration   `json:"expiry"`
	XRateRemaining   int             `json:"rate_limit"`
	XRateLimitReset  time.Duration   `json:"rate_limit_reset"`
}

func ShortenURL(c *gin.Context) {
	req := request{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format"
		})
		return
	}

	url := req.URL
	expiry := req.Expiry
	customShort := req.CustomShort

	if expiry != "" {
		expiry, err := time.ParseDuration(expiry)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid expiry format"
			})
			return
		}
	}

	serviceResponse, err := service.CreateShortURL(url, c.IP(), expiry, customShort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not call short URL service"
		})
		return
	}


	shortURL := fmt.Sprintf("https://%s/%s",cfg.DDOMAIN, code)

	resp := response{
		URL: url,
		CustomShort: serviceResponse[0],
		Expiry: expiry,
		XRateRemaining: serviceResponse[1],
		XRateLimitReset: serviceResponse[2],
	}

	c.JSON(http.StatusOK, resp)
}


func ResolveURL(c *gin.Context) {
	code := c.Param("code")

	url, err := service.GetOriginalURL(code)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":"URL not found"
		})

		return
	}

	c.Redirect(http.StatusMovedPermanently, url)
}

