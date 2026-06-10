export interface SigninResponseDTO {
  code: 'SETUP_TOTP_REQUIRED' | 'MFA_REQUIRED';
  token: string;
  flowId: string;
  clientId: string;
}