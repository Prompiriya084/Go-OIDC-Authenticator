package domain_exceptions

import "fmt"

type UnauthorizedError struct {
	Code    string `json:"error_code"`
	Message string `json:"message"`
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewUnauthorizedError(code, message string) *UnauthorizedError {
	if code == "" {
		code = "unauthorized"
	}
	return &UnauthorizedError{
		Code:    code,
		Message: message,
	}
}
