export const SESSION_KEYS = {
  PRE_MFA: 'pre_mfa_session',    // ตั๋วชั่วคราวช่วง Password -> MFA
  AUTH: 'auth_session',          // ตั๋วตัวจริงระยะยาวหลังผ่านด่านครบ
} as const; // 🌟 ใช้ as const เพื่อป้องกันไม่ให้โค้ดส่วนอื่นแอบมาแก้ไขค่าได้ (Read-Only)

// หรือถ้าชอบสไตล์ Enum:
export enum SessionKey {
  PreMfa = 'pre_mfa_session',
  Mfa = "mfa_session",
  Auth = 'auth_session',
}