package adapters_http_handlers

import (
	"OIDCAuthenticator/internal/adapters/middleware"
	domain_exceptions "OIDCAuthenticator/internal/core/domain/exceptions"
	"OIDCAuthenticator/internal/core/dto"
	"OIDCAuthenticator/internal/core/services"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HttpMfaHandler struct {
	service        services.MfaService
	authMiddleware *middleware.AuthMiddleware
}

// อัปเดต Constructor ให้รับตัวแปรน้อยลง เพราะตัดเรื่องชื่อคุกกี้ MFA ชั่วคราวออกไปแล้ว (เปลี่ยนไปใช้ JWT)
func NewHttpMfaHandler(
	service services.MfaService,
	authMiddleware *middleware.AuthMiddleware, // 🚀 ฉีด Middleware ที่คุณสร้างเข้ามาร่วมงานด้วย
) *HttpMfaHandler {
	return &HttpMfaHandler{
		service:        service,
		authMiddleware: authMiddleware,
	}

}
func (h *HttpMfaHandler) RegisterRoutes(router *gin.Engine) {
	// สร้างกลุ่ม API สำหรับ MFA
	mfaGroup := router.Group("/api/mfa")
	{
		// 🔒 เคสที่ 1: หน้า Setup และ Confirm ต้องผ่านด่าน "pre-mfa" เท่านั้น
		preMfaRoutes := mfaGroup.Group("")
		preMfaRoutes.Use(h.authMiddleware.PreMfaAuthMiddleware()) // 🚀 ใช้ตัวแปรที่คุณปั้นไว้
		{
			preMfaRoutes.POST("/setup-totp", h.Setup)
			preMfaRoutes.POST("/confirm-totp", h.ConfirmTotp)
		}

		// 🔒 เคสที่ 2: หน้า Verify ตัวตนปกติ ต้องผ่านด่าน "mfa" เท่านั้น
		mfaRoutes := mfaGroup.Group("")
		mfaRoutes.Use(h.authMiddleware.MfaAuthMiddleware()) // 🚀 ใช้ตัวแปรที่คุณปั้นไว้
		{
			mfaRoutes.POST("/verify-totp", h.VerifyTotp)
		}
	}
}

func (h *HttpMfaHandler) Setup(c *gin.Context) {
	userID, err := h.getUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uuidUserID, _ := uuid.Parse(userID)

	qr, err := h.service.StartSetup(c.Request.Context(), uuidUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "error_description": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"qr": qr})
}

func (h *HttpMfaHandler) ConfirmTotp(c *gin.Context) {
	flowId := c.Query("flowId")
	clientId := c.Query("clientId")

	var req dto.ConfirmTotpRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request_body",
			"error_description": "Please input the verification code.",
		})
		return
	}

	if len(req.Code) != 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request_body",
			"error_description": "The verification length must be 6 digits.",
		})
		return
	}

	// 🚀 ดึง userID ออกมาจากไอเทมที่คุณเอาฝังไว้ในคีย์ c.Set("userID", claims.UserID) ตอนอยู่ใน Middleware
	userId, err := h.getUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// 🚀 แก้ไขจุดนี้: ถ้า Parse ไม่ผ่าน แปลว่าระบหลังบ้านทำงานผิดพลาด ไม่ใช่เพราะลูกค้าไม่ล็อกอิน
	uuidUserID, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "server_error",
			"error_description": "User identity configuration error.",
		})
		return
	}

	// ตรวจสอบความถูกต้องของ OTP
	result, err := h.service.ConfirmTotp(c.Request.Context(), uuidUserID, req.Code)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// // 1. ผูกคุกกี้ชิ้นใหญ่ลงเครื่องลูกค้าเสร็จสรรพ
	// h.setAuthCookie(c, sessionId)

	// 🔥 (ข้อควรจำ) บรรทัด SignOutAsync ของคุกกี้ชั่วคราวตัดทิ้งได้เลยครับ!
	// ปล่อยให้หน้าที่เคลียร์สิทธิ์ของสายชั่วคราวเป็นเรื่องของหน้าบ้าน (Frontend) ไปล้างทิ้งแทน

	// 2. ดึงแคชของ /auth มารันกระบวนการ Redirect ต่อ (โค้ดลอจิกเดิมทำงานต่ออย่างสมบูรณ์แบบ...)
	oidcFlowState, err := h.service.GetOIDCFlowState(c.Request.Context(), flowId)
	if err != nil {
		uuidClientID, _ := uuid.Parse(clientId)
		defaultRedirectURI, _ := h.service.GetDefaultURIByClientId(c.Request.Context(), uuidClientID)
		c.JSON(http.StatusGone, gin.H{
			"error":             "session_expired",
			"error_description": "Session expired, Please sign in again.",
			"redirectUrl":       defaultRedirectURI,
		})
		return
	}

	redirectUrl := h.buildAuthorizeUrl(c, oidcFlowState)

	c.JSON(http.StatusOK, gin.H{
		"isVerified":             true,
		"ssoToken":               result.SessionId,
		"redirectUrl":            redirectUrl,
		"sessionName":            result.SessionName,
		"sessionExpiryInSeconds": result.SessionExpirySeconds,
	})
}
func (h *HttpMfaHandler) VerifyTotp(c *gin.Context) {
	flowId := c.Query("flowId")
	clientId := c.Query("clientId")

	var req dto.ConfirmTotpRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request_body",
			"error_description": "Please input the verification code.",
		})
		return
	}

	if len(req.Code) != 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request_body",
			"error_description": "The verification length must be 6 digits.",
		})
		return
	}

	// 🚀 1. ดึง userID จาก Context ที่แกะมาจาก "mfa" stage token (ล้อตาม ConfirmTotp)
	userId, err := h.getUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// 🚀 แก้ไขจุดนี้: ถ้า Parse ไม่ผ่าน แปลว่าระบหลังบ้านทำงานผิดพลาด ไม่ใช่เพราะลูกค้าไม่ล็อกอิน
	uuidUserID, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "server_error",
			"error_description": "User identity configuration error.",
		})
		return
	}

	// 2. ตรวจสอบโค้ด Totp
	result, err := h.service.VerifyTotp(c.Request.Context(), uuidUserID, req.Code)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// 3. ฝังคุกกี้เซสชันตัวเต็มลงเบราว์เซอร์หน้าบ้าน
	// h.setAuthCookie(c, result.SessionId)

	// 4. ค้นหา State จาก Cache
	oidcFlowState, err := h.service.GetOIDCFlowState(c.Request.Context(), flowId)
	if err != nil {
		uuidClientID, _ := uuid.Parse(clientId)
		defaultRedirectURI, _ := h.service.GetDefaultURIByClientId(c.Request.Context(), uuidClientID)
		c.JSON(http.StatusGone, gin.H{
			"error":             "session_expired",
			"error_description": "Session expired, Please sign in again.",
			"redirectUrl":       defaultRedirectURI,
		})
		return
	}

	redirectUrl := h.buildAuthorizeUrl(c, oidcFlowState)

	// 🔥 5. ตัดคำสั่ง h.clearCookie(c, h.mfaCookieName) ออกไปเลยครับ!
	// เพราะเราย้ายมาใช้ JWT Bearer ตัว Token จะหมดค่าไปเองเมื่อหน้าบ้านลบทิ้ง
	c.JSON(http.StatusOK, gin.H{
		"isVerified":             true,
		"ssoToken":               result.SessionId,
		"redirectUrl":            redirectUrl,
		"sessionName":            result.SessionName,
		"sessionExpiryInSeconds": result.SessionExpirySeconds,
	})
}
func (h *HttpMfaHandler) buildAuthorizeUrl(c *gin.Context, state *dto.OIDCFlowState) string {
	// ปั้น Base URL โดยเช็คโครงสร้าง Scheme ของ Request ปัจจุบัน
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s/auth/authorize", scheme, c.Request.Host)

	// ประกอบ Query Parameters (เหมือน Dictionary ใน C#)
	params := url.Values{}
	params.Add("client_id", state.ClientID)
	params.Add("redirect_uri", state.RedirectURI)
	params.Add("response_type", state.ResponseType)
	params.Add("state", state.State)
	params.Add("scope", state.Scope)
	params.Add("code_challenge", state.CodeChallenge)
	params.Add("code_challenge_method", state.CodeChallengeMethod)
	params.Add("nonce", state.Nonce)

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// func (h *HttpMfaHandler) setAuthCookie(c *gin.Context, sessionId string) {
// 	// คำนวณเวลาหมดอายุเป็นหน่วยวินาทีสำหรับ Go Cookie
// 	maxAge := h.cookieExpiryMin * 60

// 	// c.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
// 	// หมายเหตุ: Gin ไม่มีตัวเลือก SameSite สดๆ ในฟังก์ชันนี้ ต้องใช้ฟังก์ชันย่อยด้านล่างช่วยเซ็ต
// 	c.SetCookie(h.cookieName, sessionId, maxAge, "/", "", true, true)

//		// บังคับให้เป็น SameSite Mode Lax แบบที่คุณเขียนใน .NET
//		c.SetSameSite(http.SameSiteLaxMode)
//	}
func (h *HttpMfaHandler) getUserIdFromContext(c *gin.Context) (string, error) {
	// ดึงข้อมูลไอดีที่ Middleware ถอดรหัสฝากไว้ในระบบ
	val, exists := c.Get("userID")
	if !exists {
		return "", errors.New("user id not found in context")
	}
	userId, ok := val.(string)
	if !ok {
		return "", errors.New("invalid user id type")
	}
	return userId, nil
}
func (h *HttpMfaHandler) handleError(c *gin.Context, err error) {
	var oauthErr *domain_exceptions.OAuthError
	if errors.As(err, &oauthErr) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             oauthErr.Code,
			"error_description": oauthErr.Message,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "error_description": err.Error()})
}
