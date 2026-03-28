package ratelimiter

import (

)

type Limiter interface {
	Allow(ctx context.Context(), email string, ip string) (bool, error)
}

type rateLimiter struct {
	client *redis.Client
	window  time.Duration
}

func NewRateLimiter(client *redis.Client, window timeDuration) Limiter{
	return &rateLimiter{
		client:client,
		window: window,
	}
}

var script = redis.NewScript(
	`local current = redis.call("INCR", KEYS[1])
	if current == 1 then
		redis.call("EXPIRE",KEYS[1],ARGV[1])
	end
	return current`
) 

func (r *rateLimiter) Allow(ctx context.Context, email, ip string)(bool, error){
	key, limit := buildKeyAndLimit(email, ip)

	result, err := script.Run(ctx, r.client, []string{key}, int(r.window.Seconds()),).Result()

	if err != nil {
		return false, err
	}

	count := result.(int64)

	if count > int64(limit){
		return false, nil
	}

	return true, nil
}