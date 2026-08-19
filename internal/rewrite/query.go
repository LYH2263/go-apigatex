package rewrite

import "net/url"

// MergeQuery 合并查询串。
func MergeQuery(raw, extra string) string {
	if extra == "" {
		return raw
	}
	if raw == "" {
		return extra
	}
	return raw + "&" + extra
}

// ParseQuery 容错解析。
func ParseQuery(raw string) url.Values {
	v, err := url.ParseQuery(raw)
	if err != nil {
		return url.Values{}
	}
	return v
}
