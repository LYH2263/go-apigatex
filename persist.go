package apigatex

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/LYH2263/go-apigatex/internal/errors"
	"github.com/LYH2263/go-apigatex/internal/route"
)

type persistFile struct {
	Routes []RouteSpec `json:"routes"`
}

// Flush 将路由表写入 persistPath；失败返回 ErrPersist 包装。
func (g *Gateway) Flush() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.flushLocked()
}

func (g *Gateway) flushLocked() error {
	if g.persistPath == "" {
		g.dirty = false
		return nil
	}
	rows := g.table.List()
	specs := make([]RouteSpec, 0, len(rows))
	for _, r := range rows {
		specs = append(specs, fromInternalRoute(r))
	}
	raw, err := json.MarshalIndent(persistFile{Routes: specs}, "", "  ")
	if err != nil {
		return errors.WrapErr(ErrPersist, err)
	}
	dir := filepath.Dir(g.persistPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return errors.WrapErr(ErrPersist, err)
	}
	tmp := g.persistPath + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return errors.WrapErr(ErrPersist, err)
	}
	if err := os.Rename(tmp, g.persistPath); err != nil {
		_ = os.Remove(tmp)
		return errors.WrapErr(ErrPersist, err)
	}
	g.dirty = false
	return nil
}

// Sync 同 Flush（Close 前落盘入口）。
func (g *Gateway) Sync() error {
	return g.Flush()
}

// LoadPersist 从磁盘加载路由。
func (g *Gateway) LoadPersist() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		return ErrClosed
	}
	if g.persistPath == "" {
		return nil
	}
	raw, err := os.ReadFile(g.persistPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.WrapErr(ErrPersist, err)
	}
	var pf persistFile
	if err := json.Unmarshal(raw, &pf); err != nil {
		return errors.WrapErr(ErrPersist, err)
	}
	g.table.Replace(nil)
	for _, spec := range pf.Routes {
		if err := validateSpec(spec); err != nil {
			continue
		}
		_ = g.table.Upsert(toInternalRoute(spec, g.clk.Now()))
	}
	g.dirty = false
	return nil
}

// ReloadRoutes 热更新：先尝试持久化候选表，成功才替换内存路由。
// persistFn 返回错误时不得生效新路由（原子性）。
func (g *Gateway) ReloadRoutes(specs []RouteSpec, persistFn func([]RouteSpec) error) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		return ErrClosed
	}
	candidates := make([]route.Route, 0, len(specs))
	views := make([]RouteSpec, 0, len(specs))
	for _, spec := range specs {
		if err := validateSpec(spec); err != nil {
			return errors.WrapErr(ErrReload, err)
		}
		candidates = append(candidates, toInternalRoute(spec, g.clk.Now()))
		views = append(views, spec)
	}
	if persistFn != nil {
		if err := persistFn(views); err != nil {
			return errors.WrapErr(ErrPersist, err)
		}
	} else if g.persistPath != "" {
		// 无自定义 persist：先写盘再替换
		old := g.table.List()
		g.table.Replace(candidates)
		if err := g.flushLocked(); err != nil {
			g.table.Replace(old)
			return err
		}
		return nil
	}
	g.table.Replace(candidates)
	g.dirty = persistFn == nil
	return nil
}
