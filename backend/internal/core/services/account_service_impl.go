package services

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	domain_exceptions "OIDCAuthenticator/internal/core/domain/exceptions"
	"OIDCAuthenticator/internal/core/dto"
	ports_authentications "OIDCAuthenticator/internal/core/ports/authentications"
	ports_database "OIDCAuthenticator/internal/core/ports/database"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"
	ports_security "OIDCAuthenticator/internal/core/ports/security"
	"context"
	"time"
)

type accountServiceImpl struct {
	txManager       ports_database.TransactionManager
	repoUserAuthen  ports_repositories.UserAuthenRepository
	repoUserMfa     ports_repositories.UserMfaRepository
	passwordHasher  ports_security.PasswordHasher
	jwtTokenService ports_authentications.JwtTokenService
}

func NewAccountService(
	txManager ports_database.TransactionManager,
	repoUserAuthen ports_repositories.UserAuthenRepository,
	repoUserMfa ports_repositories.UserMfaRepository,
	passwordHasher ports_security.PasswordHasher,
	jwtTokenService ports_authentications.JwtTokenService,
) *accountServiceImpl {
	return &accountServiceImpl{
		txManager:       txManager,
		repoUserAuthen:  repoUserAuthen,
		repoUserMfa:     repoUserMfa,
		passwordHasher:  passwordHasher,
		jwtTokenService: jwtTokenService,
	}
}

func (s *accountServiceImpl) SignIn(ctx context.Context, req dto.SignInRequestDTO) (*dto.SignInResponseDTO, error) {
	filterUserAuthen := &domain_entities.UserAuthenFilter{
		Username: &req.Username,
	}
	existingUser, err := s.repoUserAuthen.Get(ctx, filterUserAuthen)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, domain_exceptions.NewUnauthorizedError("", "The account not found")
	}
	isPasswordCorrect := s.passwordHasher.Verify(existingUser.PasswordHash, req.Password)
	if !isPasswordCorrect {
		return nil, domain_exceptions.NewUnauthorizedError("", "The username or password is incorrect.")
	}

	if !existingUser.IsActive {
		return nil, domain_exceptions.NewUnauthorizedError("", "This account hasn't signed in over expected date.")
	}
	dateNow := time.Now()
	dateUtcNow := time.Now().UTC()

	txCtx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, err
	}

	// 💡 แปะป้ายป้องกันระบบระเบิดไว้ตรงนี้เลย
	defer func() {
		if r := recover(); r != nil {
			// ✨ ถ้าเกิด panic กลางคัน โค้ดจะกระโดดมาทำตรงนี้ชัวร์ๆ!
			s.txManager.Rollback(txCtx)
			panic(r) // พ่น panic ต่อเพื่อให้ระบบรู้ว่าแอปพัง
		}
		// เคสที่ 2: ฟังก์ชันนี้จบลงโดยการรีเทิร์น err != nil (ไม่ว่าจะพังจากจุดไหนในฟังก์ชัน)
		if err != nil {
			s.txManager.Rollback(txCtx)
		}
	}()
	// 5. ตรวจสอบวันหมดอายุ (ถ้าเกิน 60 วันนับจากการล็อกอินล่าสุด)
	if existingUser.IsActive && existingUser.ExpiresAt.Before(dateUtcNow) {
		existingUser.IsActive = false

		// อัปเดตสถานะใน DB ทันทีแบบสะท้อนกลับ (Mutation)
		if err := s.repoUserAuthen.Update(txCtx, existingUser); err != nil {
			return nil, err
		}

		if err := s.txManager.Commit(txCtx); err != nil {
			return nil, err
		}

		// เพื่อป้องกัน defer rollback ซ้ำ เพราะเรา handle err != nil ใน defer ไว้ พอวิ่งไปหา defer มันจะเห็นเป็น nil และเดินผ่านไปอย่างสงบ ไม่ไป Rollback ซ้ำครับ
		bizErr := domain_exceptions.NewUnauthorizedError("", "This account hasn't signed in over expected date.")
		return nil, bizErr
	}

	filterUserMfa := &domain_entities.UserMfaFilter{
		ID: &existingUser.ID,
	}
	userMfa, err := s.repoUserMfa.Get(txCtx, filterUserMfa)
	if err != nil {
		return nil, err
	}

	// 7. อัปเดตเวลาการล็อกอินและต่ออายุการใช้งานไปอีก 60 วัน
	existingUser.SignedInAt = dateUtcNow
	existingUser.ExpiresAt = dateUtcNow.AddDate(0, 0, 60)

	existingUser.SignedInAtTH = dateNow
	existingUser.ExpiresAtTH = dateNow.AddDate(0, 0, 60)

	// บันทึกเวลาล็อกอินใหม่ลง DB
	err = s.repoUserAuthen.Update(txCtx, existingUser)
	if err != nil {
		return nil, err
	}
	if err := s.txManager.Commit(txCtx); err != nil {
		return nil, err
	}

	// 8. ประกอบผลลัพธ์ส่งกลับ (ใช้ข้อมูลจากทั่งสองตารางที่ดึงมา)
	return &dto.SignInResponseDTO{
		UserID:      existingUser.ID,
		RequireTotp: userMfa.TotpEnabled,
	}, nil
}
func (s *accountServiceImpl) GeneratePreMfaToken(userID string) (string, error) {
	return s.jwtTokenService.CreatePreMfaToken(userID)
}
func (s *accountServiceImpl) GenerateMfaToken(userID string) (string, error) {
	return s.jwtTokenService.CreateMfaToken(userID)
}
