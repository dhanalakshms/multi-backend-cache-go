package redisbackend

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements the Cache interface using Redis
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(addr string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}, nil
}

func (rc *RedisCache) Get(key string) (interface{}, bool) {
	val, err := rc.client.Get(rc.ctx, key).Result()
	if err == redis.Nil {
		return nil, false
	}
	if err != nil {
		return nil, false
	}

	// Try to unmarshal JSON first
	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err == nil {
		return result, true
	}

	// If not JSON, return as string
	return val, true
}

// Set inserts or updates a key-value pair with optional TTL
func (rc *RedisCache) Set(key string, value interface{}, ttl time.Duration) {
	// Marshal value to JSON
	jsonVal, err := json.Marshal(value)
	if err != nil {
		// Fall back to string conversion
		rc.client.Set(rc.ctx, key, value, ttl)
		return
	}

	rc.client.Set(rc.ctx, key, jsonVal, ttl)
}

// Delete removes a key from Redis
func (rc *RedisCache) Delete(key string) {
	rc.client.Del(rc.ctx, key)
}

// Close closes the Redis connection
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}

// FlushAll deletes all keys in the current database
func (rc *RedisCache) FlushAll() error {
	return rc.client.FlushAll(rc.ctx).Err()
}

// Keys returns all keys matching the pattern
func (rc *RedisCache) Keys(pattern string) ([]string, error) {
	return rc.client.Keys(rc.ctx, pattern).Val(), nil
}

// GetTTL returns the remaining TTL for a key
func (rc *RedisCache) GetTTL(key string) (time.Duration, error) {
	return rc.client.TTL(rc.ctx, key).Val(), nil
}
