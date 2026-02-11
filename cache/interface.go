package cache

import "time"

// Cache defines the unified caching interface
type Cache interface {

	// Get retrieves a value by key.
	// Returns error if key not found or expired.
	Get(key string) (interface{}, error)

	// Set stores a key-value pair with optional TTL.
	Set(key string, value interface{}, ttl time.Duration) error

	// Delete removes a key from the cache.
	Delete(key string) error

	// Clear removes all entries from the cache.
	Clear() error
}
