package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// PreMfaClaims นิยามโครงสร้างข้อมูลที่จะอยู่ใน Pre-MFA Token
type PreMfaClaims struct {
	UserID    string `json:"sub"`
	AuthStage string `json:"auth_stage"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	secretKey string
}

func NewAuthMiddleware(secretKey string) *AuthMiddleware {
	return &AuthMiddleware{
		secretKey: secretKey,
	}
}

// PreMfaAuthMiddleware จะเรียกฟังก์ชันกลางและระบุว่าต้องการตรวจ stage "pre-mfa"
func (m *AuthMiddleware) PreMfaAuthMiddleware() gin.HandlerFunc {
	return m.validateStage("pre-mfa")
}

// MfaAuthMiddleware จะเรียกฟังก์ชันกลางและระบุว่าต้องการตรวจ stage "mfa"
func (m *AuthMiddleware) MfaAuthMiddleware() gin.HandlerFunc {
	return m.validateStage("mfa")
}

// validateStage (ฟังก์ชันลับภายใน) ตัวจัดการหลักที่แชร์ Logic ร่วมกันเพื่อลด Code Duplication
func (m *AuthMiddleware) validateStage(requiredStage string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. ดึง Authorization Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "error_description": "Missing Authorization header"})
			return
		}

		// 2. ตรวจสอบ Format "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "error_description": "Authorization header format must be Bearer {token}"})
			return
		}

		tokenString := parts[1]
		claims := &PreMfaClaims{}

		// 3. Parse และ Validate Token โดยใช้ m.secretKey จาก struct
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(m.secretKey), nil
		})

		// 4. ตรวจสอบความถูกต้องของ Token
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "error_description": "Invalid or expired token"})
			return
		}

		// 5. ตรวจสอบ Dynamic Stage ตามที่พาสเข้ามา
		if claims.AuthStage != requiredStage {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden", "error_description": "Invalid authentication stage"})
			return
		}

		// 6. บันทึกข้อมูลเข้า Context
		c.Set("userID", claims.UserID)

		c.Next()
	}
}
