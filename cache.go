package lunar

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

// Cache is the cache interface
type Cache interface {
	GetItems(namespace string) Items
	SetItems(namespace string, items Items) error
	GetKeys() []string
	Delete(namespace string) error
	Drain()
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
func (c *MemoryCache) SetItems(namespace string, items Items) error {
	c.items.Store(namespace, items)

	return nil
}

// GetKeys gets all the keys (namespaces)
func (c *MemoryCache) GetKeys() []string {
	var keys []string

	c.items.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))

		return true
	})

	return keys
}

// Delete deletes given namespace
func (c *MemoryCache) Delete(namespace string) error {
	c.items.Delete(namespace)

	return nil
}

// Drain deletes the whole cache
func (c *MemoryCache) Drain() {
	c.items.Range(func(key, value interface{}) bool {
		c.items.Delete(key)

		return true
	})
}

// FileCache is cache stored in files.
type FileCache struct {
	lock   sync.Mutex
	AppID  string
	Folder string // root folder
	Perm   os.FileMode
}

// make sure FileCache implements Cache
var _ Cache = new(FileCache)

// NewFileCache creates a FileCache
func NewFileCache(appID string, folder string) *FileCache {
	if folder == "" {
		folder = "/tmp"
	}

	c := &FileCache{
		AppID:  appID,
		Folder: folder,
		Perm:   0666,
	}

	if err := c.createAppFolder(); err != nil {
		return nil
	}

	return c
}

// check and create app folder if not existing
func (c *FileCache) createAppFolder() error {
	if _, err := os.Stat(c.getAppFolder()); err != nil {
		return os.MkdirAll(c.getAppFolder(), os.FileMode(0755))
	}
	return nil
}

func (c *FileCache) getAppFolder() string {
	return filepath.Join(c.Folder, c.AppID)
}

// file path is {Folder}/{app id}/{namespace}
func (c *FileCache) getFilePath(namespace string) string {
	return filepath.Join(c.getAppFolder(), namespace)
}

// GetItems gets items from cache
func (c *FileCache) GetItems(namespace string) Items {
	c.lock.Lock()
	defer c.lock.Unlock()

	items := make(Items)

	if data, err := os.ReadFile(c.getFilePath(namespace)); err == nil {
		if IsProperties(namespace) {
			if err := json.Unmarshal(data, &items); err != nil {
				log.Printf("parse data error: %s", err.Error())
			}
		} else {
			items["content"] = string(data)
		}
	}

	return items
}

// SetItems sets items into cache
func (c *FileCache) SetItems(namespace string, items Items) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	var content string
	if IsProperties(namespace) {
		content = items.String()
	} else {
		content = items.Get("content")
	}

	return os.WriteFile(c.getFilePath(namespace), []byte(content), c.Perm)
}

// GetKeys gets all the keys (namespaces)
func (c *FileCache) GetKeys() []string {
	c.lock.Lock()
	defer c.lock.Unlock()

	f, err := os.Open(filepath.Join(c.Folder, c.AppID))
	if err != nil {
		return nil
	}

	names, _ := f.Readdirnames(-1)
	f.Close()

	return names
}

// Delete deletes given namespace
func (c *FileCache) Delete(namespace string) error {
	return syscall.Unlink(c.getFilePath(namespace))
}

// Drain deletes the whole cache
func (c *FileCache) Drain() {
	for _, namespace := range c.GetKeys() {
		if err := syscall.Unlink(c.getFilePath(namespace)); err != nil {
			log.Printf("failed to delete cache file: %s", err.Error())
		}
	}
}
