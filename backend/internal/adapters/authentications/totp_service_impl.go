package adapters_authentications

import (
	ports_authentications "OIDCAuthenticator/internal/core/ports/authentications"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type TotpService struct{}

func NewTotpService() ports_authentications.TotpService {
	return &TotpService{}
}

func (s *TotpService) GenerateSecret() (string, error) {
	// สร้าง random key ขนาด 20 bytes เหมือน KeyGeneration.GenerateRandomKey(20)
	secretBytes := make([]byte, 20)
	if _, err := io.ReadFull(rand.Reader, secretBytes); err != nil {
		return "", err
	}
	// แปลงเป็น Base32 แบบไม่เติม padding (ตามมาตรฐานที่ Google Authenticator ชอบใช้)
	// หากต้องการให้เหมือน C# เป๊ะๆ ที่มี padding ให้ใช้ base32.StdEncoding.EncodeToString
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secretBytes), nil
}

func (s *TotpService) GenerateQrCodeUri(userID string, secret string) string {
	// ใช้ fmt.Sprintf ในการต่อ String (String Interpolation)
	return fmt.Sprintf("otpauth://totp/hino-authentication:%s?secret=%s&issuer=HinoAuthenticator", userID, secret)
}

func (s *TotpService) Verify(secret string, code string) bool {
	// ปรับตรรกะให้ยืดหยุ่น: ตัวไลบรารีของ Go มักต้องการ Uppercase และอาจจะมีหรือไม่มี Padding ก็ได้
	secret = strings.ToUpper(secret)

	// VerificationWindow.RfcSpecifiedNetworkDelay ของ C# มักยอมรับเวลาเยื้องได้ 1 window (30 วินาที หน้า/หลัง)
	// ใน Go เราสามารถกำหนดด้วยการตั้งค่า Skew = 1
	opts := totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	}

	valid, _ := totp.ValidateCustom(code, secret, time.Now().UTC(), opts)
	return valid
}
