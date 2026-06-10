package services

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	domain_exceptions "OIDCAuthenticator/internal/core/domain/exceptions"
	"OIDCAuthenticator/internal/core/dto"
	ports_authentications "OIDCAuthenticator/internal/core/ports/authentications"
	ports_caching "OIDCAuthenticator/internal/core/ports/caching"
	ports_configurations "OIDCAuthenticator/internal/core/ports/configurations"
	ports_crypto "OIDCAuthenticator/internal/core/ports/crypto"
	ports_database "OIDCAuthenticator/internal/core/ports/database"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"
	ports_security "OIDCAuthenticator/internal/core/ports/security"
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type mfaServiceImpl struct {
	authConfig            ports_configurations.AuthConfiguration
	txManager             ports_database.TransactionManager
	repoAuthSession       ports_repositories.AuthSessionRepository
	repoClient            ports_repositories.ClientRepository
	repoUserMfa           ports_repositories.UserMfaRepository
	repoCache             ports_caching.CacheRepository
	totp                  ports_authentications.TotpService
	crypto                ports_crypto.EncryptionService
	jwt                   ports_authentications.JwtTokenService
	randomNumberGenerator ports_authentications.RandomNumberGenerator
	hasher                ports_security.Sha256Hasher
}

// NewTotpUsecase ทำหน้าที่เป็น Constructor แบบ DI (Dependency Injection)
func NewMfaService(
	authConfig ports_configurations.AuthConfiguration,
	txManager ports_database.TransactionManager,
	repoAuthSession ports_repositories.AuthSessionRepository,
	repoClient ports_repositories.ClientRepository,
	repoUserMfa ports_repositories.UserMfaRepository,
	repoCache ports_caching.CacheRepository,
	totp ports_authentications.TotpService,
	jwt ports_authentications.JwtTokenService,
	randomNumberGenerator ports_authentications.RandomNumberGenerator,
	crypto ports_crypto.EncryptionService,
	hasher ports_security.Sha256Hasher,
) *mfaServiceImpl {
	return &mfaServiceImpl{
		authConfig:            authConfig,
		txManager:             txManager,
		repoAuthSession:       repoAuthSession,
		repoClient:            repoClient,
		repoUserMfa:           repoUserMfa,
		repoCache:             repoCache,
		totp:                  totp,
		crypto:                crypto,
		jwt:                   jwt,
		randomNumberGenerator: randomNumberGenerator,
		hasher:                hasher,
	}
}
func (s *mfaServiceImpl) GetOIDCFlowState(
	ctx context.Context,
	flowId string,
) (*dto.OIDCFlowState, error) {
	var oidcFlowState dto.OIDCFlowState
	if err := s.repoCache.Get(ctx, flowId, &oidcFlowState); err != nil {
		log.Printf("WARN: failed to get flow state from cache: %v", err)
		return nil, err
	}

	if err := s.repoCache.Delete(ctx, flowId); err != nil {
		log.Printf("WARN: failed to delete flow state from cache: %v", err)
	}
	return &oidcFlowState, nil
}
func (s *mfaServiceImpl) GetDefaultURIByClientId(ctx context.Context, clientId uuid.UUID) (string, error) {
	// ใน Go ไม่นิยมครอบ try-catch พร่ำเพรื่อ แต่จะเช็ก error ตรงๆ ขากลับจาก Repo
	filter := &domain_entities.ClientFilter{ID: &clientId}
	client, err := s.repoClient.Get(ctx, filter)
	if err != nil {
		return "", err
	}
	return client.DefaultRedirectURI, nil
}

func (s *mfaServiceImpl) StartSetup(ctx context.Context, userID uuid.UUID) (string, error) {
	filter := &domain_entities.UserMfaFilter{ID: &userID}
	userMfa, err := s.repoUserMfa.Get(ctx, filter)
	if err != nil {
		return "", err
	}
	if userMfa == nil {
		return "", errors.New("unauthorized: the account not found") // เทียบเท่า UnauthorizedAccessException
	}

	secret, err := s.totp.GenerateSecret()
	if err != nil {
		return "", err
	}

	encryptedSecret, err := s.crypto.Encrypt(secret)
	if err != nil {
		return "", err
	}

	// เรียก Domain Method เพื่อทำการเปลี่ยนสเตตัสโมเดลด้านใน
	userMfa.TotpSecretEncrypted = encryptedSecret
	userMfa.TotpEnabled = false

	if err := s.repoUserMfa.Update(ctx, userMfa); err != nil {
		return "", err
	}

	return s.totp.GenerateQrCodeUri(userID.String(), secret), nil
}

func (s *mfaServiceImpl) ConfirmTotp(ctx context.Context, userId uuid.UUID, code string) (*dto.MfaResponseDTO, error) {
	return s.processTotpVerificationAsync(ctx, userId, code, true)
}

func (s *mfaServiceImpl) VerifyTotp(ctx context.Context, userId uuid.UUID, code string) (*dto.MfaResponseDTO, error) {
	return s.processTotpVerificationAsync(ctx, userId, code, false)
}

// 🔐 แกนกลางระบบประมวลผลธุรกรรม (Core Business Process)
func (s *mfaServiceImpl) processTotpVerificationAsync(
	ctx context.Context,
	userId uuid.UUID,
	code string,
	isConfirmationMode bool,
) (*dto.MfaResponseDTO, error) {

	// 1. เปิดสวิตช์เริ่มระบบ Transaction
	txCtx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, err
	}
	// 2. ใช้ไม้ตาย "defer" คุมพฤติกรรมกู้ชีพและ Rollback ทันทีที่ฟังก์ชันเกิดพังหรือหลุดกลางคัน
	var funcErr error
	defer func() {
		if r := recover(); r != nil {
			s.txManager.Rollback(txCtx)
			panic(r) // ปล่อย panic ต่อเพื่อให้ระบบเลเยอร์บนรับรู้
		}
		if funcErr != nil {
			s.txManager.Rollback(txCtx)
		}
	}()

	// 🟢 [STEP 1] ดึงข้อมูลและตรวจสอบโปรไฟล์ MFA
	filter := &domain_entities.UserMfaFilter{ID: &userId}
	userMfa, funcErr := s.repoUserMfa.Get(txCtx, filter)
	if funcErr != nil {
		return nil, funcErr
	}
	if userMfa == nil {
		funcErr = domain_exceptions.NewOAuthError("unauthorized", "The account not found")
		return nil, funcErr
	}

	// 🟢 [STEP 2] แกะรหัสลับและตรวจสอบรหัสสุ่ม TOTP 6 หลัก
	secret, funcErr := s.crypto.Decrypt(userMfa.TotpSecretEncrypted)
	if funcErr != nil {
		return nil, funcErr
	}

	if !s.totp.Verify(secret, code) {
		// พ่น OAuth Error รูปแบบเดียวกับ C# (สร้าง Custom Error Type เองในเลเยอร์โดเมนได้)
		txErr := domain_exceptions.NewOAuthError("invalid_verification_code", "Invalid verification code, Please try again.")
		return nil, txErr
	}

	// 🟢 [STEP 3] ล้างเซสชันเก่าทิ้งในโหมดปกติ (Verify Mode)
	if !isConfirmationMode {
		sessionFilter := &domain_entities.AuthSessionFilter{UserID: &userId}
		oldSessions, err := s.repoAuthSession.GetAll(txCtx, sessionFilter)
		if err != nil {
			funcErr = err
			return nil, funcErr
		}

		// ใช้ความสามารถเลน len เช็กค่าแบบสไตล์ Go ที่เราสรุปกันไปรอบก่อน 🚀
		if len(oldSessions) > 0 {
			if err := s.repoAuthSession.DeleteRange(txCtx, oldSessions); err != nil {
				funcErr = err
				return nil, funcErr
			}
		}
	}

	// เตรียมเวลาแบบ Local (TH) และ UTC ตรงตามต้นฉบับ
	loc, _ := time.LoadLocation("Asia/Bangkok")
	dateNow := time.Now().In(loc)
	dateUtcNow := time.Now().UTC()
	sessionId := uuid.New() // เทียบเท่า Guid.NewGuid()

	expiryDuration := time.Duration(s.authConfig.GetAuthSessionExpiryInSeconds()) * time.Second

	// 🟢 [STEP 4] ประกอบร่างเซสชันใหม่
	newAuth := &domain_entities.AuthSession{
		SessionID:   sessionId,
		UserID:      userId,
		ExpiresAt:   dateUtcNow.Add(expiryDuration),
		CreatedAt:   dateUtcNow,
		ExpiresAtTH: dateNow.Add(expiryDuration),
		CreatedAtTH: dateNow,
	}

	// 🟢 [STEP 5] อัปเดตข้อมูลประวัติลงฟาร์ม Metadata ของ Domain
	if isConfirmationMode {
		userMfa.LastMfaAt = &dateUtcNow
		userMfa.TotpConfirmedAt = &dateUtcNow
		userMfa.LastMfaAtTH = &dateNow
		userMfa.TotpConfirmedAtTH = &dateNow
		userMfa.TotpEnabled = true
	} else {
		userMfa.LastMfaAt = &dateUtcNow
		userMfa.LastMfaAtTH = &dateNow
	}

	// 🟢 [STEP 6] สั่งบันทึกผ่าน Unit of Work ค้างเอาไว้ในขวดโหลชั่วคราว
	if err := s.repoUserMfa.Update(txCtx, userMfa); err != nil {
		return nil, err
	}

	if err := s.repoAuthSession.Add(txCtx, newAuth); err != nil {
		return nil, err
	}

	// 🚀 [FINAL STEP] ผ่านฉลุยทุกด่าน สั่งจารึกข้อตกลงลงดิสก์จริงพร้อมกันทีเดียว!
	if err := s.txManager.Commit(txCtx); err != nil {
		funcErr = err
		return nil, funcErr
	}

	return &dto.MfaResponseDTO{
		SessionId:            sessionId.String(),
		SessionName:          s.authConfig.GetAuthSessionName(),
		SessionExpirySeconds: s.authConfig.GetAuthSessionExpiryInSeconds(),
	}, nil
}
