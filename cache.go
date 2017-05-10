package objcache

import (
	"container/heap"
	"sync"
	"time"
)

type KeyIndex struct {
	key        string
	Expiration int64
	index      int
}

type PriorityQueue []*KeyIndex

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Expiration < pq[j].Expiration
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[i].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*KeyIndex)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *KeyIndex, value string, priority int64) {
	item.key = value
	item.Expiration = priority
	heap.Fix(pq, item.index)
}

type Item struct {
	Object interface{}
	prior  *KeyIndex
}

func (item Item) expired() bool {
	return item.prior.Expiration != 0 && time.Now().UnixNano() > item.prior.Expiration
}

type ObjCache struct {
	items     map[string]Item
	mu        sync.RWMutex
	heap      *PriorityQueue
	itemCount int
	config    Config
}

func (c *ObjCache) removeOldest() {
	target := heap.Pop(c.heap).(*KeyIndex)
	delete(c.items, target.key)
}

func (c *ObjCache) Set(k string, x interface{}, d time.Duration) error {
	e := time.Now().Add(d).UnixNano()
	c.mu.Lock()

	if _, ok := c.items[k]; !ok {
		if c.itemCount+1 > c.config.MaxEntryLimit {
			c.itemCount = c.itemCount - 1
			c.removeOldest()
		}

		ki := &KeyIndex{
			key:        k,
			Expiration: e,
		}

		heap.Push(c.heap, ki)

		c.items[k] = Item{
			Object: x,
			prior:  ki,
		}

		c.itemCount = c.itemCount + 1
	} else {
		c.heap.update(c.items[k].prior, k, e)
		c.items[k] = Item{
			Object: x,
			prior:  c.items[k].prior,
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

func New(config Config) (*ObjCache, error) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	cache := &ObjCache{
		items:     make(map[string]Item),
		itemCount: 0,
		heap:      &pq,
		config:    config,
	}
	return cache, nil
}
