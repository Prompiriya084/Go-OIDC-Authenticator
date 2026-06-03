package ports_authentications

type TotpService interface {
	GenerateSecret() (string, error)
	GenerateQrCodeUri(userID string, secret string) string
	Verify(secret string, code string) bool
}
