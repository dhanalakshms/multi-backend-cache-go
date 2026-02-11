package redisbackend

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

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

func (rc *RedisCache) Get(key string) (interface{}, error) {
	val, err := rc.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key not found")
	}
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err == nil {
		return result, nil
	}

	return val, nil
}

func (rc *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		data = []byte(fmt.Sprintf("%v", value))
	}

	return rc.client.Set(context.Background(), key, data, ttl).Err()
}

func (rc *RedisCache) Delete(key string) error {
	return rc.client.Del(context.Background(), key).Err()
}

// 🔥 THIS MUST MATCH INTERFACE EXACTLY
func (rc *RedisCache) Clear() error {
	return rc.client.FlushDB(context.Background()).Err()
}

func (rc *RedisCache) Close() error {
	return rc.client.Close()
}
