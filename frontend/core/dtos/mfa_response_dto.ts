export interface MfaResponseDTO {
    isVerified: boolean
    ssoToken: string
    redirectUrl: string
    sessionName: string
    sessionExpirySeconds: Number
}