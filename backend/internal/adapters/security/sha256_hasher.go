package adapters_security

import (
	ports_security "OIDCAuthenticator/internal/core/ports/security"
	"crypto/sha256"
	"encoding/base64"
)

type sha256HasherImpl struct{}

func NewSha256Hasher() ports_security.Sha256Hasher {
	return &sha256HasherImpl{}
}

func (s *sha256HasherImpl) Hash(text string) string {
	// 1. แปลงสตริงเป็น []byte แล้วแฮชด้วย SHA256 ทันที (เทียบเท่า SHA256.HashData)
	hash := sha256.Sum256([]byte(text))

	// 2. แปลงผลลัพธ์ [32]byte ให้เป็น Base64 String
	// (ใช้ [:] เพื่อแปลง Array เป็น Slice ส่งให้ฟังก์ชัน Standard Base64)
	return base64.StdEncoding.EncodeToString(hash[:])
}
