package services

import (
	"context"

	"github.com/google/uuid"
)

type SignInResponseDTO struct {
	UserId      uuid.UUID `json:"userId"`
	RequireTotp bool      `json:"requireTotp"`
}

type AccountService interface {
	SignIn(ctx context.Context, username string, password string) (*SignInResponseDTO, error)
}
