package adapters_authentications

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	"OIDCAuthenticator/internal/core/dto"
	ports_authentications "OIDCAuthenticator/internal/core/ports/authentications"
	ports_database "OIDCAuthenticator/internal/core/ports/database"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"time"
)

type rsaKeyStoreServiceImpl struct {
	txManager ports_database.TransactionManager
	repo      ports_repositories.SigningKeyRepository
}

// NewRsaKeyStoreService ทำหน้าที่เป็น Constructor
func NewRsaKeyStoreService(
	txManager ports_database.TransactionManager,
	repo ports_repositories.SigningKeyRepository,
) ports_authentications.RsaKeyStoreService {
	return &rsaKeyStoreServiceImpl{
		txManager: txManager,
		repo:      repo,
	}
}

// GetActiveKey ดึงค่า Active Private Key ออกมาใช้ Sign Token
func (s *rsaKeyStoreServiceImpl) GetActiveKey(ctx context.Context) (string, *rsa.PrivateKey, error) {
	key, err := s.repo.GetActiveKey(ctx)
	if err != nil {
		return "", nil, err
	}
	if key == nil {
		return "", nil, errors.New("active key not found")
	}

	// Decode Base64 ของ PrivateKey
	privateKeyBytes, err := base64.StdEncoding.DecodeString(key.PrivateKey)
	if err != nil {
		return "", nil, err
	}

	// Parse ตัว Byte ให้กลายเป็นโครงสร้าง *rsa.PrivateKey ของ Go
	// หมายเหตุ: ถ้า .NET เซฟเป็น PKCS8 ให้เปลี่ยนไปใช้ x509.ParsePKCS8PrivateKey
	rsaPrivKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		// fallback ลอง parse เผื่อเป็น PKCS#8
		parsedKey, err8 := x509.ParsePKCS8PrivateKey(privateKeyBytes)
		if err8 != nil {
			return "", nil, err // ส่ง error เดิมกลับไป
		}
		var ok bool
		rsaPrivKey, ok = parsedKey.(*rsa.PrivateKey)
		if !ok {
			return "", nil, errors.New("not an RSA private key")
		}
	}

	return key.Kid, rsaPrivKey, nil
}

// GetPublicKeys ดึง Public Keys ทั้งหมดส่งออกไป (มักใช้ทำ JWKS endpoint)
func (s *rsaKeyStoreServiceImpl) GetPublicKeys(ctx context.Context) ([]dto.RsaKeyResult, error) {
	keys, err := s.repo.GetAllKeys(ctx)
	if err != nil {
		return nil, err
	}

	var results []dto.RsaKeyResult
	for _, x := range keys {
		pubKeyBytes, err := base64.StdEncoding.DecodeString(x.PublicKey)
		if err != nil {
			return nil, err
		}

		// Parse Public Key
		rsaPubKey, err := x509.ParsePKCS1PublicKey(pubKeyBytes)
		if err != nil {
			// fallback เผื่อเซฟเป็นแบบ x509 (SubjectPublicKeyInfo)
			parsedPub, errX509 := x509.ParsePKIXPublicKey(pubKeyBytes)
			if errX509 != nil {
				return nil, err
			}
			var ok bool
			rsaPubKey, ok = parsedPub.(*rsa.PublicKey)
			if !ok {
				return nil, errors.New("not an RSA public key")
			}
		}

		results = append(results, dto.RsaKeyResult{
			Kid: x.Kid,
			Key: rsaPubKey,
		})
	}

	return results, nil
}

// RotateKey บังคับหมุนเวียนคีย์ทันที
func (s *rsaKeyStoreServiceImpl) RotateKey(ctx context.Context) error {
	if err := s.generate(ctx); err != nil {
		return err
	}
	return s.removeExpiredKeys(ctx)
}

// RotateKeyIfNeeded ตรวจสอบอายุคีย์ ถ้าเกิน 30 วันจะหมุนเวียนให้เอง
func (s *rsaKeyStoreServiceImpl) RotateKeyIfNeeded(ctx context.Context) error {
	activeKey, err := s.repo.GetActiveKey(ctx)
	if err != nil {
		return err
	}

	if activeKey == nil {
		return s.RotateKey(ctx)
	}

	// คำนวณจำนวนวันที่ผ่านไปนับจากที่สร้างคีย์
	days := time.Since(activeKey.CreatedAt).Hours() / 24

	if days >= 30 {
		return s.RotateKey(ctx)
	}

	return nil
}

// generate เป็น internal method สำหรับสร้างคีย์ชุดใหม่
func (s *rsaKeyStoreServiceImpl) generate(ctx context.Context) error {
	// 1. ดึงคีย์เก่ามาสลับสถานะเป็น IsActive = false
	currentKeys, err := s.repo.GetAllKeys(ctx)
	if err != nil {
		return err
	}

	var keysToUpdate []*domain_entities.SigningKey
	for _, k := range currentKeys {
		if k.IsActive {
			k.IsActive = false
			keysToUpdate = append(keysToUpdate, k)
		}
	}

	// 2. สร้างคีย์ RSA 2048 บิตชุดใหม่
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// แปลง Private Key เป็น ASN.1 DER และ Base64
	privBytes := x509.MarshalPKCS1PrivateKey(rsaKey)
	privBase64 := base64.StdEncoding.EncodeToString(privBytes)

	// แปลง Public Key เป็น ASN.1 DER และ Base64
	pubBytes := x509.MarshalPKCS1PublicKey(&rsaKey.PublicKey)
	pubBase64 := base64.StdEncoding.EncodeToString(pubBytes)

	// ใน Go หากไม่มีตัวช่วยสร้าง UUID สามารถเปลี่ยนไปใช้ไลบรารีภายนอก
	// เช่น "github.com/google/uuid" เพื่อสร้าง string uuid ได้ครับ (อันนี้สมมติเป็น mock string ไว้)
	newKid := "kid_" + time.Now().Format("20060102150405")

	entity := domain_entities.SigningKey{
		Id:         newKid, // หรือใช้ uuid.New().String()
		Kid:        newKid,
		PrivateKey: privBase64,
		PublicKey:  pubBase64,
		CreatedAt:  time.Now().UTC(),
		IsActive:   true,
	}

	s.txManager.Begin(ctx)

	defer func() {
		if r := recover(); r != nil {
			// ✨ ถ้าเกิด panic กลางคัน โค้ดจะกระโดดมาทำตรงนี้ชัวร์ๆ!
			s.txManager.Rollback(ctx)
			panic(r) // พ่น panic ต่อเพื่อให้ระบบรู้ว่าแอปพัง
		}
		// เคสที่ 2: ฟังก์ชันนี้จบลงโดยการรีเทิร์น err != nil (ไม่ว่าจะพังจากจุดไหนในฟังก์ชัน)
		if err != nil {
			s.txManager.Rollback(ctx)
		}
	}()

	// 3. บันทึกลง Repository
	if err := s.repo.Add(ctx, &entity); err != nil {
		return err
	}

	for _, key := range keysToUpdate {
		if err := s.repo.Update(ctx, key); err != nil {
			return err
		}
	}

	if err := s.txManager.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// removeExpiredKeys ลบคีย์เก่าทิ้ง (อายุเกิน 60 วันและไม่ได้ถูกใช้งานอยู่)
func (s *rsaKeyStoreServiceImpl) removeExpiredKeys(ctx context.Context) error {
	cutoff := time.Now().UTC().AddDate(0, 0, -60) // ถอยหลังไป 60 วัน

	expiredEntities, err := s.repo.GetExpiredKeys(ctx, cutoff)
	if err != nil {
		return err
	}

	s.txManager.Begin(ctx)

	defer func() {
		if r := recover(); r != nil {
			// ✨ ถ้าเกิด panic กลางคัน โค้ดจะกระโดดมาทำตรงนี้ชัวร์ๆ!
			s.txManager.Rollback(ctx)
			panic(r) // พ่น panic ต่อเพื่อให้ระบบรู้ว่าแอปพัง
		}
		// เคสที่ 2: ฟังก์ชันนี้จบลงโดยการรีเทิร์น err != nil (ไม่ว่าจะพังจากจุดไหนในฟังก์ชัน)
		if err != nil {
			s.txManager.Rollback(ctx)
		}
	}()

	for _, entitty := range expiredEntities {
		if err := s.repo.Delete(ctx, entitty); err != nil {
			return err
		}
	}
	if err := s.txManager.Commit(ctx); err != nil {
		return err
	}

	return nil
}
