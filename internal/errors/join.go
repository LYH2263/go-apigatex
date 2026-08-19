package errors

import "fmt"

// Join 拼接多个错误（保留首个哨兵可用 %w）。
func Join(errs ...error) error {
	var out error
	for _, e := range errs {
		if e == nil {
			continue
		}
		if out == nil {
			out = e
			continue
		}
		out = fmt.Errorf("%w; %v", out, e)
	}
	return out
}
