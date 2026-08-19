package upstream

// Target 上游目标。
type Target struct {
	ID      string
	BaseURL string
	Weight  int
	Healthy bool
}
