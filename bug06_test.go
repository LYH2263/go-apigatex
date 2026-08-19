package apigatex_test

import (
	"errors"
	"testing"

	"github.com/LYH2263/go-apigatex"
)

func TestBug06_ReloadPersistFailureNoApply(t *testing.T) {
	g := apigatex.New()
	defer g.Close()
	_ = g.AddRoute(apigatex.RouteSpec{
		ID: "old", Path: "/old", Kind: apigatex.MatchExact,
		Upstream: "http://127.0.0.1:9",
	})
	err := g.ReloadRoutes([]apigatex.RouteSpec{{
		ID: "new", Path: "/new", Kind: apigatex.MatchExact,
		Upstream: "http://127.0.0.1:9",
	}}, func([]apigatex.RouteSpec) error {
		return errors.New("disk full")
	})
	if err == nil {
		t.Fatal("expected persist error")
	}
	if _, ok := g.GetRoute("new"); ok {
		t.Fatal("failed persist still applied new routes")
	}
	if _, ok := g.GetRoute("old"); !ok {
		t.Fatal("old routes should remain after failed reload")
	}
}
