package dto

import "github.com/google/uuid"

type SignInResponseDTO struct {
	UserID      uuid.UUID `json:"userId"`
	RequireTotp bool      `json:"requireTotp"`
}
