import { AuthQueryParameterDTO } from "@/core/dtos/auth_query_parameter_dto";
import { SigninResponseDTO } from "../../dtos/signin_response_dto";
import { SigninRequestDTO } from "@/core/dtos/signin_request_dto";

export interface AccountApiRepositoryPort {
    SignIn(queryParams: AuthQueryParameterDTO, requestDetail: SigninRequestDTO): Promise<SigninResponseDTO>
}