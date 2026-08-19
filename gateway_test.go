package apigatex_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LYH2263/go-apigatex"
)

func TestGatewayRoundTrip(t *testing.T) {
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Up", "1")
		_, _ = io.WriteString(w, "pong:"+r.URL.Path)
	}))
	defer up.Close()

	g := apigatex.New(apigatex.WithDefaultToken("secret"))
	defer g.Close()
	if err := g.AddRoute(apigatex.RouteSpec{
		ID: "r1", Path: "/api", Kind: apigatex.MatchPrefix,
		Upstream: up.URL, Methods: []string{"GET"},
	}); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "http://gw/api/v1", nil)
	rr := httptest.NewRecorder()
	g.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	if got := rr.Body.String(); got != "pong:/api/v1" {
		t.Fatalf("body %q", got)
	}
}

func TestAuthRequired(t *testing.T) {
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer up.Close()
	g := apigatex.New(apigatex.WithDefaultToken("secret"))
	defer g.Close()
	_ = g.AddRoute(apigatex.RouteSpec{
		ID: "a1", Path: "/secure", Kind: apigatex.MatchExact,
		Upstream: up.URL, Auth: true, Methods: []string{"GET"},
	})
	req := httptest.NewRequest(http.MethodGet, "http://gw/secure", nil)
	rr := httptest.NewRecorder()
	g.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 got %d", rr.Code)
	}
	req2 := httptest.NewRequest(http.MethodGet, "http://gw/secure", nil)
	req2.Header.Set("Authorization", "Bearer secret")
	rr2 := httptest.NewRecorder()
	g.ServeHTTP(rr2, req2)
	if rr2.Code != 200 {
		t.Fatalf("want 200 got %d", rr2.Code)
	}
}
