package dto

type TokenResponseDTO struct {
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}
