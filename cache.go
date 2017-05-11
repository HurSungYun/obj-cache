package objcache

import (
	"container/list"
	"sync"
	"time"
)

type Item struct {
	Object   interface{}
	listElem *list.Element
}

type ObjCache struct {
	items     map[string]Item
	mu        sync.RWMutex
	list      *list.List
	itemCount int
	config    Config
}

func (c *ObjCache) removeOldest() {
	elem := c.list.Front()
	key := c.list.Remove(elem).(string)
	delete(c.items, key)
}

func (c *ObjCache) Set(k string, x interface{}, d time.Duration) error {
	c.mu.Lock()

	if _, ok := c.items[k]; !ok {
		if c.itemCount+1 > c.config.MaxEntryLimit {
			c.itemCount = c.itemCount - 1
			c.removeOldest()
		}

		elem := c.list.PushBack(k)

		c.items[k] = Item{
			Object:   x,
			listElem: elem,
		}

		c.itemCount = c.itemCount + 1
	} else {
		c.list.MoveToBack(c.items[k].listElem)

		c.items[k] = Item{
			Object:   x,
			listElem: c.items[k].listElem,
		}
	}

	c.mu.Unlock()
	return nil
}

func (c *ObjCache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	item, ok := c.items[k]

	if !ok {
		c.mu.RUnlock()
		return nil, false
	}
	c.mu.RUnlock()
	return item.Object, true
}

func (c *ObjCache) Del(k string) bool {
	if item, ok := c.items[k]; ok {
		elem := item.listElem
		c.list.Remove(elem)
		delete(c.items, k)
		return true
	}
	return false
}

func New(config Config) (*ObjCache, error) {
	l := list.New()
	cache := &ObjCache{
		items:     make(map[string]Item),
		itemCount: 0,
		list:      l,
		config:    config,
	}
	return cache, nil
}
