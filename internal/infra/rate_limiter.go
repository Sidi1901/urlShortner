package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
	script *redis.Script
}

func NewRateLimiter(rdb *redis.Client, script string) *RateLimiter {
	return &RateLimiter{
		client: rdb,
		script: redis.NewScript(script),
	}
}

func (r *RateLimiter) IsAllowed(key string, capacity float64, refillRate float64) (bool, float64, error) {
	now := float64(time.Now().Unix())
	res, err := r.script.Run(context.Background(), r.client, []string{key}, capacity, refillRate, now, 1).Result()

	if err != nil {
		return false, 0, err
	}

	result := res.([]interface{})

	// Safely extract allowed as int64 and convert to bool
	allowedInt, ok := result[0].(int64)
	if !ok {
		return false, 0, fmt.Errorf("unexpected type for allowed: %T", result[0])
	}
	allowed := allowedInt == 1

	// Safely extract remaining as int64 and convert to float64
	remainingInt, ok := result[1].(int64)
	if !ok {
		return false, 0, fmt.Errorf("unexpected type for remaining: %T", result[1])
	}
	remaining := float64(remainingInt)

	return allowed, remaining, nil
}
