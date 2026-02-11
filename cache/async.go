package cache

import "time"
// Async operations for cache
func SetAsync(c Cache, key string, value interface{}, ttl time.Duration) <-chan error {
	result := make(chan error, 1)

	go func() {
		result <- c.Set(key, value, ttl)
		close(result)
	}()

	return result
}

func DeleteAsync(c Cache, key string) <-chan error {
	result := make(chan error, 1)

	go func() {
		result <- c.Delete(key)
		close(result)
	}()

	return result
}
