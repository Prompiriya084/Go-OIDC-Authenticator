import { AuthQueryParameterDTO } from "@/core/dtos/auth_query_parameter_dto";
import { MfaResponseDTO } from "@/core/dtos/mfa_response_dto";
import { SigninRequestDTO } from "@/core/dtos/signin_request_dto";
import { SigninResponseDTO } from "@/core/dtos/signin_response_dto";
import { TotpRequestDTO } from "@/core/dtos/totp_request_dto";
import { HttpPort } from "@/core/ports/http/http_port";
import { AccountApiRepositoryPort } from "@/core/ports/repositories/account_api_repository_port";
import { MfaApiRepositoryPort } from "@/core/ports/repositories/mfa_api_repository_port";

export class MfaApiRepositoryAdapter implements MfaApiRepositoryPort {
    constructor(private http: HttpPort) { }
    async ConfirmTotp(
        preMfaToken: string,
        queryParams: AuthQueryParameterDTO,
        requestDetail: TotpRequestDTO
    ): Promise<MfaResponseDTO> {
        return await this.http.post<MfaResponseDTO>({
            url: "mfa/confirm-totp",
            params: queryParams,
            body: requestDetail,
            headers: {
                Authorization: `Bearer ${preMfaToken}`
            }
        })
    }
    async VerifyTotp(
        mfaToken: string,
        queryParams: AuthQueryParameterDTO,
        requestDetail: TotpRequestDTO
    ): Promise<MfaResponseDTO> {
        return await this.http.post<MfaResponseDTO>({
            url: "mfa/verify-totp",
            params: queryParams,
            body: requestDetail,
            headers: {
                Authorization: `Bearer ${mfaToken}`
            }
        })
    }
}