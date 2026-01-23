package main

import (
	"fmt"

	"github.com/dhanalakshms/multi-backend-cache-go/inmemory"
)

func main() {
	cache := inmemory.NewLRUCache(2)

	cache.Set("1", 1)
	cache.Set("2", 2)

	val, _ := cache.Get("1")
	fmt.Println(val) // 1

	cache.Set("3", 3)
	_, ok := cache.Get("2")
	fmt.Println(ok) // false

	cache.Set("4", 4)
	_, ok = cache.Get("1")
	fmt.Println(ok) // false

	val, _ = cache.Get("3")
	fmt.Println(val) // 3

	val, _ = cache.Get("4")
	fmt.Println(val) // 4
}
