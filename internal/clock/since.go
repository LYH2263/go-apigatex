package clock

import "time"

// Elapsed 计算耗时毫秒。
func Elapsed(c Clock, start time.Time) int64 {
	if c == nil {
		return time.Since(start).Milliseconds()
	}
	return c.Since(start).Milliseconds()
}
