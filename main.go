package main

import (
	"fmt"
	"time"

	"github.com/dhanalakshms/multi-backend-cache-go/inmemory"
)

func main() {

	fmt.Println("===== LRU Cache test =====")

	cache := inmemory.NewLRUCache(2)

	fmt.Println("\n1) Basic SET + GET")
	cache.Set("a", "apple", 0)

	val, ok := cache.Get("a")
	fmt.Println("Get a:", val, ok)

	fmt.Println("\n2) LRU Eviction")
	cache.Set("b", "Mango", 0)
	cache.Set("c", "Orange", 0) // should evict "a"

	_, ok = cache.Get("a")
	fmt.Println("a exists after eviction:", ok)


	fmt.Println("\n3) Delete test")
	cache.Set("x", "delete-me", 0)

	cache.Delete("x")
	_, ok = cache.Get("x")
	fmt.Println("x exists after delete:", ok)


	fmt.Println("\n4) TTL expiry test")
	cache.Set("ttl", "sample", 2*time.Second)

	val, ok = cache.Get("ttl")
	fmt.Println("before expiry:", val, ok)

	time.Sleep(3 * time.Second)

	_, ok = cache.Get("ttl")
	fmt.Println("after expiry:", ok)

	fmt.Println("\n5) Missing key test")
	_, ok = cache.Get("unknown")
	fmt.Println("unknown exists:", ok)


	fmt.Println("\n===== test Complete =====")
}
