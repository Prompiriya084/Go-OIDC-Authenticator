package domain_exceptions

import "fmt"

type OAuthError struct {
	Code    string `json:"error_code"`
	Message string `json:"message"`
}

func (e *OAuthError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewOAuthError(code, message string) *OAuthError {
	if code == "" {
		code = "internal_server_error"
	}
	return &OAuthError{
		Code:    code,
		Message: message,
	}
}
