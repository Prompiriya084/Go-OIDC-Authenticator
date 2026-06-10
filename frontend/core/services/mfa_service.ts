import z from "zod";
import { ValidationError } from "../domain/exceptions/validation_error";
import { AuthQueryParameterDTO } from "../dtos/auth_query_parameter_dto";
import { SigninResponseDTO } from "../dtos/signin_response_dto";
import { TotpRequestDTO } from "../dtos/totp_request_dto";
import { MfaApiRepositoryPort } from "../ports/repositories/mfa_api_repository_port";
import { MfaResponseDTO } from "../dtos/mfa_response_dto";

const ConfirmTOTPSchema = z.object({
  code: z.string().length(6, "Please input the t-otp code 6 digits.").max(6),
})
// type validateSignIn = z.infer<typeof SignInSchema>;

export class MfaService {
    constructor(private repo: MfaApiRepositoryPort){

    }
    async handleConfirmTOTPAction(
        preMfaToken: string,
        queryParams: AuthQueryParameterDTO,
        data: TotpRequestDTO
    ): Promise<MfaResponseDTO> {
        // 1.Validate input (Zod)
        const parsed = ConfirmTOTPSchema.safeParse(data)
        if (!parsed.success) {
            const fieldErrors = parsed.error.flatten().fieldErrors;
            const mappedErrors: Record<string, string> =
                Object.fromEntries(
                    Object.entries(fieldErrors).map(([field, messages]) => [
                        field,
                        messages?.[0] ?? "Invalid",
                    ])
                )

            // ✅ throw ครั้งเดียว
            throw new ValidationError(mappedErrors)
        }

        return this.repo.ConfirmTotp(preMfaToken, queryParams, data)
    }
    async handleVefityTOTPAction(
        mfaToken: string,
        queryParams: AuthQueryParameterDTO,
        data: TotpRequestDTO
    ): Promise<MfaResponseDTO> {
        // 1.Validate input (Zod)
        const parsed = ConfirmTOTPSchema.safeParse(data)
        if (!parsed.success) {
            const fieldErrors = parsed.error.flatten().fieldErrors;
            const mappedErrors: Record<string, string> =
                Object.fromEntries(
                    Object.entries(fieldErrors).map(([field, messages]) => [
                        field,
                        messages?.[0] ?? "Invalid",
                    ])
                )

            // ✅ throw ครั้งเดียว
            throw new ValidationError(mappedErrors)
        }

        return this.repo.VerifyTotp(mfaToken, queryParams, data)
    }
}