package proxy

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type closeTracker struct {
	closed bool
	body   string
}

func (c *closeTracker) Read(p []byte) (int, error) {
	if c.body == "" {
		return 0, io.EOF
	}
	n := copy(p, c.body)
	c.body = c.body[n:]
	if c.body == "" {
		return n, io.EOF
	}
	return n, nil
}

func (c *closeTracker) Close() error {
	c.closed = true
	return nil
}

type rtFunc func(*http.Request) (*http.Response, error)

func (f rtFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestBug09_UpstreamResponseBodyClosed(t *testing.T) {
	var tr *closeTracker
	client := &http.Client{Transport: rtFunc(func(r *http.Request) (*http.Response, error) {
		tr = &closeTracker{body: "hello-upstream"}
		return &http.Response{
			StatusCode: 200,
			Body:       tr,
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})}
	eng := NewEngine(client, nil, 0, true)
	res, err := eng.Forward(context.Background(), Request{
		Method: http.MethodGet, Path: "/", Upstream: "http://example.invalid",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(res.Body), "hello") {
		t.Fatalf("body %q", res.Body)
	}
	if tr == nil || !tr.closed {
		t.Fatal("upstream response Body was not Closed")
	}
}
