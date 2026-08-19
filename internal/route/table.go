package route

import (
	"sort"
	"strings"
	"sync"
)

// Table 路由表。
type Table struct {
	mu   sync.RWMutex
	byID map[string]Route
}

// NewTable 构造空表。
func NewTable() *Table {
	return &Table{byID: make(map[string]Route)}
}

// Upsert 插入或更新（存拷贝）。
func (t *Table) Upsert(r Route) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.byID == nil {
		t.byID = make(map[string]Route)
	}
	t.byID[r.ID] = CloneRoute(r)
	return nil
}

// Remove 删除。
func (t *Table) Remove(id string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.byID[id]; !ok {
		return false
	}
	delete(t.byID, id)
	return true
}

// Get 返回拷贝。
func (t *Table) Get(id string) (Route, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	r, ok := t.byID[id]
	if !ok {
		return Route{}, false
	}
	return CloneRoute(r), true
}

// List 返回全部拷贝，按 priority 降序。
func (t *Table) List() []Route {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Route, 0, len(t.byID))
	for _, r := range t.byID {
		out = append(out, CloneRoute(r))
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Priority == out[j].Priority {
			return out[i].ID < out[j].ID
		}
		return out[i].Priority > out[j].Priority
	})
	return out
}

// Replace 整体替换。
func (t *Table) Replace(rows []Route) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.byID = make(map[string]Route, len(rows))
	for _, r := range rows {
		t.byID[r.ID] = CloneRoute(r)
	}
}

// Len 条数。
func (t *Table) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.byID)
}

// Match 匹配 method+path。
func (t *Table) Match(method, path string) (Route, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	method = strings.ToUpper(strings.TrimSpace(method))
	var best Route
	found := false
	for _, r := range t.byID {
		if !methodAllowed(r.Methods, method) {
			continue
		}
		if !pathMatch(r, path) {
			continue
		}
		if !found || r.Priority > best.Priority || (r.Priority == best.Priority && r.ID < best.ID) {
			best = r
			found = true
		}
	}
	if !found {
		return Route{}, false
	}
	return CloneRoute(best), true
}

func methodAllowed(methods []string, method string) bool {
	if len(methods) == 0 {
		return true
	}
	for _, m := range methods {
		if strings.EqualFold(m, method) || m == "*" {
			return true
		}
	}
	return false
}

func pathMatch(r Route, path string) bool {
	switch r.Kind {
	case KindPrefix:
		return strings.HasPrefix(path, r.Path)
	default:
		return path == r.Path
	}
}
