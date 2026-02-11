package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dhanalakshms/multi-backend-cache-go/cache"
	"github.com/dhanalakshms/multi-backend-cache-go/inmemory"
	"github.com/dhanalakshms/multi-backend-cache-go/memcached"
	redisbackend "github.com/dhanalakshms/multi-backend-cache-go/redis"
)

func main() {

	backend := "lru"
	if len(os.Args) > 1 {
		backend = os.Args[1]
	}

	var c cache.Cache
	var dbName string

	switch backend {

	case "lru":
		dbName = "In-Memory LRU Cache"
		c = inmemory.NewLRUCache(2)

	case "redis":
		dbName = "Redis Cache"
		rc, err := redisbackend.NewRedisCache("localhost:6379")
		if err != nil {
			fmt.Println("Redis error:", err)
			return
		}
		defer rc.Close()
		c = rc

	case "memcached":
		dbName = "Memcached Cache"
		mc, err := memcached.NewMemcachedCache("localhost:11211")
		if err != nil {
			fmt.Println("Memcached error:", err)
			return
		}
		defer mc.Close()
		c = mc

	default:
		fmt.Println("Unknown backend. Use: lru | redis | memcached")
		return
	}

	fmt.Println("======================================")
	fmt.Println("Using Backend:", dbName)
	fmt.Println("======================================")

	// ------------------------------------------------
	// 1️⃣ Basic Set & Get
	// ------------------------------------------------
	fmt.Println("\n[1] Basic Set & Get")
	c.Set("TestKey1", "TestData1", 0)

	val, err := c.Get("TestKey1")
	if err != nil {
		fmt.Println("FAILED:", err)
	} else {
		fmt.Println("PASSED: Retrieved =", val)
	}

	// ------------------------------------------------
	// 2️⃣ Overwrite Existing Key
	// ------------------------------------------------
	fmt.Println("\n[2] Overwrite Existing Key")
	c.Set("TestKey1", "UpdatedData1", 0)

	val, _ = c.Get("TestKey1")
	fmt.Println("PASSED: Updated value =", val)

	// ------------------------------------------------
	// 3️⃣ LRU Eviction Test (Capacity = 2)
	// ------------------------------------------------
	fmt.Println("\n[3] LRU Eviction Test")
	c.Set("TestKey2", "TestData2", 0)
	c.Set("TestKey3", "TestData3", 0) // Should evict TestKey1

	_, err = c.Get("TestKey1")
	if err != nil {
		fmt.Println("PASSED: TestKey1 evicted correctly")
	} else {
		fmt.Println("FAILED: Eviction did not occur")
	}

			// ------------------------------------------------
		// 4️⃣ TTL Expiry Test (6 Seconds Wait)
		// ------------------------------------------------
		fmt.Println("\n[4] TTL Expiry Test")

		ttlDuration := 5 * time.Second
		fmt.Println("Setting key 'TTLKey' with TTL =", ttlDuration)

		err = c.Set("TTLKey", "TestDataTTL", ttlDuration)
		if err != nil {
			fmt.Println("Set error:", err)
			return
		}

		// Before expiry
		val, err = c.Get("TTLKey")
		if err != nil {
			fmt.Println("FAILED: Unexpected error before expiry:", err)
		} else {
			fmt.Println("Before TTL Expiry → Retrieved:", val)
		}

		fmt.Println("Waiting 6 seconds for TTL to expire...")
		time.Sleep(6 * time.Second)

		// After expiry
		val, err = c.Get("TTLKey")
		if err != nil {
			fmt.Println("After TTL Expiry → Expected Error:", err)
		} else {
			fmt.Println("FAILED: TTL did not expire. Value:", val)
		}


	// ------------------------------------------------
	// 5️⃣ Delete Test
	// ------------------------------------------------
	fmt.Println("\n[5] Delete Test")
	c.Set("DeleteKey", "DeleteData", 0)
	c.Delete("DeleteKey")

	_, err = c.Get("DeleteKey")
	if err != nil {
		fmt.Println("PASSED: Delete successful")
	} else {
		fmt.Println("FAILED: Delete failed")
	}

	// ------------------------------------------------
	// 6️⃣ Clear Test
	// ------------------------------------------------
	fmt.Println("\n[6] Clear Test")
	c.Set("ClearKey1", "ClearData1", 0)
	c.Set("ClearKey2", "ClearData2", 0)

	c.Clear()

	_, err = c.Get("ClearKey1")
	if err != nil {
		fmt.Println("PASSED: Cache cleared successfully")
	} else {
		fmt.Println("FAILED: Clear did not work")
	}

	// ------------------------------------------------
	// 7️⃣ Missing Key Test
	// ------------------------------------------------
	fmt.Println("\n[7] Missing Key Test")
	_, err = c.Get("NonExistingKey")
	if err != nil {
		fmt.Println("PASSED: Missing key handled correctly")
	} else {
		fmt.Println("FAILED: Missing key should return error")
	}

	fmt.Println("\n======================================")
	fmt.Println("All tests completed successfully for", dbName)
	fmt.Println("======================================")
}
