package handler

import (
	"net/http"
	"time"

	"github.com/Sidi1901/urlShortner/internal/dto"
	"github.com/gin-gonic/gin"
)

// GET /:shortcode
// func (h *Handler) ResolveURL(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	shortcode := c.Param("shortcode")

// 	mappedURL, err := h.service.GetOriginalURL(ctx, shortcode)

// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "URL no found"})
// 		return
// 	}

// 	// 302
// 	c.Redirect(http.StatusFound, mappedURL)
// }

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

// // GET /api/v1/urls/:shortCode
// func GetShortURL(c *gin.Context){
// 	ctx := c.Request.Context()
// 	shortcode := c.Param("shortcode")

// 	shortURLData, err := service.GetShortURL(ctx, shortcode)

// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 		return
// 	}

// 	response := dto.ShortURLDataResponse{
// 		URL          : shortURLData.URL,
// 		ShortCode    : shortCode,
// 		ShortURL     : shortURLData.ShortURL,
// 		Expiry       : shortURLData.Expiry,
// 		CreatedAt    : shortURLData.CreatedAt,
// 		LastUpdatedAt: shortURLData.LastUpdatedAt
// 	}

// 	c.JSON(http.StatusFound, response)
// }

// // PUT /api/v1/urls/:shortcode
// func UpdateURL(c *gin.Context){
// 	ctx := c.Request.Context()
// 	var request dto.UpdateShortURLRequest

// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error", err.Error()})
// 	}

// 	err := service.UpdateURL(ctx, request)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error", err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
// }

// // DELETE /api/v1/urls/:shortCode
// func DeleteURL(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	shortCode := c.Param("shortCode")

// 	err := service.DeleteURL(ctx, shortCode)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "deleted successfully",
// 	})
// }

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
