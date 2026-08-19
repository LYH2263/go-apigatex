package upstream

// PickWeighted 简单加权选择（确定性：最大 weight，同分取 ID 小）。
func PickWeighted(targets []Target) (Target, bool) {
	var best Target
	found := false
	for _, t := range targets {
		if !t.Healthy {
			continue
		}
		if !found || t.Weight > best.Weight || (t.Weight == best.Weight && t.ID < best.ID) {
			best = t
			found = true
		}
	}
	return best, found
}
