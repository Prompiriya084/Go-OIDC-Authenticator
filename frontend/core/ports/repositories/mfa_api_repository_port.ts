import { AuthQueryParameterDTO } from "@/core/dtos/auth_query_parameter_dto";
import { TotpRequestDTO } from "@/core/dtos/totp_request_dto";
import { MfaResponseDTO } from "@/core/dtos/mfa_response_dto";

export interface MfaApiRepositoryPort {
    ConfirmTotp(preMfaToken: string, queryParams: AuthQueryParameterDTO, requestDetail: TotpRequestDTO): Promise<MfaResponseDTO>
    VerifyTotp(mfaToken: string, queryParams: AuthQueryParameterDTO, requestDetail: TotpRequestDTO): Promise<MfaResponseDTO>
}