package memstorage

import "sync"

// collection is thread-safe wrapper for map.
type collection[K comparable, V any] struct {
	mutex sync.RWMutex
	data  map[K]V
}

// newCollection creates a new initialized instance of collection.
func newCollection[K comparable, V any]() *collection[K, V] {
	return &collection[K, V]{
		data: make(map[K]V),
	}
}

// Get returns an element from the collection if the specified key exists.
func (c *collection[K, V]) Get(key K) (val V, ok bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	val, ok = c.data[key]

	return val, ok
}

// GetOrCreate returns an existing element from the collection or creates a new one using the create function.
func (c *collection[K, V]) GetOrCreate(key K, create func() V) (val V) {
	val, ok := c.Get(key)
	if !ok {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		val, ok = c.data[key]
		if !ok {
			val = create()
			c.data[key] = val
		}
	}

	return val
}

// Set creates or updates an element with the specified key.
func (c *collection[K, V]) Set(key K, val V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = val
}

// Del removes an element with the specified key.
func (c *collection[K, V]) Del(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

// Map takes a snapshot of the collection and returns it as a map.
func (c *collection[K, V]) Map() map[K]V {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	m := make(map[K]V, len(c.data))

	for k, v := range c.data {
		m[k] = v
	}

	return m
}
