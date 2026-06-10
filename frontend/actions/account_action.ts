'use server';

import { redirect } from 'next/navigation';
import { accountContainer } from '@/di/account_container';
import { AuthQueryParameterDTO } from '@/core/dtos/auth_query_parameter_dto';
import { SigninRequestDTO } from '@/core/dtos/signin_request_dto';
import { ValidationError } from '@/core/domain/exceptions/validation_error';
import { ApiError } from '@/core/domain/exceptions/api_error';
import { SessionKey } from '@/core/domain/constants/session_keys';

// ประกอบร่าง (Wire Up Dependencies)
const accountSerivce = await accountContainer.accountService;
const cookieSessionStorage = await accountContainer.cookieSessionStorage;

export async function executeSignIn(queryParams: AuthQueryParameterDTO, data: SigninRequestDTO) {
  try {
    // 1. ส่งข้อมูลเข้าสู่ระบบ Core Service ตามที่คุณวาง Flow ไว้ เป๊ะๆ !
    const result = await accountSerivce.handleSignInAction(queryParams, data);

    const sessionPayload = JSON.stringify({
      token: result.token,       // Pre-MFA Token ที่ Go ให้มา
      flowId: result.flowId,     // Flow ID สำหรับ OIDC Context
      clientId: result.clientId  // Client ID ของแอปพลิเคชันต้นทาง
    });

    // 2. แยกเส้นทางการ Redirect
    if (result.code === 'MFA_REQUIRED') {
      //เก็บ Access Token (Pre-MFA, MFA) ลงใน Secure Cookie ของ Next.js
      await cookieSessionStorage.save(SessionKey.Mfa, sessionPayload, 900);
      redirect(`/mfa/verify-totp?flowId=${result.flowId}&clientId=${result.clientId}`);
    } else if (result.code === 'SETUP_TOTP_REQUIRED') {
      //เก็บ Access Token (Pre-MFA, MFA) ลงใน Secure Cookie ของ Next.js
      await cookieSessionStorage.save(SessionKey.PreMfa, sessionPayload, 900);
      redirect(`/mfa/confirm-totp?flowId=${result.flowId}&clientId=${result.clientId}`);
    }

    return { sucess: true }

  } catch (err: any) {
    // 🚨 กฎเหล็กข้อที่ 1: ถ้าเป็น Error ที่เกิดจากคำสั่ง redirect() ของ Next.js
    // ต้องปล่อยให้มัน throw ทะลุออกไปเลย ห้ามดักไว้เด็ดขาด ไม่งั้นหน้าจอจะไม่ย้ายหน้า!
    if (err.message?.includes('NEXT_REDIRECT') || err.digest?.includes('NEXT_REDIRECT')) {
      throw err;
    }

    // 🚨 กฎเหล็กข้อที่ 2: แปลงคลาส Error ให้กลายเป็น JSON เปล่าๆ ส่งกลับไป
    if (err instanceof ValidationError) {
      return {
        success: false,
        errorType: 'ValidationError',
        errors: err.errors // ส่ง Map ของฟิลด์ที่พังไปให้ UI พ่นสีแดง
      };
    }

    if (err instanceof ApiError) {
      return {
        success: false,
        errorType: 'ApiError',
        errorCode: err.errorCode,
        message: err.message
      };
    }

    // ข้อผิดพลาดอื่นๆ ที่หลุดรอดมา (Fallback)
    return {
      success: false,
      errorType: 'UnexpectedError',
      message: err.message || 'An unexpected error occurred'
    };
  }
}