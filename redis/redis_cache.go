package redisbackend

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"github.com/go-redis/redis/v8"
)

// implementing a cache using Redis as the backend.
type RedisCache struct {
	client *redis.Client
}

// creating new RedisCache instance connected to the specified address.
func NewRedisCache(addr string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
	}, nil
}

// Get retrieves the value associated with the given key from Redis. 
func (rc *RedisCache) Get(key string) (interface{}, error) {
	val, err := rc.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key not found")
	}
	if err != nil {
		return nil, err
	}
	// Attempt to unmarshal JSON, if it fails return the raw string value.
	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err == nil {
		return result, nil
	}

	return val, nil
}

// Set stores the value with the specified key in Redis, with an optional TTL.
func (rc *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {

	// marshal the value to JSON. If it fails, store the raw string representation.
	data, err := json.Marshal(value) 
	if err != nil {
		data = []byte(fmt.Sprintf("%v", value))
	}

	return rc.client.Set(context.Background(), key, data, ttl).Err()
}

// Delete removes the specified key from Redis.
func (rc *RedisCache) Delete(key string) error {
	return rc.client.Del(context.Background(), key).Err()
}

// Clear flushes the entire Redis database, removing all keys.
func (rc *RedisCache) Clear() error {
	return rc.client.FlushDB(context.Background()).Err()
}

// Close closes the Redis client connection.
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}
