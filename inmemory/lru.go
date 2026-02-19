package inmemory

import (
	"fmt"
	"sync"
	"time"
)

// Node represents one entry in the cache 
type Node struct {
	key    string
	value  interface{}
	prev   *Node
	next   *Node
	expiry time.Time // expiry time for TTL
}

// LRUCache stores cache data with LRU eviction policy
type LRUCache struct {
	capacity        int
	cache           map[string]*Node 
	head            *Node            
	tail            *Node            
	mu              sync.Mutex       
	cleanupInterval time.Duration    
	stopCleanup     chan struct{}   
}

// NewLRUCache creates a new cache and optionally starts background cleanup
func NewLRUCache(capacity int, cleanupInterval ...time.Duration) *LRUCache {
	head := &Node{}
	tail := &Node{}
	head.next = tail
	tail.prev = head

	c := &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*Node),
		head:     head,
		tail:     tail,
	}

	// start background cleanup if interval provided
	if len(cleanupInterval) > 0 && cleanupInterval[0] > 0 {
		c.cleanupInterval = cleanupInterval[0]
		c.stopCleanup = make(chan struct{})
		go c.startCleanup()
	}

	return c
}

// add inserts node right after head (mark as most recently used)
func (c *LRUCache) add(node *Node) {
	next := c.head.next
	c.head.next = node
	node.prev = c.head
	node.next = next
	next.prev = node
}

// remove disconnects a node from linked list
func (c *LRUCache) remove(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// Get returns value for a key and marks it as recently used
func (c *LRUCache) Get(key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.cache[key]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	// check TTL expiration
	if !node.expiry.IsZero() && time.Now().After(node.expiry) {
		c.remove(node)
		delete(c.cache, key)
		return nil, fmt.Errorf("key expired")
	}

	// move node to front 
	c.remove(node)
	c.add(node)

	return node.value, nil
}

// Set inserts or updates a key with optional TTL
func (c *LRUCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// if key exists, remove old node
	if node, ok := c.cache[key]; ok {
		c.remove(node)
		delete(c.cache, key)
	}

	// evict least recently used if capacity reached
	if c.capacity > 0 && len(c.cache) >= c.capacity {
		lru := c.tail.prev
		if lru != c.head {
			c.remove(lru)
			delete(c.cache, lru.key)
		}
	}

	// calculate expiry time
	expiry := time.Time{}
	if ttl > 0 {
		expiry = time.Now().Add(ttl)
	}

	node := &Node{
		key:    key,
		value:  value,
		expiry: expiry,
	}

	// add new node as most recent
	c.add(node)
	c.cache[key] = node

	return nil
}

// Delete removes a key manually
func (c *LRUCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.cache[key]
	if !ok {
		return fmt.Errorf("key not found")
	}

	c.remove(node)
	delete(c.cache, key)
	return nil
}

// Clear removes all entries and resets the list
func (c *LRUCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*Node)
	c.head.next = c.tail
	c.tail.prev = c.head
	return nil
}

// startCleanup runs background worker to remove expired entries
func (c *LRUCache) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.removeExpired()
		case <-c.stopCleanup:
			return
		}
	}
}

// removeExpired scans cache and removes expired keys
func (c *LRUCache) removeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, node := range c.cache {
		if !node.expiry.IsZero() && now.After(node.expiry) {
			c.remove(node)
			delete(c.cache, key)
		}
	}
}

// StopCleanup stops the background cleanup goroutine
func (c *LRUCache) StopCleanup() {
	if c.stopCleanup != nil {
		close(c.stopCleanup)
	}
}
