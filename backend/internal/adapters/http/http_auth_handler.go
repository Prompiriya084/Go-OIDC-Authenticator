package adapters_http_handlers

import (
	domain_exceptions "OIDCAuthenticator/internal/core/domain/exceptions"
	"OIDCAuthenticator/internal/core/dto"
	ports_configurations "OIDCAuthenticator/internal/core/ports/configurations"
	"OIDCAuthenticator/internal/core/services"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HttpAuthHandler struct {
	service      services.AuthService
	config       ports_configurations.AuthConfiguration
	frontendHost string // เก็บไว้ที่ตัว Handler เพราะเป็นเรื่องของ URL หน้าบ้าน (Web Routing)
}

func NewHttpAuthHandler(service services.AuthService, config ports_configurations.AuthConfiguration, frontendHost string) *HttpAuthHandler {
	return &HttpAuthHandler{
		service:      service,
		config:       config,
		frontendHost: frontendHost,
	}
}

// ฟังก์ชันสำหรับผูก Route ของโมดูล Auth เข้ากับ Group ที่ส่งเข้ามา
func (h *HttpAuthHandler) RegisterRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.GET("/authorize", h.Authorize)
		authGroup.POST("/token", h.Token)
	}
}

func (h *HttpAuthHandler) Authorize(c *gin.Context) {
	var query dto.AuthorizeRequestDTO
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_query_parameters"})
		return
	}

	// 1. ตรวจสอบเงื่อนไขทาง HTTP (Cookie)
	sessionId, _ := c.Cookie(h.config.GetAuthSessionName())
	// hasSession := errCookie == nil

	flowId := strings.ReplaceAll(uuid.New().String(), "-", "")

	// 2. เรียกใช้ Usecase (ส่งข้อมูลดิบไปให้คำนวณ)
	result, err := h.service.Authorize(c.Request.Context(), query, flowId, sessionId)

	// 3. API Layer จัดการผลลัพธ์ แปลงเป็น JSON หรือ Redirect ตามความเหมาะสมของเว็บ
	if err != nil {
		// 1. ประกาศตัวแปร pointer มารองรับ Type
		var unauthErr *domain_exceptions.UnauthorizedError
		// เคสที่ 1: เกิดจาก Session Expired (ตรรกะการสร้าง Signin URL และตรวจ Ajax อยู่ตรงนี้ทั้งหมด!)
		if errors.As(err, unauthErr) {
			// สร้าง URL ปลายทางสำหรับ Signin (เป็นหน้าที่ของ API Layer)
			// หมายเหตุ: ในทางปฏิบัติจริง flowId ควรจะถูกเจนขึ้นใหม่ในเฟสถัดไป แต่หากต้องการผูก flowId ตรงนี้
			// สามารถขยับตรรกะการสร้าง flowId ใน Usecase มาทำที่ระดับอินทราสตรัคเจอร์ได้ครับ
			signinUrl := fmt.Sprintf("%s/account/signin?clientId=%s", h.frontendHost, query.ClientID)

			isAjax := c.GetHeader("X-Requested-With") == "XMLHttpRequest" ||
				strings.Contains(c.GetHeader("Accept"), "application/json") ||
				c.GetHeader("Origin") != ""

			if isAjax {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":             "signin_required",
					"error_description": err.Error(),
					"redirectUri":       signinUrl, // พ่น JSON ตอบกลับตามที่ตกลงกับเว็บไว้
				})
				return
			}
			c.Redirect(http.StatusFound, signinUrl)
			return
		}

		// เคสที่ 2: เกิดจาก OAuth Error (เช่น Client ID หาย) -> ทำการ Redirect Error กลับไปหา Client ตามมาตรฐาน OAuth2
		var oAuthErr *domain_exceptions.OAuthError
		if errors.As(err, &oAuthErr) {
			errorRedirect := fmt.Sprintf("%s?error=%s&error_description=%s&state=%s",
				query.RedirectURI, oAuthErr.Code, oAuthErr.Message, query.State)
			c.Redirect(http.StatusFound, errorRedirect)
			return
		}

		// เคสที่ 3: แตกตื่นคุมไม่ได้ (Internal Server Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "error_description": err.Error()})
		return
	}

	// กรณีทำงานสำเร็จ: API Layer ทำการสร้าง Redirect URL แล้วส่ง 302 Found ไป
	successRedirect := fmt.Sprintf("%s?code=%s&state=%s", result.RedirectURI, result.AuthorizationCode, result.State)
	c.Redirect(http.StatusFound, successRedirect)
}

func (h *HttpAuthHandler) Token(c *gin.Context) {
	var req dto.TokenRequestDTO
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_form_data"})
		return
	}

	// API Layer แกะ Basic Authentication Header
	if clientId, clientSecret, ok := h.parseBasicAuth(c.GetHeader("Authorization")); ok {
		req.ClientID = clientId
		req.ClientSecret = clientSecret
	}

	// ส่งข้อมูลเรียบร้อยไปให้ Usecase จัดการต่อ
	result, err := h.service.HandleToken(c.Request.Context(), req)
	if err != nil {
		var oAuthErr *domain_exceptions.OAuthError
		if errors.As(err, &oAuthErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":             oAuthErr.Code,
				"error_description": oAuthErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	// API Layer เป็นผู้กำหนดการ Mapping ข้อมูลเป็น JSON โครงสร้างตามที่คุณต้องการ (รวมถึงเรื่อง token_type "Bearer")
	c.JSON(http.StatusOK, gin.H{
		"id_token":      result.IDToken,
		"refresh_token": result.RefreshToken,
		"access_token":  result.AccessToken,
		"token_type":    "Bearer",
		"expires_in":    result.ExpiresIn,
	})
}

func (h *HttpAuthHandler) parseBasicAuth(authHeader string) (string, string, bool) {
	if !strings.HasPrefix(authHeader, "Basic ") {
		return "", "", false
	}
	payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
	if err != nil {
		return "", "", false
	}
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return "", "", false
	}
	return pair[0], pair[1], true
}
