package dto

type TokenRequestDTO struct {
	GrantType    string `form:"grant_type"`
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri"`
	CodeVerifier string `form:"code_verifier"`
	ClientID     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
}
