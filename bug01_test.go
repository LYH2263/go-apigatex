package apigatex_test

import (
	"bytes"
	"testing"

	"github.com/LYH2263/go-apigatex"
)

func TestBug01_CachedBodyAlias(t *testing.T) {
	g := apigatex.New()
	defer g.Close()
	body := []byte("payload-one")
	cached, err := g.CacheBody("k1", body)
	if err != nil {
		t.Fatal(err)
	}
	body[0] = 'P'
	if bytes.Equal(cached, body) && cached[0] == 'P' {
		// 若缓存别名，改写会污染返回值；再读库存也应被污染
	}
	got, ok := g.CachedBody("k1")
	if !ok {
		t.Fatal("missing cache")
	}
	if string(got) != "payload-one" {
		t.Fatalf("cache aliased with caller slice: got %q", got)
	}
	if cached[0] == 'P' {
		t.Fatal("returned cache slice shared with caller")
	}
}
