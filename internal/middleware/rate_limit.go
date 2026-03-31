package middleware

// import(
// 	"net"
// 	"strings"
// 	"github.com/Sidi1901/urlShortner/internal/ratelimiter"
// 	"github.com/gin-gonic/gin"

// )

// func GetClientIP(c *gin.Context) string {
// 	ip := c.GetHeader("X-Forwarded-For")

// 	if ip != "" {
// 		return strings.Split(ip, ",")[0]
// 	}

// 	ip = c.GetHeader("X-Real-IP")

// 	if ip != "" {
// 		return ip
// 	}

// 	host, _, _, net.SplitHostPort(c.Request.RemoteAddr)

// 	return host
// }

// func RateLimit(limiter rateLimiter.Limiter) gin.HandleFunc {
// 	return func(c *gin.Context){
// 		email := c.GetHeader("X-User-Email")
// 		ip := getClientIP(c	)

// 		allowed, err := limiter.Allow(c.Request.Context(), email, ip)

// 		if err != nil {
// 			c.AbortWithStatusJSON(500, gin.H{"error":"internal error"})
// 			return
// 		}

// 		if !allowed {
// 			c.AbortWithStatusJSON(429, gin.H{"error":"rate limit exceeded"})
// 			return
// 		}

// 		c.Next()
// 	}
// }
