package cache

import (
	"container/list"
	"sync"
	"unsafe"
)

// LRUCache is a simple LRU cache inspired from https://github.com/golang/groupcache/blob/master/lru/lru.go
type LRUCache struct {
	entries     map[string]*list.Element
	entriesList *list.List
	maxEntries  int
	sync.Mutex
}

type entry struct {
	key   string
	value interface{}
}

// NewLRUCache creates a new LRU cache
func NewLRUCache(maxEntries int) *LRUCache {
	return &LRUCache{
		entries:     make(map[string]*list.Element),
		entriesList: list.New(),
		maxEntries:  maxEntries,
	}
}

// Put adds an entry to the cacche
func (c *LRUCache) Put(key string, value interface{}) {
	c.Lock()
	if _, ok := c.entries[key]; ok {
		c.Unlock()
		return
	}
	e := c.entriesList.PushFront(&entry{key, value})
	c.entries[key] = e
	if c.entriesList.Len() > c.maxEntries {
		c.removeOldest()
	}
	c.Unlock()
}

// Get finds an entry based on the given string key
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.Lock()
	val, found := c.get(key)
	c.Unlock()
	return val, found
}

func (c *LRUCache) get(key string) (interface{}, bool) {
	if e, hit := c.entries[key]; hit {
		c.entriesList.MoveToFront(e)
		return e.Value.(*entry).value, true
	}
	return nil, false
}

// GetBytes finds an entry based on the given []byte key. The given slice
// must not be modified during this call
func (c *LRUCache) GetBytes(key []byte) (interface{}, bool) {
	c.Lock()
	val, found := c.get(*(*string)(unsafe.Pointer(&key)))
	c.Unlock()
	return val, found
}

func (c *LRUCache) removeOldest() {
	e := c.entriesList.Back()
	if e != nil {
		c.entriesList.Remove(e)
		kv := e.Value.(*entry)
		delete(c.entries, kv.key)
	}
}
