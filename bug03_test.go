package apigatex_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LYH2263/go-apigatex"
)

func TestBug03_ServeAfterCloseNoPanic(t *testing.T) {
	g := apigatex.New()
	_ = g.AddRoute(apigatex.RouteSpec{
		ID: "c1", Path: "/ping", Kind: apigatex.MatchExact,
		Upstream: "http://127.0.0.1:9", Methods: []string{"GET"},
	})
	if err := g.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("ServeHTTP panicked after Close: %v", rec)
		}
	}()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://gw/ping", nil)
	g.ServeHTTP(rr, req)
	if rr.Code == 200 {
		t.Fatal("closed gateway should not proxy successfully")
	}
}
