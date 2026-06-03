package adapters_crypto

import (
	ports_crypto "OIDCAuthenticator/internal/core/ports/crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

type EncryptionService struct {
	key []byte
}

// NewEncryptionService รับค่าจาก config string แทน IConfiguration ของ C#
func NewEncryptionService(configEncryptionKey string) ports_crypto.EncryptionService {
	// SHA256.HashData(Encoding.UTF8.GetBytes(...))
	hash := sha256.Sum256([]byte(configEncryptionKey))
	return &EncryptionService{
		key: hash[:], // แปลง [32]byte เป็น []byte
	}
}

func (s *EncryptionService) Encrypt(text string) (string, error) {
	plainText := []byte(text)

	// .NET Aes.Create() ใช้ PKCS7 padding เป็นค่าเริ่มต้น
	// แต่ใน Go ของ standard library ไม่มีมาให้ ต้องเติม (pad) เอาเอง
	plainText = pkcs7Pad(plainText, aes.BlockSize)

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	// สร้าง IV ขนาด 16 bytes (aes.BlockSize) และสุ่มค่าใส่ลงไปเหมือน aes.GenerateIV()
	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// ใช้โหมด CBC (ซึ่งเป็นโหมดเริ่มต้นของ .NET Aes)
	mode := cipher.NewCBCEncrypter(block, iv)
	// ทำการเข้ารหัสและเก็บไว้ต่อท้าย IV ทันที (เหมือนกับ aes.IV.Concat(encrypted) ใน C#)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plainText)

	// แปลงเป็น Base64 String
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *EncryptionService) Decrypt(cipherTextStr string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cipherTextStr)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// แยก IV (16 bytes แรก) และ Encrypted Data (ส่วนที่เหลือ) ออกจากกัน
	iv := ciphertext[:aes.BlockSize]
	encryptedData := ciphertext[aes.BlockSize:]

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	if len(encryptedData)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	// ถอดรหัสทับลงไปในตัวแปรเดิม
	mode.CryptBlocks(encryptedData, encryptedData)

	// เอา PKCS7 padding ออกเพื่อให้ได้ข้อมูลดั้งเดิม
	plainText, err := pkcs7Unpad(encryptedData, aes.BlockSize)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

// --- PKCS7 Padding Helpers (จำเป็นต้องใช้ร่วมกับ AES-CBC เพื่อให้คุยกับ .NET รู้เรื่อง) ---

func pkcs7Pad(b []byte, blocksize int) []byte {
	padding := blocksize - (len(b) % blocksize)
	padtext := strings.Repeat(string(rune(padding)), padding)
	return append(b, padtext...)
}

func pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	if len(b)%blocksize != 0 {
		return nil, fmt.Errorf("invalid padding size")
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > blocksize {
		return nil, fmt.Errorf("invalid padding character")
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, fmt.Errorf("invalid padding text")
		}
	}
	return b[:len(b)-n], nil
}
