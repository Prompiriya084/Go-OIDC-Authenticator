'use server';

import { redirect } from 'next/navigation';
import { accountContainer } from '@/di/account_container';
import { AuthQueryParameterDTO } from '@/core/dtos/auth_query_parameter_dto';
import { SigninRequestDTO } from '@/core/dtos/signin_request_dto';

// ประกอบร่าง (Wire Up Dependencies)
const accountSerivce = await accountContainer.accountService;
const cookieSessionStorage = await accountContainer.cookieSessionStorage;

async function executeSignIn(queryParams: AuthQueryParameterDTO, data: SigninRequestDTO) {
  try {
    // 1. ส่งข้อมูลเข้าสู่ระบบ Core Service ตามที่คุณวาง Flow ไว้ เป๊ะๆ !
    const result = await accountSerivce.handleSignInAction(queryParams, data);

    // 2. เก็บ Access Token (Pre-MFA,) ลงใน Secure Cookie ของ Next.js
    await cookieSessionStorage.saveToken(result.token, 300);

    // 3. แยกเส้นทางการ Redirect
    if (result.code === 'MFA_REQUIRED') {
      redirect(`/mfa/verify-totp?flowId=${queryParams.flowId}&clientId=${queryParams.clientId}`);
    } else if (result.code === 'SETUP_TOTP_REQUIRED') {
      redirect(`/mfa/confirm-totp?flowId=${queryParams.flowId}&clientId=${queryParams.clientId}`);
    }
  
  } catch (err: any) {
    // if (err.message === 'NEXT_REDIRECT') throw err;
    // return { error: err.message || 'An unexpected error occurred' };
    throw err;
  }
}