package buffer

// AppendCopy 追加到 dst 的新底层数组。
func AppendCopy(dst, src []byte) []byte {
	out := make([]byte, len(dst)+len(src))
	copy(out, dst)
	copy(out[len(dst):], src)
	return out
}
