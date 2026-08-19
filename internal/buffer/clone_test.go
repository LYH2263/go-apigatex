package buffer

import "testing"

func TestCloneBytesIndependent(t *testing.T) {
	src := []byte("hello")
	cp := CloneBytes(src)
	src[0] = 'H'
	if string(cp) != "hello" {
		t.Fatalf("alias leaked: %q", cp)
	}
	if CloneBytes(nil) != nil {
		t.Fatal("nil should stay nil")
	}
}
