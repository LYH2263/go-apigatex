package apigatex_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LYH2263/go-apigatex"
)

func TestBug04_NilAuthenticatorNoPanic(t *testing.T) {
	g := apigatex.New(apigatex.WithAuthenticator(nil))
	defer g.Close()
	_ = g.AddRoute(apigatex.RouteSpec{
		ID: "n1", Path: "/secure", Kind: apigatex.MatchExact,
		Upstream: "http://127.0.0.1:9", Auth: true, Methods: []string{"GET"},
	})
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("nil authenticator panicked: %v", rec)
		}
	}()
	req := httptest.NewRequest(http.MethodGet, "http://gw/secure", nil)
	if err := g.AuthenticateRequest(req); err == nil {
		// 缺省应有默认鉴权或显式 ErrNilAuth，不得成功且不得 panic
		t.Fatal("expected auth error with nil authenticator")
	}
	rr := httptest.NewRecorder()
	g.ServeHTTP(rr, req)
}
