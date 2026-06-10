"use client"

import { executeSignIn } from "@/actions/account_action";
import { AuthQueryParameterDTO } from "@/core/dtos/auth_query_parameter_dto";
import { SigninRequestDTO } from "@/core/dtos/signin_request_dto";
import { useState } from "react"
import { toast } from "react-toastify";

export function useAccount() {
    const [fieldError, setFieldError] = useState<Record<string, string>>({});
    const [loading, setLoading] = useState(false)

    async function submitSignIn(
        queryParams: AuthQueryParameterDTO,
        data: SigninRequestDTO
    ) {
        setLoading(true);
        try {
            const result = await executeSignIn(queryParams, data)
            // 🌟 ถ้าระบบส่งกลับมาสำเร็จ (ในกรณีที่ไม่มีการ redirect)
            if (result?.success) {
                return { success: true };
            }

            // 🌟 ดักจับจัดการพ่นสีที่ UI ตาม "ประเภทตัวอักษร" ที่ส่งกลับมาจากเซิร์ฟเวอร์
            if (result && !result.success) {

                if (result.errorType === 'ValidationError' && result.errors) {
                    console.log(result.errors);
                    setFieldError(result.errors); // พ่นสีแดงตาม Input ต่างๆ
                    return;
                }

                if (result.errorType === 'ApiError') {
                    console.log(result);
                    if (result.errorCode === 'unauthorized') {
                        toast.error(result.message || "Username หรือ Password ไม่ถูกต้อง", {
                            theme: "colored"
                        });
                        return;
                    }
                }

                // สำหรับกรณีฉุกเฉินอื่นๆ
                if (result.errorType === 'UnexpectedError') {
                    toast.error(result.message);
                }
            }
        } catch (err: any) {

            toast.error("ไม่สามารถเชื่อมต่อกับเซิร์ฟเวอร์ได้ กรุณาลองใหม่อีกครั้ง");
        } finally {
            setLoading(false)
        }
    }

    return { submitSignIn, loading, fieldError }
}