package apigatex

import (
	"strings"
	"time"

	"github.com/LYH2263/go-apigatex/internal/route"
	"github.com/LYH2263/go-apigatex/internal/upstream"
)

// AddRoute 注册一条路由；Methods/Headers 入表前拷贝。
func (g *Gateway) AddRoute(spec RouteSpec) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		return ErrClosed
	}
	if err := validateSpec(spec); err != nil {
		return err
	}
	r := toInternalRoute(spec, g.clk.Now())
	if err := g.table.Upsert(r); err != nil {
		return err
	}
	if spec.Upstream != "" {
		g.pool.Register(upstream.Target{ID: spec.ID, BaseURL: spec.Upstream})
	}
	g.dirty = true
	return nil
}

// RemoveRoute 按 ID 删除路由。
func (g *Gateway) RemoveRoute(id string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		return ErrClosed
	}
	if !g.table.Remove(id) {
		return ErrNotFound
	}
	g.pool.Unregister(id)
	g.dirty = true
	return nil
}

// ListRoutes 返回路由快照（含 Headers 拷贝）。
func (g *Gateway) ListRoutes() []RouteSpec {
	g.mu.Lock()
	defer g.mu.Unlock()
	rows := g.table.List()
	out := make([]RouteSpec, 0, len(rows))
	for _, r := range rows {
		out = append(out, fromInternalRoute(r))
	}
	return out
}

// GetRoute 按 ID 查询；Headers 为拷贝。
func (g *Gateway) GetRoute(id string) (RouteSpec, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	r, ok := g.table.Get(id)
	if !ok {
		return RouteSpec{}, false
	}
	return fromInternalRoute(r), true
}

// Match 按 method+path 匹配最高优先级路由。
func (g *Gateway) Match(method, path string) (RouteSpec, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	r, ok := g.table.Match(method, path)
	if !ok {
		return RouteSpec{}, false
	}
	return fromInternalRoute(r), true
}

func validateSpec(spec RouteSpec) error {
	if strings.TrimSpace(spec.ID) == "" || strings.TrimSpace(spec.Path) == "" {
		return ErrBadRoute
	}
	if spec.Kind != MatchExact && spec.Kind != MatchPrefix && spec.Kind != "" {
		return ErrBadRoute
	}
	if strings.TrimSpace(spec.Upstream) == "" {
		return ErrBadRoute
	}
	return nil
}

func toInternalRoute(spec RouteSpec, now time.Time) route.Route {
	kind := route.Kind(spec.Kind)
	if kind == "" {
		kind = route.KindExact
	}
	return route.Route{
		ID:        spec.ID,
		Path:      spec.Path,
		Kind:      kind,
		Methods:   route.CloneStrings(spec.Methods),
		Upstream:  spec.Upstream,
		Headers:   route.CloneHeaders(spec.Headers),
		StripPath: spec.StripPath,
		Auth:      spec.Auth,
		Priority:  spec.Priority,
		UpdatedAt: now,
	}
}

func fromInternalRoute(r route.Route) RouteSpec {
	return RouteSpec{
		ID:        r.ID,
		Path:      r.Path,
		Kind:      MatchKind(r.Kind),
		Methods:   route.CloneStrings(r.Methods),
		Upstream:  r.Upstream,
		Headers:   route.CloneHeaders(r.Headers),
		StripPath: r.StripPath,
		Auth:      r.Auth,
		Priority:  r.Priority,
		UpdatedAt: r.UpdatedAt,
	}
}
