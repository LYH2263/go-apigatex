package apigatex

import (
	"github.com/LYH2263/go-apigatex/internal/errors"
)

// Close 先 Sync 落盘，再释放路由表与 transport；可重复调用。
func (g *Gateway) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		return nil
	}
	var first error
	// 先 Sync 落盘，再释放路由表与 transport；否则表已清空，写盘为空数组。
	if g.dirty || g.persistPath != "" {
		if err := g.flushLocked(); err != nil {
			first = err
		}
	}
	if g.table != nil {
		g.table.Replace(nil)
	}
	g.client = nil
	g.authn = nil
	g.cache = nil
	g.pool = nil
	g.closed = true
	if first != nil {
		return errors.WrapErr(ErrSync, first)
	}
	return nil
}

// Closed 是否已关闭。
func (g *Gateway) Closed() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.closed
}
