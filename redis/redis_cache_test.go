package redisbackend

import (
	"testing"
	"time"
)

func setupRedis(t *testing.T) *RedisCache {
	rc, err := NewRedisCache("localhost:6379")
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	return rc
}

func TestConnection(t *testing.T) {
	rc := setupRedis(t)
	defer rc.Close()
}

func TestSetAndGet(t *testing.T) {
	rc := setupRedis(t)
	defer rc.Close()

	err := rc.Set("TestData1", "Value1", 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := rc.Get("TestData1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if val != "Value1" {
		t.Fatalf("Expected Value1, got %v", val)
	}
}

func TestDelete(t *testing.T) {
	rc := setupRedis(t)
	defer rc.Close()

	rc.Set("TestData1", "Value1", 0)

	err := rc.Delete("TestData1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = rc.Get("TestData1")
	if err == nil {
		t.Fatalf("Expected key to be deleted")
	}
}

func TestTTLExpiry(t *testing.T) {
	rc := setupRedis(t)
	defer rc.Close()

	rc.Set("TestData1", "Value1", 2*time.Second)

	time.Sleep(3 * time.Second)

	_, err := rc.Get("TestData1")
	if err == nil {
		t.Fatalf("Expected key to expire")
	}
}

func TestClear(t *testing.T) {
	rc := setupRedis(t)
	defer rc.Close()

	rc.Set("TestData1", "Value1", 0)
	rc.Set("TestData2", "Value2", 0)

	err := rc.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	_, err = rc.Get("TestData1")
	if err == nil {
		t.Fatalf("Expected cache to be cleared")
	}
}

func TestClose(t *testing.T) {
	rc := setupRedis(t)

	err := rc.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}
