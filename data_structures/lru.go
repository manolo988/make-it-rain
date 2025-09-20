package data_structures

import (
	"container/list"
	"sync"
	"time"
)

// LRUCache is a Least Recently Used cache implementation
type LRUCache[K comparable, V any] struct {
	capacity int
	cache    map[K]*list.Element
	list     *list.List
	mu       sync.RWMutex
}

type lruEntry[K comparable, V any] struct {
	key   K
	value V
}

// NewLRUCache creates a new LRU cache with given capacity
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		cache:    make(map[K]*list.Element),
		list:     list.New(),
	}
}

// Get retrieves a value from cache
func (lru *LRUCache[K, V]) Get(key K) (V, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if elem, exists := lru.cache[key]; exists {
		// Move to front (most recently used)
		lru.list.MoveToFront(elem)
		entry := elem.Value.(lruEntry[K, V])
		return entry.value, true
	}

	var zero V
	return zero, false
}

// Put adds or updates a key-value pair
func (lru *LRUCache[K, V]) Put(key K, value V) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	// Update existing entry
	if elem, exists := lru.cache[key]; exists {
		lru.list.MoveToFront(elem)
		elem.Value = lruEntry[K, V]{key: key, value: value}
		return
	}

	// Add new entry
	entry := lruEntry[K, V]{key: key, value: value}
	elem := lru.list.PushFront(entry)
	lru.cache[key] = elem

	// Remove least recently used if over capacity
	if lru.list.Len() > lru.capacity {
		oldest := lru.list.Back()
		if oldest != nil {
			lru.list.Remove(oldest)
			oldEntry := oldest.Value.(lruEntry[K, V])
			delete(lru.cache, oldEntry.key)
		}
	}
}

// Delete removes a key from cache
func (lru *LRUCache[K, V]) Delete(key K) bool {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if elem, exists := lru.cache[key]; exists {
		lru.list.Remove(elem)
		delete(lru.cache, key)
		return true
	}
	return false
}

// Contains checks if key exists
func (lru *LRUCache[K, V]) Contains(key K) bool {
	lru.mu.RLock()
	defer lru.mu.RUnlock()
	_, exists := lru.cache[key]
	return exists
}

// Size returns current number of items
func (lru *LRUCache[K, V]) Size() int {
	lru.mu.RLock()
	defer lru.mu.RUnlock()
	return len(lru.cache)
}

// Clear removes all items
func (lru *LRUCache[K, V]) Clear() {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	lru.cache = make(map[K]*list.Element)
	lru.list.Init()
}

// Keys returns all keys in MRU to LRU order
func (lru *LRUCache[K, V]) Keys() []K {
	lru.mu.RLock()
	defer lru.mu.RUnlock()

	keys := make([]K, 0, len(lru.cache))
	for elem := lru.list.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(lruEntry[K, V])
		keys = append(keys, entry.key)
	}
	return keys
}

// LRUCacheWithTTL is an LRU cache with time-to-live support
type LRUCacheWithTTL[K comparable, V any] struct {
	capacity int
	ttl      time.Duration
	cache    map[K]*list.Element
	list     *list.List
	mu       sync.RWMutex
}

type lruEntryWithTTL[K comparable, V any] struct {
	key        K
	value      V
	expiration time.Time
}

// NewLRUCacheWithTTL creates an LRU cache with TTL
func NewLRUCacheWithTTL[K comparable, V any](capacity int, ttl time.Duration) *LRUCacheWithTTL[K, V] {
	cache := &LRUCacheWithTTL[K, V]{
		capacity: capacity,
		ttl:      ttl,
		cache:    make(map[K]*list.Element),
		list:     list.New(),
	}

	// Start cleanup goroutine
	go cache.cleanupExpired()

	return cache
}

// Get retrieves a value if not expired
func (lru *LRUCacheWithTTL[K, V]) Get(key K) (V, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if elem, exists := lru.cache[key]; exists {
		entry := elem.Value.(lruEntryWithTTL[K, V])

		// Check if expired
		if time.Now().After(entry.expiration) {
			lru.list.Remove(elem)
			delete(lru.cache, key)
			var zero V
			return zero, false
		}

		// Move to front and update expiration
		lru.list.MoveToFront(elem)
		entry.expiration = time.Now().Add(lru.ttl)
		elem.Value = entry
		return entry.value, true
	}

	var zero V
	return zero, false
}

// Put adds or updates with TTL
func (lru *LRUCacheWithTTL[K, V]) Put(key K, value V) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	expiration := time.Now().Add(lru.ttl)

	// Update existing entry
	if elem, exists := lru.cache[key]; exists {
		lru.list.MoveToFront(elem)
		elem.Value = lruEntryWithTTL[K, V]{
			key:        key,
			value:      value,
			expiration: expiration,
		}
		return
	}

	// Add new entry
	entry := lruEntryWithTTL[K, V]{
		key:        key,
		value:      value,
		expiration: expiration,
	}
	elem := lru.list.PushFront(entry)
	lru.cache[key] = elem

	// Remove LRU if over capacity
	if lru.list.Len() > lru.capacity {
		oldest := lru.list.Back()
		if oldest != nil {
			lru.list.Remove(oldest)
			oldEntry := oldest.Value.(lruEntryWithTTL[K, V])
			delete(lru.cache, oldEntry.key)
		}
	}
}

// cleanupExpired periodically removes expired entries
func (lru *LRUCacheWithTTL[K, V]) cleanupExpired() {
	ticker := time.NewTicker(lru.ttl / 2)
	defer ticker.Stop()

	for range ticker.C {
		lru.mu.Lock()
		now := time.Now()

		var next *list.Element
		for elem := lru.list.Back(); elem != nil; elem = next {
			next = elem.Prev()
			entry := elem.Value.(lruEntryWithTTL[K, V])

			if now.After(entry.expiration) {
				lru.list.Remove(elem)
				delete(lru.cache, entry.key)
			} else {
				// Since list is ordered by usage, once we find non-expired, we can stop
				break
			}
		}
		lru.mu.Unlock()
	}
}

// Size returns current number of items
func (lru *LRUCacheWithTTL[K, V]) Size() int {
	lru.mu.RLock()
	defer lru.mu.RUnlock()
	return len(lru.cache)
}

// Clear removes all items
func (lru *LRUCacheWithTTL[K, V]) Clear() {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	lru.cache = make(map[K]*list.Element)
	lru.list.Init()
}