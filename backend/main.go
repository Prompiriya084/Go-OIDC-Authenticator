package main

import (
	adapters_authentications "OIDCAuthenticator/internal/adapters/authentications"
	adapters_caching "OIDCAuthenticator/internal/adapters/caching"
	adapters_configurations "OIDCAuthenticator/internal/adapters/configurations"
	adapters_crypto "OIDCAuthenticator/internal/adapters/crypto"
	"OIDCAuthenticator/internal/adapters/dataaccess"
	adapters_http_handlers "OIDCAuthenticator/internal/adapters/http"
	adapters_repositories "OIDCAuthenticator/internal/adapters/repositories"
	adapters_security "OIDCAuthenticator/internal/adapters/security"
	"OIDCAuthenticator/internal/core/services"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	db := dataaccess.InitDB()
	r := gin.Default()
	txManager := dataaccess.NewTransactionManager(db)

	authConfig := adapters_configurations.NewAuthConfiguration()
	//caching
	cachRepository := adapters_caching.NewCacheRepository()
	//db repositories
	audienceRepository := adapters_repositories.NewAudienceRepository(db)
	authCodeRepository := adapters_repositories.NewAuthCodeRepository(db)
	authSessionRepository := adapters_repositories.NewAuthSessionRepository(db)
	clientRepository := adapters_repositories.NewClientRepository(db)
	clientGrantTypeRepository := adapters_repositories.NewClientGrantTypeRepository(db)
	clientScopeRepository := adapters_repositories.NewClientScopeRepository(db)
	grantTypeRepository := adapters_repositories.NewGrantTypeRepository(db)
	refreshTokenRepository := adapters_repositories.NewRefreshTokenRepository(db)
	refreshTokenScopeRepository := adapters_repositories.NewRefreshTokenScopeRepository(db)
	userAuthenRepository := adapters_repositories.NewUserAuthenRepository(db)
	userInformationRepository := adapters_repositories.NewUserInformationRepository(db)
	userMfaRepository := adapters_repositories.NewUserMfaRepository(db)
	scopeRepository := adapters_repositories.NewScopeRepository(db)
	signingKeyRepository := adapters_repositories.NewSigningKeyRepository(db)
	viewRefreshScopeRepository := adapters_repositories.NewViewRefreshTokenScopeRepository(db)

	//utility services
	passwordHasher := adapters_security.NewBryptPasswordHasher()
	pckeHasher := adapters_security.NewPkceHasher()
	sha256Hasher := adapters_security.NewSha256Hasher()

	//authentication services
	keyStore := adapters_authentications.NewRsaKeyStoreService(
		txManager,
		signingKeyRepository,
	)
	jwtTokenService := adapters_authentications.NewJwtTokenService(
		authConfig,
		keyStore,
	)
	totpService := adapters_authentications.NewTotpService()
	randomNumberGenerator := adapters_authentications.NewRandomNumberGenerator()

	//crypto services
	encryptionService := adapters_crypto.NewEncryptionService(authConfig.GetTotpEncryptionKey())

	//services (usecase layer)
	accountService := services.NewAccountService(
		txManager,
		userAuthenRepository,
		userMfaRepository,
		passwordHasher,
	)
	authService := services.NewAuthService(
		authConfig,
		txManager,
		audienceRepository,
		authCodeRepository,
		authSessionRepository,
		clientRepository,
		clientScopeRepository,
		clientGrantTypeRepository,
		grantTypeRepository,
		scopeRepository,
		refreshTokenRepository,
		refreshTokenScopeRepository,
		userInformationRepository,
		viewRefreshScopeRepository,
		cachRepository,
		jwtTokenService,
		randomNumberGenerator,
		sha256Hasher,
		pckeHasher,
	)
	mfaService := services.NewMfaService(
		authConfig,
		txManager,
		authSessionRepository,
		clientRepository,
		userMfaRepository,
		totpService,
		jwtTokenService,
		randomNumberGenerator,
		encryptionService,
		sha256Hasher,
	)

	//http handlers
	authHandler := adapters_http_handlers.NewHttpAuthHandler(
		r,
		authService,
		authConfig,
		os.Getenv("Frontend_Host"),
	)

	//register routes
	authHandler.RegisterRoutes()

}
