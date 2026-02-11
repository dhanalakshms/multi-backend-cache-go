package memcached

import (
	"testing"
	"time"
)

func setupMemcached(t *testing.T) *MemcachedCache {
	mc, err := NewMemcachedCache("localhost:11211")
	if err != nil {
		t.Fatalf("Failed to connect to Memcached: %v", err)
	}
	return mc
}

func TestConnection(t *testing.T) {
	mc := setupMemcached(t)
	defer mc.Close()
}

func TestSetAndGet(t *testing.T) {
	mc := setupMemcached(t)
	defer mc.Close()

	err := mc.Set("TestData1", "Value1", 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := mc.Get("TestData1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if val != "Value1" {
		t.Fatalf("Expected Value1, got %v", val)
	}
}

func TestDelete(t *testing.T) {
	mc := setupMemcached(t)
	defer mc.Close()

	mc.Set("TestData1", "Value1", 0)

	err := mc.Delete("TestData1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = mc.Get("TestData1")
	if err == nil {
		t.Fatalf("Expected key to be deleted")
	}
}

func TestTTLExpiry(t *testing.T) {
	mc := setupMemcached(t)
	defer mc.Close()

	mc.Set("TestData1", "Value1", 2*time.Second)

	time.Sleep(3 * time.Second)

	_, err := mc.Get("TestData1")
	if err == nil {
		t.Fatalf("Expected key to expire")
	}
}

func TestClear(t *testing.T) {
	mc := setupMemcached(t)
	defer mc.Close()

	mc.Set("TestData1", "Value1", 0)
	mc.Set("TestData2", "Value2", 0)

	err := mc.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	_, err = mc.Get("TestData1")
	if err == nil {
		t.Fatalf("Expected cache to be cleared")
	}
}

func TestClose(t *testing.T) {
	mc := setupMemcached(t)

	err := mc.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}
