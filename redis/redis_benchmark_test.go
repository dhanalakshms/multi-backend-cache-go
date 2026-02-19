package redisbackend

import (
	"strconv"
	"testing"
	"time"

	cacheasync "github.com/dhanalakshms/multi-backend-cache-go/cache"
)

// Create fresh Redis instance
func setupRedisBench(b *testing.B) *RedisCache {
	rc, err := NewRedisCache("localhost:6379")
	if err != nil {
		b.Fatalf("Redis connection failed: %v", err)
	}
	rc.Clear()
	return rc
}

// Set with TTL
func BenchmarkRedisSet(b *testing.B) {
	rc := setupRedisBench(b)
	defer rc.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rc.Set("key"+strconv.Itoa(i), "value", 5*time.Second)
	}
}

// Get on preloaded data
func BenchmarkRedisGet(b *testing.B) {
	rc := setupRedisBench(b)
	defer rc.Close()

	for i := 0; i < 100000; i++ {
		rc.Set("key"+strconv.Itoa(i), "value", 5*time.Second)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rc.Get("key" + strconv.Itoa(i%100000))
	}
}

// Delete benchmark
func BenchmarkRedisDelete(b *testing.B) {
	rc := setupRedisBench(b)
	defer rc.Close()

	for i := 0; i < b.N; i++ {
		rc.Set("key"+strconv.Itoa(i), "value", 5*time.Second)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rc.Delete("key" + strconv.Itoa(i))
	}
}

// Concurrent set
func BenchmarkRedisConcurrentSet(b *testing.B) {
	rc := setupRedisBench(b)
	defer rc.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			rc.Set("key"+strconv.Itoa(i%100000), "value", 5*time.Second)
			i++
		}
	})
}

// Async set
func BenchmarkRedisAsyncSet(b *testing.B) {
	rc := setupRedisBench(b)
	defer rc.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			<-cacheasync.SetAsync(rc, "key"+strconv.Itoa(i%100000), "value", 5*time.Second)
			i++
		}
	})
}

// Async delete
func BenchmarkRedisAsyncDelete(b *testing.B) {
	rc := setupRedisBench(b)
	defer rc.Close()

	for i := 0; i < 100000; i++ {
		rc.Set("key"+strconv.Itoa(i), "value", 5*time.Second)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			<-cacheasync.DeleteAsync(rc, "key"+strconv.Itoa(i%100000))
			i++
		}
	})
}

// Mixed async workload
func BenchmarkRedisAsyncMixed(b *testing.B) {
	rc := setupRedisBench(b)
	defer rc.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {

			key := "key" + strconv.Itoa(i%100000)

			<-cacheasync.SetAsync(rc, key, "value", 5*time.Second)

			rc.Get(key)

			<-cacheasync.DeleteAsync(rc, key)

			i++
		}
	})
}
