package adapters_http_handlers

import (
	"OIDCAuthenticator/internal/core/dto"
	ports_authentications "OIDCAuthenticator/internal/core/ports/authentications"
	"encoding/base64"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpJwksHandler struct {
	keyStore ports_authentications.RsaKeyStoreService
	issuer   string
}

func NewHttpJwksHandler(
	keyStore ports_authentications.RsaKeyStoreService,
	issuer string,
) *HttpJwksHandler {
	return &HttpJwksHandler{
		keyStore: keyStore,
		issuer:   issuer,
	}
}

func (h *HttpJwksHandler) RegisterRoutes(router *gin.Engine) {
	wellKnown := router.Group("/.well-known")
	{
		wellKnown.GET("/jwks.json", h.GetKeys)
		wellKnown.GET("/openid-configuration", h.GetConfig)
	}
}

func (h *HttpJwksHandler) GetKeys(c *gin.Context) {
	// 🚀 ดึงคีย์โดยตรงจาก Key Store Service ของคุณ
	keys, err := h.keyStore.GetPublicKeys(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "server_error",
			"error_description": err.Error(),
		})
		return
	}

	jwks := make([]dto.JSONWebKeyResponseDTO, 0, len(keys))

	for _, x := range keys {
		rsaPubKey := x.Key

		// แปลงตามมาตรฐานสเปก JWK
		nEncoded := base64.RawURLEncoding.EncodeToString(rsaPubKey.N.Bytes())
		eEncoded := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(rsaPubKey.E)).Bytes())

		jwks = append(jwks, dto.JSONWebKeyResponseDTO{
			Kid: x.Kid,
			Kty: "RSA",
			Use: "sig",
			Alg: "RS256",
			N:   nEncoded,
			E:   eEncoded,
		})
	}

	c.JSON(http.StatusOK, gin.H{"keys": jwks})
}

func (h *HttpJwksHandler) GetConfig(c *gin.Context) {
	config := dto.OpenIdConfigResponseDTO{
		Issuer:                            h.issuer,
		JwksURI:                           h.issuer + "/.well-known/jwks.json",
		AuthorizationEndpoint:             h.issuer + "/auth/authorize",
		TokenEndpoint:                     h.issuer + "/auth/token",
		TokenEndpointAuthMethodsSupported: []string{"client_secret_post"},
		GrantTypesSupported:               []string{"authorization_code", "refresh_token"},
		ScopesSupported:                   []string{"openid", "profile", "offline_access", "web_audience.full_access"},
	}
	c.JSON(http.StatusOK, config)
}
