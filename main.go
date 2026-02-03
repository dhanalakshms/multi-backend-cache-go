package main

import (
	"fmt"
	"time"

	"github.com/dhanalakshms/multi-backend-cache-go/cache"
	"github.com/dhanalakshms/multi-backend-cache-go/inmemory"
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
	c.Set("c", "Orange", 0) /

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


	fmt.Println("\n===== test Complete =====")
}
