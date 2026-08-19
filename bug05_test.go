package apigatex

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestBug05_UpstreamErrorWrapsSentinel(t *testing.T) {
	g := New()
	defer g.Close()
	if err := g.AddRoute(RouteSpec{ID: "r1", Methods: []string{"GET"}, Path: "/x", Upstream: "http://127.0.0.1:1"}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := g.TryRequest(ctx, "GET", "/x", nil, nil)
	if err == nil {
		t.Fatal("expected upstream error")
	}
	if !errors.Is(err, ErrUpstream) {
		t.Fatalf("TryRequest must errors.Is ErrUpstream, got %v", err)
	}
}
