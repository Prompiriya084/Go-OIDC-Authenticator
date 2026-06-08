package adapters_authentications

import (
	domain_entities "OIDCAuthenticator/internal/core/domain/entities"
	ports_authentications "OIDCAuthenticator/internal/core/ports/authentications"
	ports_configurations "OIDCAuthenticator/internal/core/ports/configurations"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtTokenServiceImpl struct {
	config   ports_configurations.AuthConfiguration
	keyStore ports_authentications.RsaKeyStoreService
}

func NewJwtTokenService(
	config ports_configurations.AuthConfiguration,
	keyStore ports_authentications.RsaKeyStoreService,
) ports_authentications.JwtTokenService {
	return &jwtTokenServiceImpl{
		config:   config,
		keyStore: keyStore,
	}
}
func (s *jwtTokenServiceImpl) CreateIdToken(
	ctx context.Context,
	userId string,
	clientId string,
	nonce *string,
	userInfo *domain_entities.UserInformation,
	expiryDateUtc time.Time,
) (string, error) {

	kid, rsaKey, err := s.keyStore.GetActiveKey(ctx)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub":            userId,
		"aud":            clientId,
		"iss":            s.config.GetTokenIssuer(),
		"exp":            jwt.NewNumericDate(expiryDateUtc),
		"iat":            jwt.NewNumericDate(time.Now().UTC()),
		"signed_in_time": time.Now().UTC().Unix(),
	}

	if nonce != nil && *nonce != "" {
		claims["nonce"] = *nonce
	}

	if userInfo != nil {
		claims["prefix_en"] = userInfo.PrefixEN
		claims["Prefix_th"] = userInfo.PrefixTH
		claims["name_en"] = userInfo.NameEN
		claims["name_th"] = userInfo.NameTH
		claims["surname_en"] = userInfo.SurnameEN
		claims["surname_th"] = userInfo.SurnameTH
		claims["address_en"] = userInfo.AddressEN
		claims["address_th"] = userInfo.AddressTH
	}

	// สร้าง Token ด้วยวิธี Signing Method RS256
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// ใส่ Key ID (kid) เข้าไปใน Header
	token.Header["kid"] = kid

	// Sign และรับ Token เป็น string
	tokenString, err := token.SignedString(rsaKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func (s *jwtTokenServiceImpl) CreateAccessToken(
	ctx context.Context,
	userId string,
	clientId string,
	audiences []string,
	scopeNames []string,
	expiryDateUtc time.Time,
) (string, error) {

	kid, rsaKey, err := s.keyStore.GetActiveKey(ctx)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub":       userId,
		"client_id": clientId,
		"iss":       s.config.GetTokenIssuer(),
		"exp":       jwt.NewNumericDate(expiryDateUtc),
		"iat":       jwt.NewNumericDate(time.Now().UTC()),
	}

	// จัดการเรื่อง Audience (กรณีมีหลายตัว)
	if len(audiences) == 1 {
		claims["aud"] = audiences[0]
	} else if len(audiences) > 1 {
		claims["aud"] = audiences
	}

	// จัดการเรื่อง Scopes
	if len(scopeNames) > 0 {
		claims["scope"] = scopeNames
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid

	tokenString, err := token.SignedString(rsaKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *jwtTokenServiceImpl) CreatePreMfaToken(userID string) (string, error) {
	expiryDateUTC := time.Now().UTC().Add(time.Duration(s.config.GetPreMfaSessionExpiryInMinutes()) * time.Minute)
	claims := jwt.MapClaims{
		"sub":        userID,
		"auth_stage": "pre-mfa",
		"iss":        s.config.GetTokenIssuer(),
		"exp":        jwt.NewNumericDate(expiryDateUTC),
		"iat":        jwt.NewNumericDate(time.Now().UTC()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(s.config.GetJwtSecret())

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
func (s *jwtTokenServiceImpl) CreateMfaToken(
	userID string,
) (string, error) {
	expiryDateUTC := time.Now().UTC().Add(time.Duration(s.config.GetMfaSessionExpiryInMinutes()) * time.Minute)
	claims := jwt.MapClaims{
		"sub":        userID,
		"auth_stage": "mfa",
		"iss":        s.config.GetTokenIssuer(),
		"exp":        jwt.NewNumericDate(expiryDateUTC),
		"iat":        jwt.NewNumericDate(time.Now().UTC()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(s.config.GetJwtSecret())

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
