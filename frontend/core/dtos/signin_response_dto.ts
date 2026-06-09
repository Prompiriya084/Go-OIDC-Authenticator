export interface SignInResponseDTO {
  code: 'SETUP_TOTP_REQUIRED' | 'MFA_REQUIRED';
  token: string;
}