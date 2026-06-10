'use server';

import { redirect } from 'next/navigation';
import { accountContainer } from '@/di/account_container';
import { AuthQueryParameterDTO } from '@/core/dtos/auth_query_parameter_dto';
import { SigninRequestDTO } from '@/core/dtos/signin_request_dto';
import { ValidationError } from '@/core/domain/exceptions/validation_error';
import { ApiError } from '@/core/domain/exceptions/api_error';
import { mfaContainer } from '@/di/mfa_container';
import { TotpRequestDTO } from '@/core/dtos/totp_request_dto';
import { SessionKey } from '@/core/domain/constants/session_keys';

// ประกอบร่าง (Wire Up Dependencies)
const mfaSerivce = await mfaContainer.mfaService;
const cookieSessionStorage = await mfaContainer.cookieSessionStorage;
export async function executeConfirmTotp(data: TotpRequestDTO) {

    const rawSession = await cookieSessionStorage.get(SessionKey.PreMfa)
    if (!rawSession) {
        return { success: false, errorType: 'SESSION_EXPIRED' };
    }
    try {
        const sessionData = JSON.parse(rawSession);
        const preMfaToken = sessionData.token;
        const flowId = sessionData.flowId;
        const clientId = sessionData.clientId;

        const queryParams: AuthQueryParameterDTO = {
            flowId: flowId,
            clientId: clientId
        };
        // 1. ส่งข้อมูลเข้าสู่ระบบ Core Service ตามที่คุณวาง Flow ไว้ เป๊ะๆ !
        const result = await mfaSerivce.handleConfirmTOTPAction(preMfaToken, queryParams, data);
        if (result.isVerified) {
            await cookieSessionStorage.save(result.sessionName, result.ssoToken, result.sessionExpirySeconds);
            await cookieSessionStorage.clearSession(SessionKey.PreMfa); // ล้างตั๋วทิ้งด้วย
        }
        return {
            success: result.isVerified,
            redirectUrl: result.redirectUrl
        }

    } catch (err: any) {
        if (err.errorCode === 'session_expired') {
            await cookieSessionStorage.clearSession(SessionKey.PreMfa); // ล้างตั๋วทิ้งด้วย
            return { success: false, errorType: 'SESSION_EXPIRED', redirectUrl: err.redirectUrl };
        }
        if (err.errorCode)

            // 🚨 กฎเหล็กข้อที่ 2: แปลงคลาส Error ให้กลายเป็น JSON เปล่าๆ ส่งกลับไป
            if (err instanceof ValidationError) {
                return {
                    success: false,
                    errorType: "ValidationError",
                    errorCode: err.errors,
                    error_description: err.message
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
export async function executeVerifyTotp(data: TotpRequestDTO) {
    const rawSession = await cookieSessionStorage.get(SessionKey.Mfa)
    if (!rawSession) {
        return { success: false, errorType: 'SESSION_EXPIRED' };
    }
    try {
        const sessionData = JSON.parse(rawSession);
        const mfaToken = sessionData.token;
        const flowId = sessionData.flowId;
        const clientId = sessionData.clientId;

        const queryParams: AuthQueryParameterDTO = {
            flowId: flowId,
            clientId: clientId
        };
        // 1. ส่งข้อมูลเข้าสู่ระบบ Core Service ตามที่คุณวาง Flow ไว้ เป๊ะๆ !
        const result = await mfaSerivce.handleVefityTOTPAction(mfaToken, queryParams, data);
        if (result.isVerified) {
            await cookieSessionStorage.clearSession(SessionKey.Mfa); // ล้างตั๋วทิ้งด้วย
        }
        return { success: result.isVerified }

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