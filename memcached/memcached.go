package memcached

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// MemcachedCache implements the Cache interface using Memcached
type MemcachedCache struct {
	client *memcache.Client
}

// NewMemcachedCache initializes a Memcached cache connection
// servers format: []string{"localhost:11211"} or {"server1:11211", "server2:11211"}
func NewMemcachedCache(servers ...string) (*MemcachedCache, error) {
	if len(servers) == 0 {
		servers = []string{"localhost:11211"}
	}

	client := memcache.New(servers...)

	// Test connection using a lightweight set operation
	err := client.Set(&memcache.Item{
		Key:        "__connection_test__",
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


// Get retrieves a value from Memcached
func (mc *MemcachedCache) Get(key string) (interface{}, bool) {
	item, err := mc.client.Get(key)
	if err == memcache.ErrCacheMiss {
		return nil, false
	}
	if err != nil {
		return nil, false
	}

	// Try to unmarshal JSON first
	var result interface{}
	if err := json.Unmarshal(item.Value, &result); err == nil {
		return result, true
	}

	// If not JSON, return as string
	return string(item.Value), true
}

// Set inserts or updates a key-value pair with optional TTL
// Note: Memcached has max expiration time of ~30 days
func (mc *MemcachedCache) Set(key string, value interface{}, ttl time.Duration) {
	// Convert TTL to seconds for Memcached
	expiration := int32(0) // 0 means no expiration
	if ttl > 0 {
		expiration = int32(ttl.Seconds())
	}

	// Marshal value to JSON
	jsonVal, err := json.Marshal(value)
	if err != nil {
		// Fall back to string conversion
		jsonVal = []byte(fmt.Sprintf("%v", value))
	}

	item := &memcache.Item{
		Key:        key,
		Value:      jsonVal,
		Expiration: expiration,
	}

	mc.client.Set(item)
}

// Delete removes a key from Memcached
func (mc *MemcachedCache) Delete(key string) {
	mc.client.Delete(key)
}

// Close closes the Memcached connection
func (mc *MemcachedCache) Close() error {
	return mc.client.Close()
}

// FlushAll deletes all items from Memcached
func (mc *MemcachedCache) FlushAll() error {
	return mc.client.FlushAll()
}

// GetMultiple retrieves multiple values at once
func (mc *MemcachedCache) GetMultiple(keys ...string) map[string]interface{} {
	items, err := mc.client.GetMulti(keys)
	if err != nil && err != memcache.ErrCacheMiss {
		return make(map[string]interface{})
	}

	result := make(map[string]interface{})
	for key, item := range items {
		// Try to unmarshal JSON
		var val interface{}
		if err := json.Unmarshal(item.Value, &val); err == nil {
			result[key] = val
		} else {
			result[key] = string(item.Value)
		}
	}

	return result
}
