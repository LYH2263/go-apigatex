package buffer

import "bytes"

// Equal 比较两切片。
func Equal(a, b []byte) bool { return bytes.Equal(a, b) }
