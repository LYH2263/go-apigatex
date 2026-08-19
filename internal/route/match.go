package route

import "strings"

// NormalizePath 去掉重复斜杠（保留前导 /）。
func NormalizePath(p string) string {
	if p == "" {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	for strings.Contains(p, "//") {
		p = strings.ReplaceAll(p, "//", "/")
	}
	return p
}

// LongestPrefix 在候选中选最长前缀。
func LongestPrefix(path string, prefixes []string) string {
	best := ""
	for _, p := range prefixes {
		if strings.HasPrefix(path, p) && len(p) > len(best) {
			best = p
		}
	}
	return best
}
