package dto

// OpenIdConfigResponse ข้อมูลสำหรับทำ OIDC Discovery (.well-known/openid-configuration)
type OpenIdConfigResponseDTO struct {
	Issuer                            string   `json:"issuer"`
	JwksURI                           string   `json:"jwks_uri"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	GrantTypesSupported               []string `json:"grant_types_supported"`
	ScopesSupported                   []string `json:"scopes_supported"`
}
