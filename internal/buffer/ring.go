package buffer

// Ring 简易字节环缓冲（试请求历史）。
type Ring struct {
	slots [][]byte
	pos   int
	size  int
}

// NewRing 构造。
func NewRing(n int) *Ring {
	if n < 1 {
		n = 8
	}
	return &Ring{slots: make([][]byte, n)}
}

// Push 推入拷贝。
func (r *Ring) Push(b []byte) {
	r.slots[r.pos%len(r.slots)] = CloneBytes(b)
	r.pos++
	if r.size < len(r.slots) {
		r.size++
	}
}

// Snapshot 返回从旧到新的拷贝列表。
func (r *Ring) Snapshot() [][]byte {
	out := make([][]byte, 0, r.size)
	start := 0
	if r.pos >= len(r.slots) {
		start = r.pos % len(r.slots)
	}
	for i := 0; i < r.size; i++ {
		idx := (start + i) % len(r.slots)
		out = append(out, CloneBytes(r.slots[idx]))
	}
	return out
}
