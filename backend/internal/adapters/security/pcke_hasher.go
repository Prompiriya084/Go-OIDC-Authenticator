package adapters_security

import (
	ports_security "OIDCAuthenticator/internal/core/ports/security"
	"crypto/sha256"
	"encoding/base64"
)

type pckeHasherImpl struct{}

func NewPkceHasher() ports_security.PkceHasher {
	return &pckeHasherImpl{}
}

func (p *pckeHasherImpl) Validate(verifier string, challenge string) bool {
	// 1. แปลง string verifier เป็น []byte และทำ SHA256 Hash
	hash := sha256.Sum256([]byte(verifier))

	// 2. แปลงผลลัพธ์เป็น Base64 URL แบบไม่มี Padding (ตรงกับ Base64UrlEncoder ใน C#)
	// เนื่องด้วย sha256.Sum256 คืนค่ากลับมาเป็น Array ขนาด fixed [32]byte
	// เราจึงต้องใส่ [:] เพื่อแปลงให้มันเป็น Slice ก่อนส่งให้ฟังก์ชัน Encode
	computed := base64.RawURLEncoding.EncodeToString(hash[:])

	// 3. ตรวจสอบเงื่อนไขเทียบค่าความถูกต้อง
	return computed == challenge
}
