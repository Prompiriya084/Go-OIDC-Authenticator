// src/app/hooks/use-confirm-totp.ts
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { toast } from 'react-toastify';
import { executeConfirmTotp } from '@/actions/mfa_action';
import { AuthQueryParameterDTO } from '@/core/dtos/auth_query_parameter_dto';
import { TotpRequestDTO } from '@/core/dtos/totp_request_dto';

export function useMfa() {
    const router = useRouter();
    const [fieldError, setFieldError] = useState<Record<string, string>>({});
    const [loading, setLoading] = useState(false);

    async function submitConfirm(data: TotpRequestDTO) {
        setLoading(true);
        try {
            const response = await executeConfirmTotp(data);

            if (response.errorCode === 'session_expired' && response.redirectUrl) {
                // 🚨 แจ้งเตือนผู้ใช้ด้วย UX ที่ดีก่อนว่าเซสชันขาดตอน
                toast.error("Session expired, Redirecting to the requested application.", {
                    position: "top-center",
                    autoClose: 3000 // ให้เวลาผู้ใช้อ่าน 3 วินาที
                });

                // ⏱️ หน่วงเวลาเล็กน้อยเพื่อให้ Toast แสดงผล จากนั้นสั่งเตะผู้ใช้กลับไปยัง Default URI ของ Client ทันที
                setTimeout(() => {
                    window.location.href = response.redirectUrl;
                    // แนะนำใช้ window.location.href หากต้องกระโดดข้าม Domain ออกไปยัง Client Audience
                }, 3000);

                return;
            }

            if (response.errorType == "ValidationError") {
                
            }

            if (response.success && response.nextStepPath) {
                router.push(response.nextStepPath);
            }
        } catch (err) {
            toast.error("เกิดข้อผิดพลาดในการเชื่อมต่อระบบ");
        } finally {
            setLoading(false);
        }
    }

    return { submitConfirm, loading, fieldError};
}