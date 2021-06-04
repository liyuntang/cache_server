package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
	GetStat() Stat
	NewScanner() Scanner
}

type Stat struct {
	Count int64
	KeySize int64
	ValueSize int64
}

func (s *Stat)add(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
	//fmt.Println("count is", s.Count, "key size is", s.KeySize, "value size is", s.ValueSize)
}

func (s *Stat)del(k string, v []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}

func New(typ string, ttl int) Cache {
	var c Cache
	if typ == "inmemory" {
		c = newInMemoryCache(ttl)
	}

	if c == nil {
		panic("unknow cache type "+ typ)
	}
	log.Println(typ, "redy to serve")
	return c
}
type value struct {
	v []byte
	created time.Time
}
type inMemoryCache struct {
	c map[string]value
	mutex sync.RWMutex
	Stat
	ttl time.Duration
}

func (c *inMemoryCache)Set(k string, v []byte) error {
	//fmt.Println("set operation is runing.................")
	c.mutex.Lock()
	defer c.mutex.Unlock()
	//tmp, exist := c.c[k]
	//if exist {
	//	c.del(k, tmp)
	//}
	c.c[k] = value{v: v, created: time.Now()}
	fmt.Println("set is ok, data is", c.c)
	c.add(k, v)
	return nil
}

func (c *inMemoryCache)Get(k string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.c[k].v, nil
}

func (c *inMemoryCache)Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, exist := c.c[k]
	if exist {
		delete(c.c, k)
		c.del(k, v.v)
	}
	return nil
}

func (c *inMemoryCache)GetStat() Stat {
	return c.Stat
}

func newInMemoryCache(ttl int) *inMemoryCache {
	c :=  &inMemoryCache{make(map[string]value), sync.RWMutex{}, Stat{}, time.Duration(ttl) * time.Second}
	if ttl >0 {
		go c.expirer()
	}
	return c
}

func (c *inMemoryCache)expirer()  {
	for {
		time.Sleep(c.ttl)
		c.mutex.RLock()
		for k, v := range c.c {
			c.mutex.RUnlock()
			if v.created.Add(c.ttl).Before(time.Now()) {
				c.Del(k)
			}
			c.mutex.RLock()
		}
		c.mutex.RUnlock()
	}
}


















































