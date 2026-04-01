package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sidi1901/urlShortner/internal/dto"
	"github.com/gin-gonic/gin"
)

// GET /:shortcode
func (h *Handler) ResolveURL(c *gin.Context) {
	ctx := c.Request.Context()
	shortcode := c.Param("shortcode")

	mappedURL, err := h.service.ResolveShortURL(ctx, shortcode)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// 302
	c.Redirect(http.StatusFound, mappedURL)
}

// POST /api/v1/urls
func (h *Handler) CreateShortURL(c *gin.Context) {
	ctx := c.Request.Context()
	var request dto.CreateShortURLRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	uRL := request.URL
	shortCode := request.ShortCode
	ip := c.ClientIP()

	var expirySec int
	if request.ExpiryDuration == nil {
		expirySec = 24 * 3600 // default expiry time of 24 hours
	} else {
		expirySec = *request.ExpiryDuration
	}

	shortURL, err := h.service.CreateShortURL(ctx, uRL, ip, expirySec, *shortCode)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.CreateShortURLResponse{
		URL:            uRL,
		ShortCode:      *shortCode,
		ShortURL:       shortURL,
		ExpiryDuration: expirySec,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	c.JSON(http.StatusCreated, response)
}

// GET /api/v1/urls/:shortCode
func (h *Handler) GetShortURL(c *gin.Context) {
	ctx := c.Request.Context()
	shortcode := c.Param("shortcode")

	shortURLData, err := h.service.GetShortURLInfo(ctx, shortcode)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := dto.ShortURLInfoResponse{
		URL:            shortURLData.URL,
		ShortCode:      shortcode,
		ShortURL:       shortURLData.ShortURL,
		ExpiryDuration: shortURLData.ExpiryDuration,
		CreatedAt:      shortURLData.CreatedAt,
		LastUpdatedAt:  shortURLData.LastUpdatedAt,
	}

	c.JSON(http.StatusFound, response)
}

// PUT /api/v1/urls/:shortcode
func (h *Handler) UpdateShortURLInfo(c *gin.Context) {

	ctx := c.Request.Context()

	var request dto.UpdateShortURLRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	fmt.Printf("1st In handler %s", request.Shortcode)

	err := h.service.UpdateShortURLInfo(ctx, request.Shortcode, request.URL, request.ExpiryDuration, request.IsActive)

	fmt.Printf("2nd In handler %s", request.Shortcode)

	if err != nil {
		fmt.Printf("In handler %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("3rd In handler %s", request.Shortcode)

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

// DELETE /api/v1/urls/:shortCode
func (h *Handler) DeleteShortURL(c *gin.Context) {
	ctx := c.Request.Context()
	shortCode := c.Param("shortcode")

	err := h.service.DeleteShortCode(ctx, shortCode)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted successfully",
	})
}

// // GET /api/v1/urls/:shortcode/stats
// func GetStats(c *gin.Context){
// 	// shortcode := c.Param("shortcode")

// 	// response, err := service.GetStats(shortcode)

// 	// if err!= nil{
// 	// 	c.JSON(http.StatusNotFound, gin.H{"error":err.Error()})
// 	// 	return
// 	// }

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "stats endpoint",
// 	})
// }

// GET /health
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "UP",
		"message": "Service is healthy",
	})
}
