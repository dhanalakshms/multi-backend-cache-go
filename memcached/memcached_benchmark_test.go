package memcached

import (
	"strconv"
	"testing"
	"time"
	cacheasync "github.com/dhanalakshms/multi-backend-cache-go/cache"
)

// Create fresh memcached instance
func setupMemcachedBench(b *testing.B) *MemcachedCache {
	mc, err := NewMemcachedCache("localhost:11211")
	if err != nil {
		b.Fatalf("Connection failed: %v", err)
	}
	mc.Clear()
	return mc
}

// Measure Set with TTL
func BenchmarkMemcachedSet(b *testing.B) {
	mc := setupMemcachedBench(b)
	defer mc.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Set("key"+strconv.Itoa(i), "value", 5*time.Second)
	}
}

// Measure Get on preloaded data
func BenchmarkMemcachedGet(b *testing.B) {
	mc := setupMemcachedBench(b)
	defer mc.Close()

	for i := 0; i < 100000; i++ {
		mc.Set("key"+strconv.Itoa(i), "value", 5*time.Second)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Get("key" + strconv.Itoa(i%100000))
	}
}

// Measure Delete
func BenchmarkMemcachedDelete(b *testing.B) {
	mc := setupMemcachedBench(b)
	defer mc.Close()

	for i := 0; i < b.N; i++ {
		mc.Set("key"+strconv.Itoa(i), "value", 5*time.Second)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Delete("key" + strconv.Itoa(i))
	}
}

// Concurrent Set
func BenchmarkMemcachedConcurrentSet(b *testing.B) {
	mc := setupMemcachedBench(b)
	defer mc.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			mc.Set("key"+strconv.Itoa(i%100000), "value", 5*time.Second)
			i++
		}
	})
}

// Async Set
func BenchmarkMemcachedAsyncSet(b *testing.B) {
	mc := setupMemcachedBench(b)
	defer mc.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			<-cacheasync.SetAsync(mc, "key"+strconv.Itoa(i%100000), "value", 5*time.Second)
			i++
		}
	})
}

// Async Delete
func BenchmarkMemcachedAsyncDelete(b *testing.B) {
	mc := setupMemcachedBench(b)
	defer mc.Close()

	for i := 0; i < 100000; i++ {
		mc.Set("key"+strconv.Itoa(i), "value", 5*time.Second)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			<-cacheasync.DeleteAsync(mc, "key"+strconv.Itoa(i%100000))
			i++
		}
	})
}

// Mixed async workload
func BenchmarkMemcachedAsyncMixed(b *testing.B) {
	mc := setupMemcachedBench(b)
	defer mc.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {

			key := "key" + strconv.Itoa(i%100000)

			<-cacheasync.SetAsync(mc, key, "value", 5*time.Second)

			mc.Get(key)

			<-cacheasync.DeleteAsync(mc, key)

			i++
		}
	})
}
