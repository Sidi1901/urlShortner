package handler

import(
	"fmt"
	"net/http"
	"time"
	"github.com/Sidi1901/urlShortner/internal/service"
	"github.com/Sidi1901/urlShortner/internal/config"
	"github.com/Sidi1901/urlShortner/internal/dto"
	"github.com/gin-gonic/gin"
)

var cfg *config.Config



// GET /:shortcode
func ResolveURL(c *gin.Context){
	shortcode := c.Param("shortcode")

	mappedURL, err := service.GetOriginalURL(shortcode)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL no found"})
		return
	} 

	// 302
	c.Redirect(http.StatusFound, mappedURL)
}

// POST /api/v1/urls
func CreateShortURL(c *gin.Context) {
	var request dto.CreateShortURLRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	uRL := request.URL
	ShortCode := request.ShortCode
	ExpirySec := time.Duration(*request.ExpirySec) * time.Second
	ip := c.ClientIP()
	

	shortURL, err := service.CreateShortURL(uRL, ip, ExpirySec, *ShortCode)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.CreateShortURLResponse{
		URL: uRL,
		ShortCode : shortCode,
		ShortURL  : shortURL,
		Expiry    : ExpirySec,
		CreatedAt : time.Now(),
		UpdatedAt : time.Now()
	}

	c.JSON(http.StatusCreated, response)
}

// GET /api/v1/urls/:shortCode
func GetURLDetails(c *gin.Context){
	shortcode := c.Param("shortcode")

	shortURLData, err := service.GetURL(shortcode)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := dto.ShortURLDataResponse{
		URL          : shortURLData.URL,
		ShortCode    : shortCode,
		ShortURL     : shortURLData.ShortURL,
		Expiry       : shortURLData.Expiry,
		CreatedAt    : shortURLData.CreatedAt,
		LastUpdatedAt: shortURLData.LastUpdatedAt
	}

	c.JSON(http.StatusFound, response)
}


// PUT /api/v1/urls/:shortcode
func UpdateURL(c *gin.Context){
	var request dto.UpdateShortURLRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error", err.Error()})
	}

	err := service.UpdateURL(request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error", err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

// DELETE /api/v1/urls/:shortCode
func DeleteURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	err := service.DeleteURL(shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted successfully",
	})
}

// GET /api/v1/urls/:shortcode/stats
func GetStats(c *gin.Context){
	// shortcode := c.Param("shortcode")

	// response, err := service.GetStats(shortcode)

	// if err!= nil{
	// 	c.JSON(http.StatusNotFound, gin.H{"error":err.Error()})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "stats endpoint",
	})
}

// GET /health
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "UP",
		"message": "Service is healthy",
	})
}



