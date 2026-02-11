package memcached

import (
	"strconv"
	"testing"
	"time"
)

func setupMemcachedBench(b *testing.B) *MemcachedCache {
	mc, err := NewMemcachedCache("localhost:11211")
	if err != nil {
		b.Fatalf("Memcached connection failed: %v", err)
	}
	mc.Clear()
	return mc
}

func BenchmarkMemcachedSet(b *testing.B) {
	mc := setupMemcachedBench(b)
	defer mc.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Set("key"+strconv.Itoa(i), "value", 0)
	}
}

func BenchmarkMemcachedGet(b *testing.B) {
	mc := setupMemcachedBench(b)
	defer mc.Close()

	for i := 0; i < 100000; i++ {
		mc.Set("key"+strconv.Itoa(i), "value", 0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Get("key" + strconv.Itoa(i%100000))
	}
}

func BenchmarkMemcachedDelete(b *testing.B) {
	mc := setupMemcachedBench(b)
	defer mc.Close()

	for i := 0; i < b.N; i++ {
		mc.Set("key"+strconv.Itoa(i), "value", 0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Delete("key" + strconv.Itoa(i))
	}
}

func BenchmarkMemcachedSetWithTTL(b *testing.B) {
	mc := setupMemcachedBench(b)
	defer mc.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Set("key"+strconv.Itoa(i), "value", 5*time.Second)
	}
}
