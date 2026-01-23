package inmemory

// Node represents a doubly linked list node
type Node struct {
	key   string
	value interface{}
	prev  *Node
	next  *Node
}

// LRUCache represents the LRU cache
type LRUCache struct {
	capacity int
	cache    map[string]*Node
	head     *Node
	tail     *Node
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
	if node, ok := c.cache[key]; ok {
		c.remove(node)
		c.add(node)
		return node.value, true
	}
	return nil, false
}

// Set inserts or updates a key-value pair
func (c *LRUCache) Set(key string, value interface{}) {
	// If key exists, remove old node
	if node, ok := c.cache[key]; ok {
		c.remove(node)
		delete(c.cache, key)
	}

	// If capacity exceeded, remove least recently used
	if len(c.cache) >= c.capacity {
		lru := c.tail.prev
		c.remove(lru)
		delete(c.cache, lru.key)
	}

	// Add new node
	node := &Node{key: key, value: value}
	c.add(node)
	c.cache[key] = node
}
