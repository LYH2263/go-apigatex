package route

// ComparePriority 更高 priority 优先；同分 ID 字典序小者优先。
func ComparePriority(a, b Route) int {
	if a.Priority != b.Priority {
		if a.Priority > b.Priority {
			return -1
		}
		return 1
	}
	if a.ID < b.ID {
		return -1
	}
	if a.ID > b.ID {
		return 1
	}
	return 0
}

// SortByPriority 原地排序。
func SortByPriority(rows []Route) {
	for i := 0; i < len(rows); i++ {
		for j := i + 1; j < len(rows); j++ {
			if ComparePriority(rows[i], rows[j]) > 0 {
				rows[i], rows[j] = rows[j], rows[i]
			}
		}
	}
}
