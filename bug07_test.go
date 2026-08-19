package apigatex_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LYH2263/go-apigatex"
)

func TestBug07_ProxyHonorsRequestContext(t *testing.T) {
	started := make(chan struct{})
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		select {
		case <-r.Context().Done():
			return
		case <-time.After(3 * time.Second):
			_, _ = io.WriteString(w, "late")
		}
	}))
	defer up.Close()

	g := apigatex.New(apigatex.WithProxyTimeout(5 * time.Second))
	defer g.Close()
	_ = g.AddRoute(apigatex.RouteSpec{
		ID: "slow", Path: "/slow", Kind: apigatex.MatchExact,
		Upstream: up.URL, Methods: []string{"GET"},
	})

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "http://gw/slow", nil).WithContext(ctx)
	rr := httptest.NewRecorder()
	done := make(chan struct{})
	go func() {
		g.ServeHTTP(rr, req)
		close(done)
	}()
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("upstream not hit")
	}
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("proxy ignored request context cancel")
	}
	_ = errors.New("ok")
}
