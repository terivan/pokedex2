package PokeCache

import(
	"fmt"
	"time"
	"sync"
)

type Cache struct{
	entries map[string]cacheEntry
	interval time.Duration
	mutex sync.Mutex

}

type cacheEntry struct{
	createdAt time.Time 
	val []byte
}

func NewCache(interval time.Duration) (*Cache){
	fmt.Println("New cache created!")
	var mu sync.Mutex
	cache := &Cache{
		entries: make(map[string]cacheEntry),
		mutex: mu,
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	createdAt := time.Now() 
	var entry cacheEntry
	entry.createdAt = createdAt
	entry.val = val
	c.entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	if c.entries[key].val != nil{
		return c.entries[key].val, true
	} else {
		return []byte{}, false
	}
}

func (c *Cache) reapLoop(interval time.Duration) error{
	ticker := time.NewTicker(interval)
	for {
		<- ticker.C 
		
		c.mutex.Lock()
		now := time.Now()

		for key, entry := range(c.entries){
			if now.Sub(entry.createdAt) > interval{
				delete(c.entries, key)
			}
		}

		c.mutex.Unlock()

	}
}