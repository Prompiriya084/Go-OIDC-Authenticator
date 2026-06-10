import z from "zod";
import { AccountApiRepositoryPort } from "../ports/repositories/account_api_repository_port";
import { ValidationError } from "../domain/exceptions/validation_error";
import { SigninRequestDTO } from "../dtos/signin_request_dto";
import { AuthQueryParameterDTO } from "../dtos/auth_query_parameter_dto";
import { SigninResponseDTO } from "../dtos/signin_response_dto";

const SignInSchema = z.object({
  username: z.string().min(1, "Please input the employee code."),
  password: z.string().min(1, "Please input the password."),
})

// type validateSignIn = z.infer<typeof SignInSchema>;

export class AccountService {
    constructor(private repo: AccountApiRepositoryPort){

    }
    async handleSignInAction(
        queryParams: AuthQueryParameterDTO,
        data: SigninRequestDTO
    ): Promise<SigninResponseDTO> {
        // 1.Validate input (Zod)
        const parsed = SignInSchema.safeParse(data)
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

        return this.repo.SignIn(queryParams, data)
    }
}