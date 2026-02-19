package redisbackend

import (
	"strconv"
	"sync"
	"testing"
	"time"
	cacheasync "github.com/dhanalakshms/multi-backend-cache-go/cache"
)

// Create Redis connection and start with empty DB
func setupRedis(t *testing.T) *RedisCache {
	rc, err := NewRedisCache("localhost:6379")
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	rc.Clear()
	return rc
}

// Basic connection test
func TestConnection(t *testing.T) {
	rc := setupRedis(t)
	defer rc.Close()
}

// Set + Get
func TestSetAndGet(t *testing.T) {
	rc := setupRedis(t)
	defer rc.Close()

	rc.Set("TestData1", "Value1", 5*time.Second)

	val, err := rc.Get("TestData1")
	if err != nil || val != "Value1" {
		t.Fatalf("Set/Get failed")
	}
}

// Delete test
func TestDelete(t *testing.T) {
	rc := setupRedis(t)
	defer rc.Close()

	rc.Set("TestData1", "Value1", 5*time.Second)
	rc.Delete("TestData1")

	_, err := rc.Get("TestData1")
	if err == nil {
		t.Fatalf("Expected delete")
	}
}

// TTL expiry
func TestTTLExpiry(t *testing.T) {
	rc := setupRedis(t)
	defer rc.Close()

	rc.Set("ttl", "value", 2*time.Second)

	time.Sleep(3 * time.Second)

	_, err := rc.Get("ttl")
	if err == nil {
		t.Fatalf("Expected expiry")
	}
}

// Clear DB
func TestClear(t *testing.T) {
	rc := setupRedis(t)
	defer rc.Close()

	for i := 0; i < 5; i++ {
		rc.Set("k"+strconv.Itoa(i), i, 5*time.Second)
	}

	rc.Clear()

	_, err := rc.Get("k0")
	if err == nil {
		t.Fatalf("Clear failed")
	}
}

// Close connection
func TestClose(t *testing.T) {
	rc := setupRedis(t)

	if err := rc.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

// Concurrent access
func TestRedisConcurrentAccess(t *testing.T) {
	rc := setupRedis(t)
	defer rc.Close()

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := "c_" + strconv.Itoa(i)

			rc.Set(key, i, 5*time.Second)
			rc.Get(key)
			rc.Delete(key)

		}(i)
	}

	wg.Wait()
}

// Async operations
func TestRedisAsyncOperations(t *testing.T) {
	rc := setupRedis(t)
	defer rc.Close()

	// set async to test non-blocking behavior
	setCh := cacheasync.SetAsync(rc, "async", "value", 5*time.Second)
	if err := <-setCh; err != nil {
		t.Fatal(err)
	}

	val, err := rc.Get("async")
	if err != nil || val != "value" {
		t.Fatal("Async set failed")
	}

	delCh := cacheasync.DeleteAsync(rc, "async")
	<-delCh

	_, err = rc.Get("async")
	if err == nil {
		t.Fatal("Async delete failed")
	}
}
