package middleware

import (
	"fmt"
	"net/http"

	"github.com/Sidi1901/urlShortner/internal/service"
	"github.com/gin-gonic/gin"
)

type RateLimitMiddleware struct {
	limiter service.RateLimiter
}

func NewRateLimitMiddlware(l service.RateLimiter) *RateLimitMiddleware {
	return &RateLimitMiddleware{limiter: l}
}

func (m *RateLimitMiddleware) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {

		// ip based rate limiting
		ip := c.ClientIP()
		userID, _ := c.Get(ContextUserIDKey)

		ipKey := fmt.Sprintf("rate_limit:ip:%s", ip)
		ipCapacity := 5.0
		ipRefill := 1.0

		allowed, _, err := m.limiter.IsAllowed(ipKey, ipCapacity, ipRefill)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests (ip limit)",
			})
			return
		}

		// User ID based rate limiting
		if userID != "" {

			userKey := fmt.Sprintf("rate_limit:user:%s", userID)

			isPremium := false // TODO: fetch from DB

			capacity := 10.0
			refill := 0.4

			if isPremium {
				capacity = 100
				refill = 10
			}

			allowed, remaining, err := m.limiter.IsAllowed(userKey, capacity, refill)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
				return
			}

			if !allowed {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"error": "too many requests (user limit)",
				})
				return
			}

			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%.0f", remaining))
		}

		c.Next()

	}
}
