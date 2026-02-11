package inmemory

import (
	"strconv"
	"testing"
	"time"
)

func BenchmarkLRUSet(b *testing.B) {
	cache := NewLRUCache(100000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key"+strconv.Itoa(i), "value", 0)
	}
}

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

func BenchmarkLRUDelete(b *testing.B) {
	cache := NewLRUCache(b.N)

	for i := 0; i < b.N; i++ {
		cache.Set("key"+strconv.Itoa(i), "value", 0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Delete("key" + strconv.Itoa(i))
	}
}

func BenchmarkLRUSetWithTTL(b *testing.B) {
	cache := NewLRUCache(100000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key"+strconv.Itoa(i), "value", 5*time.Second)
	}
}

func BenchmarkLRUConcurrentSet(b *testing.B) {
	cache := NewLRUCache(100000)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Set("key"+strconv.Itoa(i), "value", 0)
			i++
		}
	})
}
