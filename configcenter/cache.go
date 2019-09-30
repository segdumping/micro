package configcenter

import (
	"sync"
	"sync/atomic"
)

var cacheClient *cache

type cache struct {
	cache        map[string]*value //kv cache
	invalidCount uint32            //invalid cache key count
	*config
	*dispatch
	sync.RWMutex
}

type value struct {
	value   []byte
	version int64
	isValid bool
}

func init() {
	cacheClient = &cache{
		config:   NewConfigCenter(),
		cache:    make(map[string]*value),
		dispatch: newDispatch(),
	}

	go cacheClient.config.watch() //etcd watch
	go cacheClient.watch()        //configcenter watch
}

func (c *cache) watch() {
	for v := range c.eventChan {
		c.Lock()
		e, ok := c.cache[v.key]
		if !ok {
			e = &value{}
		}

		if e.version == v.version {
			c.Unlock()
			continue
		}

		switch v.action {
		case "create":
			c.cache[v.key] = &value{
				value:   v.value,
				version: v.version,
				isValid: true,
			}
			c.fire(v)
		case "modify":
			e.version = v.version
			e.value = v.value
			c.fire(v)
		case "delete":
			e.isValid = false
			//move when invalid count greater than 30% of total len
			count := atomic.LoadUint32(&c.invalidCount)
			if count*3 >= uint32(len(c.cache)) {
				c.cleanAndMove()
				atomic.StoreUint32(&c.invalidCount, 0)
			} else {
				atomic.AddUint32(&c.invalidCount, 1)
			}
		}
		c.Unlock()
	}
}

//just move
func (c *cache) cleanAndMove() {
	m := make(map[string]*value, len(c.cache)/3)
	for k := range c.cache {
		val := c.cache[k]
		if val.isValid {
			m[k] = val
		}
	}

	c.cache = m
}

// get from etcd and store in cache
func (c *cache) get(key string) (map[string][]byte, error) {
	c.RLock()
	v, ok := c.cache[key]
	c.RUnlock()
	if ok {
		return map[string][]byte{key: v.value}, nil
	}

	vals, err := c.config.Get(key)
	if err != nil {
		return nil, err
	}

	//just store
	c.Lock()
	for k, val := range vals {
		c.cache[k] = &value{
			value:   val,
			isValid: true,
		}
	}
	c.Unlock()

	return vals, nil
}

func Get(key string) (map[string][]byte, error) {
	return cacheClient.get(key)
}

func Listen(this Watch, key string) {
	cacheClient.AddListen(key, this)
}
