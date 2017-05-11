package objcache

import (
	"container/list"
	"sync"
	"time"
)

type pair struct {
	Object interface{}
	expire int64
	key    string
}

type ObjCache struct {
	mu        sync.RWMutex
	items     map[string]*list.Element
	list      *list.List
	itemCount int
	config    Config
}

func (c *ObjCache) removeExpired() {
	e := time.Now().UnixNano()
	for {
		elem := c.list.Front()
		if elem == nil {
			break
		}
		v := elem.Value.(pair)
		if v.expire < e {
			c.itemCount = c.itemCount - 1
			delete(c.items, v.key)
			c.list.Remove(elem)
		} else {
			break
		}
	}
}

func (c *ObjCache) removeOldest() {
	c.itemCount = c.itemCount - 1
	elem := c.list.Front()
	v := elem.Value.(pair)
	delete(c.items, v.key)
	c.list.Remove(elem)
}

func (c *ObjCache) Set(k string, x interface{}, d time.Duration) error {
	if d == 0 {
		d = c.config.Expiration
	}
	c.mu.Lock()

	if _, ok := c.items[k]; !ok {

		c.removeExpired()

		if c.itemCount >= c.config.MaxEntryLimit {
			c.removeOldest()
		}

		p := pair{
			Object: x,
			key:    k,
			expire: time.Now().Add(d).UnixNano(),
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
		//TODO: delete if expired
		c.mu.RUnlock()
		return nil, false
	}
	c.mu.RUnlock()
	return elem.Value.(pair).Object, true
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
