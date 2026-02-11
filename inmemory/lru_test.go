package inmemory

import (
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	cache := NewLRUCache(2)

	err := cache.Set("TestData1", "Value1", 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := cache.Get("TestData1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if val != "Value1" {
		t.Fatalf("Expected Value1, got %v", val)
	}
}

func TestDelete(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("TestData1", "Value1", 0)
	err := cache.Delete("TestData1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = cache.Get("TestData1")
	if err == nil {
		t.Fatalf("Expected key to be deleted")
	}
}

func TestTTlExpiry(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("TestData1", "Value1", 2*time.Second)

	time.Sleep(3 * time.Second)

	_, err := cache.Get("TestData1")
	if err == nil {
		t.Fatalf("Expected key to expire")
	}
}

func TestLRUEviction(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("A", "1", 0)
	cache.Set("B", "2", 0)
	cache.Set("C", "3", 0) // Should evict A

	_, err := cache.Get("A")
	if err == nil {
		t.Fatalf("Expected A to be evicted")
	}
}

func TestClear(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("TestData1", "Value1", 0)
	cache.Set("TestData2", "Value2", 0)

	cache.Clear()

	_, err := cache.Get("TestData1")
	if err == nil {
		t.Fatalf("Expected cache to be cleared")
	}
}

func TestBackgroundCleanup(t *testing.T) {
	cache := NewLRUCache(2, 1*time.Second)

	cache.Set("TestData1", "Value1", 1*time.Second)

	time.Sleep(3 * time.Second)

	_, err := cache.Get("TestData1")
	if err == nil {
		t.Fatalf("Expected background cleanup to remove expired key")
	}
}
