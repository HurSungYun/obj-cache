package objcache

import (
	"container/list"
	"sync"
	"time"
)

type Pair struct {
	Object interface{}
	key    string
}

type ObjCache struct {
	mu        sync.RWMutex
	items     map[string]*list.Element
	list      *list.List
	itemCount int
	config    Config
}

func (c *ObjCache) removeOldest() {
	c.itemCount = c.itemCount - 1
	elem := c.list.Front()
	v := elem.Value.(Pair)
	delete(c.items, v.key)
	c.list.Remove(elem)
}

func (c *ObjCache) Set(k string, x interface{}, d time.Duration) error {
	c.mu.Lock()

	if _, ok := c.items[k]; !ok {
		if c.itemCount >= c.config.MaxEntryLimit {
			c.removeOldest()
		}

		p := Pair{
			Object: x,
			key:    k,
		}
		c.items[k] = c.list.PushBack(p)
		c.itemCount = c.itemCount + 1
	} else {
		c.list.MoveToBack(c.items[k])
	}

	c.mu.Unlock()
	return nil
}

func (c *ObjCache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	elem, ok := c.items[k]
	if !ok {
		c.mu.RUnlock()
		return nil, false
	}
	c.mu.RUnlock()
	return elem.Value.(Pair).Object, true
}

func (c *ObjCache) Del(k string) bool {
	c.mu.Lock()
	item, ok := c.items[k]
	if ok {
		c.itemCount = c.itemCount - 1
		delete(c.items, k)
		c.list.Remove(item)
	}
	c.mu.Unlock()
	return ok
}

func New(config Config) (*ObjCache, error) {
	l := list.New()
	cache := &ObjCache{
		items:     make(map[string]*list.Element),
		itemCount: 0,
		list:      l,
		config:    config,
	}
	return cache, nil
}
