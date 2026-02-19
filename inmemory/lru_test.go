package inmemory

import (
	"sync"
	"testing"
	"time"
	cacheasync "github.com/dhanalakshms/multi-backend-cache-go/cache"
)

// TestLRU_AllFunctionalities checks all basic behaviors of LRU cache
func TestLRU_AllFunctionalities(t *testing.T) {

	// Create cache with capacity 2 and background cleanup every 1 second
	cache := NewLRUCache(2, 1*time.Second)

	// Insert value and ensure it can be retrieved
	err := cache.Set("t", "test1", 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	val, err := cache.Get("t")
	if err != nil || val != "test1" {
		t.Fatal("Get success failed")
	}

	// Fetching unknown key should return error
	_, err = cache.Get("unknown")
	if err == nil {
		t.Fatal("Expected key not found")
	}

	// Key expires and is removed during Get
	cache.Set("ttl", "value", 1*time.Second)
	time.Sleep(2 * time.Second)

	_, err = cache.Get("ttl")
	if err == nil {
		t.Fatal("Expected key expired")
	}

	// Delete should remove entry completely
	cache.Set("x", "delete", 5*time.Second)
	cache.Delete("x")

	_, err = cache.Get("x")
	if err == nil {
		t.Fatal("Delete failed")
	}

	// Clear removes all entries and resets cache
	cache.Set("c1", "1", 5*time.Second)
	cache.Set("c2", "2", 5*time.Second)

	cache.Clear()

	_, err = cache.Get("c1")
	if err == nil {
		t.Fatal("Clear failed")
	}

	// Capacity=2 then the oldest key should be removed
	cache.Set("a", "A", 5*time.Second)
	cache.Set("b", "B", 5*time.Second)
	cache.Set("c", "C", 5*time.Second)

	_, err = cache.Get("a")
	if err == nil {
		t.Fatal("LRU eviction failed")
	}

	// Expired key removed automatically by cleanup goroutine
	cache.Set("bg", "value", 1*time.Second)
	time.Sleep(3 * time.Second)

	_, err = cache.Get("bg")
	if err == nil {
		t.Fatal("Background cleanup failed")
	}

	cache.StopCleanup()
}

// TestLRU_ConcurrentAccess checks thread safety
func TestLRU_ConcurrentAccess(t *testing.T) {

	cache := NewLRUCache(100, 1*time.Second)

	var wg sync.WaitGroup

	// Multiple goroutines access cache simultaneously
	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := "k" + string(rune(i))

			cache.Set(key, i, 5*time.Second)
			cache.Get(key)
			cache.Delete(key)

		}(i)
	}

	wg.Wait()
	cache.StopCleanup()
}

// TestLRU_AsyncOperations checks async Set and Delete wrappers
func TestLRU_AsyncOperations(t *testing.T) {

	cache := NewLRUCache(10)

	setResult := cacheasync.SetAsync(cache, "async", "value", 5*time.Second)

	if err := <-setResult; err != nil {
		t.Fatal("Async set failed")
	}

	// Verify async set worked
	val, err := cache.Get("async")
	if err != nil || val != "value" {
		t.Fatal("Async set verification failed")
	}

	delResult := cacheasync.DeleteAsync(cache, "async")

	if err := <-delResult; err != nil {
		t.Fatal("Async delete failed")
	}

	// Verify async delete worked
	_, err = cache.Get("async")
	if err == nil {
		t.Fatal("Async delete verification failed")
	}
}

// Async + concurrency 
func TestLRU_AsyncConcurrentWrites(t *testing.T) {

	cache := NewLRUCache(200)

	var wg sync.WaitGroup

	// Launch many async writes
	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := "async_" + string(rune(i))

			errCh := cacheasync.SetAsync(cache, key, i, 5*time.Second)

			if err := <-errCh; err != nil {
				t.Errorf("Async set failed: %v", err)
			}

		}(i)
	}

	wg.Wait()

	// Verify values exist after concurrent async writes
	for i := 0; i < 100; i++ {
		key := "async_" + string(rune(i))

		val, err := cache.Get(key)
		if err != nil {
			t.Fatalf("Missing key after async write: %s", key)
		}

		if val != i {
			t.Fatalf("Wrong value for %s", key)
		}
	}
}
