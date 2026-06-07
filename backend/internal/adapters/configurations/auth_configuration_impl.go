package adapters_configurations

import (
	ports_configurations "OIDCAuthenticator/internal/core/ports/configurations"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type authConfigurationImpl struct {
	TokenIssuer                  string
	PreMfaSessionName            string
	PreMfaSessionExpiryInMinutes int
	MfaSessionName               string
	MfaSessionExpiryInMinutes    int
	AuthSessionName              string
	AuthSessionExpiryInMinutes   int
	AuthCodeExpiryInMinutes      int
	TokenExpiryInMinutes         int
	TotpEncryptionKey            string
	JwtSecret                    string
}

func NewAuthConfiguration() ports_configurations.AuthConfiguration {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	preMfaSessionExpiryInMinutes, _ := strconv.Atoi(os.Getenv("Auth_pre_mfa_session_expiry_minutes"))
	mfaSessionExpiryInMinutes, _ := strconv.Atoi(os.Getenv("Auth_mfa_session_expiry_minutes"))
	authSessionExpiryInMinutes, _ := strconv.Atoi(os.Getenv("Auth_auth_session_expiry_minutes"))
	authCodeExpiryInMinutes, _ := strconv.Atoi(os.Getenv("Auth_auth_code_expiry_minutes"))
	tokenExpiryInMinutes, _ := strconv.Atoi(os.Getenv("Auth_auth_token_expiry_minutes"))
	return &authConfigurationImpl{
		TokenIssuer:                  os.Getenv("Auth_token_issuer"),
		PreMfaSessionName:            os.Getenv("Auth_pre_mfa_session_name"),
		PreMfaSessionExpiryInMinutes: preMfaSessionExpiryInMinutes,
		MfaSessionName:               os.Getenv("Auth_mfa_session_name"),
		MfaSessionExpiryInMinutes:    mfaSessionExpiryInMinutes,
		AuthSessionName:              os.Getenv("Auth_auth_session_name"),
		AuthSessionExpiryInMinutes:   authSessionExpiryInMinutes,
		AuthCodeExpiryInMinutes:      authCodeExpiryInMinutes,
		TokenExpiryInMinutes:         tokenExpiryInMinutes,
		TotpEncryptionKey:            os.Getenv("TOTP_EncryptionKey"),
		JwtSecret:                    os.Getenv("Jwt_Secret"),
	}
}
func (c *authConfigurationImpl) GetTokenIssuer() string {
	return c.TokenIssuer
}
func (c *authConfigurationImpl) GetPreMfaSessionName() string {
	return c.PreMfaSessionName
}
func (c *authConfigurationImpl) GetPreMfaSessionExpiryInMinutes() int {
	return c.PreMfaSessionExpiryInMinutes
}
func (c *authConfigurationImpl) GetMfaSessionName() string {
	return c.MfaSessionName
}
func (c *authConfigurationImpl) GetMfaSessionExpiryInMinutes() int {
	return c.MfaSessionExpiryInMinutes
}

func (c *authConfigurationImpl) GetAuthSessionName() string {
	return c.AuthSessionName
}
func (c *authConfigurationImpl) GetAuthSessionExpiryInMinutes() int {
	return c.AuthSessionExpiryInMinutes
}
func (c *authConfigurationImpl) GetAuthCodeExpiryInMinutes() int {
	return c.AuthCodeExpiryInMinutes
}
func (c *authConfigurationImpl) GetTokenExpiryInMinutes() int {
	return c.TokenExpiryInMinutes
}
func (c *authConfigurationImpl) GetTotpEncryptionKey() string {
	return c.TotpEncryptionKey
}
func (c *authConfigurationImpl) GetJwtSecret() string {
	return c.JwtSecret
}
