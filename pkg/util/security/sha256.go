package security

import (
	"crypto/sha256"
)

// SHA256Digest 获取摘要
func SHA256Digest(data []byte) []byte {
	s := sha256.New()
	s.Write(data)
	return s.Sum(nil)
}
