package service

type RateLimiter interface {
	IsAllowed(key string, capacity float64, refillRate float64) (bool, float64, error)
}
