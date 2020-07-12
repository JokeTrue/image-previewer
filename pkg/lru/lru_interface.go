// Package lru provides simple LRU implementation based on build-in container/list.
package lru

// Cache is the interface for LRU cache.
type Cache interface {
	// Adds a value to the cache and updates the "recently used"-ness of the key.
	Add(key, value interface{}) bool

	// Returns key's value from the cache and
	// updates the "recently used"-ness of the key.
	Get(key interface{}) (value interface{}, ok bool)

	// Checks if a key exists in cache without updating the recent-ness.
	Contains(key interface{}) (ok bool)

	// Returns key's value without updating the "recently used"-ness of the key.
	Peek(key interface{}) (value interface{}, ok bool)

	// Removes a key from the cache.
	Remove(key interface{}) bool

	// Removes the oldest entry from cache.
	RemoveOldest() (interface{}, interface{}, bool)

	// Returns the oldest entry from the cache.
	GetOldest() (interface{}, interface{}, bool)

	// Returns a slice of the keys in the cache, from oldest to newest.
	Keys() []interface{}

	// Returns the number of items in the cache.
	Len() int

	// Clears all cache entries.
	Purge()
}
