package ports_configurations

type AuthConfiguration interface {
	GetTokenIssuer() string
	GetPreMfaSessionName() string
	GetPreMfaSessionExpiryInMinutes() int
	GetMfaSessionName() string
	GetMfaSessionExpiryInMinutes() int
	GetAuthSessionName() string
	GetAuthSessionExpiryInMinutes() int
	GetAuthCodeExpiryInMinutes() int
	GetTotpEncryptionKey() string
	GetTokenExpiryInMinutes() int
}
