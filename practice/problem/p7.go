package problem

import "sync"

type Counter interface {
	Inc()
	Get() int64
	Reset()
}
type MCounter struct {
	mu sync.Mutex
	n  int64
}
type RWCounter struct {
	mu sync.RWMutex
	n  int64
}

func (c *MCounter) Inc() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}
func (c *MCounter) Get() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}
func (c *MCounter) Reset() {
	c.mu.Lock()
	c.n = 0
	c.mu.Unlock()
}
func (c *RWCounter) Inc() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}
func (c *RWCounter) Get() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.n
}
func (c *RWCounter) Reset() {
	c.mu.Lock()
	c.n = 0
	c.mu.Unlock()
}
