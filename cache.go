package lunar

import "sync"

// Cache is the cache interface
type Cache interface {
	GetItems(namespace string) Items
	SetItems(namespace string, items Items)
}

// MemoryCache is cache stored in memory, it's the default cache for use
type MemoryCache struct {
	items sync.Map // key: namespace, value: items
}

// make sure MemoryCache implements Cache
var _ Cache = new(MemoryCache)

// GetItems gets items from cache
func (c *MemoryCache) GetItems(namespace string) Items {
	if v, ok := c.items.Load(namespace); ok {
		return v.(Items)
	}

	return Items{}
}

// SetItems sets items into cache
func (c *MemoryCache) SetItems(namespace string, items Items) {
	c.items.Store(namespace, items)
}
