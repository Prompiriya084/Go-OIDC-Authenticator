package adapters_http_handlers

import (
	domain_exceptions "OIDCAuthenticator/internal/core/domain/exceptions"
	"OIDCAuthenticator/internal/core/dto"
	ports_configurations "OIDCAuthenticator/internal/core/ports/configurations"
	"OIDCAuthenticator/internal/core/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AccountHandler ทำหน้าที่เสมือน AccountController
type HttpAccountHandler struct {
	r          *gin.Engine
	service    services.AccountService
	authConfig ports_configurations.AuthConfiguration
	// เพิ่มตัวแปรอื่นๆ เช่น cache, env หากจำเป็นต้องใช้งานในระดับ HTTP
}

// NewAccountHandler เป็น Constructor สำหรับสร้าง Handler
func NewAccountHandler(
	r *gin.Engine,
	service services.AccountService,
	authConfig ports_configurations.AuthConfiguration,
) *HttpAccountHandler {
	return &HttpAccountHandler{
		r:          r,
		service:    service,
		authConfig: authConfig,
	}
}

// ฟังก์ชันสำหรับผูก Route ของโมดูล Auth เข้ากับ Group ที่ส่งเข้ามา
func (h *HttpAccountHandler) RegisterRoutes() {
	accountGroup := h.r.RouterGroup.Group("/account")
	{
		accountGroup.POST("/signin", h.SignIn)
	}
}

// SignIn จัดการ Request POST /api/account/signin
func (h *HttpAccountHandler) SignIn(c *gin.Context) {
	// 1. รับ Query Parameters (เหมือน [FromQuery])
	flowID := c.Query("flowId")
	clientID := c.Query("clientId")

	// 2. รับ Body (เหมือน [FromBody])
	var req dto.SignInRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request", "error_description": err.Error()})
		return
	}

	input := dto.SignInRequestDTO{
		Username: req.Username,
		Password: req.Password,
	}

	// 3. เรียก Usecase (Business Logic)
	authResult, err := h.service.SignIn(c.Request.Context(), input)
	if err != nil {
		var unauthErr *domain_exceptions.UnauthorizedError
		// จัดการ Exception ต่างๆ เหมือนบล็อก catch
		if errors.As(err, unauthErr) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":             "unauthorized",
				"error_description": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "server_error",
			"error_description": err.Error(),
		})
		return
	}

	// 4. จัดการ token และ response code ตามเงื่อนไข RequireTotp
	var responseCode string
	var responseToken string

	if !authResult.RequireTotp {
		responseCode = "SETUP_TOTP_REQUIRED"
		responseToken, err = h.service.GeneratePreMfaToken(authResult.UserID.String())
		// sessionName := h.authConfig.GetPreMfaSessionName()

		// // ใน Go สมมติการบันทึกข้อมูล Claims ลง Session (Cookie/Redis)
		// session.Set("user_id", authResult.)
		// session.Set("auth_stage", "pre-mfa")
		// session.Set("session_name", sessionName) // จำลองประเภท Session
		// _ = session.Save()

		// redirectPath = fmt.Sprintf("/mfa/setup-totp?flowId=%s&clientId=%s", flowID, clientID)
	} else {
		responseCode = "MFA_REQUIRED"
		responseToken, err = h.service.GenerateMfaToken(authResult.UserID.String())
		// sessionName := h.authConfig.GetMfaSessionName()

		// session.Set("user_id", authResult.UserID)
		// session.Set("auth_stage", "mfa")
		// session.Set("session_name", sessionName)
		// _ = session.Save()

		// redirectPath = fmt.Sprintf("/mfa/verify-totp?flowId=%s&clientId=%s", flowID, clientID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "server_error",
			"error_description": err.Error(),
		})
	}

	// 5. ส่งผลลัพธ์กลับแบบ Ok()
	c.JSON(http.StatusOK, gin.H{
		"code":     responseCode,
		"token":    responseToken,
		"flowId":   flowID,
		"clientId": clientID,
	})
}
