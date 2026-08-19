package apigatex

import (
	"io"
	"net/http"
	"strings"

	"github.com/LYH2263/go-apigatex/internal/auth"
	"github.com/LYH2263/go-apigatex/internal/buffer"
	"github.com/LYH2263/go-apigatex/internal/errors"
	"github.com/LYH2263/go-apigatex/internal/proxy"
	"github.com/LYH2263/go-apigatex/internal/rewrite"
)

// ServeHTTP 作为反向代理入口；Close 后返回 ErrClosed，不 panic。
func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := g.serve(w, r); err != nil {
		writeErr(w, err)
	}
}

func (g *Gateway) serve(w http.ResponseWriter, r *http.Request) error {
	g.mu.Lock()
	closed := g.closed
	client := g.client
	authn := g.authn
	table := g.table
	cache := g.cache
	pool := g.pool
	timeout := g.proxyTimeout
	allowHTTP := g.allowHTTP
	g.mu.Unlock()

	if closed || client == nil || table == nil {
		return ErrClosed
	}

	rt, ok := table.Match(r.Method, r.URL.Path)
	if !ok {
		return ErrNotFound
	}

	if rt.Auth {
		// BUG: 直接解引用 nil Authenticator
		if err := authn.Authenticate(r); err != nil {
			g.bumpAuth(false)
			return err
		}
		g.bumpAuth(true)
	}

	body, err := buffer.ReadLimited(r.Body, cache.Limit())
	if err != nil {
		return err
	}
	_ = r.Body.Close()
	cached := cache.Put(rt.ID, body)

	path := rewrite.ApplyPath(r.URL.Path, rt.StripPath, rt.Path)
	hdr := rewrite.MergeHeaders(r.Header, rt.Headers)

	eng := proxy.NewEngine(client, pool, timeout, allowHTTP)
	res, err := eng.Forward(r.Context(), proxy.Request{
		Method:   r.Method,
		Path:     path,
		Query:    r.URL.RawQuery,
		Headers:  hdr,
		Body:     cached,
		Upstream: rt.Upstream,
		RouteID:  rt.ID,
	})
	if err != nil {
		g.bumpUpstream(false)
		return err
	}
	g.bumpUpstream(true)
	g.bumpProxied()

	for k, vv := range res.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(res.StatusCode)
	if len(res.Body) > 0 {
		_, _ = w.Write(res.Body)
	}
	return nil
}

// CacheBody 缓存调用方提供的请求体（独立拷贝）。
func (g *Gateway) CacheBody(key string, body []byte) ([]byte, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed || g.cache == nil {
		return nil, ErrClosed
	}
	return g.cache.Put(key, body), nil
}

// CachedBody 返回缓存体拷贝。
func (g *Gateway) CachedBody(key string) ([]byte, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.cache == nil {
		return nil, false
	}
	return g.cache.Get(key)
}

// AuthenticateRequest 走当前鉴权器。
func (g *Gateway) AuthenticateRequest(r *http.Request) error {
	g.mu.Lock()
	authn := g.authn
	closed := g.closed
	g.mu.Unlock()
	if closed {
		return ErrClosed
	}
	// BUG: 无 nil 防护
	return authn.Authenticate(r)
}

func writeErr(w http.ResponseWriter, err error) {
	code := http.StatusBadGateway
	switch {
	case errors.Is(err, ErrClosed):
		code = http.StatusServiceUnavailable
	case errors.Is(err, ErrNotFound):
		code = http.StatusNotFound
	case errors.Is(err, ErrUnauthorized), errors.Is(err, auth.ErrUnauthorized):
		code = http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		code = http.StatusForbidden
	case errors.Is(err, ErrCanceled):
		code = 499
	case errors.Is(err, ErrBadRoute):
		code = http.StatusBadRequest
	}
	http.Error(w, err.Error(), code)
}

func (g *Gateway) bumpAuth(ok bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if ok {
		g.authOK++
	} else {
		g.authFail++
	}
}

func (g *Gateway) bumpUpstream(ok bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if ok {
		g.upstreamOK++
	} else {
		g.upstreamErr++
	}
}

func (g *Gateway) bumpProxied() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.proxied++
}

// DrainBody 辅助读取剩余 body。
func DrainBody(r io.ReadCloser) {
	if r == nil {
		return
	}
	_, _ = io.Copy(io.Discard, r)
	_ = r.Close()
}

// NormalizeMethod 统一大写 method。
func NormalizeMethod(m string) string {
	return strings.ToUpper(strings.TrimSpace(m))
}
