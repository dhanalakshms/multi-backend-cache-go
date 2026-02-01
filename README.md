
# Multi Backend Cache Library 

Implementation of a lightweight in-memory cache with a Least Recently Used (LRU) eviction policy.

---

## Objective Covered

1. Develop an in-memory cache with LRU eviction policy
2. Implement Set, Get, Delete operations
3. Add TTL and Mutex support
4. Handle automatic eviction when capacity is exceeded

---

## What this code does

This project implements an **in-memory LRU cache** using:

* A **hash map** for fast key lookup
* A **doubly linked list** to track usage order
* **TTL timers** for automatic expiration
* **Mutex locking** for safe concurrent access

---

## Project Structure

```
multi-backend-cache-go
│
├── inmemory/
│   └── lru.go        → LRU cache implementation
│
├── main.go           → Demo & testing all operations
├── go.mod
└── README.md
```

---

## Supported Operations

### Create Cache

```go
cache := inmemory.NewLRUCache(2)
```

---

### Set Value (with TTL)

```go
cache.Set("a", 1, 5*time.Second)
```

---

### Get Value

```go
value, ok := cache.Get("a")
```

---

### Delete Key

```go
cache.Delete("a")
```

---

## Example Usage (main.go)

```go
cache := inmemory.NewLRUCache(2)

cache.Set("d", "dhanalakshmi", 5*time.Second)

val, ok := cache.Get("d")
fmt.Println(val, ok)

cache.Delete("d")

_, ok = cache.Get("d")
fmt.Println(ok)
```

---

## Sample Output

```
dhanalakshmi true
after delete: false
after ttl: false
```


## Current Progress

### Week 1

* Implemented doubly linked list
* Designed node structure
* Built base LRU logic
* Added Set and Get operations

### Week 2

* Implemented Delete
* Added TTL support
* Improved eviction handling
* Tested all operations in main.go

---

## Next Improvements - Planned

* Redis backend
* Memcached backend
* Common cache interface
* Unit tests

---

## How to Run

```bash
git clone https://github.com/dhanalakshms/multi-backend-cache-go.git
cd multi-backend-cache-go
go run main.go
```



