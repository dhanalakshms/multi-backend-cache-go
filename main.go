package main

import (
	"fmt"
	"time"

	"github.com/dhanalakshms/multi-backend-cache-go/cache"
	"github.com/dhanalakshms/multi-backend-cache-go/inmemory"
	memcachedbackend "github.com/dhanalakshms/multi-backend-cache-go/memcached"
	redisbackend "github.com/dhanalakshms/multi-backend-cache-go/redis"
)

func main() {

	fmt.Println("===== LRU Cache test =====")

	var c cache.Cache
	c = inmemory.NewLRUCache(2)

	fmt.Println("\n1) Basic SET + GET")
	c.Set("a", "apple", 0)

	val, ok := c.Get("a")
	fmt.Println("Get a:", val, ok)

	fmt.Println("\n2) LRU Eviction")
	c.Set("b", "Mango", 0)
	c.Set("c", "Orange", 0) 

	_, ok = c.Get("a")
	fmt.Println("a exists after eviction:", ok)


	fmt.Println("\n3) Delete test")
	c.Set("x", "delete-me", 0)

	c.Delete("x")
	_, ok = c.Get("x")
	fmt.Println("x exists after delete:", ok)


	fmt.Println("\n4) TTL expiry test")
	c.Set("ttl", "sample", 2*time.Second)

	val, ok = c.Get("ttl")
	fmt.Println("before expiry:", val, ok)

	time.Sleep(3 * time.Second)

	_, ok = c.Get("ttl")
	fmt.Println("after expiry:", ok)

	fmt.Println("\n5) Missing key test")
	_, ok = c.Get("unknown")
	fmt.Println("unknown exists:", ok)


	fmt.Println("\n===== LRU test Complete =====")

	// Redis Cache Testing
	testRedisCache()

	// Memcached Cache Testing
	testMemcachedCache()
}

func testRedisCache() {
	fmt.Println("\n\n===== Redis Cache test =====")

	// Create Redis cache connection
	rc, err := redisbackend.NewRedisCache("localhost:6379")
	if err != nil {
		fmt.Println("ERROR: Cannot connect to Redis:", err)
		fmt.Println("Make sure Redis is running on localhost:6379")
		return
	}
	defer rc.Close()

	// Clear all keys before testing
	rc.FlushAll()

	var c cache.Cache
	c = rc

	fmt.Println("\n1) Redis Basic SET + GET")
	c.Set("name", "Dhanalakshmi", 0)

	val, ok := c.Get("name")
	fmt.Println("Get name:", val, ok)

	fmt.Println("\n2) Redis SET with TTL")
	c.Set("session", "user123", 5*time.Second)

	val, ok = c.Get("session")
	fmt.Println("before expiry:", val, ok)

	fmt.Println("Waiting 6 seconds for TTL expiry...")
	time.Sleep(6 * time.Second)

	_, ok = c.Get("session")
	fmt.Println("after expiry:", ok)

	fmt.Println("\n3) Redis Delete test")
	c.Set("temp", "temporary-value", 0)

	c.Delete("temp")
	_, ok = c.Get("temp")
	fmt.Println("temp exists after delete:", ok)

	fmt.Println("\n4) Redis Multiple values")
	c.Set("user1", "Alice", 0)
	c.Set("user2", "Bob", 0)
	c.Set("user3", "Charlie", 0)

	for i := 1; i <= 3; i++ {
		key := fmt.Sprintf("user%d", i)
		val, ok := c.Get(key)
		fmt.Println(fmt.Sprintf("Get %s: %v, exists: %v", key, val, ok))
	}

	fmt.Println("\n5) Redis Advanced - GetTTL")
	c.Set("expiring", "value", 10*time.Second)
	ttl, _ := rc.GetTTL("expiring")
	fmt.Println("Remaining TTL for 'expiring':", ttl.Seconds(), "seconds")

	fmt.Println("\n6) Redis Advanced - Keys pattern")
	keys, _ := rc.Keys("user*")
	fmt.Println("Keys matching 'user*':", keys)

	fmt.Println("\n===== Redis test Complete =====")
}

func testMemcachedCache() {
	fmt.Println("\n\n===== Memcached Cache test =====")

	// Create Memcached cache connection
	mc, err := memcachedbackend.NewMemcachedCache("localhost:11211")
	if err != nil {
		fmt.Println("ERROR: Cannot connect to Memcached:", err)
		fmt.Println("Make sure Memcached is running on localhost:11211")
		return
	}
	defer mc.Close()

	// Clear all keys before testing
	mc.FlushAll()

	var c cache.Cache
	c = mc

	fmt.Println("\n1) Memcached Basic SET + GET")
	c.Set("product", "Laptop", 0)

	val, ok := c.Get("product")
	fmt.Println("Get product:", val, ok)

	fmt.Println("\n2) Memcached SET with TTL")
	c.Set("cache_key", "temp_data", 5*time.Second)

	val, ok = c.Get("cache_key")
	fmt.Println("before expiry:", val, ok)

	fmt.Println("Waiting 6 seconds for TTL expiry...")
	time.Sleep(6 * time.Second)

	_, ok = c.Get("cache_key")
	fmt.Println("after expiry:", ok)

	fmt.Println("\n3) Memcached Delete test")
	c.Set("delete_test", "removable", 0)

	c.Delete("delete_test")
	_, ok = c.Get("delete_test")
	fmt.Println("delete_test exists after delete:", ok)

	fmt.Println("\n4) Memcached Multiple values")
	c.Set("item1", "Apple", 0)
	c.Set("item2", "Banana", 0)
	c.Set("item3", "Cherry", 0)

	for i := 1; i <= 3; i++ {
		key := fmt.Sprintf("item%d", i)
		val, ok := c.Get(key)
		fmt.Println(fmt.Sprintf("Get %s: %v, exists: %v", key, val, ok))
	}

	fmt.Println("\n5) Memcached GetMultiple (batch retrieval)")
	mc.Set("id:1", "user_1", 0)
	mc.Set("id:2", "user_2", 0)
	mc.Set("id:3", "user_3", 0)

	items := mc.GetMultiple("id:1", "id:2", "id:3")
	fmt.Println("Batch retrieved items:", items)

	fmt.Println("\n6) Memcached connection successful (basic validation)")


	fmt.Println("\n===== Memcached test Complete =====")
}
