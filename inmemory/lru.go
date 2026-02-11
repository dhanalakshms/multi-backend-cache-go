package inmemory

import (
	"fmt"
	"sync"
	"time"
)

type Node struct {
	key    string
	value  interface{}
	prev   *Node
	next   *Node
	expiry time.Time
}

type LRUCache struct {
	capacity        int
	cache           map[string]*Node
	head            *Node
	tail            *Node
	mu              sync.Mutex
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

// NewLRUCache creates LRU cache with optional background cleanup.
// If cleanupInterval > 0, expired entries are removed periodically.
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

	if len(cleanupInterval) > 0 && cleanupInterval[0] > 0 {
		c.cleanupInterval = cleanupInterval[0]
		c.stopCleanup = make(chan struct{})
		go c.startCleanup()
	}

	return c
}

func (c *LRUCache) add(node *Node) {
	next := c.head.next
	c.head.next = node
	node.prev = c.head
	node.next = next
	next.prev = node
}

func (c *LRUCache) remove(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (c *LRUCache) Get(key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.cache[key]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	if !node.expiry.IsZero() && time.Now().After(node.expiry) {
		c.remove(node)
		delete(c.cache, key)
		return nil, fmt.Errorf("key expired")
	}

	c.remove(node)
	c.add(node)

	return node.value, nil
}

func (c *LRUCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, ok := c.cache[key]; ok {
		c.remove(node)
		delete(c.cache, key)
	}

	if c.capacity > 0 && len(c.cache) >= c.capacity {
		lru := c.tail.prev
		if lru != c.head {
			c.remove(lru)
			delete(c.cache, lru.key)
		}
	}

	expiry := time.Time{}
	if ttl > 0 {
		expiry = time.Now().Add(ttl)
	}

	node := &Node{
		key:    key,
		value:  value,
		expiry: expiry,
	}

	c.add(node)
	c.cache[key] = node

	return nil
}

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

func (c *LRUCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*Node)
	c.head.next = c.tail
	c.tail.prev = c.head
	return nil
}

// ------------------------------
// Background Cleanup
// ------------------------------

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

// StopCleanup safely stops background cleanup.
func (c *LRUCache) StopCleanup() {
	if c.stopCleanup != nil {
		close(c.stopCleanup)
	}
}



