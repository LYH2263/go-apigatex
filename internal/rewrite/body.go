package rewrite

import (
	"bytes"
	"strings"

	"github.com/LYH2263/go-apigatex/internal/buffer"
)

// ReplaceBodyString 简单字符串替换（返回新切片）。
func ReplaceBodyString(body []byte, old, new string) []byte {
	if old == "" || body == nil {
		return buffer.CloneBytes(body)
	}
	s := strings.ReplaceAll(string(body), old, new)
	return []byte(s)
}

// PrefixBody 前缀注入。
func PrefixBody(body []byte, prefix string) []byte {
	if prefix == "" {
		return buffer.CloneBytes(body)
	}
	var buf bytes.Buffer
	buf.WriteString(prefix)
	buf.Write(body)
	return buf.Bytes()
}
