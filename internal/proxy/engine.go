package proxy

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/LYH2263/go-apigatex/internal/errors"
	"github.com/LYH2263/go-apigatex/internal/rewrite"
	"github.com/LYH2263/go-apigatex/internal/upstream"
)

// Request 上游转发请求。
type Request struct {
	Method   string
	Path     string
	Query    string
	Headers  http.Header
	Body     []byte
	Upstream string
	RouteID  string
}

// Response 上游响应缓冲。
type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// Engine 反向代理引擎。
type Engine struct {
	Client    *http.Client
	Pool      *upstream.Pool
	Timeout   time.Duration
	AllowHTTP bool
}

// NewEngine 构造。
func NewEngine(c *http.Client, pool *upstream.Pool, timeout time.Duration, allowHTTP bool) *Engine {
	return &Engine{Client: c, Pool: pool, Timeout: timeout, AllowHTTP: allowHTTP}
}

// Forward 使用调用方 ctx（通常为 req.Context()），取消时中断上游请求。
func (e *Engine) Forward(ctx context.Context, in Request) (*Response, error) {
	if e == nil || e.Client == nil {
		return nil, errors.ErrNilTransport
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil, errors.WrapErr(errors.ErrCanceled, err)
	}
	base := in.Upstream
	if e.Pool != nil {
		base = e.Pool.ResolveURL(in.RouteID, base)
	}
	if base == "" {
		return nil, errors.Wrap(errors.ErrUpstream, "empty upstream")
	}
	if strings.HasPrefix(base, "http://") && !e.AllowHTTP {
		return nil, errors.Wrap(errors.ErrUpstream, "insecure upstream denied")
	}
	urlStr := rewrite.JoinUpstream(base, in.Path)
	if in.Query != "" {
		urlStr += "?" + in.Query
	}
	var rdr io.Reader
	if in.Body != nil {
		rdr = bytes.NewReader(in.Body)
	}
	req, err := http.NewRequestWithContext(ctx, in.Method, urlStr, rdr)
	if err != nil {
		return nil, errors.WrapErr(errors.ErrUpstream, err)
	}
	rewrite.DropHopByHop(in.Headers)
	for k, vv := range in.Headers {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}
	resp, err := e.Client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, errors.WrapErr(errors.ErrCanceled, ctx.Err())
		}
		return nil, errors.WrapErr(errors.ErrUpstream, err)
	}
	// BUG: 未 Close 上游响应 Body
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, errors.WrapErr(errors.ErrUpstream, err)
	}
	return &Response{
		StatusCode: resp.StatusCode,
		Header:     resp.Header.Clone(),
		Body:       body,
	}, nil
}
