package inmemory

import (
	"strconv"
	"testing"
	"time"
	cacheasync "github.com/dhanalakshms/multi-backend-cache-go/cache"
)

// Set without TTL
func BenchmarkLRUSet(b *testing.B) {
	cache := NewLRUCache(100000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i%100000) 
		cache.Set(key, "value", 0)
	}
}

// Get on populated cache
func BenchmarkLRUGet(b *testing.B) {
	cache := NewLRUCache(100000)

	for i := 0; i < 100000; i++ {
		cache.Set("key"+strconv.Itoa(i), "value", 0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("key" + strconv.Itoa(i%100000))
	}
}

// Delete operations
func BenchmarkLRUDelete(b *testing.B) {
	cache := NewLRUCache(100000)

	for i := 0; i < 100000; i++ {
		cache.Set("key"+strconv.Itoa(i), "value", 0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Delete("key" + strconv.Itoa(i%100000))
	}
}

// Set with TTL
func BenchmarkLRUSetWithTTL(b *testing.B) {
	cache := NewLRUCache(100000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i%100000)
		cache.Set(key, "value", 5*time.Second)
	}
}

// background cleanup impact
func BenchmarkLRUWithCleanup(b *testing.B) {
	cache := NewLRUCache(100000, 1*time.Second)

	for i := 0; i < 100000; i++ {
		cache.Set("key"+strconv.Itoa(i), "value", 1*time.Second)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Get("key" + strconv.Itoa(i%100000))
	}

	cache.StopCleanup()
}

// concurrent Set operations
func BenchmarkLRUConcurrentSet(b *testing.B) {
	cache := NewLRUCache(100000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Set("key"+strconv.Itoa(i%100000), "value", 0)
			i++
		}
	})
}

// concurrent workload 
func BenchmarkLRUMixedConcurrent(b *testing.B) {
	cache := NewLRUCache(100000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "key" + strconv.Itoa(i%100000)

			cache.Set(key, "value", 0)
			cache.Get(key)

			i++
		}
	})
}

// async set with TTL
func BenchmarkLRUAsyncSetWithTTL(b *testing.B) {
	cache := NewLRUCache(100000)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "key" + strconv.Itoa(i%100000)

			errCh := cacheasync.SetAsync(cache, key, "value", 5*time.Second)
			<-errCh

			i++
		}
	})
}

// async mixed workload
func BenchmarkLRUAsyncMixedWithTTL(b *testing.B) {
	cache := NewLRUCache(100000, 1*time.Second)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {

			key := "key" + strconv.Itoa(i%100000)

			setCh := cacheasync.SetAsync(cache, key, "value", 5*time.Second)
			<-setCh

			cache.Get(key)

			delCh := cacheasync.DeleteAsync(cache, key)
			<-delCh

			i++
		}
	})

	cache.StopCleanup()
}

