package database

import (
	"context"

	"github.com/Sidi1901/urlShortner/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg *config.Config) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDBNo,
	})

	_, err := rdb.Ping(context.Background()).Result()

	if err != nil {
		panic("Failed to Create Redis Client: " + err.Error())
	}

	return &RedisClient{Client: rdb}
}
