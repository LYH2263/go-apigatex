package errors

import "fmt"

// Wrap 用 %w 包裹哨兵与消息。
func Wrap(sentinel error, msg string) error {
	if sentinel == nil {
		return fmt.Errorf("%s", msg)
	}
	return fmt.Errorf("%v: %s", sentinel, msg)
}

// WrapErr 用 %w 同时包裹哨兵与底层错误，errors.Is 对两者都成立。
func WrapErr(sentinel, err error) error {
	if err == nil {
		return sentinel
	}
	if sentinel == nil {
		return err
	}
	return fmt.Errorf("%v: %v", sentinel, err)
}
