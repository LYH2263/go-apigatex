package apigatex_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/LYH2263/go-apigatex"
)

func TestBug10_CloseSyncsBeforeDropRoutes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "routes.json")
	g := apigatex.New(apigatex.WithPersistPath(path))
	if err := g.AddRoute(apigatex.RouteSpec{
		ID: "keep", Path: "/keep", Kind: apigatex.MatchExact,
		Upstream: "http://127.0.0.1:9",
	}); err != nil {
		t.Fatal(err)
	}
	if err := g.Close(); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(raw, []byte("keep")) {
		t.Fatalf("persist missing routes after Close (dropped before Sync): %s", raw)
	}
}
