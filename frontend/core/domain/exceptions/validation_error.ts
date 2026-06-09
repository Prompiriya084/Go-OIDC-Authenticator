export abstract class AppError extends Error {}

export class ValidationError extends AppError{
    constructor(
        public readonly errors: Record<string,string>
    ){
        super("Validation Error")
    }
}