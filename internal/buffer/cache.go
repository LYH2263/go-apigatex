package buffer

import "sync"

// Cache 按 key 缓存请求体拷贝。
type Cache struct {
	mu    sync.Mutex
	limit int
	data  map[string][]byte
}

// NewCache 构造。
func NewCache(limit int) *Cache {
	if limit < 1 {
		limit = 1 << 20
	}
	return &Cache{limit: limit, data: make(map[string][]byte)}
}

// SetLimit 更新上限。
func (c *Cache) SetLimit(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if n > 0 {
		c.limit = n
	}
}

// Limit 返回上限。
func (c *Cache) Limit() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.limit
}

// Put 存入 body 的独立拷贝并返回该拷贝。
func (c *Cache) Put(key string, body []byte) []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	cp := CloneBytes(body)
	if c.data == nil {
		c.data = make(map[string][]byte)
	}
	c.data[key] = cp
	return CloneBytes(cp)
}

// Get 返回缓存体的拷贝。
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.data[key]
	if !ok {
		return nil, false
	}
	return CloneBytes(v), true
}

// Delete 删除。
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Clear 清空。
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string][]byte)
}
