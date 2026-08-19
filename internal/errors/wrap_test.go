package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestWrapErrIs(t *testing.T) {
	base := fmt.Errorf("boom")
	err := WrapErr(ErrUpstream, base)
	if !Is(err, ErrUpstream) || !errors.Is(err, base) {
		t.Fatalf("wrap broken: %v", err)
	}
}
