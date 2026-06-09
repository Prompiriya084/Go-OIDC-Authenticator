package main

import (
	adapters_authentications "OIDCAuthenticator/internal/adapters/authentications"
	adapters_caching "OIDCAuthenticator/internal/adapters/caching"
	adapters_configurations "OIDCAuthenticator/internal/adapters/configurations"
	adapters_crypto "OIDCAuthenticator/internal/adapters/crypto"
	"OIDCAuthenticator/internal/adapters/dataaccess"
	adapters_http_handlers "OIDCAuthenticator/internal/adapters/http"
	"OIDCAuthenticator/internal/adapters/middleware"
	adapters_repositories "OIDCAuthenticator/internal/adapters/repositories"
	adapters_security "OIDCAuthenticator/internal/adapters/security"
	"OIDCAuthenticator/internal/core/services"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	db := dataaccess.InitDB()

	r := gin.Default()
	txManager := dataaccess.NewTransactionManager(db)

	authConfig := adapters_configurations.NewAuthConfiguration()

	redisClient, err := dataaccess.InitRedis(ctx)
	if err != nil {
		log.Fatalf("Fatal Redis connection: %v", err)
	}

	// 👈 ย้าย defer มาไว้ตรงนี้แทน! มันจะปิดตัวก็ต่อเมื่อฟังก์ชัน main() นี้จบลงเท่านั้น
	defer redisClient.Close()
	//Middleware
	authMiddleware := middleware.NewAuthMiddleware(authConfig.GetJwtSecret())
	//caching
	cachRepository := adapters_caching.NewCacheRepository(redisClient)
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
		jwtTokenService,
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
		cachRepository,
		totpService,
		jwtTokenService,
		randomNumberGenerator,
		encryptionService,
		sha256Hasher,
	)

	//http handlers
	authHandler := adapters_http_handlers.NewHttpAuthHandler(
		authService,
		authConfig,
		os.Getenv("Frontend_Host"),
	)
	accountHandler := adapters_http_handlers.NewAccountHandler(
		accountService,
		authConfig,
	)
	mfaHandler := adapters_http_handlers.NewHttpMfaHandler(
		mfaService,
		authMiddleware,
		authConfig.GetAuthSessionName(),
		authConfig.GetAuthSessionExpiryInMinutes(),
	)

	//register routes
	authHandler.RegisterRoutes(r)
	accountHandler.RegisterRoutes(r)
	mfaHandler.RegisterRoutes(r)

	r.Run(":8080")

}
