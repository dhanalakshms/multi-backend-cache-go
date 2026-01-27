# In-Memory LRU Cache (Go)

Implementation of **in-memory cache with a Least Recently Used (LRU) eviction policy**.

## Objective Covered
1. Develop an in-memory cache with LRU eviction policy.

## What this code does
This code implements an in-memory LRU cache using:
- A **hash map** for fast key lookup
- A **doubly linked list** to track usage order

It ensures:
- O(1) time complexity for `Get` and `Set`
- Automatic eviction of the least recently used item when capacity is reached

---

## How it works
- Each cache entry is stored as a node in a doubly linked list
- The most recently used item is placed near the head
- The least recently used item is placed near the tail
- On every `Get` or `Set`, the accessed node is moved to the front
- When the cache exceeds its capacity, the node at the tail is removed

---

## Supported Operations

### Create Cache

```go
cache := inmemory.NewLRUCache(2)
```

### Set Value

```go
cache.Set("a", 1)
```
### Get Value

```go
value, ok := cache.Get("a")
```
---

## Eviction Policy

* The cache size is configurable during initialization
* When the cache reaches its maximum size, 
The least recently used entry is evicted automatically



