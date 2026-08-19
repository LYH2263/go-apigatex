package apigatex

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/LYH2263/go-apigatex/internal/buffer"
	"github.com/LYH2263/go-apigatex/internal/proxy"
	"github.com/LYH2263/go-apigatex/internal/rewrite"
)

// TryRequest 管理页试请求：匹配路由后转发，尊重 ctx 取消。
func (g *Gateway) TryRequest(ctx context.Context, method, path string, headers map[string]string, body []byte) (*ProxyResult, error) {
	g.mu.Lock()
	if g.closed || g.client == nil {
		g.mu.Unlock()
		return nil, ErrClosed
	}
	client := g.client
	table := g.table
	pool := g.pool
	timeout := g.proxyTimeout
	allowHTTP := g.allowHTTP
	cache := g.cache
	g.mu.Unlock()

	rt, ok := table.Match(method, path)
	if !ok {
		return nil, ErrNotFound
	}
	bodyCopy := buffer.CloneBytes(body)
	_ = cache // 试请求不强制入缓存
	outPath := rewrite.ApplyPath(path, rt.StripPath, rt.Path)
	hdr := rewrite.MergeHeaders(http.Header{}, rt.Headers)
	for k, v := range headers {
		hdr.Set(k, v)
	}
	start := time.Now()
	eng := proxy.NewEngine(client, pool, timeout, allowHTTP)
	res, err := eng.Forward(ctx, proxy.Request{
		Method:   method,
		Path:     outPath,
		Headers:  hdr,
		Body:     bodyCopy,
		Upstream: rt.Upstream,
		RouteID:  rt.ID,
	})
	if err != nil {
		return nil, err
	}
	out := &ProxyResult{
		Status:     res.StatusCode,
		Upstream:   rt.Upstream,
		BytesOut:   len(res.Body),
		DurationMs: time.Since(start).Milliseconds(),
		Headers:    map[string]string{},
	}
	for k, vv := range res.Header {
		if len(vv) > 0 {
			out.Headers[k] = vv[0]
		}
	}
	return out, nil
}

// TryRequestHTTP 从 http.Request 试请求。
func (g *Gateway) TryRequestHTTP(r *http.Request) (*ProxyResult, error) {
	var body []byte
	if r.Body != nil {
		b, err := io.ReadAll(io.LimitReader(r.Body, int64(g.maxBodyOrDefault())))
		if err != nil {
			return nil, err
		}
		body = b
	}
	hdr := map[string]string{}
	for k, vv := range r.Header {
		if len(vv) > 0 {
			hdr[k] = vv[0]
		}
	}
	return g.TryRequest(r.Context(), r.Method, r.URL.Path, hdr, body)
}

func (g *Gateway) maxBodyOrDefault() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.maxBody > 0 {
		return g.maxBody
	}
	return defaultMaxBody
}

// BuildProbeRequest 构造本地探测请求。
func BuildProbeRequest(method, urlStr string, body []byte) (*http.Request, error) {
	var rdr io.Reader
	if body != nil {
		rdr = bytes.NewReader(body)
	}
	return http.NewRequest(method, urlStr, rdr)
}
