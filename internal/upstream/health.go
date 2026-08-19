package upstream

import (
	"context"
	"net/http"
	"time"

	"github.com/LYH2263/go-apigatex/internal/errors"
)

// Probe 对 baseURL 发 HEAD/GET 探测；尊重 ctx。
func Probe(ctx context.Context, client *http.Client, baseURL string) error {
	if client == nil {
		return errors.ErrNilTransport
	}
	if ctx == nil {
		ctx = context.Background()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return errors.WrapErr(errors.ErrUpstream, err)
	}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return errors.WrapErr(errors.ErrCanceled, ctx.Err())
		}
		return errors.WrapErr(errors.ErrUpstream, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		return errors.Wrap(errors.ErrUpstream, "unhealthy status")
	}
	return nil
}

// Mark 更新健康状态。
func (p *Pool) Mark(id string, healthy bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	t, ok := p.byID[id]
	if !ok {
		return
	}
	t.Healthy = healthy
	p.byID[id] = t
}

// ProbeLoop 周期性探测，直到 ctx 结束。
func (p *Pool) ProbeLoop(ctx context.Context, client *http.Client, every time.Duration) {
	if every <= 0 {
		every = time.Second
	}
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			for _, tgt := range p.List() {
				err := Probe(ctx, client, tgt.BaseURL)
				p.Mark(tgt.ID, err == nil)
			}
		}
	}
}
