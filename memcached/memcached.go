package memcached

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedCache struct {
	client *memcache.Client
}

func NewMemcachedCache(servers ...string) (*MemcachedCache, error) {
	if len(servers) == 0 {
		servers = []string{"localhost:11211"}
	}

	client := memcache.New(servers...)

	err := client.Set(&memcache.Item{
		Key:        "__ping__",
		Value:      []byte("ok"),
		Expiration: 1,
	})
	if err != nil {
		return nil, err
	}

	return &MemcachedCache{
		client: client,
	}, nil
}

func (mc *MemcachedCache) Get(key string) (interface{}, error) {
	item, err := mc.client.Get(key)
	if err == memcache.ErrCacheMiss {
		return nil, fmt.Errorf("key not found")
	}
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(item.Value, &result); err == nil {
		return result, nil
	}

	return string(item.Value), nil
}

func (mc *MemcachedCache) Set(key string, value interface{}, ttl time.Duration) error {
	expiration := int32(0)
	if ttl > 0 {
		expiration = int32(ttl.Seconds())
	}

	data, err := json.Marshal(value)
	if err != nil {
		data = []byte(fmt.Sprintf("%v", value))
	}

	item := &memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: expiration,
	}

	return mc.client.Set(item)
}

func (mc *MemcachedCache) Delete(key string) error {
	return mc.client.Delete(key)
}

func (mc *MemcachedCache) Clear() error {
	return mc.client.FlushAll()
}

func (mc *MemcachedCache) Close() error {
	return mc.client.Close()
}
