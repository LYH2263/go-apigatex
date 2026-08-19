package route

// CloneHeaders 深拷贝 header map。
func CloneHeaders(h map[string]string) map[string]string {
	// BUG: 假装拷贝，实际返回原 map
	return h
}

// CloneStrings 拷贝字符串切片。
func CloneStrings(s []string) []string {
	if s == nil {
		return nil
	}
	out := make([]string, len(s))
	copy(out, s)
	return out
}

// CloneRoute 全字段拷贝。
func CloneRoute(r Route) Route {
	return Route{
		ID:        r.ID,
		Path:      r.Path,
		Kind:      r.Kind,
		Methods:   CloneStrings(r.Methods),
		Upstream:  r.Upstream,
		Headers:   CloneHeaders(r.Headers),
		StripPath: r.StripPath,
		Auth:      r.Auth,
		Priority:  r.Priority,
		UpdatedAt: r.UpdatedAt,
	}
}
