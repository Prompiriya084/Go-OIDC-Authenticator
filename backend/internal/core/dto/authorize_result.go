package dto

type AuthorizeResult struct {
	AuthorizationCode string
	RedirectURI       string
	State             string
}
