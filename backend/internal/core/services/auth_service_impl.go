package services

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	domain_exceptions "OIDCAuthenticator/internal/core/domain/exceptions"
	"OIDCAuthenticator/internal/core/dto"
	ports_authentications "OIDCAuthenticator/internal/core/ports/authentications"
	ports_caching "OIDCAuthenticator/internal/core/ports/caching"
	ports_configurations "OIDCAuthenticator/internal/core/ports/configurations"
	ports_database "OIDCAuthenticator/internal/core/ports/database"
	ports_repositories "OIDCAuthenticator/internal/core/ports/repositories"
	ports_security "OIDCAuthenticator/internal/core/ports/security"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AuthUsecase implements IAuthUsecase
type authServiceImpl struct {
	authConfig                ports_configurations.AuthConfiguration
	txManager                 ports_database.TransactionManager
	repoAudience              ports_repositories.AudienceRepository
	repoAuthCode              ports_repositories.AuthCodeRepository
	repoAuthSession           ports_repositories.AuthSessionRepository
	repoClient                ports_repositories.ClientRepository
	repoClientScope           ports_repositories.ClientScopeRepository
	repoClientGrantType       ports_repositories.ClientGrantTypeRepository
	repoGrantType             ports_repositories.GrantTypeRepository
	repoScope                 ports_repositories.ScopeRepository
	repoRefreshToken          ports_repositories.RefreshTokenRepository
	repoRefreshTokenScope     ports_repositories.RefreshTokenScopeRepository
	repoUserInfo              ports_repositories.UserInformationRepository
	repoViewRefreshTokenScope ports_repositories.ViewRefreshTokenScopeRepository
	repoCache                 ports_caching.CacheRepository
	jwtToken                  ports_authentications.JwtTokenService
	randomNumberGenerator     ports_authentications.RandomNumberGenerator
	sha256Hasher              ports_security.Sha256Hasher
	pckeHasher                ports_security.PkceHasher
}

// NewAuthUsecase is the constructor (แทน Constructor ใน C#)
func NewAuthService(
	authConfig ports_configurations.AuthConfiguration,
	txManager ports_database.TransactionManager,
	repoAudience ports_repositories.AudienceRepository,
	repoAuthCode ports_repositories.AuthCodeRepository,
	repoAuthSession ports_repositories.AuthSessionRepository,
	repoClient ports_repositories.ClientRepository,
	repoClientScope ports_repositories.ClientScopeRepository,
	repoClientGrantType ports_repositories.ClientGrantTypeRepository,
	repoGrantType ports_repositories.GrantTypeRepository,
	repoScope ports_repositories.ScopeRepository,
	repoRefreshToken ports_repositories.RefreshTokenRepository,
	repoRefreshTokenScope ports_repositories.RefreshTokenScopeRepository,
	repoUserInfo ports_repositories.UserInformationRepository,
	repoViewRefreshTokenScope ports_repositories.ViewRefreshTokenScopeRepository,
	repoCache ports_caching.CacheRepository,
	jwtToken ports_authentications.JwtTokenService,
	randomNumberGenerator ports_authentications.RandomNumberGenerator,
	sha256Hasher ports_security.Sha256Hasher,
	pckeHasher ports_security.PkceHasher,
) *authServiceImpl {
	return &authServiceImpl{
		authConfig:                authConfig,
		txManager:                 txManager,
		repoAudience:              repoAudience,
		repoAuthCode:              repoAuthCode,
		repoAuthSession:           repoAuthSession,
		repoClient:                repoClient,
		repoClientScope:           repoClientScope,
		repoClientGrantType:       repoClientGrantType,
		repoGrantType:             repoGrantType,
		repoScope:                 repoScope,
		repoRefreshToken:          repoRefreshToken,
		repoRefreshTokenScope:     repoRefreshTokenScope,
		repoUserInfo:              repoUserInfo,
		repoViewRefreshTokenScope: repoViewRefreshTokenScope,
		repoCache:                 repoCache,
		jwtToken:                  jwtToken,
		randomNumberGenerator:     randomNumberGenerator,
		sha256Hasher:              sha256Hasher,
		pckeHasher:                pckeHasher,
	}
}

func (s *authServiceImpl) ValidateGrantType(ctx context.Context, grantType string) bool {
	// var currentGrantTypes = u.unitOfWork.Grantty.Get(...)
	return false
}

func (s *authServiceImpl) Authorize(ctx context.Context, req dto.AuthorizeRequestDTO, flowID string, sessionID string) (*dto.AuthorizeResult, error) {
	// 1. Business Validation (กฎของ OAuth)
	if req.ResponseType != "code" {
		return nil, domain_exceptions.NewOAuthError("unsupported_response_type", "Invalid response type.")
	}
	if req.ClientID == "" {
		return nil, domain_exceptions.NewOAuthError("client_id_required", "The client is null or empty.")
	}
	if req.RedirectURI == "" {
		return nil, domain_exceptions.NewOAuthError("redirect_uri_required", "The redirect uri is null or empty.")
	}
	if req.Scope == "" {
		return nil, domain_exceptions.NewOAuthError("scope_required", "The scopes is null or empty.")
	}
	oidcFlowState := dto.OIDCFlowState{
		ClientID:            req.ClientID,
		RedirectURI:         req.RedirectURI,
		ResponseType:        req.ResponseType,
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
		State:               req.State,
		Scope:               req.Scope,
		Nonce:               req.Nonce,
	}

	// บันทึกลง Cache (ผ่าน Port Interface)
	err := s.repoCache.Set(ctx, flowID, oidcFlowState, 15*time.Minute)
	if err != nil {
		return nil, err
	}
	// 2. ตรวจสอบสถานะของธุรกิจ (มี Session ไหม)
	// ให้ฝั่ง API เป็นคนส่งสถานะมา Usecase จะได้ไม่ยึดติดกับคำว่า Cookie
	if sessionID == "" {
		return nil, domain_exceptions.NewUnauthorizedError("session_expired", "Session has expired.")
	}
	uuidSessionId, _ := uuid.Parse(sessionID)

	// flowId := strings.ReplaceAll(uuid.New().String(), "-", "")

	// สมมติตรรกะการเจน Code

	authCode, err := s.createAuthorizationCode(
		ctx,
		uuidSessionId,
		req,
	)
	if err != nil {
		return nil, err
	}

	// ส่งกลับแค่ข้อมูลดิบๆ ที่ระบบทำงานสำเร็จ
	return &dto.AuthorizeResult{
		AuthorizationCode: authCode,
		RedirectURI:       req.RedirectURI,
		State:             req.State,
	}, nil
}

func (s *authServiceImpl) createAuthorizationCode(ctx context.Context, sessionId uuid.UUID, req dto.AuthorizeRequestDTO) (string, error) {

	filterAuthSession := &domain_entities.AuthSessionFilter{
		SessionID: &sessionId,
	}
	storedAuthSession, err := s.repoAuthSession.Get(ctx, filterAuthSession)

	if err != nil || storedAuthSession == nil {
		return "", domain_exceptions.NewUnauthorizedError("", "Session not found.")
	}

	if storedAuthSession.ExpiresAt.Before(time.Now().UTC()) {
		return "", domain_exceptions.NewUnauthorizedError("", "Session has expired.")
	}

	clientUUID, err := uuid.Parse(req.ClientID)
	if err != nil {
		return "", domain_exceptions.NewOAuthError("invalid_client", "Invalid client ID format.")
	}
	filterClient := &domain_entities.ClientFilter{
		ID: &clientUUID,
	}
	client, err := s.repoClient.Get(ctx, filterClient)
	if err != nil || client == nil {
		return "", domain_exceptions.NewOAuthError("invalid_client", "Invalid client.")
	}

	if client.RedirectURI != req.RedirectURI {
		return "", domain_exceptions.NewOAuthError("invalid_redirect_uri", "Invalid redirect uri.")
	}

	requestedScopes := strings.Split(req.Scope, " ")
	filterClientScope := &domain_entities.ClientScopeFilter{
		ClientID: &clientUUID,
	}
	// LINQ -> Go Filter/Map logic
	clientScopes, err := s.repoClientScope.GetAll(ctx, filterClientScope)
	if err != nil {
		return "", err
	}

	var allowedScopes []uuid.UUID
	for _, x := range clientScopes {
		allowedScopes = append(allowedScopes, x.ScopeID)
	}

	scopes, err := s.repoScope.GetAllByIDs(ctx, allowedScopes)
	if err != nil {
		return "", err
	}

	allowedScopeNames := make(map[string]bool)
	for _, x := range scopes {
		allowedScopeNames[x.Name] = true
	}

	// Check if all requested scopes are allowed
	for _, reqScope := range requestedScopes {
		if !allowedScopeNames[reqScope] {
			return "", domain_exceptions.NewOAuthError("invalid_scope", "The requested scope is invalid or exceeds permission.")
		}
	}

	base64Str, err := s.randomNumberGenerator.ToBase64String()
	if err != nil {
		return "", err
	}

	authorizationCode := strings.NewReplacer("+", "-", "/", "_", "=", "").Replace(base64Str)

	expiryMinutes := s.authConfig.GetAuthCodeExpiryInMinutes()
	loc, _ := time.LoadLocation("Asia/Bangkok") // TH Time

	newAuthCode := &domain_entities.AuthCode{
		Code:            authorizationCode,
		UserID:          storedAuthSession.UserID,
		SessionID:       sessionId,
		ClientID:        clientUUID,
		CodeChallenge:   req.CodeChallenge,
		ChallengeMethod: req.CodeChallengeMethod,
		RequiredScopes:  strings.Join(requestedScopes, " "),
		RedirectURI:     &req.RedirectURI,
		Nonce:           &req.Nonce,
		ExpiresAt:       time.Now().UTC().Add(time.Duration(expiryMinutes) * time.Minute),
		ExpiresAtTH:     time.Now().In(loc).Add(time.Duration(expiryMinutes) * time.Minute),
	}

	s.txManager.Begin(ctx)

	// ใช้ defer สำหรับการ Rollback อัตโนมัติเมื่อเกิด panic หรือ error ก่อนที่จะ Commit
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

	if err := s.repoAuthCode.Add(ctx, newAuthCode); err != nil {
		return "", err
	}
	filterAuthCode := &domain_entities.AuthCodeFilter{
		SessionID: &sessionId,
	}
	oldAuthorizationCode, err := s.repoAuthCode.Get(ctx, filterAuthCode)
	if err == nil && oldAuthorizationCode != nil {
		if err := s.repoAuthCode.Delete(ctx, oldAuthorizationCode); err != nil {
			return "", err
		}
		// if err := s.txManager.Commit(ctx); err != nil {
		// 	return "", err
		// }
	}

	if err := s.txManager.Commit(ctx); err != nil {
		return "", err
	}

	return authorizationCode, nil
}

func (s *authServiceImpl) HandleToken(ctx context.Context, req dto.TokenRequestDTO) (*dto.TokenResponseDTO, error) {
	uuidClientID, err := uuid.Parse(req.ClientID)
	if err != nil {
		return nil, err
	}
	filterClientGrant := &domain_entities.ClientGrantTypeFilter{
		ClientID: &uuidClientID,
	}
	clientGrants, err := s.repoClientGrantType.GetAllWithMaster(ctx, filterClientGrant)
	if err != nil {
		return nil, err
	}
	if len(clientGrants) == 0 {
		return nil, domain_exceptions.NewOAuthError("invalid_grant", "Invalid grant type.")
	}
	//validate grant type
	IsValidGrant := false
	for _, cg := range clientGrants {
		if cg.Grant.Type == req.GrantType {
			IsValidGrant = true
			break
		}
	}

	if !IsValidGrant {
		return nil, domain_exceptions.NewOAuthError("invalid_grant", "Invalid grant type.")
	}

	var tokenResult *dto.TokenResult
	var refreshTokenStr string
	if req.RefreshToken != nil {
		// 2. Dereference the pointer to get the actual normal string value
		refreshTokenStr = *req.RefreshToken
	}
	switch req.GrantType {
	case "authorization_code":
		tokenResult, err = s.handleTokenAuthorizationCode(
			ctx,
			req.Code,
			req.ClientID,
			&req.ClientSecret,
			req.RedirectURI,
			req.CodeVerifier,
		)
	case "refresh_token":
		tokenResult, err = s.handleTokenRefreshToken(
			ctx,
			refreshTokenStr,
			req.ClientID,
			&req.ClientSecret,
		)

	default:
		return nil, domain_exceptions.NewOAuthError("invalid_grant", "unsupported grant type")
	}

	return &dto.TokenResponseDTO{
		IDToken:      tokenResult.IdToken,
		RefreshToken: tokenResult.RefreshToken,
		AccessToken:  tokenResult.AccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.authConfig.GetTokenExpiryInMinutes() * 60,
	}, nil
	// return nil, nil

}

func (s *authServiceImpl) handleTokenAuthorizationCode(
	ctx context.Context,
	authorizationCode string,
	clientId string,
	clientSecret *string,
	redirectUri string,
	codeVerifier string,
) (*dto.TokenResult, error) {

	dateUTCNow := time.Now().UTC()
	dateNow := time.Now()

	filterAuthcode := &domain_entities.AuthCodeFilter{
		Code: &authorizationCode,
	}

	storedAuthCode, err := s.repoAuthCode.Get(ctx, filterAuthcode)
	if err != nil || storedAuthCode == nil {
		return nil, domain_exceptions.NewOAuthError("invalid_grant", "Invalid authorization code.")
	}

	if storedAuthCode.ClientID.String() != clientId {
		return nil, domain_exceptions.NewOAuthError("invalid_client", "Invalid client.")
	}
	if *storedAuthCode.RedirectURI != redirectUri {
		return nil, domain_exceptions.NewOAuthError("invalid_grant", "Invalid redirect uri.")
	}

	filterClient := &domain_entities.ClientFilter{
		ID: &storedAuthCode.ClientID,
	}
	client, err := s.repoClient.Get(ctx, filterClient)
	if err != nil || client == nil {
		return nil, domain_exceptions.NewOAuthError("invalid_client", "Invalid client.")
	}

	if client.RequiredClientSecret {
		if clientSecret == nil || strings.TrimSpace(*clientSecret) == "" {
			return nil, domain_exceptions.NewOAuthError("invalid_client", "Client secret required.")
		}
		if client.HashSecret == nil || *client.HashSecret != *clientSecret {
			return nil, domain_exceptions.NewOAuthError("invalid_client", "Invalid client secret.")
		}
	}

	if client.RequirePCKE && strings.TrimSpace(codeVerifier) == "" {
		return nil, domain_exceptions.NewOAuthError("invalid_request", "Code verifier required.")
	}

	txCtx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, err
	}

	// ใช้ defer สำหรับการ Rollback อัตโนมัติเมื่อเกิด panic หรือ error ก่อนที่จะ Commit
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

	if storedAuthCode.ExpiresAt.Before(dateUTCNow) || storedAuthCode.ExpiresAt.Equal(dateUTCNow) {
		if err := s.repoAuthCode.Delete(txCtx, storedAuthCode); err != nil {
			return nil, err
		}

		if err := s.txManager.Commit(txCtx); err != nil {
			return nil, err
		}

		txErr := domain_exceptions.NewOAuthError("invalid_grant", "Authorization code expired.")
		return nil, txErr
	}

	if !s.pckeHasher.Validate(codeVerifier, storedAuthCode.CodeChallenge) {
		if err := s.repoAuthCode.Delete(txCtx, storedAuthCode); err != nil {
			return nil, err
		}
		if err := s.txManager.Commit(txCtx); err != nil {
			return nil, err
		}
		txErr := domain_exceptions.NewOAuthError("invalid_grant", "Invalid PCKE")
		// เพื่อไม่ให้ defer จับ ได้ว่า err = nil (ไม่ต้องการ roll back)
		err = nil
		return nil, txErr
	}

	requestedScopeNames := strings.Split(storedAuthCode.RequiredScopes, " ")
	accessTokenExpiry := dateUTCNow.Add(time.Duration(client.AccessTokenLifeTimeMinutes) * time.Minute)
	refreshExpiry := dateUTCNow.Add(time.Duration(client.RefreshTokenLifeTimeMinutes) * time.Minute)

	tokenResult, err := s.createTokenResult(
		txCtx,
		storedAuthCode.UserID,
		storedAuthCode.ClientID,
		storedAuthCode.SessionID,
		requestedScopeNames,
		storedAuthCode.Nonce,
		dateUTCNow,
		dateNow,
		accessTokenExpiry,
		accessTokenExpiry,
		refreshExpiry,
		true,
	)
	if err != nil {
		return nil, err
	}
	//ลบทันที่ที่ใช้เสร็จแล้ว
	if err := s.repoAuthCode.Delete(txCtx, storedAuthCode); err != nil {
		return nil, err
	}

	if err := s.txManager.Commit(txCtx); err != nil {
		return nil, err
	}

	return tokenResult, nil
}

func (s *authServiceImpl) handleTokenRefreshToken(
	ctx context.Context,
	refreshToken string,
	clientId string,
	clientSecret *string,
) (*dto.TokenResult, error) {

	dateUTCNow := time.Now().UTC()
	refreshTokenHash := s.sha256Hasher.Hash(refreshToken)

	filterRefreshToken := &domain_entities.RefreshTokenFilter{
		TokenHash: &refreshTokenHash,
	}
	storedRefreshToken, err := s.repoRefreshToken.Get(ctx, filterRefreshToken)
	if err != nil || storedRefreshToken == nil {
		return nil, domain_exceptions.NewOAuthError("invalid_grant", "Invalid or revoked refresh token.")
	}

	if storedRefreshToken.ClientID.String() != clientId {
		return nil, domain_exceptions.NewOAuthError("invalid_client", "Invalid client id.")
	}

	if storedRefreshToken.IsRevoked {
		return nil, domain_exceptions.NewOAuthError("invalid_grant", "Invalid or revoked refresh token.")
	}

	if storedRefreshToken.ExpiresAt.Before(dateUTCNow) {
		return nil, domain_exceptions.NewOAuthError("invalid_grant", "The refresh token has expired.")
	}

	filterClient := &domain_entities.ClientFilter{
		ID: &storedRefreshToken.ClientID,
	}
	client, err := s.repoClient.Get(ctx, filterClient)
	if err != nil || client == nil {
		return nil, domain_exceptions.NewOAuthError("invalid_client", "Invalid client.")
	}

	if client.RequiredClientSecret {
		if clientSecret == nil || strings.TrimSpace(*clientSecret) == "" {
			return nil, domain_exceptions.NewOAuthError("invalid_client", "Client secret required.")
		}
		if client.HashSecret == nil || *client.HashSecret != *clientSecret {
			return nil, domain_exceptions.NewOAuthError("invalid_client", "Invalid client secret.")
		}
	}

	absoluteRefreshExpiry := storedRefreshToken.InitialSignInDate.Add(time.Duration(client.RefreshTokenLifeTimeMinutes) * time.Minute)
	slidingRefreshExpiry := dateUTCNow.Add(time.Duration(client.RefreshTokenLifeTimeMinutes) * time.Minute)

	refreshTokenExpiry := slidingRefreshExpiry
	if slidingRefreshExpiry.Before(absoluteRefreshExpiry) {
		refreshTokenExpiry = slidingRefreshExpiry
	}

	if refreshTokenExpiry.Before(dateUTCNow) {
		return nil, domain_exceptions.NewOAuthError("invalid_grant", "Session expired, Please sign in again.")
	}

	// 4. เริ่ม Transaction ทำงานกับ DB
	txCtx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, err
	}

	// ใช้ defer สำหรับการ Rollback อัตโนมัติเมื่อเกิด panic หรือ error ก่อนที่จะ Commit
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

	storedRefreshToken.IsRevoked = true
	if err := s.repoRefreshToken.Update(txCtx, storedRefreshToken); err != nil {
		return nil, err
	}

	filterRefreshTokenScope := &domain_entities.RefreshTokenScopeFilter{
		TokenID: &storedRefreshToken.ID,
	}
	refreshScopes, err := s.repoRefreshTokenScope.GetAllWithMaster(ctx, filterRefreshTokenScope)
	if err != nil {
		return nil, err
	}

	var scopeNames []string
	for _, x := range refreshScopes {
		scopeNames = append(scopeNames, x.Scope.Name)
	}

	tokenResult, err := s.createTokenResult(
		txCtx,
		storedRefreshToken.UserID,
		storedRefreshToken.ClientID,
		storedRefreshToken.SessionID,
		scopeNames,
		nil,
		storedRefreshToken.InitialSignInDate,
		storedRefreshToken.InitialSignInDateTH,
		dateUTCNow.Add(time.Duration(client.AccessTokenLifeTimeMinutes)*time.Minute),
		dateUTCNow.Add(time.Duration(client.AccessTokenLifeTimeMinutes)*time.Minute),
		refreshTokenExpiry,
		true,
	)
	if err != nil {
		return nil, err
	}

	if err := s.txManager.Commit(txCtx); err != nil {
		return nil, err
	}

	return tokenResult, nil
}

func (s *authServiceImpl) createTokenResult(
	ctx context.Context,
	userID, clientID, sessionID uuid.UUID,
	scopeNames []string,
	nonce *string,
	initialSignInDateUtc, initialSignInDateTh,
	idTokenExpiryDateUTC, accessTokenExpiryDateUTC,
	refreshTokenExpiryDateUTC time.Time,
	issueRefreshToken bool,
) (*dto.TokenResult, error) {

	scopes, err := s.repoScope.GetAllByNames(ctx, scopeNames)
	if err != nil {
		return nil, err
	}

	audiences, err := s.getAudienceNames(ctx, scopes)
	if err != nil {
		return nil, err
	}

	var scopeIds []uuid.UUID
	for _, scope := range scopes {
		scopeIds = append(scopeIds, scope.ID)
	}

	accessToken, err := s.jwtToken.CreateAccessToken(
		ctx,
		userID.String(),
		clientID.String(),
		audiences,
		scopeNames,
		accessTokenExpiryDateUTC,
	)

	idToken, err := s.createIdToken(ctx, userID, clientID, nonce, scopeNames, idTokenExpiryDateUTC)
	if err != nil {
		return nil, err
	}

	var refreshToken string
	hasOfflineAccess := false
	for _, name := range scopeNames {
		if name == "offline_access" {
			hasOfflineAccess = true
			break
		}
	}

	if issueRefreshToken && hasOfflineAccess {
		refreshResult, err := s.issueRefreshToken(
			userID,
			clientID,
			sessionID,
			initialSignInDateUtc,
			initialSignInDateTh,
			refreshTokenExpiryDateUTC,
		)
		if err != nil {
			return nil, err
		}

		if err := s.repoRefreshToken.Add(ctx, &refreshResult.Entity); err != nil {
			return nil, err
		}

		refreshScopes := s.createRefreshTokenScopes(refreshResult.Entity.ID, scopeIds)
		if err := s.repoRefreshTokenScope.AddRange(ctx, refreshScopes); err != nil {
			return nil, err
		}

		if err := s.txManager.Commit(ctx); err != nil {
			return nil, err
		}

		refreshToken = refreshResult.PlainTextToken
	}

	return &dto.TokenResult{
		AccessToken:  accessToken,
		IdToken:      idToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authServiceImpl) getUserInfo(ctx context.Context, userID uuid.UUID, scopes []string) (*domain_entities.UserInformation, error) {
	hasProfileScope := false
	for _, s := range scopes {
		if s == "profile" {
			hasProfileScope = true
			break
		}
	}
	if !hasProfileScope {
		return nil, nil
	}
	filterUserInfo := &domain_entities.UserInformationFilter{
		ID: &userID,
	}
	userInfo, err := s.repoUserInfo.Get(ctx, filterUserInfo)
	if err != nil || userInfo == nil {
		return nil, err
	}
	return userInfo, nil
}

func (s *authServiceImpl) getAudienceNames(ctx context.Context, scopes []*domain_entities.Scope) ([]string, error) {
	var audienceIds []uuid.UUID
	seenIds := make(map[uuid.UUID]bool)

	for _, x := range scopes {
		if x.AudienceID != nil {
			if !seenIds[*x.AudienceID] {
				seenIds[*x.AudienceID] = true
				audienceIds = append(audienceIds, *x.AudienceID)
			}
		}
	}

	audiences, err := s.repoAudience.GetAllByIDs(ctx, audienceIds)
	if err != nil {
		return nil, err
	}

	var names []string
	seenNames := make(map[string]bool)
	for _, x := range audiences {
		if !seenNames[x.Name] {
			seenNames[x.Name] = true
			names = append(names, x.Name)
		}
	}

	return names, nil
}

func (s *authServiceImpl) createIdToken(
	ctx context.Context,
	userID, clientId uuid.UUID,
	nonce *string,
	scopes []string,
	expiryDateUtc time.Time,
) (string, error) {
	userInfo, err := s.getUserInfo(ctx, userID, scopes)
	if err != nil {
		return "", err
	}

	idToken, err := s.jwtToken.CreateIdToken(
		ctx,
		userID.String(),
		clientId.String(),
		nonce,
		userInfo,
		expiryDateUtc,
	)
	if err != nil {
		return "", err
	}

	return idToken, nil
}

func (u *authServiceImpl) createRefreshTokenScopes(tokenId uuid.UUID, scopeIds []uuid.UUID) []*domain_entities.RefreshTokenScope {
	var scopes []*domain_entities.RefreshTokenScope
	for _, scopeId := range scopeIds {
		scopes = append(scopes, &domain_entities.RefreshTokenScope{
			TokenID: tokenId,
			ScopeID: scopeId,
		})
	}
	return scopes
}

func (s *authServiceImpl) issueRefreshToken(
	userId, clientId, sessionId uuid.UUID,
	initialSignInDateUtc, initialSignInDateTh,
	expiryDateUtc time.Time,
) (*dto.RefreshTokenResult, error) {
	dateUTCNow := time.Now().UTC()
	dateNow := time.Now()

	base64Str, err := s.randomNumberGenerator.ToBase64String()
	if err != nil {
		return nil, err
	}

	plainTextToken := strings.NewReplacer("+", "-", "/", "_", "=", "").Replace(base64Str)

	tokenHash := s.sha256Hasher.Hash(plainTextToken)

	entity := domain_entities.RefreshToken{
		ID:                  uuid.New(),
		UserID:              userId,
		ClientID:            clientId,
		SessionID:           sessionId,
		TokenHash:           tokenHash,
		CreatedAt:           dateUTCNow,
		ExpiresAt:           expiryDateUtc,
		InitialSignInDate:   initialSignInDateUtc,
		CreatedAtTH:         dateNow,
		ExpiresAtTH:         expiryDateUtc.Local(), // หรือใช้ Location ของไทยแทน .Local()
		InitialSignInDateTH: initialSignInDateTh,
		IsRevoked:           false,
	}

	return &dto.RefreshTokenResult{
		PlainTextToken: plainTextToken,
		Entity:         entity,
	}, nil
}

// // --- Custom Error Helpers (เพื่อความง่ายในการตรวจสอบ Error Type แบบ OAuth) ---

// type OAuthError struct {
// 	ErrorType        string
// 	ErrorDescription string
// }

// func (e *OAuthError) Error() string {
// 	return fmt.Sprintf("%s: %s", e.ErrorType, e.ErrorDescription)
// }

// func NewOAuthError(errType, desc string) error {
// 	return &OAuthError{ErrorType: errType, ErrorDescription: desc}
// }

// func NewUnauthorizedError(msg string) error {
// 	return errors.New("unauthorized: " + msg)
// }
