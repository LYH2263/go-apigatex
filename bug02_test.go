package apigatex_test

import (
	"testing"

	"github.com/LYH2263/go-apigatex"
)

func TestBug02_RouteHeadersShared(t *testing.T) {
	g := apigatex.New()
	defer g.Close()
	if err := g.AddRoute(apigatex.RouteSpec{
		ID: "h1", Path: "/x", Kind: apigatex.MatchExact,
		Upstream: "http://127.0.0.1:9", Headers: map[string]string{"X-A": "1"},
	}); err != nil {
		t.Fatal(err)
	}
	spec, ok := g.GetRoute("h1")
	if !ok {
		t.Fatal("missing")
	}
	spec.Headers["X-A"] = "mutated"
	spec.Headers["X-B"] = "leak"
	again, _ := g.GetRoute("h1")
	if again.Headers["X-A"] != "1" {
		t.Fatalf("headers map shared: %+v", again.Headers)
	}
	if _, ok := again.Headers["X-B"]; ok {
		t.Fatalf("caller insert leaked into store: %+v", again.Headers)
	}
}
