package proxy

// BodyLimit 默认响应体上限。
const BodyLimit = 8 << 20

// ClampLimit 夹紧读取上限。
func ClampLimit(n int) int {
	if n < 1 {
		return BodyLimit
	}
	if n > BodyLimit {
		return BodyLimit
	}
	return n
}
