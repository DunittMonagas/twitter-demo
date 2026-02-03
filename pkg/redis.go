package pkg

import (
	"context"
	"time"

	"twitter-demo/internal/config"

	"github.com/redis/go-redis/v9"
)

// Cache defines the interface for cache operations.
// External code should depend on this interface, not on the concrete implementation.
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

// redisCache is the concrete implementation of Cache using go-redis.
type redisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache instance.
func NewRedisCache(cfg config.RedisConfig) Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &redisCache{
		client: client,
	}
}

// Get retrieves a value from Redis by key.
// Returns redis.Nil error if the key does not exist.
func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set stores a key-value pair in Redis with an optional expiration time.
// Use expiration = 0 for no expiration.
func (r *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Delete removes a key from Redis.
func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
