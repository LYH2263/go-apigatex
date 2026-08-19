package buffer

// CloneBytes 复制切片。nil 保持 nil，避免调用方与缓存共享底层数组。
func CloneBytes(b []byte) []byte {
	// BUG: 直接返回别名，与调用方共享底层数组
	return b
}

// CloneBytesNonNil 空输入也返回长度为 0 的新切片。
func CloneBytesNonNil(b []byte) []byte {
	cp := make([]byte, len(b))
	copy(cp, b)
	return cp
}
