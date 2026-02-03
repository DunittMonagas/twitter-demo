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

	// List operations for timeline caching
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	LPush(ctx context.Context, key string, values ...interface{}) error
	RPush(ctx context.Context, key string, values ...interface{}) error
	LTrim(ctx context.Context, key string, start, stop int64) error
	LLen(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
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

// LRange retrieves a range of elements from a Redis list.
// Start and stop are zero-based indexes. Use -1 for the last element.
func (r *redisCache) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.LRange(ctx, key, start, stop).Result()
}

// LPush inserts values at the head of the list (prepends).
// Use this for Fan-Out to add new tweets at the beginning of the timeline.
func (r *redisCache) LPush(ctx context.Context, key string, values ...interface{}) error {
	return r.client.LPush(ctx, key, values...).Err()
}

// RPush inserts values at the tail of the list (appends, maintains order as provided).
// Use this to build the initial cache from database results.
func (r *redisCache) RPush(ctx context.Context, key string, values ...interface{}) error {
	return r.client.RPush(ctx, key, values...).Err()
}

// LTrim trims the list to only keep elements within the specified range.
// Use this to limit cache size (e.g., keep only the latest 1000 tweets).
func (r *redisCache) LTrim(ctx context.Context, key string, start, stop int64) error {
	return r.client.LTrim(ctx, key, start, stop).Err()
}

// LLen returns the length of the list.
func (r *redisCache) LLen(ctx context.Context, key string) (int64, error) {
	return r.client.LLen(ctx, key).Result()
}

// Expire sets the expiration time for a key.
func (r *redisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}
