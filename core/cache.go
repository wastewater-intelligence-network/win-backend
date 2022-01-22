package core

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	cache *cache.Cache
}

func NewCache(defaultExpiration time.Duration) *Cache {
	return &Cache{
		cache: cache.New(defaultExpiration, time.Hour*1),
	}
}

func (c *Cache) Load(key interface{}) (interface{}, bool) {
	return c.cache.Get(key.(string))
}

func (c *Cache) Store(key interface{}, value interface{}) {
	c.cache.Set(key.(string), value, cache.DefaultExpiration)
}

func (c *Cache) StoreWithTTL(key interface{}, value interface{}, ttl time.Duration) {
	c.cache.Set(key.(string), value, ttl)
}

func (c *Cache) Delete(key interface{}) {
	c.cache.Delete(key.(string))
}
