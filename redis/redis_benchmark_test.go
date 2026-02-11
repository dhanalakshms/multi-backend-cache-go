package redisbackend

import (
	"strconv"
	"testing"
	"time"
)

func setupRedisBench(b *testing.B) *RedisCache {
	rc, err := NewRedisCache("localhost:6379")
	if err != nil {
		b.Fatalf("Redis connection failed: %v", err)
	}
	rc.Clear()
	return rc
}

func BenchmarkRedisSet(b *testing.B) {
	rc := setupRedisBench(b)
	defer rc.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rc.Set("key"+strconv.Itoa(i), "value", 0)
	}
}

func BenchmarkRedisGet(b *testing.B) {
	rc := setupRedisBench(b)
	defer rc.Close()

	for i := 0; i < 100000; i++ {
		rc.Set("key"+strconv.Itoa(i), "value", 0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rc.Get("key" + strconv.Itoa(i%100000))
	}
}

func BenchmarkRedisDelete(b *testing.B) {
	rc := setupRedisBench(b)
	defer rc.Close()

	for i := 0; i < b.N; i++ {
		rc.Set("key"+strconv.Itoa(i), "value", 0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rc.Delete("key" + strconv.Itoa(i))
	}
}

func BenchmarkRedisSetWithTTL(b *testing.B) {
	rc := setupRedisBench(b)
	defer rc.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rc.Set("key"+strconv.Itoa(i), "value", 5*time.Second)
	}
}
