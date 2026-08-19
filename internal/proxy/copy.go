package proxy

import (
	"io"
	"net/http"
)

// CopyHeader 复制响应头。
func CopyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// CopyBody 将 src 写入 dst 并关闭 src。
func CopyBody(dst io.Writer, src io.ReadCloser) (int64, error) {
	if src == nil {
		return 0, nil
	}
	defer src.Close()
	return io.Copy(dst, src)
}
