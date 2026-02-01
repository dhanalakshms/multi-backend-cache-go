package inmemory

import (
	"sync"
	"time"
)

// LL node creation
type Node struct {
	key   string
	value interface{}
	prev  *Node
	next  *Node
	expiry time.Time
}

// LRU cache
type LRUCache struct {
	capacity int
	cache    map[string]*Node
	head     *Node
	tail     *Node
	mu sync.RWMutex  //Mutex
}

// NewLRUCache initializes the LRU cache
func NewLRUCache(capacity int) *LRUCache {
	head := &Node{}
	tail := &Node{}

	head.next = tail
	tail.prev = head

	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*Node),
		head:     head,
		tail:     tail,
	}
}

// add inserts node right after head (most recently used)
func (c *LRUCache) add(node *Node) {
	nextNode := c.head.next
	c.head.next = node
	node.prev = c.head
	node.next = nextNode
	nextNode.prev = node
}

// remove removes a node from the doubly linked list
func (c *LRUCache) remove(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}


// Get retrieves a value and marks it as recently used
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	// TTL check
	if !node.expiry.IsZero() && time.Now().After(node.expiry) {
		c.remove(node)
		delete(c.cache, key)
		return nil, false
	}

	c.remove(node)
	c.add(node)
	return node.value, true
}


// Set inserts or updates a key-value pair
// Set inserts or updates a key-value pair with optional TTL
func (c *LRUCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If exists, remove old
	if node, ok := c.cache[key]; ok {
		c.remove(node)
		delete(c.cache, key)
	}

	// Evict if full
	if len(c.cache) >= c.capacity {
		lru := c.tail.prev
		c.remove(lru)
		delete(c.cache, lru.key)
	}

	exp := time.Time{}
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}

	node := &Node{
		key:    key,
		value:  value,
		expiry: exp,
	}

	c.add(node)
	c.cache[key] = node
}


// Delete removes a key manually
func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, ok := c.cache[key]; ok {
		c.remove(node)
		delete(c.cache, key)
	}
}

