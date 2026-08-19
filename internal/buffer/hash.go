package buffer

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256Hex 计算十六进制摘要。
func SHA256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// Fingerprint 短指纹（前 16 hex）。
func Fingerprint(b []byte) string {
	h := SHA256Hex(b)
	if len(h) > 16 {
		return h[:16]
	}
	return h
}
