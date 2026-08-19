package upstream

import "sync"

// Pool 上游注册表。
type Pool struct {
	mu   sync.RWMutex
	byID map[string]Target
}

// NewPool 构造。
func NewPool() *Pool {
	return &Pool{byID: make(map[string]Target)}
}

// Register 注册/更新。
func (p *Pool) Register(t Target) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.byID == nil {
		p.byID = make(map[string]Target)
	}
	if t.Weight <= 0 {
		t.Weight = 1
	}
	t.Healthy = true
	p.byID[t.ID] = t
}

// Unregister 删除。
func (p *Pool) Unregister(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.byID, id)
}

// Lookup 查询。
func (p *Pool) Lookup(id string) (Target, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	t, ok := p.byID[id]
	return t, ok
}

// List 列表。
func (p *Pool) List() []Target {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]Target, 0, len(p.byID))
	for _, t := range p.byID {
		out = append(out, t)
	}
	return out
}

// ResolveURL 优先用池中 BaseURL，否则回落传入。
func (p *Pool) ResolveURL(routeID, fallback string) string {
	if p == nil {
		return fallback
	}
	if t, ok := p.Lookup(routeID); ok && t.BaseURL != "" {
		return t.BaseURL
	}
	return fallback
}
