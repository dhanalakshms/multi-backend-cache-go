package memcached

import (
	"strconv"
	"sync"
	"testing"
	"time"
	cacheasync "github.com/dhanalakshms/multi-backend-cache-go/cache"
)

// helper to create new memcached instance
func setupMemcached(t *testing.T) *MemcachedCache {
	mc, err := NewMemcachedCache("localhost:11211")
	if err != nil {
		t.Fatalf("Failed to connect to Memcached: %v", err)
	}
	mc.Clear() // ensure clean state
	return mc
}

// connection sanity check
func TestConnection(t *testing.T) {
	mc := setupMemcached(t)
	defer mc.Close()
}

// basic set/get functionality
func TestSetAndGet(t *testing.T) {
	mc := setupMemcached(t)
	defer mc.Close()

	if err := mc.Set("TestData1", "Value1", 5*time.Second); err != nil {
		t.Fatal(err)
	}

	val, err := mc.Get("TestData1")
	if err != nil || val != "Value1" {
		t.Fatalf("Set/Get failed")
	}
}

// delete operation validation
func TestDelete(t *testing.T) {
	mc := setupMemcached(t)
	defer mc.Close()

	mc.Set("TestData1", "Value1", 5*time.Second)
	mc.Delete("TestData1")

	_, err := mc.Get("TestData1")
	if err == nil {
		t.Fatalf("Expected delete")
	}
}

// TTL expiry validation
func TestTTLExpiry(t *testing.T) {
	mc := setupMemcached(t)
	defer mc.Close()

	mc.Set("ttl", "value", 2*time.Second)
	time.Sleep(3 * time.Second)

	_, err := mc.Get("ttl")
	if err == nil {
		t.Fatalf("Expected expiry")
	}
}

// clear / flush validation
func TestClear(t *testing.T) {
	mc := setupMemcached(t)
	defer mc.Close()

	mc.Set("a", "1", 5*time.Second)
	mc.Set("b", "2", 5*time.Second)

	mc.Clear()

	if _, err := mc.Get("a"); err == nil {
		t.Fatalf("Clear failed")
	}
}

// concurrent sync operations
func TestMemcachedConcurrentAccess(t *testing.T) {
	mc := setupMemcached(t)
	defer mc.Close()

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := "c_" + strconv.Itoa(i)

			mc.Set(key, strconv.Itoa(i), 5*time.Second)
			mc.Get(key)
			mc.Delete(key)
		}(i)
	}

	wg.Wait()
}

// async set/delete validation
func TestMemcachedAsyncOperations(t *testing.T) {
	mc := setupMemcached(t)
	defer mc.Close()

	setCh := cacheasync.SetAsync(mc, "async", "value", 5*time.Second)
	if err := <-setCh; err != nil {
		t.Fatal(err)
	}

	val, err := mc.Get("async")
	if err != nil || val != "value" {
		t.Fatal("Async set failed")
	}

	delCh := cacheasync.DeleteAsync(mc, "async")
	<-delCh

	_, err = mc.Get("async")
	if err == nil {
		t.Fatal("Async delete failed")
	}
}

// async + concurrency stress
func TestMemcachedAsyncConcurrentWrites(t *testing.T) {
	mc := setupMemcached(t)
	defer mc.Close()

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := "async_" + strconv.Itoa(i)
			<-cacheasync.SetAsync(mc, key, strconv.Itoa(i), 5*time.Second)

		}(i)
	}

	wg.Wait()

	for i := 0; i < 50; i++ {
		key := "async_" + strconv.Itoa(i)

		if _, err := mc.Get(key); err != nil {
			t.Fatalf("Missing key %s", key)
		}
	}
}
