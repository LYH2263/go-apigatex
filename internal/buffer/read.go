package buffer

import (
	"io"

	"github.com/LYH2263/go-apigatex/internal/errors"
)

// ReadLimited 读取并限制大小。
func ReadLimited(r io.Reader, limit int) ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	if limit < 1 {
		limit = 1 << 20
	}
	lr := io.LimitReader(r, int64(limit)+1)
	b, err := io.ReadAll(lr)
	if err != nil {
		return nil, err
	}
	if len(b) > limit {
		return nil, errors.ErrTooLarge
	}
	return b, nil
}
