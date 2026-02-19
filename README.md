
# Multi-Backend Caching Library in Go

A high-performance, pluggable caching library in Go supporting multiple backends:

- In-Memory LRU Cache
- Redis
- Memcached

The library provides a unified API for cache operations and allows seamless switching between backends without modifying application logic.

---

## 🚀 Features

- In-Memory LRU Cache with O(1) operations
- Redis integration using go-redis/v8
- Memcached integration using gomemcache
- Unified Cache Interface
- TTL (Time-To-Live) support
- Manual invalidation (Delete, Clear)
- Thread-safe implementation
- supports sync and async operations
- Performance benchmark suite
- Pluggable backend selection
- Clean and intuitive API design

---

## 🏗 Architecture Overview

The library follows a unified interface design:

```go
type Cache interface {
    Set(key string, value interface{}, ttl time.Duration) error
    Get(key string) (interface{}, error)
    Delete(key string) error
    Clear() error
}
````

Each backend implements this interface:

* `inmemory` → In-Memory LRU implementation
* `redis` → Redis backend
* `memcached` → Memcached backend

This design ensures seamless backend switching.

---

## 📘 Comprehensive API Documentation

### 1️⃣ Set

```go
Set(key string, value interface{}, ttl time.Duration) error
```
Stores a key-value pair with optional TTL.

* `key` → Unique cache key
* `value` → Any Go type
* `ttl` → Expiration duration (0 = no expiration)

---

### 2️⃣ Get

```go
Get(key string) (interface{}, error)
```

Retrieves value by key.

* Returns error if:

  * Key not found
  * Key expired

---

### 3️⃣ Delete

```go
Delete(key string) error
```

Manually removes a key from cache.

---

### 4️⃣ Clear

```go
Clear() error
```

Clears entire cache.

---

## 🐳 Running Redis & Memcached using Docker

### Redis

```bash
docker run -d -p 6379:6379 redis
```

### Memcached

```bash
docker run -d -p 11211:11211 memcached
```

Your library connects using:

```
localhost:6379
localhost:11211
```

Containers ensure consistent testing and benchmarking environments.

---

## 🧪 Usage Guide & Examples

### Using In-Memory LRU

```go
import "github.com/dhanalakshms/multi-backend-cache-go/inmemory"

cache := inmemory.NewLRUCache(100)
cache.Set("user", "data", 5*time.Second)

value, err := cache.Get("user")
if err != nil {
    fmt.Println("Error:", err)
}

fmt.Println("Value:", value)
```

---

### Using Redis
install Redis server and go-redis/v8 client library before running this example.

```go
import redisbackend "github.com/dhanalakshms/multi-backend-cache-go/redis"

rc, _ := redisbackend.NewRedisCache("localhost:6379")
defer rc.Close()
rc.Set("key", "value", 5*time.Second)
val, _ := rc.Get("key")
```

---

### Using Memcached
install memcached server and client library before running this example.

```go
import "github.com/dhanalakshms/multi-backend-cache-go/memcached"

mc, _ := memcached.NewMemcachedCache("localhost:11211")
defer mc.Close()
mc.Set("key", "value", 5*time.Second)
val, _ := mc.Get("key")
```

---

### Switching Backend Easily

```go
var cache cache.Cache

cache = inmemory.NewLRUCache(100)
// OR
cache = redisCache
// OR
cache = memcachedCache
```

No application logic changes required.

---

## 📊 Benchmark Results

### In-Memory LRU

| Operation      | Latency |
| -------------- | ------- |
| Set            | ~250 ns |
| Get            | ~90 ns  |
| Delete         | ~145 ns |
| Concurrent Set | ~436 ns |

Sub-microsecond latency confirms O(1) design.

---

### Redis (Localhost)

| Operation      | Latency |
| -------------- | ------- |
| Set            | ~365 µs |
| Get            | ~302 µs |
| Delete         | ~405 µs |
| Concurrent Set | ~86 µs  |
| Async Mixed    | ~308 µs |

Network latency dominates performance.

---

### Memcached (Localhost)

| Operation      | Latency  |
| -------------- | -------- |
| Set            | ~2.10 ms |
| Get            | ~1.13 ms |
| Delete         | ~0.76 ms |
| Concurrent Set | ~0.18 ms |
| Async Mixed    | ~0.66 ms |

Performance consistent with network-based caching systems.

---

## 🧵 Thread Safety

* LRU uses mutex locking
* External caches rely on client concurrency
* Benchmarks validate stability under parallel workloads

---

## 🔄 Cache Invalidation & Expiration

* TTL support for automatic expiration
* Manual invalidation using `Delete`
* Complete cache reset using `Clear`
* Background cleanup (LRU)

---

## 🛠 Best Practices for Integration

1. Use In-Memory LRU for:

   * High-frequency, low-latency local caching
   * Single-instance applications

2. Use Redis for:

   * Distributed systems
   * Shared caching across services
   * Persistence requirements

3. Use Memcached for:

   * High-speed distributed caching
   * Stateless microservices

4. Keep TTL values meaningful:

   * Avoid extremely short TTLs
   * Balance freshness and performance

5. Always handle errors from `Get` properly:

   * Distinguish between cache miss and system error

6. Use capacity limits carefully in LRU to prevent memory overuse.

---

## 📂 Project Structure

```
/cache        → Interface definition
/inmemory     → LRU implementation
/redis        → Redis integration
/memcached    → Memcached integration
/main.go      → Example usage
```

---

## 🧪 Running Tests

```
go test ./...
```

---

## 📈 Running Benchmarks

```
go test ./... -bench="." -benchmem
```

---

## 🎯 Design Highlights

* Doubly linked list + hash map for O(1) LRU eviction
* Pluggable backend architecture
* Unified API abstraction
* TTL-based expiration policy
* Benchmark-driven performance validation

---

## 📌 Current Stable Release

**v1.2.1**

---

