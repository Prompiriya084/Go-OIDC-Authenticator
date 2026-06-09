import { AuthQueryParameterDTO } from "@/core/dtos/auth_query_parameter_dto";
import { SigninRequestDTO } from "@/core/dtos/signin_request_dto";
import { SignInResponseDTO } from "@/core/dtos/signin_response_dto";
import { HttpPort } from "@/core/ports/http/http_port";
import { AccountApiRepositoryPort } from "@/core/ports/repositories/account_api_repository_port";

export class AccountApiRepositoryAdapter implements AccountApiRepositoryPort {
    constructor(private http: HttpPort) { }
    async SignIn(queryParams: AuthQueryParameterDTO, requestDetail: SigninRequestDTO): Promise<SignInResponseDTO> {
        return await this.http.post<SignInResponseDTO>({
            url: "account/signin",
            params: queryParams,
            body: requestDetail
        })
    }
}