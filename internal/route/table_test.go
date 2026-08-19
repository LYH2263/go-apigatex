package route

import "testing"

func TestTableMatchPriority(t *testing.T) {
	tb := NewTable()
	_ = tb.Upsert(Route{ID: "a", Path: "/api", Kind: KindPrefix, Priority: 1, Upstream: "u1"})
	_ = tb.Upsert(Route{ID: "b", Path: "/api/v1", Kind: KindExact, Priority: 10, Upstream: "u2"})
	r, ok := tb.Match("GET", "/api/v1")
	if !ok || r.ID != "b" {
		t.Fatalf("got %+v ok=%v", r, ok)
	}
}
