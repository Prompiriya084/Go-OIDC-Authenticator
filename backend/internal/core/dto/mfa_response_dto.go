package dto

type MfaResponseDTO struct {
	SessionId            string `json:"session_id"`
	SessionName          string `json:"session_name"`
	SessionExpirySeconds int    `json:"session_expiry_seconds"`
}
